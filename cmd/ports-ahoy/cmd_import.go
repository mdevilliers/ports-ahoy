package main

import (
	"context"
	"os"

	"github.com/mdevilliers/ports-ahoy/internal/importer"
	"github.com/mdevilliers/ports-ahoy/internal/store"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func registerImportCommand(root *cobra.Command) {

	pathToFile := "./ports.json"

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Imports some data from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			client, err := store.NewCacheServiceClientFromEnvironment(ctx)

			if err != nil {
				return errors.Wrap(err, "error creating cache service client")
			}

			store := store.New(client)

			file, err := os.Open(pathToFile)

			if err != nil {
				return errors.Wrap(err, "error opening import file")
			}

			i, err := importer.New(file, store)

			if err != nil {
				return errors.Wrap(err, "error creating importer")
			}

			err = i.Import()

			if err != nil {
				return errors.Wrap(err, "error importing records")
			}

			log.Info().Msg("import finished")

			return nil
		},
	}

	cmd.Flags().StringVar(&pathToFile, "path", pathToFile, "location of JSON file to import")

	root.AddCommand(cmd)
}
