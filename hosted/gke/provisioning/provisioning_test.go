package provisioning_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/gke"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"k8s.io/client-go/rest"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
)

var _ = Describe("Provisioning", func() {
	It("should be able to create the cluster, successfully provision it and then delete it", func() {
		clusterName := namegen.AppendRandomString("gkehostcluster")
		var cluster *management.Cluster

		By("creating the cluster", func() {
			var err error
			cluster, err = gke.CreateGKEHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
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

		By("deleting the cluster", func() {
			err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
			Expect(err).To(BeNil())
		})
	})
	Context("Importing a Cluster", func() {
		var (
			clusterName string
			restConfig  *rest.Config
			location    = "us-central1-c"
			project     = "container-project-qa"
		)
		BeforeEach(func() {
			clusterName = namegen.AppendRandomString("imported-gke")
			// TODO: Allow project/zone values to be fetched from a config
			err := helper.CreateGKEClusterUsingGCloud(clusterName, location, project, "", "1", []string{})
			Expect(err).To(BeNil())
			restConfig, err = helper.GetGKEClusterKubeConfigUsingGCloud(clusterName, location, project)
			Expect(err).To(BeNil())
			Expect(restConfig).ToNot(BeNil())
		})
		AfterEach(func() {
			err := helper.DeleteGKEClusterUsingGCloud(clusterName, location, project, false)
			Expect(err).To(BeNil())
		})

		It("should be able to import a cluster, successfully validate it and delete it", func() {
			var cluster *management.Cluster

			By("importing the cluster", func() {
				var err error
				cluster, err = helper.ImportCluster(ctx.RancherClient, clusterName, restConfig)
				Expect(err).To(BeNil())
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

			By("deleting the cluster", func() {
				err := helper.DeleteGKEHostCluster(cluster, ctx.RancherClient)
				Expect(err).To(BeNil())
			})
		})
	})

})
