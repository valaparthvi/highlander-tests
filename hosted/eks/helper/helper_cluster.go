package helper

import (
	. "github.com/onsi/gomega"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"
	"github.com/rancher/rancher/tests/framework/pkg/wait"

	"github.com/Masterminds/semver/v3"
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

func ListSingleVariantEKSAvailableVersions(client *rancher.Client) (availableVersions []string, err error) {
	availableVersions, err = kubernetesversions.ListEKSAllVersions(client)
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
