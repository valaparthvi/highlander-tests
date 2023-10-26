package helper

import (
	"github.com/Masterminds/semver/v3"

	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/pkg/wait"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/defaults"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// WaitUntilClusterIsReady waits until the cluster is in a Ready state,
// fetch the cluster again once it's ready so that it has everything up to date and then return it.
// For e.g. once the cluster has been updated, it contains information such as Version.GitVersion which it does not have before it's ready
func WaitUntilClusterIsReady(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	opts := metav1.ListOptions{FieldSelector: "metadata.name=" + cluster.ID, TimeoutSeconds: &defaults.WatchTimeoutSeconds}
	watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
	if err != nil {
		return nil, err
	}
	watchFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, watchFunc)
	if err != nil {
		return nil, err
	}
	return client.Management.Cluster.ByID(cluster.ID)

}

// UpgradeClusterKubernetesVersion upgrades the k8s version to the value defined by upgradeToVersion.
func UpgradeClusterKubernetesVersion(cluster *management.Cluster, upgradeToVersion *string, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.AKSConfig = cluster.AKSConfig
	upgradedCluster.AKSConfig.KubernetesVersion = upgradeToVersion

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// UpgradeNodeKubernetesVersion upgrades the k8s version of nodepool to the value defined by upgradeToVersion.
func UpgradeNodeKubernetesVersion(cluster *management.Cluster, upgradeToVersion *string, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.AKSConfig = cluster.AKSConfig
	for i := range upgradedCluster.AKSConfig.NodePools {
		upgradedCluster.AKSConfig.NodePools[i].OrchestratorVersion = upgradeToVersion
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// DeleteAKSHostCluster deletes the AKS cluster
func DeleteAKSHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}

func ListSingleVariantAKSAvailableVersions(client *rancher.Client, cloudCredentialID, region string) (availableVersions []string, err error) {
	// TODO: passing a cluster is a temporary workaround until https://github.com/rancher/rancher/pull/43034 is merged
	cluster := &management.Cluster{AKSConfig: &management.AKSClusterConfigSpec{AzureCredentialSecret: cloudCredentialID, ResourceLocation: region}}
	availableVersions, err = kubernetesversions.ListAKSAllVersions(client, cluster)
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

// AddNodePool adds a nodepool to the list
// TODO: Modify this method to add a custom qty of AddNodePool, perhaps by adding an `increaseBy int` arg
func AddNodePool(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.AKSConfig = cluster.AKSConfig

	existingNodepool := cluster.AKSConfig.NodePools[0]
	newNodepool := management.AKSNodePool{
		Count:               existingNodepool.Count,
		VMSize:              existingNodepool.VMSize,
		Mode:                "User",
		Name:                pointer.String(namegen.RandStringLower(5)),
		OrchestratorVersion: existingNodepool.OrchestratorVersion,
	}
	upgradedCluster.AKSConfig.NodePools = append(upgradedCluster.AKSConfig.NodePools, newNodepool)

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// DeleteNodePool deletes a nodepool from the list
// TODO: Modify this method to delete a custom qty of DeleteNodePool, perhaps by adding an `decreaseBy int` arg
func DeleteNodePool(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.AKSConfig = cluster.AKSConfig
	upgradedCluster.AKSConfig.NodePools = cluster.AKSConfig.NodePools[:1]

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
	upgradedCluster.AKSConfig = cluster.AKSConfig
	for i := range upgradedCluster.AKSConfig.NodePools {
		upgradedCluster.AKSConfig.NodePools[i].Count = pointer.Int64(nodeCount)
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// ListAKSAvailableVersions is a function to list and return only available AKS versions for a specific cluster.
func ListAKSAvailableVersions(client *rancher.Client, clusterID string) (availableVersions []string, err error) {
	// kubernetesversions.ListAKSAvailableVersions expects cluster.Version.GitVersion to be available, which it is not sometimes, so we fetch the cluster again to ensure it has all the available data
	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return nil, err
	}
	return kubernetesversions.ListAKSAvailableVersions(client, cluster)
}
