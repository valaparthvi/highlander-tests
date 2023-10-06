package support_matrix_test

import (
	"github.com/valaparthvi/highlander-tests/hosted/aks/helper"
	"testing"

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
	availableVersionList, err = helper.ListSingleVariantAKSAvailableVersions(ctx.RancherClient, ctx.CloudCred.ID, "eastus")
	Expect(err).To(BeNil())
	RunSpecs(t, "SupportMatrix Suite")
}
