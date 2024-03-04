package utils

import (
	"bytes"
	"fmt"
	"log"
	"mime"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

func BuildProject(repoName string, gitRepoURL string) error {
	cmd := exec.Command("docker", "build", "--build-arg", fmt.Sprintf("GIT_REPOSITORY_URL=%s", gitRepoURL), "-t", repoName, ".")
	cmd.Dir = filepath.Join(".", "docker-files", "golang")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build command execution error: %v\n%s", err, stderr.String())
	}
	return nil
}

func UploadImage(dockerImage string) error {
	// Get login password for ECR Public
	getLoginPasswordCmd := exec.Command("aws", "ecr-public", "get-login-password", "--region", "us-east-1")
	password, err := getLoginPasswordCmd.Output()
	if err != nil {
		return err
	}
	log.Printf("AWS Log in succesful")

	// Login to ECR Public
	dockerLoginCmd := exec.Command("docker", "login", "--username", "AWS", "--password-stdin", "public.ecr.aws/z7u8q2p8")
	dockerLoginCmd.Stdin = bytes.NewReader(password)
	err = dockerLoginCmd.Run()
	if err != nil {
		return err
	}
	log.Println("Docker Log in succesful")
	tagCmdff := exec.Command("docker", "version")
	outt, err := tagCmdff.Output()
	if err != nil {
		return err
	}
	log.Printf("Docker version:\n%s\n", outt)
	log.Println("Docker Log in succesful")
	tagCmdf := exec.Command("docker", "images")
	out, err := tagCmdf.Output()
	if err != nil {
		return err
	}
	log.Printf("Docker Images:\n%s\n", out)
	// Tag Docker image
	tagCmd := exec.Command("docker", "tag", fmt.Sprintf("%s:latest", dockerImage), "public.ecr.aws/z7u8q2p8/nikhil-build-server:latest")
	err = tagCmd.Run()
	if err != nil {
		return err
	}
	log.Println("Tag succesful")

	// Push Docker image to ECR Public
	pushCmd := exec.Command("docker", "push", "public.ecr.aws/z7u8q2p8/nikhil-build-server:latest")
	err = pushCmd.Run()
	if err != nil {
		return err
	}
	log.Println("Push succesful")
	return nil
}

func GetRepoName(gitURL string) (string, error) {
	repoName, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}
	return strings.ToLower(path.Base(repoName.Path)), nil
}
