// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	backup "github.com/riccap/provider-arubacloud/internal/controller/namespaced/backup/backup"
	blockstorage "github.com/riccap/provider-arubacloud/internal/controller/namespaced/blockstorage/blockstorage"
	cloudserver "github.com/riccap/provider-arubacloud/internal/controller/namespaced/cloudserver/cloudserver"
	containerregistry "github.com/riccap/provider-arubacloud/internal/controller/namespaced/containerregistry/containerregistry"
	database "github.com/riccap/provider-arubacloud/internal/controller/namespaced/database/database"
	databasebackup "github.com/riccap/provider-arubacloud/internal/controller/namespaced/databasebackup/databasebackup"
	databasegrant "github.com/riccap/provider-arubacloud/internal/controller/namespaced/databasegrant/databasegrant"
	dbaas "github.com/riccap/provider-arubacloud/internal/controller/namespaced/dbaas/dbaas"
	dbaasuser "github.com/riccap/provider-arubacloud/internal/controller/namespaced/dbaasuser/dbaasuser"
	elasticip "github.com/riccap/provider-arubacloud/internal/controller/namespaced/elasticip/elasticip"
	kaas "github.com/riccap/provider-arubacloud/internal/controller/namespaced/kaas/kaas"
	keypair "github.com/riccap/provider-arubacloud/internal/controller/namespaced/keypair/keypair"
	kms "github.com/riccap/provider-arubacloud/internal/controller/namespaced/kms/kms"
	project "github.com/riccap/provider-arubacloud/internal/controller/namespaced/project/project"
	providerconfig "github.com/riccap/provider-arubacloud/internal/controller/namespaced/providerconfig"
	restore "github.com/riccap/provider-arubacloud/internal/controller/namespaced/restore/restore"
	schedulejob "github.com/riccap/provider-arubacloud/internal/controller/namespaced/schedulejob/schedulejob"
	securitygroup "github.com/riccap/provider-arubacloud/internal/controller/namespaced/securitygroup/securitygroup"
	securityrule "github.com/riccap/provider-arubacloud/internal/controller/namespaced/securityrule/securityrule"
	snapshot "github.com/riccap/provider-arubacloud/internal/controller/namespaced/snapshot/snapshot"
	subnet "github.com/riccap/provider-arubacloud/internal/controller/namespaced/subnet/subnet"
	vpc "github.com/riccap/provider-arubacloud/internal/controller/namespaced/vpc/vpc"
	vpcpeering "github.com/riccap/provider-arubacloud/internal/controller/namespaced/vpcpeering/vpcpeering"
	vpcpeeringroute "github.com/riccap/provider-arubacloud/internal/controller/namespaced/vpcpeeringroute/vpcpeeringroute"
	vpnroute "github.com/riccap/provider-arubacloud/internal/controller/namespaced/vpnroute/vpnroute"
	vpntunnel "github.com/riccap/provider-arubacloud/internal/controller/namespaced/vpntunnel/vpntunnel"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		backup.Setup,
		blockstorage.Setup,
		cloudserver.Setup,
		containerregistry.Setup,
		database.Setup,
		databasebackup.Setup,
		databasegrant.Setup,
		dbaas.Setup,
		dbaasuser.Setup,
		elasticip.Setup,
		kaas.Setup,
		keypair.Setup,
		kms.Setup,
		project.Setup,
		providerconfig.Setup,
		restore.Setup,
		schedulejob.Setup,
		securitygroup.Setup,
		securityrule.Setup,
		snapshot.Setup,
		subnet.Setup,
		vpc.Setup,
		vpcpeering.Setup,
		vpcpeeringroute.Setup,
		vpnroute.Setup,
		vpntunnel.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		backup.SetupGated,
		blockstorage.SetupGated,
		cloudserver.SetupGated,
		containerregistry.SetupGated,
		database.SetupGated,
		databasebackup.SetupGated,
		databasegrant.SetupGated,
		dbaas.SetupGated,
		dbaasuser.SetupGated,
		elasticip.SetupGated,
		kaas.SetupGated,
		keypair.SetupGated,
		kms.SetupGated,
		project.SetupGated,
		providerconfig.SetupGated,
		restore.SetupGated,
		schedulejob.SetupGated,
		securitygroup.SetupGated,
		securityrule.SetupGated,
		snapshot.SetupGated,
		subnet.SetupGated,
		vpc.SetupGated,
		vpcpeering.SetupGated,
		vpcpeeringroute.SetupGated,
		vpnroute.SetupGated,
		vpntunnel.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
