package auth_test

import (
	"context"
	"testing"

	"go.breu.io/ctrlplane/internal/testutils"
)

func TestHandler(t *testing.T) {
	dbcon, err := testutils.StartDBContainer(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = dbcon.Stop()
	})
}
