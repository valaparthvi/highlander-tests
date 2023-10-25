package support_matrix_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/clusters/eks"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/rancher/rancher/tests/framework/extensions/pipeline"
	"github.com/rancher/rancher/tests/framework/extensions/provisioninginput"
	"github.com/rancher/rancher/tests/framework/extensions/workloads/pods"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"

	"github.com/valaparthvi/highlander-tests/hosted/eks/helper"
)

var _ = Describe("SupportMatrix", func() {

	for _, version := range availableVersionList {
		version := version

		When(fmt.Sprintf("a cluster is created with kubernetes version %s", version), func() {
			var (
				clusterName string
				cluster     *management.Cluster
			)
			BeforeEach(func() {
				clusterName = namegen.AppendRandomString("ekshostcluster")
				//TODO: Rever the changes
				pipeline.UpdateHostedKubernetesVField(provisioninginput.AWSProviderName.String(), version)
				var err error
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
		})
	}
})
