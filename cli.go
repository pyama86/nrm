package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/atotto/clipboard"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

const layout = "2006-01-02T15:04:05"

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		owner string
		repo  string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&owner, "owner", "", "owner")
	flags.StringVar(&owner, "o", "", "owner(Short)")

	flags.StringVar(&repo, "repo", "", "repo")
	flags.StringVar(&repo, "r", "", "repo(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if err := validate(owner, repo, os.Getenv("GITHUB_AUTH_TOKEN")); err != nil {
		logrus.Error(err)
		return ExitCodeError
	}

	data, err := makeScreenShot()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	var client *github.Client
	if os.Getenv("GITHUB_API_ENDPOINT") != "" {
		c, err := github.NewEnterpriseClient(os.Getenv("GITHUB_API_ENDPOINT"), "", tc)
		if err != nil {
			logrus.Fatal(err)
		}
		client = c
	} else {
		client = github.NewClient(tc)
	}

	name := time.Now().Format(layout) + ".png"
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message: github.String(name),
		Content: data,
	}

	createResponse, _, err := client.Repositories.CreateFile(context.Background(), owner, repo, name, repositoryContentsOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	assined := regexp.MustCompile("(.*)/commit/")
	group := assined.FindStringSubmatch(*createResponse.HTMLURL)
	uploadPath := fmt.Sprintf(`![](%s/blob/master/%s?raw=true)`, group[1], url.QueryEscape(name))
	if err := clipboard.WriteAll(uploadPath); err != nil {
		logrus.Fatal(err)
	}
	return ExitCodeOK
}

func makeScreenShot() ([]byte, error) {
	f, err := ioutil.TempFile(os.TempDir(), "nrm")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := exec.Command("/usr/sbin/screencapture", "-i", f.Name()).Run(); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func validate(owner, repo, token string) error {
	if owner == "" {
		return errors.New("owner parameter is required")
	}
	if repo == "" {
		return errors.New("repo parameter is required")
	}
	if token == "" {
		return errors.New("Please set GITHUB_AUTH_TOKEN to environment variable")
	}
	return nil
}
