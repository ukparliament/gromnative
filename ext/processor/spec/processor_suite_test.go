package spec

import (
  "io/ioutil"
  "log"
  "testing"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

func TestProcessor(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  RegisterFailHandler(Fail)
  RunSpecs(t, "Processor Suite")
}
