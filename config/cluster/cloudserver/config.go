package cloudserver

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("arubacloud_cloudserver", func(r *config.Resource) {
		r.ShortGroup = "cloudserver"
		r.Kind = "CloudServer"
	})
}
