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
  client := &http.Client{}
  resp, err := client.Do(request)

  // Store the response code
  if resp.StatusCode != 0 {
    output.StatusCode = int32(resp.StatusCode)
  }

  // Handle any request errors
  if err != nil {
    defer resp.Body.Close()

    output.Error = err.Error()
    return output, err
  }

  // Read the body into a []byte
  body, err := ioutil.ReadAll(resp.Body)
  output.Body = body

  defer resp.Body.Close()

  // Handle non-200 responses
  if resp.StatusCode != 200 {
    errorMessage := fmt.Sprintf("Recieved %v status code from %v: %s", resp.StatusCode, input.Uri, body)

    output.Error = err.Error()
    return output, errors.New(errorMessage)
  }

  return output, nil
}
