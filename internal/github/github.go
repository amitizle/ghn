package github

import (
	"context"

	"github.com/amitizle/ghn/pkg/logger"
	"github.com/google/go-github/v33/github"
	"github.com/rs/zerolog"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Client is a struct that holds a Github
// v4 (graphql) client
type Client struct {
	v3Client *github.Client // needed for at least user notifications
	client   *githubv4.Client
	log      zerolog.Logger
}

// NewClient returns a Client with the given token
func NewClient(ctx context.Context, token string) (*Client, error) {
	log, err := logger.GetContext(ctx)
	if err != nil {
		return &Client{}, err
	}
	log.Debug().Msg("Initializing Github V4 client")
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	return &Client{
		client:   githubv4.NewClient(httpClient),
		v3Client: github.NewClient(httpClient),
		log:      log,
	}, nil
}

// ValidateToken validates that the token is valid by
// making a user query (see https://docs.github.com/en/graphql/reference/queries#user)
func (client *Client) ValidateToken() error {
	var query struct {
		Viewer struct {
			Login githubv4.String
		}
	}
	err := client.client.Query(context.Background(), &query, nil)

	client.log.Debug().Msgf("user: %s", query.Viewer.Login)
	return err
}
