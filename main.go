package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func main() {
	regex, _ := regexp.Compile("([A-Za-z0-9]+\\.)+[a-z]+")
	regex2, _ := regexp.Compile("https")

	var droneGitHttpUrl string = os.Getenv("DRONE_GIT_HTTP_URL")
	var droneRepo string = os.Getenv("DRONE_REPO")
	var gitlabToken string = os.Getenv("PLUGIN_GITLAB_TOKEN")
	var pluginSourceBranch string = os.Getenv("PLUGIN_SOURCE_BRANCH")
	var pluginTargetBranch string = os.Getenv("PLUGIN_TARGET_BRANCH")
	var droneProto string = os.Getenv("DRONE_SYSTEM_PROTO")

	if len(gitlabToken) == 0 {
		log.Fatalf("Not provided gitlabToken")
	}

	if len(pluginSourceBranch) == 0 {
		log.Fatalf("Not provided sourceBranch")
	}

	if len(pluginTargetBranch) == 0 {
		log.Fatalf("Not provided targetBranch")
	}

	var arrayBaseUrl []string = regex.FindAllString(droneGitHttpUrl, 1)

	if len(arrayBaseUrl) == 0 {
		log.Fatalf("Cannot parse git remote url")
	}

	var protoArray []string = regex2.FindAllString(droneGitHttpUrl, 1)

	if len(protoArray) > 0 {
		droneProto = "https"
	}

	var gitlabToken string = os.Getenv("PLUGIN_GITLAB_TOKEN")
	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(baseUrl))


	if err != nil {
		log.Fatalf("Failed to create gitlab client: %v", err)
	}

	listProjectOption := &gitlab.ListProjectsOptions{
		Search: gitlab.String(strings.Split(droneRepo, "/")[1]),
	}

	projects, _, _ := git.Projects.ListProjects(listProjectOption)

	var project *gitlab.Project = nil

	for idx := range projects {
		if projects[idx].PathWithNamespace == droneRepo {
			project = projects[idx]
		}
	}

	var projectId = project.ID

	listProjectMergeRequestOption := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(pluginSourceBranch),
		TargetBranch: gitlab.String(pluginTargetBranch),
	}

	mrList, _, _ := git.MergeRequests.ListProjectMergeRequests(projectId, listProjectMergeRequestOption)

	for _, mr := range mrList {
		acceptMROption := &gitlab.AcceptMergeRequestOptions{
			Squash:                    gitlab.Bool(pluginSquash),
			ShouldRemoveSourceBranch:  gitlab.Bool(removeSourceBranch),
			MergeWhenPipelineSucceeds: gitlab.Bool(false),
		}
		_, res, _ := git.MergeRequests.AcceptMergeRequest(projectId, mr.IID, acceptMROption)
		fmt.Printf("[%d] %s: %d", mr.IID, mr.Title, res.StatusCode)
	}
}
