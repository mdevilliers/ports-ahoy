package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/mdevilliers/ports-ahoy/internal/api"
	"github.com/mdevilliers/ports-ahoy/internal/env"
	"github.com/mdevilliers/ports-ahoy/internal/healthcheck"
	"github.com/mdevilliers/ports-ahoy/internal/store"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

func registerServerCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Runs the service as a HTTP API service",
		RunE: func(cmd *cobra.Command, args []string) error {

			port := env.LookUpWithDefaultStr("PORT", "8000")
			binding := fmt.Sprintf(":%s", port)

			ctx := context.Background()

			client, err := store.NewCacheServiceClientFromEnvironment(ctx)

			if err != nil {
				return errors.Wrap(err, "error creating cache service client")
			}

			store := store.New(client)

			go healthcheck.Start(log)

			apiLogger := log.With().Fields(map[string]interface{}{
				"server": "api",
			}).Logger()

			restful.Add(api.New(apiLogger, store).WebService())

			apiLogger.Info().Fields(map[string]interface{}{
				"binding": binding,
			}).Msg("starting up...")

			srv := &http.Server{Addr: binding}

			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					apiLogger.Fatal().Err(err)
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			<-stop

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				apiLogger.Fatal().Err(err)
				return err
			}

			apiLogger.Info().Msg("api-server closed")

			return nil
		},
	}
	root.AddCommand(cmd)
}
