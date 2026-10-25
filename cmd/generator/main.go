package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/pipeline"

	"github.com/riccap/provider-arubacloud/config"
	"github.com/riccap/provider-arubacloud/internal/genfix"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("root directory is required to be given as argument")
	}
	rootDir := os.Args[1]
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(fmt.Sprintf("cannot calculate the absolute path with %s", rootDir))
	}

	clusterProvider := config.GetProvider()
	namespacedProvider := config.GetProviderNamespaced()

	pipeline.Run(clusterProvider, namespacedProvider, absRootDir)

	// Record the Terraform resource names we generate managed resources for.
	// scripts/version_diff.py reads this list so that the native state schema
	// version check covers exactly the resources this provider ships, and not
	// the ones left out of the include list.
	if err := writeGeneratedList(filepath.Join(absRootDir, "config", "generated.lst"), clusterProvider); err != nil {
		panic(fmt.Sprintf("cannot write the generated resource list: %v", err))
	}

	// Workaround for https://github.com/crossplane/upjet/issues/239: upjet
	// v2.2.0 demotes top-level required fields to
	// +kubebuilder:validation:Optional and substitutes a CEL rule that
	// no-ops when spec.managementPolicies is unset, so kubectl apply
	// silently accepts manifests missing TF-required fields. Promote those
	// back to Required so controller-gen emits the right OpenAPI
	// `required:` arrays in the generated CRDs.
	if err := genfix.PromoteRequiredFields(absRootDir, genfix.Cluster, clusterProvider); err != nil {
		panic(fmt.Sprintf("genfix (cluster): %v", err))
	}
	if err := genfix.PromoteRequiredFields(absRootDir, genfix.Namespaced, namespacedProvider); err != nil {
		panic(fmt.Sprintf("genfix (namespaced): %v", err))
	}
}

// writeGeneratedList writes the Terraform resource names configured on pc to
// path, as a sorted JSON array. The cluster and namespaced providers share an
// include list, so either one yields the same set.
func writeGeneratedList(path string, pc *ujconfig.Provider) error {
	names := make([]string, 0, len(pc.Resources))
	for name := range pc.Resources {
		names = append(names, name)
	}
	sort.Strings(names)

	out, err := json.MarshalIndent(names, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644) //nolint:gosec // generated path
}
