package cfn

type FloatStats struct {
	Min float64
	Max float64
	Avg float64
}

type IntStats struct {
	Min int
	Max int
	Avg float64
}

type FloatMetricComparison struct {
	Unit       string
	Current    float64
	TopPlayers FloatStats
}

type IntMetricComparison struct {
	Unit       string
	Current    int
	TopPlayers IntStats
}

type BattleStatsFloatMetric struct {
	Key   string
	Name  string
	Unit  string
	Scale float64
	Metric func(*BattleStatsComparison) *FloatMetricComparison
	Get   func(*BattleStats) float64
}

type BattleStatsIntMetric struct {
	Key   string
	Name  string
	Unit  string
	Metric func(*BattleStatsComparison) *IntMetricComparison
	Get   func(*BattleStats) int
}

type BattleStatsComparison struct {
	BattleHubMatchPlayCount          IntMetricComparison
	CasualMatchPlayCount             IntMetricComparison
	CustomRoomMatchPlayCount         IntMetricComparison
	RankMatchPlayCount               IntMetricComparison
	CornerTime                       FloatMetricComparison
	CorneredTime                     FloatMetricComparison
	DriveImpact                      FloatMetricComparison
	DriveImpactToDriveImpact         FloatMetricComparison
	DriveParry                       FloatMetricComparison
	DriveReversal                    FloatMetricComparison
	JustParry                        FloatMetricComparison
	PunishCounter                    FloatMetricComparison
	ThrowCount                       FloatMetricComparison
	ThrowDriveParry                  FloatMetricComparison
	ThrowTech                        FloatMetricComparison
	Stun                             FloatMetricComparison
	ReceivedDriveImpact              FloatMetricComparison
	ReceivedDriveImpactToDriveImpact FloatMetricComparison
	ReceivedPunishCounter            FloatMetricComparison
	ReceivedStun                     FloatMetricComparison
	ReceivedThrowCount               FloatMetricComparison
	ReceivedThrowDriveParry          FloatMetricComparison
	GaugeRateCA                      FloatMetricComparison
	GaugeRateDriveArts               FloatMetricComparison
	GaugeRateDriveGuard              FloatMetricComparison
	GaugeRateDriveImpact             FloatMetricComparison
	GaugeRateDriveOther              FloatMetricComparison
	GaugeRateDriveReversal           FloatMetricComparison
	GaugeRateDriveRushFromCancel     FloatMetricComparison
	GaugeRateDriveRushFromParry      FloatMetricComparison
	GaugeRateSaLv1                   FloatMetricComparison
	GaugeRateSaLv2                   FloatMetricComparison
	GaugeRateSaLv3                   FloatMetricComparison
	TargetClearCount                 IntMetricComparison
	TotalAllCharacterPlayPoint       IntMetricComparison
	RivalAiAchievedChallengeCount    IntMetricComparison
	RivalAiHighestLeagueRank         IntMetricComparison
}

