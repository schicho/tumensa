package tumensa

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	endpoint = "https://backend.mensen.at/api"
)

type gqlMensaReqVariables struct {
	LocationUri string `json:"locationUri"`
}

type gqlMensaReqBody struct {
	OperationName string               `json:"operationName"`
	Query         string               `json:"query"`
	Variables     gqlMensaReqVariables `json:"variables"`
}

var gqlMensaRequest = gqlMensaReqBody{
	OperationName: "Location",
	Query: `query Location($locationUri: String!) {
  nodeByUri(uri: $locationUri) {
    ... on Location {
      menuplanCurrentWeek
    }
  }
}`,
	Variables: gqlMensaReqVariables{
		LocationUri: "standort/mensa-tu-wien/",
	},
}

func RequestMenuPlan() (*http.Response, error) {
	client := http.DefaultClient
	reqBuf := &bytes.Buffer{}
	encoder := json.NewEncoder(reqBuf)

	err := encoder.Encode(gqlMensaRequest)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(endpoint, "application/json", reqBuf)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New("request failed to graphql endpoint failed")
	}
	return resp, nil
}
