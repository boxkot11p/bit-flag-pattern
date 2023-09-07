package sample_test

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/boxkot11p/big-flag-pattern/sample"
	"github.com/google/go-cmp/cmp"
)

func Test_Sample(t *testing.T) {
	client, err := sample.NewClient()
	if err != nil {
		t.Fatalf("failed new spanner client: %v", err)
	}

	const userID = "sample-id"

	type userData struct {
		userID string
		name   string
		entranceFlag int64		
	}
	tests := map[string]struct {
		data     *userData
		reqFlag sample.EntranceFlag
		wantRes  bool
	}{
		"sample-1": {
			data: &userData{
				userID: userID, 
				name: "sample-name",
				entranceFlag: sample.MergeFlag([]sample.EntranceFlag{sample.EntranceFlag_NORMAL}),
			},
			reqFlag: sample.EntranceFlag_NORMAL,
			wantRes: true,
		},
		"sample-2": {
			data: &userData{
				userID: userID, 
				name: "sample-name",
				entranceFlag: sample.MergeFlag([]sample.EntranceFlag{
					sample.EntranceFlag_NORMAL,
					sample.EntranceFlag_PREMIUM,
				}),
			},
			reqFlag: sample.EntranceFlag_NORMAL,
			wantRes: true,
		},
		"sample-3": {
			data: &userData{
				userID: userID, 
				name: "sample-name",
				entranceFlag: sample.MergeFlag([]sample.EntranceFlag{
					sample.EntranceFlag_NORMAL,
					sample.EntranceFlag_PREMIUM,
				}),
			},
			reqFlag: sample.EntranceFlag_PREMIUM,
			wantRes: true,
		},
		"sample-4": {
			data: &userData{
				userID: userID, 
				name: "sample-name",
				entranceFlag: sample.MergeFlag([]sample.EntranceFlag{
					sample.EntranceFlag_NORMAL,
					sample.EntranceFlag_PREMIUM,
				}),
			},
			reqFlag: sample.EntranceFlag_SPECIAL,
			wantRes: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				ctx := context.Background()
				m := spanner.Delete("Users", spanner.Key{userID})
				if _, err := client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
					log.Fatal(err)
				}
			})

			ctx := context.Background()
			if _, err := client.Apply(ctx, []*spanner.Mutation{
				spanner.InsertOrUpdate(
					"Users", 
					[]string{"UserId", "Name", "EntranceFlag"}, 
					[]interface{}{tt.data.userID, tt.data.name, tt.data.entranceFlag},
				),
			}); err != nil {
				t.Fatalf("failed to apply mutation: %v", err)
			}

			stmt := spanner.Statement{SQL: "SELECT * FROM Users"}
			iter := client.Single().Query(ctx, stmt)
			defer iter.Stop()
			row, err := iter.Next()
			if err != nil {
				t.Fatalf("failed select users: %v", err)
			}
			user := &sample.User{}
			if err := row.Columns(&user.UserID, &user.Name, &user.EntranceFlag); err != nil {
				t.Fatalf("failed to get columns: %v", err)
			}

			res := sample.HasEntranceFlag(user.EntranceFlag, tt.reqFlag)
			if diff := cmp.Diff(tt.wantRes, res); diff != "" {
				t.Errorf("mismatch: (-want +got): %s", diff)	
			}
		})
	}
}