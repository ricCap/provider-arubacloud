package databasegrant

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("arubacloud_databasegrant", func(r *config.Resource) {
		r.ShortGroup = "databasegrant"
		r.Kind = "DatabaseGrant"
	})
}
