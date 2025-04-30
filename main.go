package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"github.com/k0kubun/go-ansi"
	"github.com/t3rm1n4l/go-mega"
)

func uploadFileWithProgress(m *mega.Mega, filePath string, parent *mega.Node) (*mega.Node, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// Create a progress bar
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("Uploading %s", filepath.Base(filePath))),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Create a progress channel
	progressChan := make(chan int)
	go func() {
		for progress := range progressChan {
			bar.Add(progress)
		}
	}()

	// Upload the file with progress tracking
	node, err := m.UploadFile(filePath, parent, "", &progressChan)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	bar.Finish()
	fmt.Println() 

	return node, nil
}

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
			// Upload files with progress tracking
			_, err := uploadFileWithProgress(m, fullPath, megaFolder)
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
		// Upload single file with progress tracking
		node, err := uploadFileWithProgress(m, *path, root)
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}
		fmt.Printf("File uploaded successfully! Node ID: %s\n", node.GetHash())
	}
}
