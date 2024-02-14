package utils

import (
	"bytes"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func BuildGoProject(outDirPath string) error {
	cmd := exec.Command("go", "build", "-o", "app")
	cmd.Dir = outDirPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build command execution error: %v\n%s", err, stderr.String())
	}
	return nil
}

func CheckBinaryExists(binaryPath string) error {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func UploadToS3(binaryPath string, projectID string, bucketName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return err
	}
	s3Client := s3.New(sess)

	fileContent, err := os.ReadFile(binaryPath)
	if err != nil {
		return err
	}

	contentType := mime.TypeByExtension(filepath.Ext(binaryPath))

	if _, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fmt.Sprintf("__outputs/%s/app", projectID)),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
	}); err != nil {
		return err
	}
	return nil
}
