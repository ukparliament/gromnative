package spec

import (
  "gopkg.in/jarcoal/httpmock.v1"

  "testing"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

func TestNet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Net Suite")
}

var _ = BeforeSuite(func() {
  // block all HTTP requests
  httpmock.Activate()
})

var _ = BeforeEach(func() {
  // remove any mocks
  httpmock.Reset()
})

var _ = AfterSuite(func() {
  httpmock.DeactivateAndReset()
})