package validation_test

import (
	"testing"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx helper.Context

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	ctx = helper.CommonBeforeSuite()
	RunSpecs(t, "Validation Suite")
}
