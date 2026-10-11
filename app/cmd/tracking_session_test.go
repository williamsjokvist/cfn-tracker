package cmd_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn/cfntest"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu/wavutest"
)

type liveTracking struct {
	t       *testing.T
	handler *cmd.TrackingHandler
	matches chan model.Match
	done    chan error
}

func newGameHandler(t *testing.T, wavuClient *wavutest.Client, cfnClient *cfntest.Client, game model.GameType) *cmd.TrackingHandler {
	t.Helper()
	handler := cmd.NewTrackingHandler(wavuClient, cfnClient, testSuite.sqlDb, testSuite.cfgDb, testSuite.txtDb, &config.BuildConfig{})
	if err := handler.SelectGame(game); err != nil {
		t.Fatalf("SelectGame(%s) error = %v", game, err)
	}
	return handler
}

func startLiveTracking(t *testing.T, handler *cmd.TrackingHandler, userCode string, restore bool) *liveTracking {
	t.Helper()
	lt := &liveTracking{t: t, handler: handler, matches: make(chan model.Match, 16), done: make(chan error, 1)}
	handler.SetEventEmitter(func(name string, data ...interface{}) {
		if name == "match" {
			lt.matches <- data[0].(model.Match)
		}
	})
	go func() { lt.done <- handler.StartTracking(userCode, restore) }()
	t.Cleanup(lt.stop)

	lt.waitForMatch()
	if restore {
		lt.waitForMatch()
	}
	return lt
}

func (lt *liveTracking) next() model.Match {
	lt.t.Helper()
	lt.handler.ForcePoll()
	return lt.waitForMatch()
}

func (lt *liveTracking) waitForMatch() model.Match {
	lt.t.Helper()
	select {
	case match := <-lt.matches:
		return match
	case err := <-lt.done:
		lt.t.Fatalf("tracking stopped: %v", err)
	case <-time.After(5 * time.Second):
		lt.t.Fatal("no match event")
	}
	return model.Match{}
}

func (lt *liveTracking) stop() {
	lt.handler.StopTracking()
	select {
	case err := <-lt.done:
		if err != nil {
			lt.t.Errorf("StartTracking() error = %v", err)
		}
		lt.done <- err
	case <-time.After(5 * time.Second):
		lt.t.Fatal("StartTracking did not stop")
	}
	select {
	case match := <-lt.matches:
		lt.t.Errorf("unexpected match event after the last one: %+v", match)
	default:
	}
}

type wantStats struct {
	character                       string
	wins, losses, winStreak, lpGain int
}

func assertStats(t *testing.T, match model.Match, want wantStats) {
	t.Helper()
	got := wantStats{match.Character, match.Wins, match.Losses, match.WinStreak, match.LPGain}
	if got != want {
		t.Fatalf("match stats = %+v, want %+v (match %+v)", got, want, match)
	}
}

func uniqueUserCode(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

func TestTrackingT8(t *testing.T) {
	polarisId := uniqueUserCode("t8")
	wavuClient := wavutest.NewClient(polarisId, "Tekk!")
	handler := newGameHandler(t, wavuClient, nil, model.GameTypeT8)

	lt := startLiveTracking(t, handler, polarisId, false)
	wavuClient.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 9})
	assertStats(t, lt.next(), wantStats{"Steve", 1, 0, 1, 0})
	// Polling again before the next match must not record the same replay twice.
	handler.ForcePoll()
	handler.ForcePoll()
	wavuClient.Play(wavutest.Match{Opponent: "lowhigh", Character: 9, AsP2: true})
	assertStats(t, lt.next(), wantStats{"Steve", 1, 1, 0, 0})
	wavuClient.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 8})
	assertStats(t, lt.next(), wantStats{"Kazuya", 1, 0, 1, 0})
	lt.stop()

	stored, err := testSuite.sqlDb.GetLatestSession(context.Background(), polarisId)
	if err != nil || stored == nil || len(stored.Matches) != 3 {
		t.Fatalf("GetLatestSession() = %+v, %v, want 3 matches", stored, err)
	}

	lt = startLiveTracking(t, handler, polarisId, true)
	wavuClient.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 9})
	assertStats(t, lt.next(), wantStats{"Steve", 2, 1, 1, 0})
}

func TestTrackingSF6(t *testing.T) {
	userCode := uniqueUserCode("")
	cfnClient := cfntest.NewClient("Kakeru", "Ryu", 10000, 1500)
	cfnClient.SetPoints("Ken", 7000, 1500)
	handler := newGameHandler(t, nil, cfnClient, model.GameTypeSF6)

	lt := startLiveTracking(t, handler, userCode, false)
	cfnClient.Play(cfntest.Match{Won: true, Opponent: "Punk", LPChange: 100})
	assertStats(t, lt.next(), wantStats{"Ryu", 1, 0, 1, 100})
	handler.ForcePoll()
	handler.ForcePoll()
	cfnClient.Play(cfntest.Match{Opponent: "Punk", LPChange: -80, AsP2: true})
	assertStats(t, lt.next(), wantStats{"Ryu", 1, 1, 0, 20})
	cfnClient.Play(cfntest.Match{Won: true, Opponent: "Punk", Character: "Ken", LPChange: 20})
	assertStats(t, lt.next(), wantStats{"Ken", 1, 0, 1, 0})
	cfnClient.Play(cfntest.Match{Won: true, Opponent: "Punk", Character: "Ryu", LPChange: 50})
	assertStats(t, lt.next(), wantStats{"Ryu", 2, 1, 1, 70})
	lt.stop()

	stored, err := testSuite.sqlDb.GetLatestSession(context.Background(), userCode)
	if err != nil || stored == nil || len(stored.Matches) != 4 || stored.EndingLP != 10070 {
		t.Fatalf("GetLatestSession() = %+v, %v, want 4 matches ending at 10070 LP", stored, err)
	}

	lt = startLiveTracking(t, handler, userCode, true)
	cfnClient.Play(cfntest.Match{Won: true, Opponent: "Punk", Character: "Ken", LPChange: 40})
	assertStats(t, lt.next(), wantStats{"Ken", 2, 0, 2, 40})
}
