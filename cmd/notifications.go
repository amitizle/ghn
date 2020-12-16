package cmd

import (
	"context"
	"os"
	"time"

	"github.com/amitizle/ghn/internal/github"
	"github.com/amitizle/ghn/internal/scheduler"
	"github.com/amitizle/ghn/pkg/logger"
	"github.com/amitizle/ghn/pkg/notifiers"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	notificationsRootCmd = &cobra.Command{
		Use:     "notifications",
		Aliases: []string{"notif", "n"},
		Short:   "notifications related commands",
	}

	pollCmd = &cobra.Command{
		Use:   "poll",
		Short: "start polling for notifications",
		Run:   startScheduler,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list notifications one time",
		Run:   listNotifications,
	}
)

func init() {
	// listCmd.Flags().S
	rootCmd.AddCommand(notificationsRootCmd)
	notificationsRootCmd.AddCommand(pollCmd)
	notificationsRootCmd.AddCommand(listCmd)
}

func listNotifications(cmd *cobra.Command, args []string) {
	t, _ := time.Parse(time.RFC3339, "2020-12-01T15:04:05Z")
	validateUserLog := log.With().Logger()
	ctx := context.Background()
	ctxWithLog := logger.StoreContext(ctx, validateUserLog)

	ghClient, err := github.NewClient(ctxWithLog, cfg.Github.Token)
	if err != nil {
		exitWithError(err)
	}
	notifications, err := ghClient.ListNotifications(&github.ListNotificationOpts{
		All:   true,
		Since: t,
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Reason", "Repo", "URL", "Type"})
	for _, notification := range notifications {
		table.Append([]string{
			notification.Title,
			notification.Reason,
			notification.Repository,
			notification.RepositoryURL,
			notification.Type,
		})
	}

	table.Render()

}

// TODO implement
func startScheduler(cmd *cobra.Command, args []string) {
	if err := initializeNotifiers(); err != nil {
		exitWithError(err)
	}

	s := scheduler.New()

	if err := s.Start(); err != nil {
		exitWithError(err)
	}
	select {}
}

func initializeNotifiers() error {
	for _, cfgNotifier := range cfg.Notifiers {
		notifierLogger := log.With().Str("notifier", cfgNotifier.Name).Str("notifier_type", cfgNotifier.Type).Logger()
		notifier, err := notifiers.FromString(cfgNotifier.Type)
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctxWithLog := logger.StoreContext(ctx, notifierLogger)

		if err := notifier.Initialize(ctxWithLog); err != nil {
			return err
		}

		if err := notifier.Configure(cfgNotifier.Config); err != nil {
			return err
		}

		cfgNotifier.Notifier = notifier
	}
	return nil
}
