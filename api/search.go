package main

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func GetMeilisearch() *meilisearch.Client {
	meiliSearchHost := os.Getenv("MEILISEARCH_HOST")

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: meiliSearchHost,
	})

	taskInfo, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "people",
		PrimaryKey: "id",
	})

	if err != nil {
		panic("Failed to create index: " + err.Error())
	}

	_, err = client.WaitForTask(taskInfo.TaskUID)

	if err != nil {
		panic("Failed to create index: " + err.Error())
	}

	_, err = client.Index("people").UpdateSearchableAttributes(&[]string{"nome", "apelido", "stack"})

	if err != nil {
		panic("Failed define searchable attributes: " + err.Error())
	}

	return client
}
