package fcapi

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestDecode(t *testing.T) {
	tests := map[string]struct {
		status     int
		text       string
		wantStatus int
		wantParse  bool
		wantErr    bool
	}{
		"ok":            {status: 200, text: `{"res":"OK","user":{"name":"naggamurasang"}}`},
		"left the site": {status: 0, wantStatus: 401},
		"check":         {status: 403, text: "<!DOCTYPE html><title>Just a moment...</title>", wantStatus: 403},
		"check as 200":  {status: 200, text: "  <!DOCTYPE html>", wantStatus: 403},
		"server error":  {status: 502, text: "bad gateway", wantStatus: 502},
		"not json":      {status: 200, text: "nope", wantParse: true},
		"api error":     {status: 200, text: `{"res":"error","error":"user not found"}`, wantErr: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var res userResponse
			err := decode("get user", tt.status, tt.text, &res)
			var statusErr *model.HTTPStatusError
			var parseErr *model.ParseError
			switch {
			case tt.wantStatus != 0:
				if !errors.As(err, &statusErr) || statusErr.StatusCode != tt.wantStatus {
					t.Fatalf("decode() error = %v, want http status %d", err, tt.wantStatus)
				}
			case tt.wantParse:
				if !errors.As(err, &parseErr) {
					t.Fatalf("decode() error = %v, want a parse error", err)
				}
			case tt.wantErr:
				if err == nil || errors.As(err, &statusErr) || errors.As(err, &parseErr) {
					t.Fatalf("decode() error = %v, want an api error", err)
				}
			default:
				if err != nil || res.User.Name != "naggamurasang" {
					t.Fatalf("decode() = %+v, %v, want user naggamurasang", res.User, err)
				}
			}
		})
	}
}

func TestDecodeReplay(t *testing.T) {
	text := `{"res":"OK","results":{"count":1,"results":[{"quarkid":"1","channelname":"3rd Strike","date":1791794018000,"duration":312.5,"emulator":"fbneo","gameid":"sfiii3nr1","num_matches":5,"ranked":"cancelled","players":[{"name":"a","country":{"iso_code":"SE","full_name":"Sweden"},"rank":5,"score":3},{"name":"b","country":"JP","rank":null,"score":null}]}]}}`
	var res replaysResponse
	if err := decode("get replays", 200, text, &res); err != nil {
		t.Fatalf("decode() error = %v", err)
	}
	r := res.Results.Results[0]
	if !r.Cancelled() || r.Players[0].Country.IsoCode != "SE" || r.Players[1].Country.IsoCode != "JP" || r.Players[1].Score != nil {
		t.Fatalf("decode() = %+v, want a cancelled set between SE and JP", r)
	}
	if got := r.Players[0].Rank.String(); got != "A" {
		t.Fatalf("rank = %s, want A", got)
	}
}

func TestTimestampAcceptsSecondsAndMilliseconds(t *testing.T) {
	want := time.Unix(1791794018, 0)
	for _, raw := range []string{"1791794018", "1791794018000"} {
		var ts Timestamp
		if err := json.Unmarshal([]byte(raw), &ts); err != nil {
			t.Fatalf("unmarshal %s: %v", raw, err)
		}
		if !ts.Time().Equal(want) {
			t.Fatalf("Timestamp(%s) = %s, want %s", raw, ts.Time(), want)
		}
	}
}
