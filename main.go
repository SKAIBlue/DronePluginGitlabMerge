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
	var gitlabToken string = os.Getenv("PLUGIN_GITLAB_TOKEN")
	git, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(baseUrl))
}
