package cmd

import (
	"context"

	"github.com/amitizle/ghn/internal/github"
	"github.com/amitizle/ghn/pkg/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Github user related commands",
	}

	userValidateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate the user and access based on configured token",
		Run:   validateUser,
	}
)

func init() {
	userCmd.AddCommand(userValidateCmd)
	rootCmd.AddCommand(userCmd)
}

func validateUser(cmd *cobra.Command, args []string) {
	validateUserLog := log.With().Logger()
	ctx := context.Background()
	ctxWithLog := logger.StoreContext(ctx, validateUserLog)

	ghClient, err := github.NewClient(ctxWithLog, cfg.Github.Token)
	if err != nil {
		exitWithError(err)
	}

	if err := ghClient.ValidateToken(); err != nil {
		validateUserLog.Error().Msgf("user could not be validated: (original error: %v)", err)
		exitWithError(err)
	}
}
