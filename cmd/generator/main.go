package main

import (
	"fmt"
	"os"
	"path/filepath"

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
