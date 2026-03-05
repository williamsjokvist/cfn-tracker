package tracker

import (
	"context"
	"errors"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestMockGameTrackerPollOrder(t *testing.T) {
	match1 := &model.Match{UserId: "user1", ReplayID: "replay1"}
	match2 := &model.Match{UserId: "user2", ReplayID: "replay2"}

	mock := NewMockGameTracker(nil,
		NewMockPollResult(match1),
		NewMockPollResult(match2),
	)

	ctx := context.Background()
	session := &model.Session{}

	// First call should return match1
	m, err := mock.Poll(ctx, session)
	if err != nil {
		t.Fatalf("Poll() returned error: %v", err)
	}
	if m != match1 {
		t.Errorf("Poll() returned %v, want %v", m, match1)
	}

	// Second call should return match2
	m, err = mock.Poll(ctx, session)
	if err != nil {
		t.Fatalf("Poll() returned error: %v", err)
	}
	if m != match2 {
		t.Errorf("Poll() returned %v, want %v", m, match2)
	}

	// Third call should return (nil, nil)
	m, err = mock.Poll(ctx, session)
	if err != nil {
		t.Fatalf("Poll() returned error: %v", err)
	}
	if m != nil {
		t.Errorf("Poll() returned %v, want nil", m)
	}
}

func TestMockGameTrackerPollError(t *testing.T) {
	testErr := errors.New("test error")
	mock := NewMockGameTracker(nil,
		NewMockPollError(testErr),
	)

	ctx := context.Background()
	session := &model.Session{}

	m, err := mock.Poll(ctx, session)
	if err != testErr {
		t.Errorf("Poll() returned error %v, want %v", err, testErr)
	}
	if m != nil {
		t.Errorf("Poll() returned match %v, want nil", m)
	}
}

func TestMockGameTrackerGetUser(t *testing.T) {
	user := &model.User{
		Id:          1,
		DisplayName: "TestUser",
		Code:        "ABC123",
	}
	mock := NewMockGameTracker(user)

	ctx := context.Background()
	u, err := mock.GetUser(ctx, "user1")
	if err != nil {
		t.Fatalf("GetUser() returned error: %v", err)
	}
	if u != user {
		t.Errorf("GetUser() returned %v, want %v", u, user)
	}
}

func TestMockGameTrackerGetUserError(t *testing.T) {
	userErr := errors.New("user not found")
	mock := NewMockGameTrackerWithError(userErr)

	ctx := context.Background()
	u, err := mock.GetUser(ctx, "user1")
	if err != userErr {
		t.Errorf("GetUser() returned error %v, want %v", err, userErr)
	}
	if u != nil {
		t.Errorf("GetUser() returned user %v, want nil", u)
	}
}
