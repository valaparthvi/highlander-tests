package aks

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/aks"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/valaparthvi/highlander-tests/hosted/aks/helper"
)

var (
	clusterName = namegen.AppendRandomString("akshostcluster")
	ctx         helper.Context
	cluster     *management.Cluster
)

var _ = BeforeSuite(func() {
	ctx = helper.CommonBeforeEach(helper.ContextOpts{})
})

var _ = AfterSuite(func() {
	err := helper.DeleteAKSHostCluster(cluster, ctx.RancherClient)
	Expect(err).To(BeNil())
})

var _ = Describe("ProvisioningAks", func() {

	When("a cluster is created", func() {
		It("should successfully provision the cluster", func() {

			By("provisioning the cluster", func() {
				var err error
				cluster, err = aks.CreateAKSHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
				Expect(err).To(BeNil())
				helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
			})

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
				versions, err := helper.ListAKSAvailableVersions(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(versions).ToNot(BeEmpty())
				upgradeToVersion = &versions[len(versions)-1]
			})

			When(fmt.Sprintf("the k8s version of the cluster is upgraded"), func() {

				It("the k8s version of the cluster is upgraded", func() {
					var err error
					cluster, err = helper.UpgradeClusterKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
					Expect(cluster.AKSConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))
				})

				It("the k8s version of the cluster nodegroups is upgraded", func() {
					var err error
					cluster, err = helper.UpgradeNodeKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
					for _, np := range cluster.AKSConfig.NodePools {
						Expect(np.OrchestratorVersion).To(BeEquivalentTo(upgradeToVersion))
					}
				})
			})
		})

		Context("Scaling Up/Down NodeGroup", func() {
			var initialNodeCount int64
			BeforeEach(func() {
				initialNodeCount = *cluster.AKSConfig.NodePools[0].Count
			})

			When("a nodegroup is scaled up", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodepool(cluster, ctx.RancherClient, initialNodeCount+1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled up the node Group", func() {
					for i := range cluster.AKSConfig.NodePools {
						Expect(*cluster.AKSConfig.NodePools[i].Count).To(BeNumerically(">", initialNodeCount))
					}
				})
			})

			When("a nodegroup is scaled down", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodepool(cluster, ctx.RancherClient, initialNodeCount-1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled down the node Group", func() {
					for i := range cluster.AKSConfig.NodePools {
						Expect(*cluster.AKSConfig.NodePools[i].Count).To(BeNumerically("<", initialNodeCount))
					}
				})
			})
		})

		Context("Adding/Deleting NodeGroups", func() {
			var currentNodeGroupNumber int
			BeforeEach(func() {
				currentNodeGroupNumber = len(cluster.AKSConfig.NodePools)
			})

			When("a nodeGroup is added", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.AddNodepool(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully increased the node Group qty", func() {
					Expect(len(cluster.AKSConfig.NodePools)).To(BeNumerically(">", currentNodeGroupNumber))
				})
			})

			When("a nodeGroup is deleted", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.DeleteNodepool(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully reduced the node Group qty", func() {
					Expect(len(cluster.AKSConfig.NodePools)).To(BeNumerically("<", currentNodeGroupNumber))
				})
			})
		})
	})

})
