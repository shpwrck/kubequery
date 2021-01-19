/**
 * Copyright (c) 2020-present, The kubequery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Uptycs/kubequery/internal/k8s"
	"github.com/Uptycs/kubequery/internal/k8s/admissionregistration"
	"github.com/Uptycs/kubequery/internal/k8s/apps"
	"github.com/Uptycs/kubequery/internal/k8s/autoscaling"
	"github.com/Uptycs/kubequery/internal/k8s/batch"
	core "github.com/Uptycs/kubequery/internal/k8s/core"
	"github.com/Uptycs/kubequery/internal/k8s/discovery"
	"github.com/Uptycs/kubequery/internal/k8s/networking"
	"github.com/Uptycs/kubequery/internal/k8s/policy"
	"github.com/Uptycs/kubequery/internal/k8s/rbac"
	"github.com/Uptycs/kubequery/internal/k8s/storage"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
)

var (
	socket   = flag.String("socket", "", "Path to the extensions UNIX domain socket")
	timeout  = flag.Int("timeout", 3, "Seconds to wait for autoloaded extensions")
	interval = flag.Int("interval", 3, "Seconds delay between connectivity checks")
)

func main() {
	flag.Parse()
	if *socket == "" {
		panic("Missing required --socket argument")
	}

	err := k8s.Init()
	if err != nil {
		panic(err.Error())
	}

	serverTimeout := osquery.ServerTimeout(
		time.Second * time.Duration(*timeout),
	)
	serverPingInterval := osquery.ServerPingInterval(
		time.Second * time.Duration(*interval),
	)

	// TODO: Version and SDK version
	server, err := osquery.NewExtensionManagerServer(
		"kubequery",
		*socket,
		serverTimeout,
		serverPingInterval,
	)

	if err != nil {
		panic(fmt.Sprintf("Error launching kubequery: %s\n", err))
	}

	// Admission Registration
	server.RegisterPlugin(table.NewPlugin("kubernetes_mutating_webhooks", admissionregistration.MutatingWebhookColumns(), admissionregistration.MutatingWebhooksGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_validating_webhooks", admissionregistration.ValidatingWebhookColumns(), admissionregistration.ValidatingWebhooksGenerate))

	// Apps
	server.RegisterPlugin(table.NewPlugin("kubernetes_daemon_sets", apps.DaemonSetColumns(), apps.DaemonSetsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_daemon_set_containers", apps.DaemonSetContainerColumns(), apps.DaemonSetContainersGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_daemon_set_volumes", apps.DaemonSetVolumeColumns(), apps.DaemonSetVolumesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_deployments", apps.DeploymentColumns(), apps.DeploymentsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_deployments_containers", apps.DeploymentContainerColumns(), apps.DeploymentContainersGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_deployments_volumes", apps.DeploymentVolumeColumns(), apps.DeploymentVolumesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_replica_sets", apps.ReplicaSetColumns(), apps.ReplicaSetsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_replica_set_containers", apps.ReplicaSetContainerColumns(), apps.ReplicaSetContainersGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_replica_set_volumes", apps.ReplicaSetVolumeColumns(), apps.ReplicaSetVolumesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_stateful_sets", apps.StatefulSetColumns(), apps.StatefulSetsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_stateful_set_containers", apps.StatefulSetContainerColumns(), apps.StatefulSetContainersGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_stateful_set_volumes", apps.StatefulSetVolumeColumns(), apps.StatefulSetVolumesGenerate))

	// Autoscaling
	server.RegisterPlugin(table.NewPlugin("kubernetes_horizontal_pod_autoscalers", autoscaling.HorizontalPodAutoscalersColumns(), autoscaling.HorizontalPodAutoscalerGenerate))

	// Batch
	server.RegisterPlugin(table.NewPlugin("kubernetes_cron_jobs", batch.CronJobColumns(), batch.CronJobsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_jobs", batch.JobColumns(), batch.JobsGenerate))

	// Core
	server.RegisterPlugin(table.NewPlugin("kubernetes_component_statuses", core.ComponentStatusColumns(), core.ComponentStatusesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_config_maps", core.ConfigMapColumns(), core.ConfigMapsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_endpoint_subsets", core.EndpointSubsetColumns(), core.EndpointSubsetsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_limit_ranges", core.LimitRangeColumns(), core.LimitRangesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_namespaces", core.NamespaceColumns(), core.NamespacesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_nodes", core.NodeColumns(), core.NodesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_persistent_volume_claims", core.PersistentVolumeClaimColumns(), core.PersistentVolumeClaimsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_persistent_volumes", core.PersistentVolumeColumns(), core.PersistentVolumesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_pod_templates", core.PodTemplateColumns(), core.PodTemplatesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_pods", core.PodColumns(), core.PodsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_pod_containers", core.PodContainerColumns(), core.PodContainersGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_pod_volumes", core.PodVolumeColumns(), core.PodVolumesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_resource_quotas", core.ResourceQuotaColumns(), core.ResourceQuotasGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_secrets", core.SecretColumns(), core.SecretsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_service_accounts", core.ServiceAccountColumns(), core.ServiceAccountsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_services", core.ServiceColumns(), core.ServicesGenerate))

	// Discovery
	server.RegisterPlugin(table.NewPlugin("kubernetes_api_resources", discovery.APIResourceColumns(), discovery.APIResourcesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_info", discovery.InfoColumns(), discovery.InfoGenerate))

	// Networking
	server.RegisterPlugin(table.NewPlugin("kubernetes_ingress_classes", networking.IngressClassColumns(), networking.IngressClassesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_ingresses", networking.IngressColumns(), networking.IngressesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_network_policies", networking.NetworkPolicyColumns(), networking.NetworkPoliciesGenerate))

	// Policy
	server.RegisterPlugin(table.NewPlugin("kubernetes_pod_disruption_budget", policy.PodDisruptionBudgetColumns(), policy.PodDisruptionBudgetsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_pod_security_policies", policy.PodSecurityPolicyColumns(), policy.PodSecurityPoliciesGenerate))

	// RBAC
	server.RegisterPlugin(table.NewPlugin("kubernetes_cluster_role_binding_subjects", rbac.ClusterRoleBindingSubjectColumns(), rbac.ClusterRoleBindingSubjectsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_cluster_role_policy_rule", rbac.ClusterRolePolicyRuleColumns(), rbac.ClusterRolePolicyRulesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_role_binding_subjects", rbac.RoleBindingSubjectColumns(), rbac.RoleBindingSubjectsGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_role_policy_rule", rbac.RolePolicyRuleColumns(), rbac.RolePolicyRulesGenerate))

	// Storage
	server.RegisterPlugin(table.NewPlugin("kubernetes_csi_drivers", storage.CSIDriverColumns(), storage.CSIDriversGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_csi_node_drivers", storage.CSINodeDriverColumns(), storage.CSINodeDriversGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_storage_capacities", storage.CSIStorageCapacityColumns(), storage.CSIStorageCapacitiesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_storage_classes", storage.SGClassColumns(), storage.SGClassesGenerate))
	server.RegisterPlugin(table.NewPlugin("kubernetes_volume_attachments", storage.VolumeAttachmentColumns(), storage.VolumeAttachmentsGenerate))

	if err := server.Run(); err != nil {
		panic(err)
	}
}