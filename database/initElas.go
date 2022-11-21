package database

import (
	"fmt"

	"github.com/olivere/elastic/v7"
)

func NewInitEsClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true))

	fmt.Println("ES initialized...")

	return client, err

}
