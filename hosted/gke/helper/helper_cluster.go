package helper

import (
	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/wait"
	"k8s.io/utils/pointer"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//TODO: Add documentation

func WaitUntilClusterIsReady(cluster *management.Cluster, client *rancher.Client) {
	opts := metav1.ListOptions{FieldSelector: "metadata.name=" + cluster.ID, TimeoutSeconds: &defaults.WatchTimeoutSeconds}
	watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
	Expect(err).To(BeNil())

	watchFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, watchFunc)
	Expect(err).To(BeNil())
}

func UpgradeKubernetesVersion(cluster *management.Cluster, upgradeToVersion *string, client *rancher.Client, upgradeNodePool bool) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.GKEConfig = cluster.GKEConfig
	upgradedCluster.GKEConfig.KubernetesVersion = upgradeToVersion
	if upgradeNodePool {
		for i := range upgradedCluster.GKEConfig.NodePools {
			upgradedCluster.GKEConfig.NodePools[i].Version = upgradeToVersion
		}
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func DeleteGKEHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}

// TODO: Modify this method to add a custom qty of nodepool, perhaps by adding an `increaseBy int` arg
func AddNodePool(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.GKEConfig = cluster.GKEConfig

	existingNodePool := cluster.GKEConfig.NodePools[0]
	newNodePool := management.GKENodePoolConfig{
		Autoscaling:       existingNodePool.Autoscaling,
		Config:            existingNodePool.Config,
		InitialNodeCount:  existingNodePool.InitialNodeCount,
		Management:        existingNodePool.Management,
		MaxPodsConstraint: existingNodePool.MaxPodsConstraint,
		Name:              pointer.String(namegen.AppendRandomString("nodepool")),
		Version:           existingNodePool.Version,
	}
	upgradedCluster.GKEConfig.NodePools = append(upgradedCluster.GKEConfig.NodePools, newNodePool)

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// TODO: Modify this method to delete a custom qty of nodepool, perhaps by adding an `decreaseBy int` arg
func DeleteNodePool(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.GKEConfig = cluster.GKEConfig

	upgradedCluster.GKEConfig.NodePools = cluster.GKEConfig.NodePools[1:]

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func ScaleNodePool(cluster *management.Cluster, client *rancher.Client, nodeCount int64) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.GKEConfig = cluster.GKEConfig
	for i := range upgradedCluster.GKEConfig.NodePools {
		upgradedCluster.GKEConfig.NodePools[i].InitialNodeCount = pointer.Int64(nodeCount)
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}
