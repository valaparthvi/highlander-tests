package support_matrix_test

import (
	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
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
	availableVersionList, err = helper.ListSingleVariantGKEAvailableVersions(ctx.RancherClient, "container-project-qa", ctx.CloudCred.ID, "", "us-central1")
	Expect(err).To(BeNil())
	RunSpecs(t, "SupportMatrix Suite")
}
