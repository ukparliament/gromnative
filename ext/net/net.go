package net

import (
  "errors"
  "fmt"
  . "github.com/ukparliament/gromnative/ext/types/net"
  "io/ioutil"
  "net/http"
)

func Get(input *GetInput) (*GetOutput, error) {
  output := &GetOutput{ Uri: input.Uri }

  // Build a new get request object
  request, err := http.NewRequest("GET", input.Uri, nil)
  if err != nil {
    output.Error = err.Error()
    return output, err
  }

  // Add any header objects to our request
  for i := 0; i < len(input.Headers); i++ {
    request.Header.Add(input.Headers[i].Key, input.Headers[i].Value)
  }

  // Perform our request
  client := http.Client{}
  resp, err := client.Do(request)
  if resp != nil {
    defer resp.Body.Close()
  }

  if err != nil {
    output.Error = err.Error()
    return output, err
  }

  // Store the response code
  if resp.StatusCode != 0 {
    output.StatusCode = int32(resp.StatusCode)
  }

  // Read the body into a []byte
  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    errorMessage := fmt.Sprintf("Error reading body from %v: %v", input.Uri, err)
    output.Error = errorMessage
    return output, err
  }

  output.Body = body

  // Handle non-200 responses
  if resp.StatusCode != 200 {
    errorMessage := fmt.Sprintf("Recieved %v status code from %v: %s", resp.StatusCode, input.Uri, body)

    output.Error = errorMessage
    return output, errors.New(errorMessage)
  }

  return output, err
}
