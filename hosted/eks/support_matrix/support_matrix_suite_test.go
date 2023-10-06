package support_matrix_test

import (
	"testing"

	"github.com/rancher/rancher/tests/framework/extensions/clusters/kubernetesversions"

	"github.com/valaparthvi/highlander-tests/hosted/eks/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	availableVersionList []string
	ctx                  helper.Context
)

func TestSupportMatrix(t *testing.T) {
	RegisterFailHandler(Fail)
	ctx = helper.CommonBeforeSuite()

	var err error
	availableVersionList, err = kubernetesversions.ListEKSAllVersions(ctx.RancherClient)
	Expect(err).To(BeNil())
	Expect(availableVersionList).ToNot(BeEmpty())
	RunSpecs(t, "SupportMatrix Suite")
}
