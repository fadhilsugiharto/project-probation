package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"project-probation/model"
	"strings"

	"github.com/nsqio/go-nsq"
)

var nsqProducer *nsq.Producer

func initNSQProducer() (*nsq.Config, error) {
	nsqConfig := nsq.NewConfig()

	var err error
	nsqProducer, err = nsq.NewProducer("127.0.0.1:4150", nsqConfig)
	if err != nil {
		return nil, err
	}

	return nsqConfig, nil
}

func ConsumeNSQMessages() {
	// Initialize NSQ Producer
	nsqConfig, err := initNSQProducer()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize Elasticsearch client
	esClient, err := InitESClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	consumer, err := nsq.NewConsumer("project_probation", "products", nsqConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle NSQ messages
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {

		var product model.Product
		err := json.Unmarshal(message.Body, &product)
		if err != nil {
			log.Println("Error decoding NSQ message:", err)
			return err
		}

		// Store the product in Elasticsearch using esClient
		indexName := "product" // Choose an index name
		jsonBody := fmt.Sprintf(`{
            "id": %d,
            "name": "%s",
            "price": %d
        }`, product.ID, product.Name, product.Price)

		_, err = esClient.Index(indexName, strings.NewReader(jsonBody))
		if err != nil {
			log.Println("Error ES - indexing product:", err)
			return err
		}

		return nil
	}))

	err = consumer.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Fatal(err)
	}

}
