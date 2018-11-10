package spec

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ukparliament/gromnative/ext/net"
	. "github.com/ukparliament/gromnative/ext/types/net"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) { return 0, errors.New("test error") }
func (errReader) Close() error                     { return nil }

var _ = Describe("Net", func() {
	Describe("Get", func() {
		Context("with a valid URI", func() {
			BeforeEach(func() {
				httpmock.RegisterResponder("GET", "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
					httpmock.NewStringResponder(200, `<https://id.parliament.uk/43RHonMf> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/Person> .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personGivenName> "Diane" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personOtherNames> "Julie" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personFamilyName> "Abbott" .\n\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/oppositionPersonHasOppositionIncumbency> <https://id.parliament.uk/wE8Hq016> .\n\r<https://id.parliament.uk/wE8Hq016> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/OppositionIncumbency> .\n\r<https://id.parliament.uk/wE8Hq016> <https://id.parliament.uk/schema/incumbencyStartDate> "2016-06-27+01:00"^^<http://www.w3.org/2001/XMLSchema#date> .`))
			})

			It("makes the expected request", func() {
				resp, err := net.Get(&GetInput{Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf"})

				Expect(resp).To(Equal(&GetOutput{
					Uri:        "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
					Body:       []byte("<https://id.parliament.uk/43RHonMf> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/Person> .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personGivenName> \"Diane\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personOtherNames> \"Julie\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/personFamilyName> \"Abbott\" .\\n\\r<https://id.parliament.uk/43RHonMf> <https://id.parliament.uk/schema/oppositionPersonHasOppositionIncumbency> <https://id.parliament.uk/wE8Hq016> .\\n\\r<https://id.parliament.uk/wE8Hq016> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <https://id.parliament.uk/schema/OppositionIncumbency> .\\n\\r<https://id.parliament.uk/wE8Hq016> <https://id.parliament.uk/schema/incumbencyStartDate> \"2016-06-27+01:00\"^^<http://www.w3.org/2001/XMLSchema#date> ."),
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
					headers = append(headers, &GetInput_Header{Key: "Foo", Value: "Bar"}, &GetInput_Header{Key: "Bar", Value: "Baz"})
					resp, err := net.Get(&GetInput{Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf", Headers: headers})

					Expect(resp).To(Equal(&GetOutput{
						Uri:        "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						Body:       []byte("done"),
						StatusCode: 200,
					}))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("with an error making the request", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder(
						"GET",
						"https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						func(req *http.Request) (*http.Response, error) {
							return httpmock.NewStringResponse(200, "done"), errors.New("There was a problem")
						},
					)
				})

				It("returns the expected response and error", func() {
					resp, err := net.Get(&GetInput{Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf"})

					Expect(resp).To(Equal(&GetOutput{
						Uri:   "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						Error: "Get https://api.parliament.uk/query/person_by_id?person_id=43RHonMf: There was a problem",
					}))

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Get https://api.parliament.uk/query/person_by_id?person_id=43RHonMf: There was a problem"))
				})
			})

			Context("with an error reading the response body", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder(
						"GET",
						"https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						func(req *http.Request) (*http.Response, error) {
							return &http.Response{StatusCode: 200, Body: errReader(0)}, nil
						},
					)
				})

				It("returns an error", func() {
					resp, err := net.Get(&GetInput{Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf"})

					Expect(resp).To(Equal(&GetOutput{
						Uri:        "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						StatusCode: int32(200),
						Error:      "Error reading body from https://api.parliament.uk/query/person_by_id?person_id=43RHonMf: test error",
					}))

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("test error"))
				})
			})

			Context("with a non-200 status code", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder(
						"GET",
						"https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						func(req *http.Request) (*http.Response, error) {
							return httpmock.NewStringResponse(500, "Error"), nil
						},
					)
				})

				It("returns an error", func() {
					resp, err := net.Get(&GetInput{Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf"})

					Expect(resp).To(Equal(&GetOutput{
						Uri:        "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
						Body:       []byte("Error"),
						StatusCode: int32(500),
						Error:      "Received 500 status code from https://api.parliament.uk/query/person_by_id?person_id=43RHonMf: Error",
					}))

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Received 500 status code from https://api.parliament.uk/query/person_by_id?person_id=43RHonMf: Error"))
				})
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

				resp, err := net.Get(&GetInput{Uri: uri})

				Expect(resp).To(Equal(&GetOutput{Uri: "some_invalid-value.foo", Error: errorString}))

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(errorString))
			})
		})
	})
})
