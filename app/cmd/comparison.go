package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

type ComparisonHandler struct {
	cfnClient    *cfn.Client
	exporter     *Exporter
	eventEmitter EventEmitFn
}

func NewSF6ComparisonHandler(cfnClient *cfn.Client, exporter *Exporter) *ComparisonHandler {
	return &ComparisonHandler{
		cfnClient: cfnClient,
		exporter:  exporter,
	}
}

type ComparisonProgress struct {
	LocalizationKey string         `json:"localizationKey"`
	Params          map[string]any `json:"params,omitempty"`
}

func (h *ComparisonHandler) progress(localizationKey string, params map[string]any) {
	h.emit("comparison-progress", ComparisonProgress{LocalizationKey: localizationKey, Params: params})
}

func (h *ComparisonHandler) SetEventEmitter(eventEmitter EventEmitFn) { h.eventEmitter = eventEmitter }

// ExportSF6BattleStatsComparison takes the report the GUI already shows, as fetching
// it again takes a while.
func (h *ComparisonHandler) ExportSF6BattleStatsComparison(report model.SF6BattleStatsComparisonReport) error {
	records := [][]string{{"metric", "unit", report.UserName, "top avg", "top min", "top max", "delta"}}
	for _, m := range report.Metrics {
		records = append(records, []string{
			m.Key,
			m.Unit,
			formatMetric(m.Current),
			formatMetric(m.TopAvg),
			formatMetric(m.TopMin),
			formatMetric(m.TopMax),
			formatMetric(m.Current - m.TopAvg),
		})
	}
	return h.exporter.CSV(report.UserName+"-comparison.csv", records)
}

func formatMetric(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func (h *ComparisonHandler) emit(name string, data ...interface{}) {
	if h.eventEmitter != nil {
		h.eventEmitter(name, data...)
	}
}

func (h *ComparisonHandler) GetSF6BattleStatsComparison(userCode string) (*model.SF6BattleStatsComparisonReport, error) {
	if userCode == "" {
		return nil, fmt.Errorf("userCode is empty")
	}
	const topN = 5

	tabClient, cleanup, err := h.cfnClient.NewTabClient()
	if err != nil {
		return nil, err
	}
	defer func() { cleanup() }()

	// The tab shares the tracking browser's profile and so its buckler cookie.
	h.progress("comparisonCheckingLogin", nil)
	if !tabClient.HasSession(context.Background()) {
		// Logging in closes and relaunches the browser, which takes the tab with it.
		cleanup()
		done, err := authenticate(h.cfnClient.Authenticate, func(action tracker.AuthAction) {
			h.emit("comparison-auth-action-required", action)
		}, func(step string) {
			h.progress(step, nil)
		})
		if err != nil {
			return nil, err
		}
		if !done {
			return nil, model.WrapError(model.ErrAuth, errors.New("login stopped without a buckler session"))
		}
		h.emit("comparison-auth-success")
		tabClient, cleanup, err = h.cfnClient.NewTabClient()
		if err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	h.progress("comparisonFetchingBattleLog", nil)
	bl, err := tabClient.GetBattleLog(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if bl == nil {
		return nil, fmt.Errorf("battle log is nil")
	}

	characterToolName := bl.FighterBannerInfo.FavoriteCharacterToolName
	if characterToolName == "" {
		return nil, fmt.Errorf("favorite_character_tool_name is empty")
	}

	h.progress("comparisonFetchingStats", nil)
	playData, err := tabClient.GetPlayData(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if playData == nil {
		return nil, fmt.Errorf("play data is nil")
	}

	character := bl.FighterBannerInfo.FavoriteCharacterName
	h.progress("comparisonFetchingRanking", map[string]any{"character": character})
	topPlayers, err := tabClient.GetTopMRPlayersByCharacter(ctx, characterToolName)
	if err != nil {
		return nil, err
	}

	topStats := make([]*cfn.BattleStats, 0, topN)
	for _, p := range topPlayers {
		if len(topStats) >= topN {
			break
		}
		profileID, err := p.ProfileID()
		if err != nil {
			continue
		}
		h.progress("comparisonFetchingTopPlayer", map[string]any{"current": len(topStats) + 1, "total": topN})
		pd, err := tabClient.GetPlayData(ctx, profileID)
		if err != nil {
			continue
		}
		if pd == nil {
			continue
		}
		topStats = append(topStats, &pd.Play.BattleStats)
	}

	if len(topStats) == 0 {
		return nil, fmt.Errorf("no top player stats were fetched")
	}

	h.progress("comparisonComparing", nil)
	comparison, err := tabClient.CompareBattleStats(&playData.Play.BattleStats, topStats)
	if err != nil {
		return nil, err
	}

	metrics := make([]model.SF6BattleStatsMetric, 0, len(cfn.BattleStatsIntMetrics())+len(cfn.BattleStatsFloatMetrics()))
	for _, m := range cfn.BattleStatsIntMetrics() {
		metric := m.Metric(comparison)
		metrics = append(metrics, model.SF6BattleStatsMetric{
			Key:      m.Key,
			Unit:     metric.Unit,
			Kind:     "int",
			Polarity: cfn.MetricPolarity(m.Key),
			Current:  float64(metric.Current),
			TopAvg:   metric.TopPlayers.Avg,
			TopMin:   float64(metric.TopPlayers.Min),
			TopMax:   float64(metric.TopPlayers.Max),
		})
	}
	for _, m := range cfn.BattleStatsFloatMetrics() {
		metric := m.Metric(comparison)
		metrics = append(metrics, model.SF6BattleStatsMetric{
			Key:      m.Key,
			Unit:     metric.Unit,
			Kind:     "float",
			Polarity: cfn.MetricPolarity(m.Key),
			Current:  metric.Current,
			TopAvg:   metric.TopPlayers.Avg,
			TopMin:   metric.TopPlayers.Min,
			TopMax:   metric.TopPlayers.Max,
		})
	}

	return &model.SF6BattleStatsComparisonReport{
		UserCode:          userCode,
		UserName:          bl.FighterBannerInfo.PersonalInfo.FighterID,
		CharacterName:     character,
		CharacterToolName: characterToolName,
		TopN:              len(topStats),
		Metrics:           metrics,
	}, nil
}
