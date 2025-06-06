package model

type Localization struct {
	AppVersion               string `json:"appVersion"`
	Source                   string `json:"source"`
	About                    string `json:"about"`
	Settings                 string `json:"settings"`
	Tracking                 string `json:"tracking"`
	History                  string `json:"history"`
	Language                 string `json:"language"`
	Changelog                string `json:"changelog"`
	StartTracking            string `json:"startTracking"`
	CFNName                  string `json:"cfnName"`
	TekkenId                 string `json:"tekkenId"`
	Start                    string `json:"start"`
	Opponent                 string `json:"opponent"`
	Character                string `json:"character"`
	LPGain                   string `json:"lpGain"`
	MRGain                   string `json:"mrGain"`
	DeleteLog                string `json:"deleteLog"`
	GoBack                   string `json:"goBack"`
	Loading                  string `json:"loading"`
	Wins                     string `json:"wins"`
	Losses                   string `json:"losses"`
	WinRate                  string `json:"winRate"`
	Stop                     string `json:"stop"`
	Files                    string `json:"files"`
	EnterCfnName             string `json:"enterCfnName"`
	EnterTekkenId            string `json:"enterTekkenId"`
	Result                   string `json:"result"`
	ReplayId                 string `json:"replayId"`
	Copy                     string `json:"copy"`
	Time                     string `json:"time"`
	WinStreak                string `json:"winStreak"`
	UpdateAvailable          string `json:"newVersionAvailable"`
	Pause                    string `json:"pause"`
	Unpause                  string `json:"unpause"`
	Statistics               string `json:"statistics"`
	Date                     string `json:"date"`
	Minimize                 string `json:"minimize"`
	RestoreSession           string `json:"restoreSession"`
	ExportLog                string `json:"exportLog"`
	League                   string `json:"league"`
	PickGame                 string `json:"pickGame"`
	Follow                   string `json:"follow"`
	ContinueStep             string `json:"continueStep"`
	Placement                string `json:"placement"`
	CFNError                 string `json:"cfnError"`
	Output                   string `json:"output"`
	Theme                    string `json:"theme"`
	DefaultTheme             string `json:"defaultTheme"`
	PickTheme                string `json:"pickTheme"`
	UsingBrowserSource       string `json:"usingBrowserSource"`
	CopyBrowserSourceLink    string `json:"copyBrowserSourceLink"`
	BrowserSourceDescription string `json:"browserSourceDescription"`
	DisplayStats             string `json:"displayStats"`
	ImportFiles              string `json:"importFiles"`
	OBSCustomize             string `json:"obsCustomize"`
	LastMatch                string `json:"lastMatch"`
	StatsWillBeDisplayed     string `json:"statsWillBeDisplayed"`
	SelectTheme              string `json:"selectTheme"`
	Started                  string `json:"started"`
	User                     string `json:"user"`
	MatchesWon               string `json:"matchesWon"`
	MatchesLost              string `json:"matchesLost"`
	Sessions                 string `json:"sessions"`
	Refresh                  string `json:"refresh"`
	Cooldown                 string `json:"cooldown"`
	T8Beginner               string `json:"T8_BEGINNER"`
	T81stDan                 string `json:"T8_1ST_DAN"`
	T82ndDan                 string `json:"T8_2ND_DAN"`
	T8Fighter                string `json:"T8_FIGHTER"`
	T8Strategist             string `json:"T8_STRATEGIST"`
	T8Combatant              string `json:"T8_COMBATANT"`
	T8Brawler                string `json:"T8_BRAWLER"`
	T8Ranger                 string `json:"T8_RANGER"`
	T8Cavalry                string `json:"T8_CAVALRY"`
	T8Warrior                string `json:"T8_WARRIOR"`
	T8Assailant              string `json:"T8_ASSAILANT"`
	T8Dominator              string `json:"T8_DOMINATOR"`
	T8Vanquisher             string `json:"T8_VANQUISHER"`
	T8Destroyer              string `json:"T8_DESTROYER"`
	T8Eliminator             string `json:"T8_ELIMINATOR"`
	T8Garyu                  string `json:"T8_GARYU"`
	T8Shinryu                string `json:"T8_SHINRYU"`
	T8Tenryu                 string `json:"T8_TENRYU"`
	T8MightyRuler            string `json:"T8_MIGHTY_RULER"`
	T8FlameRuler             string `json:"T8_FLAME_RULER"`
	T8BattleRuler            string `json:"T8_BATTLE_RULER"`
	T8Fujin                  string `json:"T8_FUJIN"`
	T8Raijin                 string `json:"T8_RAIJIN"`
	T8Kishin                 string `json:"T8_KISHIN"`
	T8Bushin                 string `json:"T8_BUSHIN"`
	T8TekkenKing             string `json:"T8_TEKKEN_KING"`
	T8TekkenEmperor          string `json:"T8_TEKKEN_EMPEROR"`
	T8TekkenGod              string `json:"T8_TEKKEN_GOD"`
	T8TekkenGodSupreme       string `json:"T8_TEKKEN_GOD_SUPREME"`
	T8GodOfDestruction       string `json:"T8_GOD_OF_DESTRUCTION"`
	ErrUnknown               string `json:"errUnknown"`
	ErrSelectGame            string `json:"errSelectGame"`
	ErrAuth                  string `json:"errAuth"`
	ErrGetLatestSession      string `json:"errGetLatestSession"`
	ErrGetUser               string `json:"errGetUser"`
	ErrGetMatches            string `json:"errGetMatches"`
	ErrSaveLocale            string `json:"errSaveLocale"`
	ErrCheckForUpdate        string `json:"errCheckForUpdate"`
	ErrGetGuiConfig          string `json:"errGetGuiConfig"`
	ErrSaveTheme             string `json:"errSaveTheme"`
	ErrSaveUser              string `json:"errSaveUser"`
	ErrSaveSidebar           string `json:"errSaveSidebar"`
	ErrGetSessions           string `json:"errGetSessions"`
	ErrGetTranslations       string `json:"errGetTranslations"`
	ErrGetSessionStatistics  string `json:"errGetSessionStatistics"`
	ErrCreateSession         string `json:"errCreateSession"`
	ErrOpenResultsDirectory  string `json:"errOpenResultsDirectory"`
	ErrReadThemeCSS          string `json:"errReadThemeCSS"`
}
