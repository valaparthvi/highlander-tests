package gke

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/gke"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
	"k8s.io/utils/pointer"
)

var _ = Describe("ProvisioningGke", func() {
	var (
		clusterName string
		ctx         helper.Context
	)
	var _ = BeforeEach(func() {
		clusterName = namegen.AppendRandomString("gkehostcluster")
		ctx = helper.CommonBeforeEach(helper.ContextOpts{})
	})

	var _ = AfterEach(func() {
		helper.CommonAfterEach(ctx)
	})
	When("a cluster is created", func() {
		var cluster *management.Cluster

		BeforeEach(func() {
			var err error
			//TODO: Create GKE cluster only once
			//TODO(contd.): Currently a new GKE cluster is created for every new test, this significantly increases the test time for the entire suite by 6x (it takes 6 minutes for a cluster to be setup)
			cluster, err = gke.CreateGKEHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
		})
		AfterEach(func() {
			err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
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
			// TODO: Programmatically obtain the the version
			var version = pointer.String("1.27.4-gke.900")
			When("the k8s version of the cluster is upgraded", func() {
				BeforeEach(func() {
					cluster, err := helper.UpgradeKubernetesVersion(cluster, version, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})
				It("should have upgraded the cluster's kubernetes version", func() {
					Expect(cluster.GKEConfig.KubernetesVersion).To(BeEquivalentTo(version))
					for _, np := range cluster.GKEConfig.NodePools {
						Expect(np.Version).To(BeEquivalentTo(version))
					}
				})
			})
			//When("the k8s version of nodepools is upgraded", func() {
			//	BeforeEach(func() {
			//		cluster, err := helper.UpgradeNodePoolKubernetesVersion(cluster, version, ctx.RancherClient)
			//		Expect(err).To(BeNil())
			//		err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
			//		Expect(err).To(BeNil())
			//	})
			//	It("should have upgraded the nodepools' kubernetes version", func() {
			//		for _, np := range cluster.GKEConfig.NodePools {
			//			Expect(np.Version).To(BeEquivalentTo(version))
			//		}
			//	})
			//})

		})
		Context("Scaling NodePools", func() {
			var currentNodePoolNumber int
			BeforeEach(func() {
				currentNodePoolNumber = len(cluster.GKEConfig.NodePools)
			})

			When("a nodepool is added", func() {
				BeforeEach(func() {
					cluster, err := helper.AddNodePool(cluster, ctx.RancherClient)
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
					cluster, err := helper.DeleteNodePool(cluster, ctx.RancherClient)
					Expect(err).To(BeNil())
					err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
					Expect(err).To(BeNil())
				})
				It("should have successfully reduced the node pool qty", func() {
					Expect(len(cluster.GKEConfig.NodePools)).To(BeNumerically("<", currentNodePoolNumber))
				})
			})
		})
	})

})
