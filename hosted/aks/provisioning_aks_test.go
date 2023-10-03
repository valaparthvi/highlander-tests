package aks

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
)

var _ = Describe("ProvisioningAks", func() {
	var (
		clusterName = namegen.AppendRandomString("akshostcluster")
		dnsPrefix   = clusterName + "-dns"
		ctx         helper.Context
	)
	var _ = BeforeEach(func() {
		ctx = helper.CommonBeforeSuite()
	})

	When("a cluster is created", func() {
		var cluster *management.Cluster

		BeforeEach(func() {
			var err error

			aksConfig := new(aks.ClusterConfig)
			config.LoadAndUpdateConfig(aks.AKSClusterConfigConfigurationFileKey, aksConfig, func() {
				aksConfig.ResourceGroup = clusterName
				aksConfig.DNSPrefix = &dnsPrefix
			})

			cluster, err = aks.CreateAKSHostedCluster(ctx.RancherClient, clusterName, ctx.CloudCred.ID, false, false, false, false, map[string]string{})
			Expect(err).To(BeNil())
			helper.WaitUntilClusterIsReady(cluster, ctx.RancherClient)
		})
		//AfterEach(func() {
		//	err := gke.DeleteGKEHostCluster(ctx.RancherClient, cluster)
		//	Expect(err).To(BeNil())
		//})

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
		When("the cluster is upgraded", func() {

		})
	})

})
