package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"arubacloud_backup":            config.IdentifierFromProvider,
	"arubacloud_blockstorage":      config.IdentifierFromProvider,
	"arubacloud_cloudserver":       config.IdentifierFromProvider,
	"arubacloud_containerregistry": config.IdentifierFromProvider,
	"arubacloud_database":          config.IdentifierFromProvider,
	"arubacloud_databasebackup":    config.IdentifierFromProvider,
	"arubacloud_databasegrant":     config.IdentifierFromProvider,
	"arubacloud_dbaas":             config.IdentifierFromProvider,
	"arubacloud_dbaasuser":         config.IdentifierFromProvider,
	"arubacloud_elasticip":         config.IdentifierFromProvider,
	"arubacloud_kaas":              config.IdentifierFromProvider,
	"arubacloud_keypair":           config.IdentifierFromProvider,
	"arubacloud_kms":               config.IdentifierFromProvider,
	"arubacloud_project":           config.IdentifierFromProvider,
	"arubacloud_restore":           config.IdentifierFromProvider,
	"arubacloud_schedulejob":       config.IdentifierFromProvider,
	"arubacloud_securitygroup":     config.IdentifierFromProvider,
	"arubacloud_securityrule":      config.IdentifierFromProvider,
	"arubacloud_snapshot":          config.IdentifierFromProvider,
	"arubacloud_subnet":            config.IdentifierFromProvider,
	"arubacloud_vpc":               config.IdentifierFromProvider,
	"arubacloud_vpcpeering":        config.IdentifierFromProvider,
	"arubacloud_vpcpeeringroute":   config.IdentifierFromProvider,
	"arubacloud_vpnroute":          config.IdentifierFromProvider,
	"arubacloud_vpntunnel":         config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
