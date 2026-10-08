package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// MockPollResult represents a single Poll() return value for testing.
type MockPollResult struct {
	Match *model.Match
	Err   error
}

// NewMockPollResult creates a MockPollResult that returns the given match.
func NewMockPollResult(match *model.Match) MockPollResult {
	return MockPollResult{Match: match}
}

// NewMockPollError creates a MockPollResult that returns the given error.
func NewMockPollError(err error) MockPollResult {
	return MockPollResult{Err: err}
}

// MockGameTracker is a test double for GameTracker.
// It returns pre-configured results in order from Poll().
// When the queue is exhausted, Poll() returns (nil, nil).
type MockGameTracker struct {
	user        *model.User
	userErr     error
	pollResults []MockPollResult
	pollIndex   int
}

var _ GameTracker = (*MockGameTracker)(nil)

// NewMockGameTracker creates a MockGameTracker with the given user and poll results.
func NewMockGameTracker(user *model.User, results ...MockPollResult) *MockGameTracker {
	return &MockGameTracker{
		user:        user,
		pollResults: results,
	}
}

// NewMockGameTrackerWithError creates a MockGameTracker whose GetUser returns an error.
func NewMockGameTrackerWithError(userErr error) *MockGameTracker {
	return &MockGameTracker{userErr: userErr}
}

func (m *MockGameTracker) GetUser(ctx context.Context, userId string) (*model.User, error) {
	if m.userErr != nil {
		return nil, m.userErr
	}
	return m.user, nil
}

func (m *MockGameTracker) Poll(ctx context.Context, session *model.Session) (*model.Match, error) {
	if m.pollIndex >= len(m.pollResults) {
		return nil, nil
	}
	result := m.pollResults[m.pollIndex]
	m.pollIndex++
	return result.Match, result.Err
}
