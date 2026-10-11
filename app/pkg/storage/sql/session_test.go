package sql

import (
	"context"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "123456789", DisplayName: "Kakeru"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	created, err := storage.CreateSession(ctx, "123456789")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := storage.GetSession(ctx, created.Id)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if got.Id != created.Id || got.UserId != "123456789" || got.UserName != "Kakeru" || got.CreatedAt != created.CreatedAt {
		t.Fatalf("GetSession() = %+v, want session %d by Kakeru created at %s", got, created.Id, created.CreatedAt)
	}

	if _, err := storage.GetSession(ctx, created.Id+1); err == nil {
		t.Fatal("GetSession() of a missing session returned no error")
	}
}
