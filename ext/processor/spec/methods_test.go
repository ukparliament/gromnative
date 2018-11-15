package spec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/ukparliament/gromnative/ext/processor"
  "github.com/wallix/triplestore"
  "io/ioutil"
)

var _ = Describe("Processor", func() {
  Describe("NetTriple", func() {
    Context("object is a BNode", func() {
      It("creates the expected triple", func() {
        expected := processor.Triple{
          Subject: "https://id.parliament.uk/12345678",
          Predicate: "https://id.parliament.uk/shema/PredicateName",
          Object: "_:node39387803",
        }

        triple := triplestore.SubjPredBnode("https://id.parliament.uk/12345678", "https://id.parliament.uk/shema/PredicateName", "node39387803")

        result := processor.NewTriple(triple)

        Expect(result).To(Equal(expected))
      })
    })

    Context("object is a Literal", func() {
      It("creates the expected triple", func() {
        expected := processor.Triple{
          Subject: "https://id.parliament.uk/12345678",
          Predicate: "https://id.parliament.uk/shema/PredicateName",
          Object:  "\"12\"^^<xsd:integer>",
        }

        triple, err := triplestore.SubjPredLit("https://id.parliament.uk/12345678", "https://id.parliament.uk/shema/PredicateName", 12)

        result := processor.NewTriple(triple)

        Expect(result).To(Equal(expected))
        Expect(err).NotTo(HaveOccurred())
      })
    })

    Context("object is a Resource", func() {
      It("creates the expected triple", func() {
        expected := processor.Triple{
          Subject: "https://id.parliament.uk/12345678",
          Predicate: "https://id.parliament.uk/shema/PredicateName",
          Object:  "https://id.parliament.uk/23456789",
        }

        triple := triplestore.SubjPredRes("https://id.parliament.uk/12345678", "https://id.parliament.uk/shema/PredicateName", "https://id.parliament.uk/23456789")

        result := processor.NewTriple(triple)

        Expect(result).To(Equal(expected))
      })
    })
  })

  Describe("Process", func() {
    Context("with no edges", func() {
      It("returns empty edges object", func() {
        fixture, _ := ioutil.ReadFile("../../../spec/fixtures/no_edges.nt")

        res, err := processor.Process(&processor.ProcessorInput{ Body: fixture })

        Expect(len(res.StatementsBySubject)).To(Equal(1))
        Expect(len(res.EdgesBySubject)).To(Equal(0))
        Expect(err).NotTo(HaveOccurred())
      })
    })

    Context("with empty data", func() {
      It("returns empty objects", func() {
        expected := &processor.ProcessorOutput{
          StatementsBySubject: make(map[string][]processor.Triple),
          EdgesBySubject: make(map[string]map[string][]string),
        }

        res, err := processor.Process(&processor.ProcessorInput{ Body: []byte("") })

        Expect(res).To(Equal(expected))
        Expect(err).NotTo(HaveOccurred())
      })
    })

    Context("with full data", func() {
      It("returns the expected objects", func() {
        fixture, _ := ioutil.ReadFile("../../../spec/fixtures/full.nt")

        res, err := processor.Process(&processor.ProcessorInput{ Body: fixture })

        Expect(len(res.StatementsBySubject)).To(Equal(43))
        Expect(len(res.EdgesBySubject)).To(Equal(26))
        Expect(err).NotTo(HaveOccurred())
      })
    })
  })
})