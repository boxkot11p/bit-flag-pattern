package sample

import (
	"context"
	"os"

	"cloud.google.com/go/spanner"
)

func NewClient() (*spanner.Client, error) {
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	ctx := context.Background()
	cli, err := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-database")
	if err != nil {
		return nil, err
	}
	return cli, nil
}
