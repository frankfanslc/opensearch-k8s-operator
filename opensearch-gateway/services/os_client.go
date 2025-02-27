package services

import (
	"context"
	"encoding/json"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"opensearch-k8-operator/opensearch-gateway/responses"
)

type OsClusterClient struct {
	client   *opensearch.Client
	MainPage responses.MainResponse
}

func NewOsClusterClient(config opensearch.Config) (*OsClusterClient, error) {
	service := new(OsClusterClient)
	client, err := opensearch.NewClient(config)
	if err == nil {
		service.client = client
	}
	pingReq := opensearchapi.PingRequest{}
	pingRes, err := pingReq.Do(context.Background(), client)
	if err == nil && pingRes.StatusCode == 200 {
		mainPageResponse, err := MainPage(client)
		if err == nil {
			service.MainPage = mainPageResponse
		}
	}
	return service, err
}

func MainPage(client *opensearch.Client) (responses.MainResponse, error) {
	req := opensearchapi.InfoRequest{}
	infoRes, err := req.Do(context.Background(), client)
	var response responses.MainResponse
	if err == nil {
		defer infoRes.Body.Close()
		err = json.NewDecoder(infoRes.Body).Decode(&response)
	}
	return response, err
}

func (client *OsClusterClient) CatNodes() (responses.CatNodesResponse, error) {
	req := opensearchapi.CatNodesRequest{Format: "json"}
	catNodesRes, err := req.Do(context.Background(), client.client)
	var response responses.CatNodesResponse
	if err == nil {
		defer catNodesRes.Body.Close()
		err = json.NewDecoder(catNodesRes.Body).Decode(&response)
	}
	return response, err
}

func (client *OsClusterClient) NodesStats() (responses.NodesStatsResponse, error) {
	req := opensearchapi.NodesStatsRequest{}
	catNodesRes, err := req.Do(context.Background(), client.client)
	var response responses.NodesStatsResponse
	if err == nil {
		defer catNodesRes.Body.Close()
		err = json.NewDecoder(catNodesRes.Body).Decode(&response)
	}
	return response, err
}

func (client *OsClusterClient) CatIndices() ([]responses.CatIndicesResponse, error) {
	req := opensearchapi.CatIndicesRequest{Format: "json"}
	indicesRes, err := req.Do(context.Background(), client.client)
	var response []responses.CatIndicesResponse
	if err != nil {
		return response, err
	}
	defer indicesRes.Body.Close()
	err = json.NewDecoder(indicesRes.Body).Decode(&response)
	return response, err
}
