package spec

import (
  "github.com/ukparliament/gromnative/ext/net"
  . "github.com/ukparliament/gromnative/ext/types/net"
  "gopkg.in/jarcoal/httpmock.v1"
  "net/http"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Net", func() {
  Describe("Get", func() {
    Context("with a valid URI", func() {
      BeforeEach(func() {
        httpmock.RegisterResponder("GET", "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          httpmock.NewStringResponder(200, `<https://id.parliament.uk/43RHonMf> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/Person> .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personGivenName> "Diane" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personOtherNames> "Julie" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personFamilyName> "Abbott" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/oppositionPersonHasOppositionIncumbency> <https://id.parliament.uk/wE8Hq016> .\n\r<https://id.parliament.uk/wE8Hq016> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/OppositionIncumbency> .\n\r<https://id.parliament.uk/wE8Hq016> <https://id.parliament.uk/schema/incumbencyStartDate> "2016-06-27+01:00"^^<http://www.w3.org/2001/XMLSchema#date> .`))
      })

      It("makes the expected request", func() {
        resp, err := net.Get(&GetInput{ Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf" })

        Expect(resp).To(Equal(&GetOutput{
          Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          Body: []byte("<https://id.parliament.uk/43RHonMf> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/Person> .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personGivenName> \"Diane\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personOtherNames> \"Julie\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personFamilyName> \"Abbott\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/oppositionPersonHasOppositionIncumbency> <https://id.parliament.uk/wE8Hq016> .\\n\\r<https://id.parliament.uk/wE8Hq016> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/OppositionIncumbency> .\\n\\r<https://id.parliament.uk/wE8Hq016> <https://id.parliament.uk/schema/incumbencyStartDate> \"2016-06-27+01:00\"^^<http://www.w3.org/2001/XMLSchema#date> ."),
          StatusCode: 200,
        }))
        Expect(err).NotTo(HaveOccurred())
      })

      Context("with headers", func() {
        It("includes the headers in a request", func() {
          httpmock.RegisterResponder(
            "GET",
            "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
            func(req *http.Request) (*http.Response, error) {
              httpHeaders := http.Header{}
              httpHeaders.Add("Foo", "Bar")
              httpHeaders.Add("Bar", "Baz")

              Expect(req.Header).To(Equal(httpHeaders))

              return httpmock.NewStringResponse(200, "done"), nil
            },
          )

          var headers []*GetInput_Header
          headers = append(headers, &GetInput_Header{ Key: "Foo", Value: "Bar" }, &GetInput_Header{ Key: "Bar", Value: "Baz" })
          resp, err := net.Get(&GetInput{ Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf", Headers: headers })

          Expect(resp).To(Equal(&GetOutput{
            Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
            Body: []byte("done"),
            StatusCode: 200,
          }))
          Expect(err).NotTo(HaveOccurred())
        })
      })

      Context("with an invalid URI", func() {
        BeforeEach(func() {
          httpmock.DeactivateAndReset()
        })

        AfterEach(func() {
          httpmock.Activate()
        })

        It("returns the expected values", func() {
          uri := "some_invalid-value.foo"
          errorString := "Get some_invalid-value.foo: unsupported protocol scheme \"\""

          resp, err := net.Get(&GetInput{ Uri: uri })

          Expect(resp).To(Equal(&GetOutput{ Uri: "some_invalid-value.foo", Error: errorString }))

          Expect(err).To(HaveOccurred())
          Expect(err.Error()).To(Equal(errorString))
        })
      })

      Context("with an error making the request", func() {
        PIt("returns the expected response and error", func() {

        })
      })

      Context("with an error reading the response body", func() {
        PIt("returns an error", func() {

        })
      })

      Context("with a non-200 status code", func() {
        PIt("returns an error", func() {

        })
      })
    })

    Context("With an invalid URI", func() {
      PIt("returns with an error", func() {

      })
    })
  })
})