package github

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-github/github"
	"github.com/google/go-querystring/query"
)

// Notification is a struct that represents a single notification
type Notification struct {
	Reason          string
	Repository      string
	RepositoryShort string
	RepositoryURL   string
	Title           string
	Type            string
	URL             string
	UpdatedAt       *time.Time
}

// ListNotificationOpts is a struct that should be passed
// when listing notifications
type ListNotificationOpts struct {
	All           bool
	Participating bool
	Since         time.Time
	Before        time.Time
}

// ListNotifications querying the Github API for users' notifications.
// It receives a "since" timestamp and returns a list of *Notification
func (client *Client) ListNotifications(opts *ListNotificationOpts) ([]*Notification, error) {
	var ghNotifications []*github.Notification
	var notifications []*Notification
	u, err := url.Parse("notifications")
	if err != nil {
		return nil, err
	}
	qs, err := query.Values(opts)
	if err != nil {
		return nil, err
	}
	u.RawQuery = qs.Encode()
	req, err := client.v3Client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	_, err = client.v3Client.Do(context.Background(), req, &ghNotifications)
	if err != nil {
		return nil, err
	}
	for _, notification := range ghNotifications {
		notifications = append(notifications, &Notification{
			Repository:      *notification.Repository.FullName,
			RepositoryShort: fmt.Sprintf("%s/%s", *notification.Repository.Owner, *notification.Repository.Name),
			RepositoryURL:   *notification.Repository.HTMLURL,
			Reason:          *notification.Reason,
			Title:           *notification.Subject.Title,
			Type:            *notification.Subject.Type,
			URL:             *notification.Subject.URL,
			UpdatedAt:       notification.UpdatedAt,
		})
	}
	return notifications, nil
}
