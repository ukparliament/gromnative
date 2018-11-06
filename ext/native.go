package main

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/wallix/triplestore"
  "io/ioutil"
  "log"
  "net/http"
  "C"
)

type RequestResponse struct {
  Body        []byte
  StatusCode  int
}

type Response struct {
  StatementsBySubject map[string][]GromTriple         `json:"statementsBySubject"`
  EdgesBySubject      map[string]map[string][]string  `json:"edgesBySubject"`
  StatusCode          int                             `json:"statusCode"`
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

func MakeRequest(httpRequest *http.Request, httpError error, uri string, includeAuth bool, authToken string) (RequestResponse, error) {
  client := &http.Client{}

  if httpError != nil {
    errorMessage := fmt.Sprintf("error creating request object for: %s\n", uri)

    log.Print(errorMessage)

    return RequestResponse{}, errors.New(errorMessage)
  }

  if includeAuth {
    log.Println("Adding auth header")
    authHeaderValue := fmt.Sprintf("Bearer %s", authToken)
    httpRequest.Header.Add("Authorization", authHeaderValue)
  }

  log.Println("Making request")
  resp, err := client.Do(httpRequest)
  if err != nil {
    errorMessage := fmt.Sprintf("error making request to: %s\n", uri)

    fmt.Print(errorMessage)

    defer resp.Body.Close()

    return RequestResponse{}, errors.New(errorMessage)
  }

  log.Printf("Recieved status code: %v\n", resp.StatusCode)

  log.Println("Reading body")
  body, err := ioutil.ReadAll(resp.Body)

  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    errorMessage := fmt.Sprintf("Non-200 Status code. Status code: (%v), body: %s", resp.StatusCode, body)

    log.Println(errorMessage)

    return RequestResponse{}, errors.New(errorMessage)
  }

  requestResponse := RequestResponse {
    Body: body,
    StatusCode: resp.StatusCode,
  }

  return requestResponse, err
}

func GetRequest(uri string, includeAuth bool, authToken string) (RequestResponse, error) {
  // Construct a new Request
  req, err := http.NewRequest("GET", uri, nil)

  return MakeRequest(req, err, uri, includeAuth, authToken)
}

func GetData(uri string) (Response, error) {
  // Placeholder response object
  response := Response { Err: "" }
  // Used to group all statements under a shared subject
  statementsBySubject := make(map[string][]GromTriple)
  // Used to show how one object, through a predicate, links to one or more object
  edgesBySubject := make(map[string]map[string][]string)

  log.Printf("Requesting: %v\n", uri)

  requestResponse, err := GetRequest(uri, false, "")
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
