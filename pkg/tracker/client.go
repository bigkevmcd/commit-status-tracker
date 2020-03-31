package tracker

import (
	"context"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/github"
	"golang.org/x/oauth2"
)

// SCMClientFactory implementations create clients with the
// correct authentication.
//
// TODO: this should accept a URL and use the factory mechanism in go-scm to
// create the correct scm.Client.
type SCMClientFactory func(string) *scm.Client

// TODO: fix this to determine the type of scm Client to create.
func CreateSCMClient(token string) *scm.Client {
	client := github.NewDefault()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client.Client = oauth2.NewClient(context.Background(), ts)
	return client
}
