package main

import (
	"fmt"
	"log"
	"os"

	utils "github.com/czarnikhil/RepoDeployer.git/internal"
)

func main() {
	programmingLanguage := os.Getenv("LANGUAGE")
	gitURL := os.Getenv("GIT_REPOSITORY_URL")
	repoName, err := utils.GetRepoName(gitURL)
	if err != nil {
		fmt.Errorf("Unable to parse git url :%s", gitURL)
	}
	switch programmingLanguage {
	case "golang":
		if err := utils.BuildProject(repoName, gitURL); err != nil {
			log.Fatalf("Error building Go project: %v", err)
		}
		if err := utils.UploadImage(repoName); err != nil {
			log.Fatalf("Unable to upload image to ECR: %v", err)
		}
	}

}
