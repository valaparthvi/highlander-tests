package p0_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/gke"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
	"github.com/valaparthvi/highlander-tests/hosted/helpers"
)

var _ = Describe("P0Importing", func() {
	var (
		clusterName string
		ctx         helpers.Context
		zone        = "us-central1-c"
		project     = "<project>"
		k8sVersion  = "1.26.5-gke.2700"
		increaseBy  = 1
	)
	var _ = BeforeEach(func() {
		clusterName = namegen.AppendRandomString("gkehostcluster")
		ctx = helpers.CommonBeforeSuite("gke")
	})

	When("a cluster is created", func() {
		var cluster *management.Cluster

		BeforeEach(func() {
			var err error
			gkeConfig := new(helper.ImportClusterConfig)
			config.LoadAndUpdateConfig(gke.GKEClusterConfigConfigurationFileKey, gkeConfig, func() {
				gkeConfig.ProjectID = project
			})
			err = helper.CreateGKEClusterOnGCloud(zone, clusterName, project, k8sVersion)
			Expect(err).To(BeNil())
			cluster, err = helper.ImportGKEHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			cluster, err = helpers.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
			// Workaround to add new Nodegroup till https://github.com/rancher/aks-operator/issues/251 is fixed
			cluster.GKEConfig = cluster.GKEStatus.UpstreamSpec
		})
		AfterEach(func() {
			err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
			err = helper.DeleteGKEClusterOnGCloud(zone, clusterName)
			Expect(err).To(BeNil())
		})
		It("should successfully provision the cluster", func() {

			By("checking cluster name is same", func() {
				Expect(cluster.Name).To(BeEquivalentTo(clusterName))
			})

			By("checking service account token secret", func() {
				success, err := clusters.CheckServiceAccountTokenSecret(ctx.RancherClient, clusterName)
				Expect(err).To(BeNil())
				Expect(success).To(BeTrue())
			})

			By("checking all management nodes are ready", func() {
				err := nodestat.AllManagementNodeReady(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
			})

			By("checking all pods are ready", func() {
				podResults, errs := pods.StatusPods(ctx.RancherClient, cluster.ID)
				Expect(errs).To(BeEmpty())
				Expect(podResults).ToNot(BeEmpty())
			})

		})
		Context("Upgrading K8s version", func() {
			var upgradeToVersion *string
			BeforeEach(func() {
				versions, err := helper.ListGKEAvailableVersions(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(versions).ToNot(BeEmpty())
				upgradeToVersion = &versions[0]
			})

			It("should be able to upgrade k8s version of the cluster", func() {
				By("upgrading the Controlplane & NodePools", func() {
					var err error
					cluster, err = helper.UpgradeKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient, true)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())

					Expect(cluster.GKEConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))
					for _, np := range cluster.GKEConfig.NodePools {
						Expect(np.Version).To(BeEquivalentTo(upgradeToVersion))
					}
				})
			})
		})

		It("should be possible to add or delete the nodepools", func() {
			currentNodePoolNumber := len(cluster.GKEConfig.NodePools)

			By("adding a nodepool", func() {
				var err error
				cluster, err = helper.AddNodePool(cluster, increaseBy, ctx.RancherClient)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically("==", currentNodePoolNumber+1))
			})
			By("deleting the nodepool", func() {
				var err error
				cluster, err = helper.DeleteNodePool(cluster, ctx.RancherClient)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically("==", currentNodePoolNumber))

			})

		})

		It("should be possible to scale up/down the nodepool", func() {
			initialNodeCount := *cluster.GKEConfig.NodePools[0].InitialNodeCount

			By("scaling up the nodepool", func() {
				var err error
				cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount+1)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				for i := range cluster.GKEConfig.NodePools {
					Expect(*cluster.GKEConfig.NodePools[i].InitialNodeCount).To(BeNumerically("==", initialNodeCount+1))
				}
			})

			By("scaling down the nodepool", func() {
				var err error
				cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				for i := range cluster.GKEConfig.NodePools {
					Expect(*cluster.GKEConfig.NodePools[i].InitialNodeCount).To(BeNumerically("==", initialNodeCount))
				}
			})
		})
	})
})
