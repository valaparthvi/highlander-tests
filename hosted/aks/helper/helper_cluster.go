package helper

import (
	"github.com/Masterminds/semver/v3"

	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/pkg/wait"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	client "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WaitUntilClusterIsReady(cluster *client.Cluster, client *rancher.Client) {
	opts := metav1.ListOptions{FieldSelector: "metadata.name=" + cluster.ID, TimeoutSeconds: &defaults.WatchTimeoutSeconds}
	watchInterface, err := client.GetManagementWatchInterface(management.ClusterType, opts)
	Expect(err).To(BeNil())

	watchFunc := clusters.IsHostedProvisioningClusterReady

	err = wait.WatchWait(watchInterface, watchFunc)
	Expect(err).To(BeNil())
}

// DeleteAKSHostCluster deletes the AKS cluster
func DeleteAKSHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}

func ListSingleVariantAKSAvailableVersions(client *rancher.Client, cloudCredentialID, region string) (availableVersions []string, err error) {
	//TODO: passing a cluster is a temporary workaround until https://github.com/rancher/rancher/pull/43034 is merged
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

// DeleteEKSHostCluster deletes the EKS cluster
func DeleteEKSHostCluster(cluster *management.Cluster, client *rancher.Client) error {
	return client.Management.Cluster.Delete(cluster)
}
