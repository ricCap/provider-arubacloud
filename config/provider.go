package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	backupCluster "github.com/riccap/provider-arubacloud/config/cluster/backup"
	blockstorageCluster "github.com/riccap/provider-arubacloud/config/cluster/blockstorage"
	cloudserverCluster "github.com/riccap/provider-arubacloud/config/cluster/cloudserver"
	containerregistryCluster "github.com/riccap/provider-arubacloud/config/cluster/containerregistry"
	databaseCluster "github.com/riccap/provider-arubacloud/config/cluster/database"
	databasebackupCluster "github.com/riccap/provider-arubacloud/config/cluster/databasebackup"
	databasegrantCluster "github.com/riccap/provider-arubacloud/config/cluster/databasegrant"
	dbaasCluster "github.com/riccap/provider-arubacloud/config/cluster/dbaas"
	dbaasuserCluster "github.com/riccap/provider-arubacloud/config/cluster/dbaasuser"
	elasticipCluster "github.com/riccap/provider-arubacloud/config/cluster/elasticip"
	kaasCluster "github.com/riccap/provider-arubacloud/config/cluster/kaas"
	keypairCluster "github.com/riccap/provider-arubacloud/config/cluster/keypair"
	kmsCluster "github.com/riccap/provider-arubacloud/config/cluster/kms"
	projectCluster "github.com/riccap/provider-arubacloud/config/cluster/project"
	restoreCluster "github.com/riccap/provider-arubacloud/config/cluster/restore"
	schedulejobCluster "github.com/riccap/provider-arubacloud/config/cluster/schedulejob"
	securitygroupCluster "github.com/riccap/provider-arubacloud/config/cluster/securitygroup"
	securityruleCluster "github.com/riccap/provider-arubacloud/config/cluster/securityrule"
	snapshotCluster "github.com/riccap/provider-arubacloud/config/cluster/snapshot"
	subnetCluster "github.com/riccap/provider-arubacloud/config/cluster/subnet"
	vpcCluster "github.com/riccap/provider-arubacloud/config/cluster/vpc"
	vpcpeeringCluster "github.com/riccap/provider-arubacloud/config/cluster/vpcpeering"
	vpcpeeringrouteCluster "github.com/riccap/provider-arubacloud/config/cluster/vpcpeeringroute"
	vpnrouteCluster "github.com/riccap/provider-arubacloud/config/cluster/vpnroute"
	vpntunnelCluster "github.com/riccap/provider-arubacloud/config/cluster/vpntunnel"

	backupNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/backup"
	blockstorageNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/blockstorage"
	cloudserverNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/cloudserver"
	containerregistryNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/containerregistry"
	databaseNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/database"
	databasebackupNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/databasebackup"
	databasegrantNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/databasegrant"
	dbaasNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/dbaas"
	dbaasuserNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/dbaasuser"
	elasticipNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/elasticip"
	kaasNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/kaas"
	keypairNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/keypair"
	kmsNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/kms"
	projectNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/project"
	restoreNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/restore"
	schedulejobNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/schedulejob"
	securitygroupNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/securitygroup"
	securityruleNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/securityrule"
	snapshotNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/snapshot"
	subnetNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/subnet"
	vpcNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/vpc"
	vpcpeeringNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/vpcpeering"
	vpcpeeringrouteNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/vpcpeeringroute"
	vpnrouteNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/vpnroute"
	vpntunnelNamespaced "github.com/riccap/provider-arubacloud/config/namespaced/vpntunnel"
)

const (
	resourcePrefix = "arubacloud"
	modulePath     = "github.com/riccap/provider-arubacloud"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("arubacloud.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		backupCluster.Configure,
		blockstorageCluster.Configure,
		cloudserverCluster.Configure,
		containerregistryCluster.Configure,
		databaseCluster.Configure,
		databasebackupCluster.Configure,
		databasegrantCluster.Configure,
		dbaasCluster.Configure,
		dbaasuserCluster.Configure,
		elasticipCluster.Configure,
		kaasCluster.Configure,
		keypairCluster.Configure,
		kmsCluster.Configure,
		projectCluster.Configure,
		restoreCluster.Configure,
		schedulejobCluster.Configure,
		securitygroupCluster.Configure,
		securityruleCluster.Configure,
		snapshotCluster.Configure,
		subnetCluster.Configure,
		vpcCluster.Configure,
		vpcpeeringCluster.Configure,
		vpcpeeringrouteCluster.Configure,
		vpnrouteCluster.Configure,
		vpntunnelCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("arubacloud.m.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		backupNamespaced.Configure,
		blockstorageNamespaced.Configure,
		cloudserverNamespaced.Configure,
		containerregistryNamespaced.Configure,
		databaseNamespaced.Configure,
		databasebackupNamespaced.Configure,
		databasegrantNamespaced.Configure,
		dbaasNamespaced.Configure,
		dbaasuserNamespaced.Configure,
		elasticipNamespaced.Configure,
		kaasNamespaced.Configure,
		keypairNamespaced.Configure,
		kmsNamespaced.Configure,
		projectNamespaced.Configure,
		restoreNamespaced.Configure,
		schedulejobNamespaced.Configure,
		securitygroupNamespaced.Configure,
		securityruleNamespaced.Configure,
		snapshotNamespaced.Configure,
		subnetNamespaced.Configure,
		vpcNamespaced.Configure,
		vpcpeeringNamespaced.Configure,
		vpcpeeringrouteNamespaced.Configure,
		vpnrouteNamespaced.Configure,
		vpntunnelNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
