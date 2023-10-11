package imports_test

import (
	"testing"

	"github.com/valaparthvi/highlander-tests/hosted/gke/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx helper.Context

func TestImports(t *testing.T) {
	RegisterFailHandler(Fail)
	ctx = helper.CommonBeforeSuite()
	RunSpecs(t, "Imports Suite")
}
