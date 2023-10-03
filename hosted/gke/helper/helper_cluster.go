package helper

import (
	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/extensions/defaults"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/wait"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// WaitUntilClusterIsReady waits until the cluster is in a Ready state.
func WaitUntilClusterIsReady(cluster *management.Cluster, client *rancher.Client) {
	opts := metav1.ListOptions{FieldSelector: "metadata.name=" + cluster.ID, TimeoutSeconds: &defaults.WatchTimeoutSeconds}
	watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
	Expect(err).To(BeNil())

	watchFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, watchFunc)
	Expect(err).To(BeNil())
}

// UpgradeKubernetesVersion upgrades the k8s version to the value defined by upgradeToVersion; if upgradeNodePool is true, it also upgrades nodepools' k8s version
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

// DeleteGKEHostCluster deletes the GKE cluster
func DeleteGKEHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}

// AddNodePool adds a nodepool to the list
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

// DeleteNodePool deletes a nodepool from the list
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

// ScaleNodePool modifies the number of initialNodeCount of all the nodepools as defined by nodeCount
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

// ListGKEAvailableVersions is a function to list and return only available GKE versions for a specific cluster.
func ListGKEAvailableVersions(client *rancher.Client, clusterID string) (availableVersions []string, err error) {
	// kubernetesversions.ListGKEAvailableVersions expects cluster.Version.GitVersion to be available, which it is not sometimes, so we fetch the cluster again to ensure it has all the available data
	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return nil, err
	}
	return kubernetesversions.ListGKEAvailableVersions(client, cluster)
}

func ListSingleVariantGKEAvailableVersions(client *rancher.Client, projectID, cloudCredentialID, zone, region string) (availableVersions []string, err error) {
	availableVersions, err = kubernetesversions.ListGKEAllVersions(client, projectID, cloudCredentialID, zone, region)
	if err != nil {
		return nil, err
	}
	var singleVersionList []string
	var oldMinor uint64
	for _, version := range availableVersions {
		semVersion := semver.MustParse(version)
		if currentMinor := semVersion.Minor(); oldMinor != currentMinor {
			singleVersionList = append(singleVersionList, version)
			oldMinor = currentMinor
		}
	}
	return singleVersionList, nil
}
