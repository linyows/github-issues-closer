package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

const (
	defaultEndpoint = "https://api.github.com"
	searchQuery     = "%s is:issue is:open repo:%s/%s"
)

func main() {
	var (
		endpoint      = flag.String("e", defaultEndpoint, "github api endpoint")
		duplicateOnly = flag.Bool("d", false, "duplicate only")
	)
	flag.Parse()
	args := flag.Args()

	argRepo := args[0]
	if argRepo == "" {
		panic("repo is empty")
	}
	argWord := args[1]
	if argWord == "" {
		panic("search word is empty")
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("token is empty")
	}

	repo := strings.Split(argRepo, "/")

	ctx := context.Background()
	i := New(ctx, &Config{
		endpoint:      *endpoint,
		token:         token,
		owner:         repo[0],
		repo:          repo[1],
		word:          argWord,
		duplicateOnly: *duplicateOnly,
		perpage:       100,
	})
	i.Do(ctx)
}

type Config struct {
	endpoint      string
	token         string
	owner         string
	repo          string
	perpage       int
	word          string
	duplicateOnly bool
}

type IssueCloser struct {
	client *github.Client
	config *Config
}

func New(ctx context.Context, config *Config) *IssueCloser {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.token},
	)
	c := oauth2.NewClient(ctx, ts)
	cl := github.NewClient(c)
	if defaultEndpoint != config.endpoint {
		var err error
		cl, err = github.NewEnterpriseClient(config.endpoint, config.endpoint, c)
		if err != nil {
			panic(err)
		}
	}

	return &IssueCloser{client: cl, config: config}
}

func (ic *IssueCloser) Do(ctx context.Context) {
	page := 1
	total := 0
	q := fmt.Sprintf(searchQuery, ic.config.word, ic.config.owner, ic.config.repo)

	for {
		fmt.Printf("Fetch page %d...", page)
		opts := &github.SearchOptions{
			Sort:  "created",
			Order: "desc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: ic.config.perpage,
			},
		}
		result, _, err := ic.client.Search.Issues(ctx, q, opts)
		if err != nil {
			panic(err)
		}
		total = result.GetTotal()
		fmt.Printf(" Total: %d\n", total)

		if total > 0 && len(result.Issues) == 0 {
			fmt.Printf("Return is strange, Please retry!\n")
			break
		}

		for i, issue := range result.Issues {
			if ic.config.duplicateOnly && i == 0 {
				continue
			}
			ic.Close(ctx, issue.GetNumber())
		}

		if total > ic.config.perpage {
			page += 1
			fmt.Printf("Sleep 5 sec...\n")
			time.Sleep(5 * time.Second)
		} else {
			fmt.Printf("Finish!\n")
			break
		}
	}
}

func (ic *IssueCloser) Close(ctx context.Context, issueNumber int) {
	state := "closed"
	req := &github.IssueRequest{State: &state}
	issue, _, err := ic.client.Issues.Edit(ctx, ic.config.owner, ic.config.repo, issueNumber, req)
	if err != nil {
		fmt.Printf("%#v\n", err)
	} else {
		fmt.Printf("%s\n", issue.GetHTMLURL())
	}
}
