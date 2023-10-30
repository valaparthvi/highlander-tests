package support_matrix_test

import (
	"testing"

	"github.com/valaparthvi/highlander-tests/hosted/aks/helper"
	"github.com/valaparthvi/highlander-tests/hosted/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	availableVersionList []string
	ctx                  helpers.Context
)

func TestSupportMatrix(t *testing.T) {
	RegisterFailHandler(Fail)
	ctx = helpers.CommonBeforeSuite("aks")
	var err error
	availableVersionList, err = helper.ListSingleVariantAKSAvailableVersions(ctx.RancherClient, ctx.CloudCred.ID, "eastus")
	Expect(err).To(BeNil())
	RunSpecs(t, "SupportMatrix Suite")
}
