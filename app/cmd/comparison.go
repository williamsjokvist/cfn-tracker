package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

type ComparisonHandler struct {
	cfnClient cfn.CFNClient
}

func NewSF6ComparisonHandler(cfnClient cfn.CFNClient) *ComparisonHandler {
	return &ComparisonHandler{
		cfnClient: cfnClient,
	}
}

func (h *ComparisonHandler) GetSF6BattleStatsComparison(userCode string) (*model.SF6BattleStatsComparisonReport, error) {
	if userCode == "" {
		return nil, fmt.Errorf("userCode is empty")
	}
	const topN = 5

	baseClient, ok := h.cfnClient.(*cfn.Client)
	if !ok {
		return nil, fmt.Errorf("unexpected CFN client type")
	}

	tabClient, cleanup, err := baseClient.NewTabClient()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

	playData, err := tabClient.GetPlayData(ctx, userCode)
	if err != nil {
		return nil, err
	}
	if playData == nil {
		return nil, fmt.Errorf("play data is nil")
	}

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
		CharacterName:     bl.FighterBannerInfo.FavoriteCharacterName,
		CharacterToolName: characterToolName,
		TopN:              len(topStats),
		Metrics:           metrics,
	}, nil
}
