package eks

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/eks"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/valaparthvi/highlander-tests/hosted/eks/helper"
)

var (
	ctx helper.Context
)

var _ = BeforeSuite(func() {
	ctx = helper.CommonBeforeSuite()
})

var _ = Describe("ProvisioningEks", Ordered, func() {
	When("a cluster is created", func() {
		var (
			cluster     *management.Cluster
			clusterName string
		)
		BeforeEach(func() {
			var err error
			clusterName = namegen.AppendRandomString("ekshostcluster")
			cluster, err = eks.CreateEKSHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
		})
		AfterEach(func() {
			err := helper.DeleteEKSHostCluster(cluster, ctx.RancherClient)
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
				versions, err := helper.ListEKSAvailableVersions(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
				Expect(versions).ToNot(BeEmpty())
				upgradeToVersion = &versions[0]
			})

			When(fmt.Sprintf("the k8s version of the cluster is upgraded"), func() {

				It("the k8s version of the cluster is upgraded", func() {
					var err error
					cluster, err = helper.UpgradeClusterKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
					Expect(cluster.EKSConfig.KubernetesVersion).To(BeEquivalentTo(upgradeToVersion))
				})

				It("the k8s version of the cluster nodegroups is upgraded", func() {
					var err error
					cluster, err = helper.UpgradeNodeKubernetesVersion(cluster, upgradeToVersion, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
					for _, ng := range cluster.EKSConfig.NodeGroups {
						Expect(ng.Version).To(BeEquivalentTo(upgradeToVersion))
					}
				})
			})
		})

		Context("Scaling Up/Down NodeGroup", func() {
			var initialNodeCount int64
			BeforeEach(func() {
				initialNodeCount = *cluster.EKSConfig.NodeGroups[0].DesiredSize
			})

			When("a nodegroup is scaled up", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodeGroup(cluster, ctx.RancherClient, initialNodeCount+1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled up the node Group", func() {
					for i := range cluster.EKSConfig.NodeGroups {
						Expect(*cluster.EKSConfig.NodeGroups[i].DesiredSize).To(BeNumerically(">", initialNodeCount))
					}
				})
			})

			When("a nodegroup is scaled down", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.ScaleNodeGroup(cluster, ctx.RancherClient, initialNodeCount-1)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully scaled down the node Group", func() {
					for i := range cluster.EKSConfig.NodeGroups {
						Expect(*cluster.EKSConfig.NodeGroups[i].DesiredSize).To(BeNumerically("<", initialNodeCount))
					}
				})
			})
		})

		Context("Adding/Deleting NodeGroups", func() {
			var currentNodeGroupNumber int
			BeforeEach(func() {
				currentNodeGroupNumber = len(cluster.EKSConfig.NodeGroups)
			})

			When("a nodeGroup is added", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.AddNodeGroup(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully increased the node Group qty", func() {
					Expect(len(cluster.EKSConfig.NodeGroups)).To(BeNumerically(">", currentNodeGroupNumber))
				})
			})

			When("a nodeGroup is deleted", func() {
				BeforeEach(func() {
					var err error
					cluster, err = helper.DeleteNodeGroup(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})

				It("should have successfully reduced the node Group qty", func() {
					Expect(len(cluster.EKSConfig.NodeGroups)).To(BeNumerically("<", currentNodeGroupNumber))
				})
			})
		})
	})

})
