package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/t3rm1n4l/go-mega"
)

func uploadDirectory(m *mega.Mega, localPath string, megaParent *mega.Node) error {
	// Create folder in MEGA with the same name as the local directory
	folderName := filepath.Base(localPath)
	megaFolder, err := m.CreateDir(folderName, megaParent)
	if err != nil {
		return fmt.Errorf("failed to create MEGA folder: %v", err)
	}

	// Read directory contents
	entries, err := os.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Process each entry
	for _, entry := range entries {
		fullPath := filepath.Join(localPath, entry.Name())

		if entry.IsDir() {
			// Recursively upload subdirectories
			if err := uploadDirectory(m, fullPath, megaFolder); err != nil {
				return err
			}
		} else {
			// Upload files
			fmt.Printf("Uploading file: %s\n", fullPath)
			_, err := m.UploadFile(fullPath, megaFolder, "", nil)
			if err != nil {
				return fmt.Errorf("failed to upload file %s: %v", fullPath, err)
			}
		}
	}
	return nil
}

func main() {
	// Define command line flags
	email := flag.String("email", "", "MEGA account email")
	password := flag.String("password", "", "MEGA account password")
	path := flag.String("path", "", "Path to the file or directory to upload")
	flag.Parse()

	// Validate required flags
	if *email == "" || *password == "" || *path == "" {
		fmt.Println("Usage: megauploader -email <email> -password <password> -path <filepath>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if path exists
	fileInfo, err := os.Stat(*path)
	if os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %s", *path)
	}

	// Create MEGA client
	m := mega.New()

	// Login to MEGA
	err = m.Login(*email, *password)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	// Get root node
	root := m.FS.GetRoot()

	if fileInfo.IsDir() {
		// Upload directory
		fmt.Printf("Uploading directory: %s\n", *path)
		if err := uploadDirectory(m, *path, root); err != nil {
			log.Fatalf("Failed to upload directory: %v", err)
		}
		fmt.Println("Directory uploaded successfully!")
	} else {
		// Upload single file
		fmt.Printf("Uploading file: %s\n", *path)
		node, err := m.UploadFile(*path, root, "", nil)
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}
		fmt.Printf("File uploaded successfully! Node ID: %s\n", node.GetHash())
	}
}
