package net_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Net Suite")
}

var _ = Describe("Net", func() {
	Describe("Get", func() {
		Context("with a valid URI", func() {
			PIt("makes the expected request", func() {

			})

			Context("with headers", func() {
				PIt("includes the headers in a request", func() {

				})
			})

			PContext("with an error in the request", func() {})
		})

		Context("With an invalid URI", func() {
			PIt("returns with an error", func() {

			})
		})
	})
})