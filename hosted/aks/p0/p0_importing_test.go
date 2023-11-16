package p0_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/aks"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"

	"github.com/valaparthvi/highlander-tests/hosted/aks/helper"
	"github.com/valaparthvi/highlander-tests/hosted/helpers"
)

var _ = Describe("P0Importing", func() {
	var (
		ctx         helpers.Context
		clusterName string
		location    = "eastus"
		k8sVersion  = "1.26.6"
		increaseBy  = 1
	)
	var _ = BeforeEach(func() {
		clusterName = namegen.AppendRandomString("akshostcluster")
		ctx = helpers.CommonBeforeSuite("aks")
	})
	When("a cluster is imported", func() {
		var cluster *management.Cluster

		BeforeEach(func() {
			var err error
			aksConfig := new(helper.ImportClusterConfig)
			config.LoadAndUpdateConfig(aks.AKSClusterConfigConfigurationFileKey, aksConfig, func() {
				aksConfig.ResourceGroup = clusterName
				aksConfig.ResourceLocation = location
			})
			err = helper.CreateAKSClusterOnAzure(location, clusterName, k8sVersion, "1")
			Expect(err).To(BeNil())
			cluster, err = helper.ImportAKSHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			cluster, err = helpers.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
			// Workaround to add new Nodegroup till https://github.com/rancher/aks-operator/issues/251 is fixed
			cluster.AKSConfig = cluster.AKSStatus.UpstreamSpec
		})
		AfterEach(func() {
			err := helper.DeleteAKSHostCluster(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
			err = helper.DeleteAKSClusteronAzure(clusterName)
			Expect(err).To(BeNil())
		})
		It("should successfully import the cluster & add, delete, scale nodepool", func() {

			By("checking cluster name is same", func() {
				Expect(cluster.Name).To(BeEquivalentTo(clusterName))
			})

			By("checking service account token secret", func() {
				success, err := clusters.CheckServiceAccountTokenSecret(ctx.RancherClient, clusterName)
				Expect(err).To(BeNil())
				Expect(success).To(BeTrue())
			})

			By("checking all management nodes are ready", func() {
				err := nodestat.AllManagementNodeReady(ctx.RancherClient, cluster.ID, helpers.Timeout)
				Expect(err).To(BeNil())
			})

			By("checking all pods are ready", func() {
				podErrors := pods.StatusPods(ctx.RancherClient, cluster.ID)
				Expect(podErrors).To(BeEmpty())
			})

			currentNodePoolNumber := len(cluster.AKSConfig.NodePools)
			initialNodeCount := *cluster.AKSConfig.NodePools[0].Count

			By("scaling up the nodepool", func() {
				var err error
				cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount+1)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				for i := range cluster.AKSConfig.NodePools {
					Expect(*cluster.AKSConfig.NodePools[i].Count).To(BeNumerically("==", initialNodeCount+1))
				}
			})

			By("scaling down the nodepool", func() {
				var err error
				cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				for i := range cluster.AKSConfig.NodePools {
					Expect(*cluster.AKSConfig.NodePools[i].Count).To(BeNumerically("==", initialNodeCount))
				}
			})

			By("adding a nodepool/s", func() {
				var err error
				cluster, err = helper.AddNodePool(cluster, increaseBy, ctx.RancherClient)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(len(cluster.AKSConfig.NodePools)).To(BeNumerically("==", currentNodePoolNumber+1))
			})
			By("deleting the nodepool", func() {
				var err error
				cluster, err = helper.DeleteNodePool(cluster, ctx.RancherClient)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(len(cluster.AKSConfig.NodePools)).To(BeNumerically("==", currentNodePoolNumber))
			})
		})

		Context("Upgrading K8s version", func() {
			var upgradeToVersion, currentVersion *string
			BeforeEach(func() {
				currentVersion = cluster.AKSConfig.KubernetesVersion
				versions, err := helper.ListAKSAvailableVersions(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(versions).ToNot(BeEmpty())
				upgradeToVersion = &versions[0]
			})

			It("should be able to upgrade k8s version of the cluster", func() {
				By("upgrading the ControlPlane", func() {
					var err error
					cluster, err = helper.UpgradeClusterKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					cluster, err = helpers.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					Expect(cluster.AKSConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))
					for _, np := range cluster.AKSConfig.NodePools {
						Expect(np.OrchestratorVersion).To(BeEquivalentTo(currentVersion))
					}
				})

				By("upgrading the NodePools", func() {
					var err error
					cluster, err = helper.UpgradeNodeKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
					Expect(cluster.AKSConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))
					for _, np := range cluster.AKSConfig.NodePools {
						Expect(np.OrchestratorVersion).To(BeEquivalentTo(upgradeToVersion))
					}
				})
			})
		})
	})

})
