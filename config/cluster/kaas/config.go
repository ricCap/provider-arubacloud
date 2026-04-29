package kaas

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("arubacloud_kaas", func(r *config.Resource) {
		r.ShortGroup = "kaas"
		r.Kind = "KaaS"
	})
}
