// SPDX-FileCopyrightText: 2026 Riccardo Capraro
//
// SPDX-License-Identifier: Apache-2.0

// Package genfix post-processes upjet-generated Go API types to work around
// upjet v2.2.0's unconditional demotion of top-level required fields to
// +kubebuilder:validation:Optional.
//
// Background: pkg/types/builder.go:400 in upjet v2.2.0 strips Required=true
// from every top-level non-identifier parameter and substitutes a CEL
// x-kubernetes-validations rule that no-ops when spec.managementPolicies is
// unset. The net effect on this provider is that kubectl apply silently
// accepts manifests missing TF-required fields; the user only learns about
// the missing fields when the runtime provider webhook fires on first
// reconcile.
//
// PromoteRequiredFields rewrites the Optional marker to Required for each
// top-level parameter whose Terraform schema has Required: true, restoring
// admission-time validation.
//
// Tracked upstream as https://github.com/crossplane/upjet/issues/239 — once
// upjet stops demoting, this package becomes a no-op (the matcher will not
// find any Optional comments to flip) and can be deleted.
package genfix

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Scope mirrors the apis/{cluster,namespaced} directory layout.
type Scope string

const (
	// Cluster targets cluster-scoped CRDs under apis/cluster.
	Cluster Scope = "cluster"
	// Namespaced targets namespaced CRDs under apis/namespaced.
	Namespaced Scope = "namespaced"
)

const (
	optionalMarker = "// +kubebuilder:validation:Optional"
	requiredMarker = "// +kubebuilder:validation:Required"
)

// PromoteRequiredFields walks every resource in provider, identifies the
// top-level <Kind>Parameters struct under
// apis/<scope>/<shortGroup>/<version>/zz_*_types.go, and rewrites the
// Optional kubebuilder marker to Required for each field whose underlying
// Terraform schema has Required: true.
//
// The rewrite is purely textual on the original source bytes — the AST is
// used only to locate comments — so unaffected lines are byte-identical
// after the edit. The function is idempotent: a second invocation finds
// no Optional markers to flip.
func PromoteRequiredFields(absRootDir string, scope Scope, provider *ujconfig.Provider) error {
	for tfName, r := range provider.Resources {
		if r == nil || r.TerraformResource == nil {
			continue
		}

		required := requiredFieldsOf(r)
		if len(required) == 0 {
			continue
		}

		dir := filepath.Join(absRootDir, "apis", string(scope), r.ShortGroup, r.Version)
		structName := r.Kind + "Parameters"

		path, err := findStructFile(dir, structName)
		if err != nil {
			return fmt.Errorf("locate %s for %s: %w", structName, tfName, err)
		}
		if path == "" {
			// No matching file — resource is registered but no generated
			// types exist yet on disk for this scope. Skip silently.
			continue
		}

		if _, err := promoteFieldsInFile(path, structName, required); err != nil {
			return fmt.Errorf("promote required fields in %s (%s): %w", path, tfName, err)
		}
	}
	return nil
}

// requiredFieldsOf collects the snake-case names of every top-level
// schema attribute on r whose Required flag is set. Computed-only and
// Optional fields are skipped.
func requiredFieldsOf(r *ujconfig.Resource) map[string]struct{} {
	out := map[string]struct{}{}
	for name, s := range r.TerraformResource.Schema {
		if s == nil {
			continue
		}
		if s.Required {
			out[name] = struct{}{}
		}
	}
	return out
}

// findStructFile returns the absolute path of the zz_*_types.go file under
// dir whose AST contains a top-level type declaration named structName.
// Returns "" (no error) when no match exists.
func findStructFile(dir, structName string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, "zz_") || !strings.HasSuffix(name, "_types.go") {
			continue
		}
		path := filepath.Join(dir, name)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", path, err)
		}
		if hasStruct(f, structName) {
			return path, nil
		}
	}
	return "", nil
}

func hasStruct(f *ast.File, name string) bool {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if ts.Name.Name != name {
				continue
			}
			if _, ok := ts.Type.(*ast.StructType); ok {
				return true
			}
		}
	}
	return false
}

// edit is a half-open byte range [start, end) inside the source file.
type edit struct{ start, end int }

// fieldTFName returns the snake-case Terraform name encoded in field's
// `tf:"<snake>,..."` struct tag, or "" if the tag is missing/unparseable.
func fieldTFName(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	tagVal, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return ""
	}
	tfTag := reflect.StructTag(tagVal).Get("tf")
	if tfTag == "" {
		return ""
	}
	return strings.SplitN(tfTag, ",", 2)[0]
}

// optionalMarkerEdit returns the byte range of the Optional marker line
// preceding field, or (edit{}, false) if no such line exists.
func optionalMarkerEdit(field *ast.Field, fset *token.FileSet) (edit, bool) {
	if field.Doc == nil {
		return edit{}, false
	}
	for _, c := range field.Doc.List {
		if strings.TrimSpace(c.Text) != optionalMarker {
			continue
		}
		return edit{
			start: fset.Position(c.Slash).Offset,
			end:   fset.Position(c.End()).Offset,
		}, true
	}
	return edit{}, false
}

// collectEdits returns the byte ranges of every Optional marker on a field
// of st whose TF tag is a key in requiredFields.
func collectEdits(st *ast.StructType, fset *token.FileSet, requiredFields map[string]struct{}) []edit {
	var edits []edit
	for _, field := range st.Fields.List {
		snake := fieldTFName(field)
		if snake == "" {
			continue
		}
		if _, ok := requiredFields[snake]; !ok {
			continue
		}
		if e, ok := optionalMarkerEdit(field, fset); ok {
			edits = append(edits, e)
		}
	}
	return edits
}

// applyEdits returns a copy of src with each edit's [start, end) range
// replaced by repl. Edits must be non-overlapping.
func applyEdits(src, repl []byte, edits []edit) []byte {
	// Apply back-to-front so earlier offsets remain valid.
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := append([]byte(nil), src...)
	for _, e := range edits {
		out = append(out[:e.start], append(append([]byte{}, repl...), out[e.end:]...)...)
	}
	return out
}

// promoteFieldsInFile rewrites the Optional kubebuilder marker to Required
// for each field of structName whose `tf:"<snake>,..."` tag matches a key
// in requiredFields. Returns true iff the file was modified.
func promoteFieldsInFile(path, structName string, requiredFields map[string]struct{}) (bool, error) {
	src, err := os.ReadFile(path) //nolint:gosec // generated path under apis/
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("parse: %w", err)
	}
	st := findStruct(file, structName)
	if st == nil {
		return false, nil
	}
	edits := collectEdits(st, fset, requiredFields)
	if len(edits) == 0 {
		return false, nil
	}
	out := applyEdits(src, []byte(requiredMarker), edits)
	if err := os.WriteFile(path, out, 0o644); err != nil { //nolint:gosec // generated path
		return false, err
	}
	return true, nil
}

func findStruct(f *ast.File, name string) *ast.StructType {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != name {
				continue
			}
			if st, ok := ts.Type.(*ast.StructType); ok {
				return st
			}
		}
	}
	return nil
}
