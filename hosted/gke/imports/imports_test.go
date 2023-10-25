package imports_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"k8s.io/client-go/rest"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
)

var _ = Describe("Imports", func() {
	var (
		clusterName string
		restConfig  *rest.Config
		location    = "us-central1-c"
		project     = "container-project-qa"
	)
	When("a cluster is created in GCloud console and imported to rancher", func() {
		var (
			cluster *management.Cluster
		)
		BeforeEach(func() {
			clusterName = namegen.AppendRandomString("imported-gke")
			// TODO: Allow project/zone values to be fetched from a config
			err := helper.CreateGKEClusterUsingGCloud(clusterName, location, project, "", "1", []string{})
			Expect(err).To(BeNil())
			restConfig, err = helper.GetGKEClusterKubeConfigUsingGCloud(clusterName, location, project)
			Expect(err).To(BeNil())
			Expect(restConfig).ToNot(BeNil())
			cluster, err = helper.ImportCluster(ctx.RancherClient, clusterName, restConfig)
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
			err = helper.DeleteGKEClusterUsingGCloud(clusterName, location, project, false)
			Expect(err).To(BeNil())
		})

		Context("Upgrading K8s version", func() {
			var upgradeToVersion, currentVersion *string
			BeforeEach(func() {
				currentVersion = cluster.GKEConfig.KubernetesVersion
				versions, err := helper.ListGKEAvailableVersions(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(versions).ToNot(BeEmpty())
				upgradeToVersion = &versions[0]
			})

			for _, upgradeNodePool := range []bool{true, false} {
				upgradeNodePool := upgradeNodePool

				When(fmt.Sprintf("the k8s version of the cluster is upgraded; upgradeNodePool=%v", upgradeNodePool), func() {

					BeforeEach(func() {
						var err error
						cluster, err = helper.UpgradeKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient, upgradeNodePool)
						Expect(err).To(BeNil())
						err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
						Expect(err).To(BeNil())
					})

					It("should have upgraded the cluster's kubernetes version", func() {
						Expect(cluster.GKEConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))

						if upgradeNodePool {
							for _, np := range cluster.GKEConfig.NodePools {
								Expect(np.Version).To(BeEquivalentTo(upgradeToVersion))
							}
						} else {
							for _, np := range cluster.GKEConfig.NodePools {
								Expect(np.Version).To(BeEquivalentTo(currentVersion))
							}
						}
					})
				})
			}
		})

		Context("Adding/Deleting NodePools", func() {
			var currentNodePoolNumber int
			BeforeEach(func() {
				currentNodePoolNumber = len(cluster.GKEConfig.NodePools)
			})

			When("a nodepool is added", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.AddNodePool(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				FIt("should have successfully increased the node pool qty", func() {
					Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically(">", currentNodePoolNumber))
				})
			})

			When("a nodepool is deleted", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.DeleteNodePool(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully reduced the node pool qty", func() {
					Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically("<", currentNodePoolNumber))
				})
			})
		})

		Context("Scaling Up/Down NodePool", func() {
			var initialNodeCount int64
			BeforeEach(func() {
				initialNodeCount = *cluster.GKEConfig.NodePools[0].InitialNodeCount
			})

			When("a node is added to the nodepool", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount+1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled up the node pool", func() {
					for i := range cluster.GKEConfig.NodePools {
						Expect(*cluster.GKEConfig.NodePools[i].InitialNodeCount).To(BeNumerically(">", initialNodeCount))
					}
				})
			})

			When("a node is deleted from the nodepool", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount-1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled down the node pool", func() {
					for i := range cluster.GKEConfig.NodePools {
						Expect(*cluster.GKEConfig.NodePools[i].InitialNodeCount).To(BeNumerically("<", initialNodeCount))
					}
				})
			})
		})
	})
})
