package main

import (
  "C"
  "bytes"
  "encoding/json"
  "github.com/ukparliament/gromnative/ext/net"
  . "github.com/ukparliament/gromnative/ext/types/net"
  "github.com/wallix/triplestore"
  "log"
)

type Response struct {
  StatementsBySubject map[string][]GromTriple         `json:"statementsBySubject"`
  EdgesBySubject      map[string]map[string][]string  `json:"edgesBySubject"`
  StatusCode          int32                           `json:"statusCode"`
  Uri                 string                          `json:"uri"`
  Err                 string                          `json:"error"`
}

type GromTriple struct {
  Subject   string  `json:"subject"`
  Predicate string  `json:"predicate"`
  Object    string  `json:"object"`
}

func NewGromTriple(t triplestore.Triple) GromTriple {
  object := ""

  bnode,      isBnode   := t.Object().Bnode()
  literalObj, isLiteral := t.Object().Literal()
  resource,   _         := t.Object().Resource()

  if isBnode {
    object = bnode
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

  return GromTriple {
    Subject: t.Subject(),
    Predicate: t.Predicate(),
    Object: object,
  }
}

func GetData(uri string) (Response, error) {
  // Placeholder response object
  response := Response { Err: "" }
  // Used to group all statements under a shared subject
  statementsBySubject := make(map[string][]GromTriple)
  // Used to show how one object, through a predicate, links to one or more object
  edgesBySubject := make(map[string]map[string][]string)

  log.Printf("Requesting: %v\n", uri)

  requestResponse, err := net.Get(&GetInput{ Uri: uri })
  if requestResponse.StatusCode != 0 {
    response.StatusCode = requestResponse.StatusCode
  }

  if err != nil {
    log.Printf("Error getting: %v\n", err)
    response.Err = err.Error()
    return response, err
  }

  log.Println("Decoding")
  dec := triplestore.NewDatasetDecoder(triplestore.NewLenientNTDecoder, bytes.NewReader(requestResponse.Body))
  tris, err := dec.Decode()
  if err != nil {
    log.Printf("Error decoding: %v\n", err)
    response.Err = err.Error()
    return response, err
  }
  log.Printf("Decoded %v triples", len(tris))


  for i := 0; i < len(tris); i++ {
    triple := tris[i]

    // log.Printf("%v", triple)

    subject := triple.Subject()
    predicate := triple.Predicate()
    object := triple.Object()

    statementsBySubject[subject] = append(statementsBySubject[subject], NewGromTriple(triple))

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
  response.StatementsBySubject = statementsBySubject
  response.EdgesBySubject = edgesBySubject

  log.Println("Done")

  return response, nil
}

//export get
func get(data *C.char) *C.char {
  uri := C.GoString(data)
  errorReturn := C.CString("{}")

  response, err := GetData(uri)
  if err != nil {
    log.Printf("Error getting data: %v\n", err)
    return errorReturn
  }

  responseJson, err := json.Marshal(response)
  if err != nil {
    log.Printf("Error marshalling data: %v\n", err)
    return errorReturn
  }

  return C.CString(string(responseJson))
}

func main() {}
