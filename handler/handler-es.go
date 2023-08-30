package handler

import "github.com/elastic/go-elasticsearch/v8"

func InitESClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}
