package postgres_test

import (
	"context"
	"errors"
	"os"
	"testing"

	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
	"github.com/BenyaChef/micro/infrastructure/postgres"
)

const envTestDSN = "POSTGRES_TEST_DSN"

func TestBuildFailsWithoutRequiredFields(t *testing.T) {
	_, err := postgres.NewBuilder().Build(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, postgres.ErrDSNIsRequired) {
		t.Errorf("error %v does not contain %v", err, postgres.ErrDSNIsRequired)
	}
}

func TestBuildFailsOnMalformedDSN(t *testing.T) {
	_, err := postgres.NewBuilder().
		DSN("://not-a-dsn").
		Build(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !apperrors.EqualByCode(err, postgres.ErrCodeInvalidDSN) {
		t.Errorf("got code %s, want %s", apperrors.CodeOf(err), postgres.ErrCodeInvalidDSN)
	}
}

func TestExecQueryAndQueryRowOnLiveDatabase(t *testing.T) {
	client := newTestClient(t)

	const table = "exec_query_probe"
	createProbeTable(t, client, table)

	ctx := t.Context()

	if _, err := client.Exec(ctx, "INSERT INTO "+table+" (n) VALUES (1), (2)"); err != nil {
		t.Fatalf("Exec: %v", err)
	}

	rows, err := client.Query(ctx, "SELECT n FROM "+table+" ORDER BY n")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	var got []int

	for rows.Next() {
		var n int
		if scanErr := rows.Scan(&n); scanErr != nil {
			t.Fatalf("Scan: %v", scanErr)
		}

		got = append(got, n)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}

	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Errorf("got %v, want [1 2]", got)
	}

	if n := countRows(t, client, table); n != 2 {
		t.Errorf("QueryRow: got %d rows, want 2", n)
	}
}

func TestPingOnLiveDatabase(t *testing.T) {
	client := newTestClient(t)

	if err := client.Ping(t.Context()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func newTestClient(t *testing.T) *postgres.Client {
	t.Helper()

	dsn := os.Getenv(envTestDSN)
	if dsn == "" {
		t.Skipf("%s is not set, skipping test against live database", envTestDSN)
	}

	client, err := postgres.NewBuilder().
		DSN(dsn).
		Build(t.Context())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	t.Cleanup(client.Close)

	return client
}

func createProbeTable(t *testing.T, client *postgres.Client, table string) {
	t.Helper()

	ctx := t.Context()

	if _, err := client.Exec(ctx, "CREATE TABLE "+table+" (n int)"); err != nil {
		t.Fatalf("create table: %v", err)
	}

	t.Cleanup(func() {
		//nolint:usetesting // t.Context() уже отменён к моменту очистки.
		cleanupCtx := context.Background()

		if _, err := client.Exec(cleanupCtx, "DROP TABLE IF EXISTS "+table); err != nil {
			t.Errorf("drop table: %v", err)
		}
	})
}

func countRows(t *testing.T, client *postgres.Client, table string) int {
	t.Helper()

	ctx := t.Context()

	var count int
	if err := client.QueryRow(ctx, "SELECT count(*) FROM "+table).Scan(&count); err != nil {
		t.Fatalf("count rows: %v", err)
	}

	return count
}
