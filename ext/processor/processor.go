package processor

import (
	"bytes"
	"github.com/wallix/triplestore"
	"log"
)

type Triple struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

type ProcessorInput struct {
	Body []byte
}

type ProcessorOutput struct {
	StatementsBySubject map[string][]Triple
	EdgesBySubject      map[string]map[string][]string
	Error               string
}

func NewTriple(t triplestore.Triple) Triple {
	object := ""

	bnode, isBnode := t.Object().Bnode()
	literalObj, isLiteral := t.Object().Literal()
	resource, _ := t.Object().Resource()

	if isBnode {
		object = "_:" + bnode
	} else if isLiteral {
		var literal string

		if literalObj.Lang() != "" {
			literal = "\"" + literalObj.Value() + "\"@" + literalObj.Lang()
		} else {
			literal = "\"" + literalObj.Value() + "\"^^<" + string(literalObj.Type()) + ">"
		}

		object = literal
	} else {
		object = resource
	}

	return Triple{
		Subject:   t.Subject(),
		Predicate: t.Predicate(),
		Object:    object,
	}
}

func Process(input *ProcessorInput) (*ProcessorOutput, error) {
	output := ProcessorOutput{}

	// Used to group all statements under a shared subject
	statementsBySubject := make(map[string][]Triple)
	// Used to show how one object, through a predicate, links to one or more object
	edgesBySubject := make(map[string]map[string][]string)

	log.Println("Decoding")
	dec := triplestore.NewDatasetDecoder(triplestore.NewLenientNTDecoder, bytes.NewReader(input.Body))
	tris, err := dec.Decode()
	if err != nil {
		log.Printf("Error decoding: %v\n", err)
		output.Error = err.Error()
		return &output, err
	}
	log.Printf("Decoded %v triples", len(tris))

	for i := 0; i < len(tris); i++ {
		triple := tris[i]

		subject := triple.Subject()
		predicate := triple.Predicate()
		object := triple.Object()

		statementsBySubject[subject] = append(statementsBySubject[subject], NewTriple(triple))

		// decide if this is an edge
		objectResource, _ := object.Resource()
		if objectResource != "" && predicate != "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			if edgesBySubject[subject] == nil {
				edgesBySubject[subject] = make(map[string][]string)
			}
			edgesBySubject[subject][predicate] = append(edgesBySubject[subject][predicate], objectResource)
		}
	}

	log.Printf("Found %v subjects\n", len(statementsBySubject))

	// Pass our statements and edges back in our response
	output.StatementsBySubject = statementsBySubject
	output.EdgesBySubject = edgesBySubject

	return &output, nil
}