var battleStatsFloatMetrics = []BattleStatsFloatMetric{
	{Key: "corneredTime", Name: "Cornered Time", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.CorneredTime }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.CorneredTime }},
	{Key: "cornerTime", Name: "Corner Time", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.CornerTime }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.CornerTime }},
	{Key: "driveImpact", Name: "Drive Impact", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.DriveImpact }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.DriveImpact }},
	{Key: "driveImpactToDriveImpact", Name: "Drive Impact (DI vs DI)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.DriveImpactToDriveImpact }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.DriveImpactToDriveImpact }},
	{Key: "driveParry", Name: "Drive Parry", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.DriveParry }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.DriveParry }},
	{Key: "driveReversal", Name: "Drive Reversal", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.DriveReversal }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.DriveReversal }},
	{Key: "justParry", Name: "Just Parry", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.JustParry }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.JustParry }},
	{Key: "punishCounter", Name: "Punish Counter", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.PunishCounter }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.PunishCounter }},
	{Key: "throwCount", Name: "Throw Count", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ThrowCount }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ThrowCount }},
	{Key: "throwDriveParry", Name: "Throw (vs Drive Parry)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ThrowDriveParry }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ThrowDriveParry }},
	{Key: "throwTech", Name: "Throw Tech", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ThrowTech }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ThrowTech }},
	{Key: "stun", Name: "Stun", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.Stun }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.Stun }},
	{Key: "receivedDriveImpact", Name: "Received Drive Impact", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedDriveImpact }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedDriveImpact }},
	{Key: "receivedDriveImpactToDriveImpact", Name: "Received DI (DI vs DI)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedDriveImpactToDriveImpact }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedDriveImpactToDriveImpact }},
	{Key: "receivedPunishCounter", Name: "Received Punish Counter", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedPunishCounter }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedPunishCounter }},
	{Key: "receivedStun", Name: "Received Stun", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedStun }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedStun }},
	{Key: "receivedThrowCount", Name: "Received Throw Count", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedThrowCount }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedThrowCount }},
	{Key: "receivedThrowDriveParry", Name: "Received Throw (vs Drive Parry)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.ReceivedThrowDriveParry }, Unit: "%", Scale: 1, Get: func(s *BattleStats) float64 { return s.ReceivedThrowDriveParry }},
	{Key: "gaugeRateCA", Name: "Gauge Rate (CA)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateCA }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateCA }},
	{Key: "gaugeRateDriveArts", Name: "Gauge Rate (Drive Arts)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveArts }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveArts }},
	{Key: "gaugeRateDriveGuard", Name: "Gauge Rate (Drive Guard)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveGuard }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveGuard }},
	{Key: "gaugeRateDriveImpact", Name: "Gauge Rate (Drive Impact)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveImpact }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveImpact }},
	{Key: "gaugeRateDriveOther", Name: "Gauge Rate (Drive Other)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveOther }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveOther }},
	{Key: "gaugeRateDriveReversal", Name: "Gauge Rate (Drive Reversal)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveReversal }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveReversal }},
	{Key: "gaugeRateDriveRushFromCancel", Name: "Gauge Rate (DR from Cancel)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveRushFromCancel }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveRushFromCancel }},
	{Key: "gaugeRateDriveRushFromParry", Name: "Gauge Rate (DR from Parry)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateDriveRushFromParry }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateDriveRushFromParry }},
	{Key: "gaugeRateSaLv1", Name: "Gauge Rate (SA Lv1)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateSaLv1 }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateSaLv1 }},
	{Key: "gaugeRateSaLv2", Name: "Gauge Rate (SA Lv2)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateSaLv2 }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateSaLv2 }},
	{Key: "gaugeRateSaLv3", Name: "Gauge Rate (SA Lv3)", Metric: func(c *BattleStatsComparison) *FloatMetricComparison { return &c.GaugeRateSaLv3 }, Unit: "%", Scale: 100, Get: func(s *BattleStats) float64 { return s.GaugeRateSaLv3 }},
}

var battleStatsIntMetrics = []BattleStatsIntMetric{
	{Key: "battleHubMatchPlayCount", Name: "Battle Hub Match Play Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.BattleHubMatchPlayCount }, Unit: "matches", Get: func(s *BattleStats) int { return s.BattleHubMatchPlayCount }},
	{Key: "casualMatchPlayCount", Name: "Casual Match Play Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.CasualMatchPlayCount }, Unit: "matches", Get: func(s *BattleStats) int { return s.CasualMatchPlayCount }},
	{Key: "customRoomMatchPlayCount", Name: "Custom Room Match Play Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.CustomRoomMatchPlayCount }, Unit: "matches", Get: func(s *BattleStats) int { return s.CustomRoomMatchPlayCount }},
	{Key: "rankMatchPlayCount", Name: "Rank Match Play Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.RankMatchPlayCount }, Unit: "matches", Get: func(s *BattleStats) int { return s.RankMatchPlayCount }},
	{Key: "targetClearCount", Name: "Target Clear Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.TargetClearCount }, Unit: "count", Get: func(s *BattleStats) int { return s.TargetClearCount }},
	{Key: "totalAllCharacterPlayPoint", Name: "Total All Character Play Point", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.TotalAllCharacterPlayPoint }, Unit: "pt", Get: func(s *BattleStats) int { return s.TotalAllCharacterPlayPoint }},
	{Key: "rivalAiAchievedChallengeCount", Name: "Rival AI Achieved Challenge Count", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.RivalAiAchievedChallengeCount }, Unit: "count", Get: func(s *BattleStats) int { return s.RivalAiAchievedChallengeCount }},
	{Key: "rivalAiHighestLeagueRank", Name: "Rival AI Highest League Rank", Metric: func(c *BattleStatsComparison) *IntMetricComparison { return &c.RivalAiHighestLeagueRank }, Unit: "rank", Get: func(s *BattleStats) int { return s.RivalAiHighestLeagueRank }},
}

func NewBattleStatsComparison() *BattleStatsComparison {
	out := &BattleStatsComparison{}
	for _, m := range battleStatsFloatMetrics {
		m.Metric(out).Unit = m.Unit
	}
	for _, m := range battleStatsIntMetrics {
		m.Metric(out).Unit = m.Unit
	}
	return out
}

var metricPolarity = map[string]int{
	"corneredTime": -1,
	"receivedStun": -1,
	"receivedDriveImpact": -1,
}

func MetricPolarity(key string) int {
	if p, ok := metricPolarity[key]; ok {
		return p
	}
	return 1
}

func BattleStatsFloatMetrics() []BattleStatsFloatMetric {
	out := make([]BattleStatsFloatMetric, len(battleStatsFloatMetrics))
	copy(out, battleStatsFloatMetrics)
	return out
}

func BattleStatsIntMetrics() []BattleStatsIntMetric {
	out := make([]BattleStatsIntMetric, len(battleStatsIntMetrics))
	copy(out, battleStatsIntMetrics)
	return out
}

func CalcStatsFloat(values []float64) (min, max, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	min, max = values[0], values[0]
	sum := 0.0
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	avg = sum / float64(len(values))
	return min, max, avg
}

func CalcStatsInt(values []int) (min, max int, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	min, max = values[0], values[0]
	sum := 0
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	avg = float64(sum) / float64(len(values))
	return min, max, avg
}
