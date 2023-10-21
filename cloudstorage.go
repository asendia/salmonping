package main

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

func writeToCloudStorage(bucketName string, onjectName string, body []byte) {
	// Create a new Google Cloud Storage client
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	objectName := "your-object-name"

	// Get a handle to the bucket
	bucket := client.Bucket(bucketName)

	// Create a new object in the bucket
	obj := bucket.Object(objectName)

	// Write the body to the object
	wc := obj.NewWriter(context.Background())
	if _, err := wc.Write(body); err != nil {
		log.Fatalf("Failed to write object: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}
}
