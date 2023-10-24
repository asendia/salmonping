package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func writeToCloudStorage(bucketName string, objectName string, body []byte) {
	// Create a new Google Cloud Storage client
	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return
	}

	// Get a handle to the bucket
	bucket := client.Bucket(bucketName)

	// Create a new object in the bucket
	obj := bucket.Object(objectName)

	// Write the body to the object
	wc := obj.NewWriter(context.Background())
	if _, err := wc.Write(body); err != nil {
		fmt.Printf("Failed to write object: %v", err)
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Failed to close writer: %v", err)
		return
	}
}
