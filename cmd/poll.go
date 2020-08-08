package cmd

import (
	"context"
	"fmt"

	"github.com/amitizle/ghn/internal/scheduler"
	"github.com/amitizle/ghn/pkg/logger"
	"github.com/amitizle/ghn/pkg/notifiers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "poll",
	Short: "start polling for notifications",
	Run:   startScheduler,
}

func init() {
	// startCmd.Flags().Bool("dry-run", false, "don't run the checks, just print that they were supposed to be running")
	rootCmd.AddCommand(startCmd)
}

func startScheduler(cmd *cobra.Command, args []string) {
	if err := initializeNotifiers(); err != nil {
		exitWithError(err)
	}

	s := scheduler.New()

	fn := func() {
		fmt.Println("Amit")
	}
	if err := s.NewTask("* * * * *", fn); err != nil {
		exitWithError(err)
	}
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
