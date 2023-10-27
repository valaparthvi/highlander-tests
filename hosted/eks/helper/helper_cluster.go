package helper

import (
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
	upgradedCluster.EKSConfig = cluster.EKSConfig
	upgradedCluster.EKSConfig.KubernetesVersion = upgradeToVersion

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// UpgradeNodeKubernetesVersion upgrades the k8s version of nodegroup to the value defined by upgradeToVersion.
func UpgradeNodeKubernetesVersion(cluster *management.Cluster, upgradeToVersion *string, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.EKSConfig = cluster.EKSConfig
	for i := range upgradedCluster.EKSConfig.NodeGroups {
		upgradedCluster.EKSConfig.NodeGroups[i].Version = upgradeToVersion
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// DeleteEKSHostCluster deletes the EKS cluster
func DeleteEKSHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}

// AddNodeGroup adds a nodegroup to the list
// TODO: Modify this method to add a custom qty of AddNodeGroup, perhaps by adding an `increaseBy int` arg
func AddNodeGroup(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.EKSConfig = cluster.EKSConfig

	existingNodeGroup := cluster.EKSConfig.NodeGroups[0]
	newNodeGroup := management.NodeGroup{
		DesiredSize:   existingNodeGroup.DesiredSize,
		DiskSize:      existingNodeGroup.DiskSize,
		InstanceType:  existingNodeGroup.InstanceType,
		MaxSize:       existingNodeGroup.MaxSize,
		MinSize:       existingNodeGroup.MinSize,
		NodegroupName: pointer.String(namegen.AppendRandomString("nodegroup")),
		Version:       existingNodeGroup.Version,
	}
	upgradedCluster.EKSConfig.NodeGroups = append(upgradedCluster.EKSConfig.NodeGroups, newNodeGroup)

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// DeleteNodeGroup deletes a nodegroup from the list
// TODO: Modify this method to delete a custom qty of DeleteNodeGroup, perhaps by adding an `decreaseBy int` arg
func DeleteNodeGroup(cluster *management.Cluster, client *rancher.Client) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.EKSConfig = cluster.EKSConfig

	upgradedCluster.EKSConfig.NodeGroups = cluster.EKSConfig.NodeGroups[1:]

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// ScaleNodeGroup modifies the number of initialNodeCount of all the nodegroups as defined by nodeCount
func ScaleNodeGroup(cluster *management.Cluster, client *rancher.Client, nodeCount int64) (*management.Cluster, error) {
	upgradedCluster := new(management.Cluster)
	upgradedCluster.Name = cluster.Name
	upgradedCluster.EKSConfig = cluster.EKSConfig
	for i := range upgradedCluster.EKSConfig.NodeGroups {
		upgradedCluster.EKSConfig.NodeGroups[i].DesiredSize = pointer.Int64(nodeCount)
		upgradedCluster.EKSConfig.NodeGroups[i].MaxSize = pointer.Int64(nodeCount)
		upgradedCluster.EKSConfig.NodeGroups[i].MinSize = pointer.Int64(nodeCount)
	}

	cluster, err := client.Management.Cluster.Update(cluster, &upgradedCluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// ListEKSAvailableVersions is a function to list and return only available EKS versions for a specific cluster.
func ListEKSAvailableVersions(client *rancher.Client, clusterID string) (availableVersions []string, err error) {
	// kubernetesversions.ListEKSAvailableVersions expects cluster.Version.GitVersion to be available, which it is not sometimes, so we fetch the cluster again to ensure it has all the available data
	cluster, err := client.Management.Cluster.ByID(clusterID)
	if err != nil {
		return nil, err
	}
	return kubernetesversions.ListEKSAvailableVersions(client, cluster)
}
