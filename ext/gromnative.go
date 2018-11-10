package main

import (
	"C"
	"encoding/json"
	"github.com/ukparliament/gromnative/ext/net"
	"github.com/ukparliament/gromnative/ext/processor"
	. "github.com/ukparliament/gromnative/ext/types/net"
	"log"
)

type Response struct {
	StatementsBySubject map[string][]processor.Triple  `json:"statementsBySubject"`
	EdgesBySubject      map[string]map[string][]string `json:"edgesBySubject"`
	StatusCode          int32                          `json:"statusCode"`
	Uri                 string                         `json:"uri"`
	Err                 string                         `json:"error"`
}

func GetandProcess(uri string) (Response, error) {
	// Placeholder response object
	response := Response{Err: ""}

	log.Printf("Requesting: %v\n", uri)

	requestResponse, err := net.Get(&GetInput{Uri: uri})
	if requestResponse.StatusCode != 0 {
		response.StatusCode = requestResponse.StatusCode
	}

	if err != nil {
		log.Printf("Error getting: %v\n", err)
		response.Err = requestResponse.Error
		return response, err
	}

	processedData, err := processor.Process(&processor.ProcessorInput{Body: requestResponse.Body})
	if err != nil {
		log.Printf("Error processing: %v\n", err)
		response.Err = processedData.Error
		return response, err
	}

	response.StatementsBySubject = processedData.StatementsBySubject
	response.EdgesBySubject = processedData.EdgesBySubject

	log.Println("Done")

	return response, nil
}

//export get
func get(data *C.char) *C.char {
	uri := C.GoString(data)
	errorReturn := C.CString("{}")

	response, err := GetandProcess(uri)
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
