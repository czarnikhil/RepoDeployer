package main

import (
	"log"
	"os"
	"path/filepath"

	utils "github.com/czarnikhil/RepoDeployer.git/internal"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	bucketName := os.Getenv("BUCKET_NAME")
	outDirPath := filepath.Join(".", "output")

	if err := utils.BuildGoProject(outDirPath); err != nil {
		log.Fatalf("Error building Go project: %v", err)
	}

	binaryPath := filepath.Join(outDirPath, "app")
	if err := utils.CheckBinaryExists(binaryPath); err != nil {
		log.Fatalf("Binary file does not exist: %v", err)
	}

	if err := utils.UploadToS3(binaryPath, projectID, bucketName); err != nil {
		log.Fatalf("Error uploading binary file to S3: %v", err)
	}
}
