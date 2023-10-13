package p0_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/gke"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
)

var (
	ctx helper.Context
)

var _ = BeforeSuite(func() {
	ctx = helper.CommonBeforeSuite()
})

var _ = AfterSuite(func() {
	helper.CommonAfterSuite(ctx)
})

var _ = Describe("P0", func() {
	When("a cluster is created", func() {
		var (
			cluster     *management.Cluster
			clusterName string
		)
		BeforeEach(func() {
			clusterName = namegen.AppendRandomString("gkehostcluster")
			var err error
			cluster, err = gke.CreateGKEHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
		})
		AfterEach(func() {
			err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
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

				It("should have successfully increased the node pool qty", func() {
					Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically(">", currentNodePoolNumber))
				})
			})

			When("a nodepool is deleted", func() {
				BeforeEach(func() {
					Expect(currentNodePoolNumber).To(BeNumerically(">", 1))
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
					Expect(initialNodeCount).To(BeNumerically(">", 1))
					var err error
					cluster, err = helper.ScaleNodePool(cluster, ctx.RancherClient, initialNodeCount-1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled down the node pool", func() {
					// TODO: Figure out why the cluster is not deleted
					for i := range cluster.GKEConfig.NodePools {
						Expect(*cluster.GKEConfig.NodePools[i].InitialNodeCount).To(BeNumerically("<", initialNodeCount))
					}
				})
			})
		})
	})
})
