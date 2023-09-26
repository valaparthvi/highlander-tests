package ginkgo_gke_test

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
		clusterName = namegen.AppendRandomString("gkehostcluster")
		ctx         helper.Context
	)
	var _ = BeforeEach(func() {
		ctx = helper.CommonBeforeEach(helper.ContextOpts{})
	})

	var _ = AfterEach(func() {
		helper.CommonAfterEach(ctx)
	})
	When("a cluster is created", func() {
		var cluster *management.Cluster

		BeforeEach(func() {
			var err error
			cluster, err = gke.CreateGKEHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
		})
		//AfterEach(func() {
		//	err := gke.DeleteGKEHostCluster(ctx.RancherClient, cluster)
		//	Expect(err).To(BeNil())
		//})

		FIt("should successfully provision the cluster", func() {

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

		When("the cluster is upgraded", func() {
			var version = pointer.String("1.27.4-gke.900")
			BeforeEach(func() {
				cluster, err := helper.UpgradeKubernetesVersion(cluster, version, ctx.RancherClient)
				Expect(err).To(BeNil())
				err = clusters.WaitClusterToBeUpgraded(ctx.RancherClient, cluster.ID)
				Expect(err).To(BeNil())
			})
			It("should have upgraded the cluster's kubernetes version", func() {
				Expect(cluster.GKEConfig.KubernetesVersion).To(BeEquivalentTo(version))
			})
		})
	})

})
