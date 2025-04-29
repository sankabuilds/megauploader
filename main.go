package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/t3rm1n4l/go-mega"
)

func main() {
	// Define command line flags
	email := flag.String("email", "", "MEGA account email")
	password := flag.String("password", "", "MEGA account password")
	filePath := flag.String("file", "", "Path to the file to upload")
	flag.Parse()

	// Validate required flags
	if *email == "" || *password == "" || *filePath == "" {
		fmt.Println("Usage: megauploader -email <email> -password <password> -file <filepath>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *filePath)
	}

	// Create MEGA client
	m := mega.New()

	// Login to MEGA
	err := m.Login(*email, *password)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	// Get root node
	root := m.FS.GetRoot()

	// Upload file
	fmt.Printf("Uploading file: %s\n", *filePath)
	node, err := m.UploadFile(*filePath, root, "", nil)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	fmt.Printf("File uploaded successfully! Node ID: %s\n", node.GetHash())
}
