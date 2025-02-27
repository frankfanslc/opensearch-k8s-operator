package services

import (
	"context"
	"crypto/tls"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func getClusterClient(t *testing.T) *OsClusterClient {
	// Initialize the client with SSL/TLS enabled.
	config := opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://localhost:9200"},
		Username:  "admin", // For testing only. Don't store credentials in code.
		Password:  "admin",
	}

	clusterClient, err := NewOsClusterClient(config)
	assert.Nil(t, err, "failed connection to cluster")
	return clusterClient
}

func createIndex(t *testing.T, clusterClient *OsClusterClient, indexName string, mapping *strings.Reader) {
	req := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  mapping,
	}
	createIndexRes, err := req.Do(context.Background(), clusterClient.client)
	assert.Nil(t, err, "failed to create index")
	assert.Equal(t, createIndexRes.StatusCode, 200)
}

func deleteIndex(clusterClient *OsClusterClient, indexName string) {
	opensearchapi.IndicesDeleteRequest{Index: []string{indexName}}.Do(context.Background(), clusterClient.client)
}
