package helper

import (
	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/pkg/wait"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WaitUntilClusterIsReady(cluster *management.Cluster, client *rancher.Client) {
	opts := metav1.ListOptions{FieldSelector: "metadata.name=" + cluster.ID, TimeoutSeconds: &defaults.WatchTimeoutSeconds}
	watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
	Expect(err).To(BeNil())

	watchFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, watchFunc)
	Expect(err).To(BeNil())
}

func UpgradeKubernetesVersion(cluster *management.Cluster, upgradeToVersion *string, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := cluster
	upgradedCluster.GKEConfig.KubernetesVersion = upgradeToVersion
	cluster, err := client.Management.Cluster.Update(cluster, upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}
