package main

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/ukparliament/gromnative/ext/processor"
  "gopkg.in/jarcoal/httpmock.v1"
  "io/ioutil"
  "log"
  "testing"
)

func TestProcessor(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  RegisterFailHandler(Fail)
  RunSpecs(t, "gromnative Suite")
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

var _ = Describe("gromnative", func() {
  Describe("GetandProcess", func() {
    Context("with expected data", func() {
      BeforeEach(func() {
        fixture, _ := ioutil.ReadFile("../spec/fixtures/one_edge.nt")

        httpmock.RegisterResponder("GET", "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          httpmock.NewStringResponder(200, string(fixture)))
      })

      It("returns the expected data", func() {
        statementsBySubject := make(map[string][]processor.Triple)
        statementsBySubject["https://id.parliament.uk/43RHonMf"] = append(
          statementsBySubject["https://id.parliament.uk/43RHonMf"],
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
            Object: "https://id.parliament.uk/schema/Person",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "https://id.parliament.uk/schema/personGivenName",
            Object: "\"Diane\"^^<xsd:string>",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "https://id.parliament.uk/schema/personOtherNames",
            Object: "\"Julie\"^^<xsd:string>",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "https://id.parliament.uk/schema/personFamilyName",
            Object: "\"Abbott\"^^<xsd:string>",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "http://example.com/F31CBD81AD8343898B49DC65743F0BDF",
            Object: "\"Ms Diane Abbott\"^^<xsd:string>",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "http://example.com/D79B0BAC513C4A9A87C9D5AFF1FC632F",
            Object: "\"Rt Hon Diane Abbott MP\"^^<xsd:string>",
          },
          processor.Triple {
            Subject: "https://id.parliament.uk/43RHonMf",
            Predicate: "https://id.parliament.uk/schema/Test",
            Object: "https://id.parliament.uk/12345678",
          })
        edgesBySubject := make(map[string]map[string][]string)
        edgesBySubject["https://id.parliament.uk/43RHonMf"] = make(map[string][]string)
        edgesBySubject["https://id.parliament.uk/43RHonMf"]["https://id.parliament.uk/schema/Test"] = append(
          edgesBySubject["https://id.parliament.uk/43RHonMf"]["https://id.parliament.uk/schema/Test"],
          "https://id.parliament.uk/12345678")

        expected := Response{
          StatementsBySubject: statementsBySubject,
          EdgesBySubject: edgesBySubject,
          StatusCode: 200,
          Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          Err: "",
        }

        res, err := GetandProcess("https://api.parliament.uk/query/person_by_id?person_id=43RHonMf")

        Expect(res).To(Equal(expected))
        Expect(err).NotTo(HaveOccurred())
      })
    })

    Context("with an error getting", func() {
      BeforeEach(func() {
        httpmock.DeactivateAndReset()
      })

      AfterEach(func() {
        httpmock.Activate()
      })

      It("returns the expected data", func() {
        expected := Response{
          StatementsBySubject: nil,
          EdgesBySubject: nil,
          StatusCode: 0,
          Uri: "foo://a_broken.url",
          Err: "Get foo://a_broken.url: unsupported protocol scheme \"foo\"",
        }

        res, err := GetandProcess("foo://a_broken.url")

        Expect(res).To(Equal(expected))
        Expect(err.Error()).To(Equal("Get foo://a_broken.url: unsupported protocol scheme \"foo\""))
      })
    })

    Context("with an error processing", func() {
      BeforeEach(func() {
        httpmock.RegisterResponder("GET", "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          httpmock.NewStringResponder(200, "{\"error\":\"Definitely not Triples\"}"))
      })

      It("returns the expected data", func() {
        expected := Response{
          StatementsBySubject: nil,
          EdgesBySubject: nil,
          StatusCode: 200,
          Uri: "https://api.parliament.uk/query/person_by_id?person_id=43RHonMf",
          Err: "lenient parsing: line 1: invalid subject in {\"error\":\"Definitely not Triples\"}",
        }

        res, err := GetandProcess("https://api.parliament.uk/query/person_by_id?person_id=43RHonMf")

        Expect(res).To(Equal(expected))
        Expect(err.Error()).To(Equal("lenient parsing: line 1: invalid subject in {\"error\":\"Definitely not Triples\"}"))
      })
    })
  })
})
