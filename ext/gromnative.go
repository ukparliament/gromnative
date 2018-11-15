package main

import (
	"C"
	"encoding/json"
	"fmt"
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
	response := Response{ Uri: uri }

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

func cStringConversion(response *Response) *C.char {
	json, err := json.Marshal(response)
	if err != nil {
		return C.CString("{\"error\": \"error creating json for ruby\"}")
	}

	return C.CString(string(json))
}

//export get
func get(data *C.char) *C.char {
	uri := C.GoString(data)
	response, err := GetandProcess(uri)
	if err != nil {
		errorResponse := &Response{ Err: fmt.Sprintf("Error getting data: %v\n", err) }
		log.Println(errorResponse.Err)
		return cStringConversion(errorResponse)
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		errorResponse := &Response{ Err: fmt.Sprintf("Error marshalling data: %v\n", err) }
		log.Println(errorResponse.Err)
		return cStringConversion(errorResponse)
	}

	return C.CString(string(responseJson))
}

func main() {}
