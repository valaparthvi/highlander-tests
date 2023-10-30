package support_matrix_test

import (
	"testing"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"
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
	ctx = helpers.CommonBeforeSuite("gke")
	var err error
	availableVersionList, err = helper.ListSingleVariantGKEAvailableVersions(ctx.RancherClient, "container-project-qa", ctx.CloudCred.ID, "", "us-central1")
	Expect(err).To(BeNil())
	RunSpecs(t, "SupportMatrix Suite")
}
