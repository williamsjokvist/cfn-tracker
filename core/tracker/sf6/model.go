package sf6

import (
	"fmt"
	"strconv"

	"github.com/williamsjokvist/cfn-tracker/core/model"
)

type PlayerInfo struct {
	AllowCrossPlay  bool `json:"allow_cross_play"`
	BattleInputType int  `json:"battle_input_type"`
	CharacterID     int  `json:"character_id"`
	HomeID          int  `json:"home_id"`
	LeaguePoint     int  `json:"league_point"`
	LeagueRank      int  `json:"league_rank"`
	MainCircle      struct {
		CircleID   string `json:"circle_id"`
		CircleName string `json:"circle_name"`
		DataExist  bool   `json:"data_exist"`
		Emblem     struct {
			EmblemBase                       int  `json:"emblem_base"`
			EmblemBaseColor                  int  `json:"emblem_base_color"`
			EmblemFrame                      int  `json:"emblem_frame"`
			EmblemFrameColor                 int  `json:"emblem_frame_color"`
			EmblemPattern                    int  `json:"emblem_pattern"`
			EmblemPatternColor               int  `json:"emblem_pattern_color"`
			EmblemPatternHorizontalInversion bool `json:"emblem_pattern_horizontal_inversion"`
			EmblemPatternVerticalInversion   bool `json:"emblem_pattern_vertical_inversion"`
			EmblemProcessing                 int  `json:"emblem_processing"`
			EmblemSymbol1                    int  `json:"emblem_symbol1"`
			EmblemSymbol1Clipping            bool `json:"emblem_symbol1_clipping"`
			EmblemSymbol1Color               int  `json:"emblem_symbol1_color"`
			EmblemSymbol1HorizontalInversion bool `json:"emblem_symbol1_horizontal_inversion"`
			EmblemSymbol1VerticalInversion   bool `json:"emblem_symbol1_vertical_inversion"`
			EmblemSymbol2                    int  `json:"emblem_symbol2"`
			EmblemSymbol2Clipping            bool `json:"emblem_symbol2_clipping"`
			EmblemSymbol2Color               int  `json:"emblem_symbol2_color"`
			EmblemSymbol2HorizontalInversion bool `json:"emblem_symbol2_horizontal_inversion"`
			EmblemSymbol2VerticalInversion   bool `json:"emblem_symbol2_vertical_inversion"`
		} `json:"emblem"`
		Leader struct {
			FighterID  string `json:"fighter_id"`
			PlatformID int    `json:"platform_id"`
			ShortID    int    `json:"short_id"`
		} `json:"leader"`
	} `json:"main_circle"`
	MasterLeague        int `json:"master_league"`
	MasterRating        int `json:"master_rating"`
	MasterRatingRanking int `json:"master_rating_ranking"`
	Player              struct {
		FighterID        string `json:"fighter_id"`
		PlatformID       int    `json:"platform_id"`
		ShortID          int64  `json:"short_id"`
		PlatformName     string `json:"platform_name"`
		PlatformToolName string `json:"platform_tool_name"`
	} `json:"player"`
	PlayingCharacterID       int    `json:"playing_character_id"`
	RoundResults             []int  `json:"round_results"`
	TitlePlate               int    `json:"title_plate"`
	CharacterName            string `json:"character_name"`
	CharacterToolName        string `json:"character_tool_name"`
	PlayingCharacterName     string `json:"playing_character_name"`
	PlayingCharacterToolName string `json:"playing_character_tool_name"`
	TitleData                struct {
		TitleDataID        int    `json:"title_data_id"`
		TitleDataGradeID   int    `json:"title_data_grade_id"`
		TitleDataGradeName string `json:"title_data_grade_name"`
		TitleDataPlateID   int    `json:"title_data_plate_id"`
		TitleDataPlateName string `json:"title_data_plate_name"`
		TitleDataVal       string `json:"title_data_val"`
	} `json:"title_data"`
	BattleInputTypeName string `json:"battle_input_type_name"`
}

type ProfilePage struct {
	Props struct {
		PageProps BattleLog `json:"pageProps"`
		NSsp      bool      `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Sid string `json:"sid"`
	} `json:"query"`
	BuildID       string   `json:"buildId"`
	AssetPrefix   string   `json:"assetPrefix"`
	IsFallback    bool     `json:"isFallback"`
	Gssp          bool     `json:"gssp"`
	Locale        string   `json:"locale"`
	Locales       []string `json:"locales"`
	DefaultLocale string   `json:"defaultLocale"`
	ScriptLoader  []any    `json:"scriptLoader"`
}

type BattleLog struct {
	FighterBannerInfo struct {
		AllowCrossPlay              bool `json:"allow_cross_play"`
		BattleInputType             int  `json:"battle_input_type"`
		CustomRoomInviteSetting     int  `json:"custom_room_invite_setting"`
		EnjoyTotalPoint             int  `json:"enjoy_total_point"`
		FavoriteCharacterID         int  `json:"favorite_character_id"`
		FavoriteCharacterLeagueInfo struct {
			LeaguePoint         int `json:"league_point"`
			LeagueRank          int `json:"league_rank"`
			MasterLeague        int `json:"master_league"`
			MasterRating        int `json:"master_rating"`
			MasterRatingRanking int `json:"master_rating_ranking"`
			LeagueRankInfo      struct {
				LeagueRankName   string `json:"league_rank_name"`
				LeagueRankNumber int    `json:"league_rank_number"`
			} `json:"league_rank_info"`
		} `json:"favorite_character_league_info"`
		FavoriteCharacterPlayPoint struct {
			BattleHub      int `json:"battle_hub"`
			FightingGround int `json:"fighting_ground"`
			WorldTour      int `json:"world_tour"`
		} `json:"favorite_character_play_point"`
		FriendRequestFlag bool `json:"friend_request_flag"`
		Friendship        int  `json:"friendship"`
		HomeID            int  `json:"home_id"`
		IsBeginner        bool `json:"is_beginner"`
		IsCircleInvite    bool `json:"is_circle_invite"`
		IsCircleMember    bool `json:"is_circle_member"`
		LastPlayAt        int  `json:"last_play_at"`
		MainCircle        struct {
			CircleID   string `json:"circle_id"`
			CircleName string `json:"circle_name"`
			DataExist  bool   `json:"data_exist"`
			Emblem     struct {
				EmblemBase                       int  `json:"emblem_base"`
				EmblemBaseColor                  int  `json:"emblem_base_color"`
				EmblemFrame                      int  `json:"emblem_frame"`
				EmblemFrameColor                 int  `json:"emblem_frame_color"`
				EmblemPattern                    int  `json:"emblem_pattern"`
				EmblemPatternColor               int  `json:"emblem_pattern_color"`
				EmblemPatternHorizontalInversion bool `json:"emblem_pattern_horizontal_inversion"`
				EmblemPatternVerticalInversion   bool `json:"emblem_pattern_vertical_inversion"`
				EmblemProcessing                 int  `json:"emblem_processing"`
				EmblemSymbol1                    int  `json:"emblem_symbol1"`
				EmblemSymbol1Clipping            bool `json:"emblem_symbol1_clipping"`
				EmblemSymbol1Color               int  `json:"emblem_symbol1_color"`
				EmblemSymbol1HorizontalInversion bool `json:"emblem_symbol1_horizontal_inversion"`
				EmblemSymbol1VerticalInversion   bool `json:"emblem_symbol1_vertical_inversion"`
				EmblemSymbol2                    int  `json:"emblem_symbol2"`
				EmblemSymbol2Clipping            bool `json:"emblem_symbol2_clipping"`
				EmblemSymbol2Color               int  `json:"emblem_symbol2_color"`
				EmblemSymbol2HorizontalInversion bool `json:"emblem_symbol2_horizontal_inversion"`
				EmblemSymbol2VerticalInversion   bool `json:"emblem_symbol2_vertical_inversion"`
			} `json:"emblem"`
			Leader struct {
				FighterID  string `json:"fighter_id"`
				PlatformID int    `json:"platform_id"`
				ShortID    int    `json:"short_id"`
			} `json:"leader"`
		} `json:"main_circle"`
		MaxContentPlayTime struct {
			ContentType int `json:"content_type"`
			PlayTime    int `json:"play_time"`
		} `json:"max_content_play_time"`
		MobileLinkage    bool `json:"mobile_linkage"`
		OnlineStatusInfo struct {
			BattlehubForBeginner                       bool   `json:"battlehub_for_beginner"`
			BattlehubID                                string `json:"battlehub_id"`
			BattlehubPlatformID                        int    `json:"battlehub_platform_id"`
			BattlehubRegionID                          int    `json:"battlehub_region_id"`
			BattlehubServerNo                          int    `json:"battlehub_server_no"`
			CustomRoomMasterShortID                    int    `json:"custom_room_master_short_id"`
			CustomRoomPlatformID                       int    `json:"custom_room_platform_id"`
			CustomRoomPublishSetting                   int    `json:"custom_room_publish_setting"`
			CustomRoomRegionID                         int    `json:"custom_room_region_id"`
			CustomRoomRequiredNetworkConnectionQuality int    `json:"custom_room_required_network_connection_quality"`
			CustomRoomRequiredPassCode                 bool   `json:"custom_room_required_pass_code"`
			CustomRoomRoomID                           string `json:"custom_room_room_id"`
			OnlineStatus                               int    `json:"online_status"`
			OnlineStatusData                           struct {
				OnlineStatusName string `json:"online_status_name"`
				OnlineStatusType int    `json:"online_status_type"`
			} `json:"online_status_data"`
			BattlehubRegionName       string `json:"battlehub_region_name"`
			BattlehubFormatedServerNo string `json:"battlehub_formated_server_no"`
		} `json:"online_status_info"`
		PersonalInfo struct {
			FighterID        string `json:"fighter_id"`
			PlatformID       int    `json:"platform_id"`
			ShortID          int64  `json:"short_id"`
			PlatformName     string `json:"platform_name"`
			PlatformToolName string `json:"platform_tool_name"`
		} `json:"personal_info"`
		PlayTimeZone struct {
			EndHour     int `json:"end_hour"`
			EndMinute   int `json:"end_minute"`
			StartHour   int `json:"start_hour"`
			StartMinute int `json:"start_minute"`
		} `json:"play_time_zone"`
		ProfileComment struct {
			ProfileTagID     int    `json:"profile_tag_id"`
			TagOptionID      int    `json:"tag_option_id"`
			ProfileTagName   string `json:"profile_tag_name"`
			ProfileTagOption string `json:"profile_tag_option"`
		} `json:"profile_comment"`
		TitlePlate                int    `json:"title_plate"`
		HomeName                  string `json:"home_name"`
		FavoriteCharacterName     string `json:"favorite_character_name"`
		FavoriteCharacterAlpha    string `json:"favorite_character_alpha"`
		FavoriteCharacterToolName string `json:"favorite_character_tool_name"`
		TitleData                 struct {
			TitleDataID        int    `json:"title_data_id"`
			TitleDataGradeID   int    `json:"title_data_grade_id"`
			TitleDataGradeName string `json:"title_data_grade_name"`
			TitleDataPlateID   int    `json:"title_data_plate_id"`
			TitleDataPlateName string `json:"title_data_plate_name"`
			TitleDataVal       string `json:"title_data_val"`
		} `json:"title_data"`
		IsMyData bool `json:"is_my_data"`
	} `json:"fighter_banner_info"`
	CurrentPage int      `json:"current_page"`
	ReplayList  []Replay `json:"replay_list"`
	TotalPage   int      `json:"total_page"`
	Sid         int64    `json:"sid"`
	WeekList    []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"week_list"`
	SeriesList []struct {
		Value string `json:"value"`
		Label string `json:"label"`
		Short string `json:"short"`
	} `json:"series_list"`
	SnsList []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	} `json:"sns_list"`
	Common struct {
		StatusCode int  `json:"statusCode"`
		IsError    bool `json:"isError"`
		LoginUser  struct {
			PlatformID int    `json:"platformId"`
			ShortID    int64  `json:"shortId"`
			FighterID  string `json:"fighterId"`
			Flg        bool   `json:"flg"`
			RegionID   int    `json:"regionId"`
		} `json:"loginUser"`
		AppEnv string `json:"appEnv"`
	} `json:"common"`
	Lang string `json:"__lang"`
}

type FighterBanner struct {
	AllowCrossPlay              bool `json:"allow_cross_play"`
	BattleInputType             int  `json:"battle_input_type"`
	CustomRoomInviteSetting     int  `json:"custom_room_invite_setting"`
	EnjoyTotalPoint             int  `json:"enjoy_total_point"`
	FavoriteCharacterID         int  `json:"favorite_character_id"`
	FavoriteCharacterLeagueInfo struct {
		LeaguePoint         int `json:"league_point"`
		LeagueRank          int `json:"league_rank"`
		MasterLeague        int `json:"master_league"`
		MasterRating        int `json:"master_rating"`
		MasterRatingRanking int `json:"master_rating_ranking"`
	} `json:"favorite_character_league_info"`
	FavoriteCharacterPlayPoint struct {
		BattleHub      int `json:"battle_hub"`
		FightingGround int `json:"fighting_ground"`
		WorldTour      int `json:"world_tour"`
	} `json:"favorite_character_play_point"`
	FriendRequestFlag bool `json:"friend_request_flag"`
	Friendship        int  `json:"friendship"`
	HomeID            int  `json:"home_id"`
	InsideRank        int  `json:"inside_rank"`
	IsCircleInvite    bool `json:"is_circle_invite"`
	IsCircleMember    bool `json:"is_circle_member"`
	LastPlayAt        int  `json:"last_play_at"`
	MainCircle        struct {
		CircleID   string `json:"circle_id"`
		CircleName string `json:"circle_name"`
		DataExist  bool   `json:"data_exist"`
		Emblem     struct {
			EmblemBase                       int  `json:"emblem_base"`
			EmblemBaseColor                  int  `json:"emblem_base_color"`
			EmblemFrame                      int  `json:"emblem_frame"`
			EmblemFrameColor                 int  `json:"emblem_frame_color"`
			EmblemPattern                    int  `json:"emblem_pattern"`
			EmblemPatternColor               int  `json:"emblem_pattern_color"`
			EmblemPatternHorizontalInversion bool `json:"emblem_pattern_horizontal_inversion"`
			EmblemPatternVerticalInversion   bool `json:"emblem_pattern_vertical_inversion"`
			EmblemProcessing                 int  `json:"emblem_processing"`
			EmblemSymbol1                    int  `json:"emblem_symbol1"`
			EmblemSymbol1Clipping            bool `json:"emblem_symbol1_clipping"`
			EmblemSymbol1Color               int  `json:"emblem_symbol1_color"`
			EmblemSymbol1HorizontalInversion bool `json:"emblem_symbol1_horizontal_inversion"`
			EmblemSymbol1VerticalInversion   bool `json:"emblem_symbol1_vertical_inversion"`
			EmblemSymbol2                    int  `json:"emblem_symbol2"`
			EmblemSymbol2Clipping            bool `json:"emblem_symbol2_clipping"`
			EmblemSymbol2Color               int  `json:"emblem_symbol2_color"`
			EmblemSymbol2HorizontalInversion bool `json:"emblem_symbol2_horizontal_inversion"`
			EmblemSymbol2VerticalInversion   bool `json:"emblem_symbol2_vertical_inversion"`
		} `json:"emblem"`
		Leader struct {
			FighterID  string `json:"fighter_id"`
			PlatformID int    `json:"platform_id"`
			ShortID    int    `json:"short_id"`
		} `json:"leader"`
	} `json:"main_circle"`
	MaxContentPlayTime struct {
		ContentType int `json:"content_type"`
		PlayTime    int `json:"play_time"`
	} `json:"max_content_play_time"`
	MobileLinkage    bool `json:"mobile_linkage"`
	OnlineStatusInfo struct {
		BattlehubAdmissionRestriction              int    `json:"battlehub_admission_restriction"`
		BattlehubID                                string `json:"battlehub_id"`
		BattlehubPlatformID                        int    `json:"battlehub_platform_id"`
		BattlehubRegionID                          int    `json:"battlehub_region_id"`
		BattlehubServerNo                          int    `json:"battlehub_server_no"`
		CustomRoomMasterShortID                    int    `json:"custom_room_master_short_id"`
		CustomRoomPlatformID                       int    `json:"custom_room_platform_id"`
		CustomRoomPublishSetting                   int    `json:"custom_room_publish_setting"`
		CustomRoomRegionID                         int    `json:"custom_room_region_id"`
		CustomRoomRequiredNetworkConnectionQuality int    `json:"custom_room_required_network_connection_quality"`
		CustomRoomRequiredPassCode                 bool   `json:"custom_room_required_pass_code"`
		CustomRoomRoomID                           string `json:"custom_room_room_id"`
		OnlineStatus                               int    `json:"online_status"`
		OnlineStatusData                           struct {
			OnlineStatusName string `json:"online_status_name"`
			OnlineStatusType int    `json:"online_status_type"`
		} `json:"online_status_data"`
		BattlehubRegionName       string `json:"battlehub_region_name"`
		BattlehubFormatedServerNo string `json:"battlehub_formated_server_no"`
	} `json:"online_status_info"`
	PersonalInfo struct {
		FighterID        string `json:"fighter_id"`
		PlatformID       int    `json:"platform_id"`
		ShortID          int64  `json:"short_id"`
		PlatformName     string `json:"platform_name"`
		PlatformToolName string `json:"platform_tool_name"`
	} `json:"personal_info"`
	PlayTimeZone struct {
		EndHour     int `json:"end_hour"`
		EndMinute   int `json:"end_minute"`
		StartHour   int `json:"start_hour"`
		StartMinute int `json:"start_minute"`
	} `json:"play_time_zone"`
	ProfileComment struct {
		ProfileTagID int `json:"profile_tag_id"`
		TagOptionID  int `json:"tag_option_id"`
	} `json:"profile_comment"`
	TitlePlate                int    `json:"title_plate"`
	FavoriteCharacterName     string `json:"favorite_character_name"`
	FavoriteCharacterToolName string `json:"favorite_character_tool_name"`
	TitleData                 struct {
		TitleDataID        int    `json:"title_data_id"`
		TitleDataGradeID   int    `json:"title_data_grade_id"`
		TitleDataGradeName string `json:"title_data_grade_name"`
		TitleDataPlateID   int    `json:"title_data_plate_id"`
		TitleDataPlateName string `json:"title_data_plate_name"`
		TitleDataVal       string `json:"title_data_val"`
	} `json:"title_data"`
	HomeName string `json:"home_name"`
}

func (bl *BattleLog) GetLP() int {
	return bl.FighterBannerInfo.FavoriteCharacterLeagueInfo.LeaguePoint
}

func (bl *BattleLog) GetMR() int {
	return bl.FighterBannerInfo.FavoriteCharacterLeagueInfo.MasterRating
}

func (bl *BattleLog) GetCharacter() string {
	return bl.FighterBannerInfo.FavoriteCharacterName
}

func (bl *BattleLog) GetCFN() string {
	return bl.FighterBannerInfo.PersonalInfo.FighterID
}

func (bl *BattleLog) GetUserCode() string {
	return strconv.FormatInt(bl.FighterBannerInfo.PersonalInfo.ShortID, 10)
}

type Replay struct {
	BattleVersion           int        `json:"battle_version"`
	IsRegistered            bool       `json:"is_registered"`
	Player1Info             PlayerInfo `json:"player1_info"`
	Player2Info             PlayerInfo `json:"player2_info"`
	ReplayBattleSubType     int        `json:"replay_battle_sub_type"`
	ReplayBattleType        int        `json:"replay_battle_type"`
	ReplayID                string     `json:"replay_id"`
	UploadedAt              int        `json:"uploaded_at"`
	Views                   int        `json:"views"`
	ReplayBattleTypeName    string     `json:"replay_battle_type_name"`
	ReplayBattleSubTypeName string     `json:"replay_battle_sub_type_name"`
}

type SearchResult struct {
	AssetPrefix   string   `json:"assetPrefix"`
	BuildID       string   `json:"buildId"`
	DefaultLocale string   `json:"defaultLocale"`
	Gssp          bool     `json:"gssp"`
	IsFallback    bool     `json:"isFallback"`
	Locale        string   `json:"locale"`
	Locales       []string `json:"locales"`
	Page          string   `json:"page"`
	Props         struct {
		NSsp      bool `json:"__N_SSP"`
		PageProps struct {
			Lang       string `json:"__lang"`
			Namespaces struct {
				Common struct {
					T403Text                         string `json:"[t]403_text"`
					T403Title                        string `json:"[t]403_title"`
					TABOUTUSER                       string `json:"[t]ABOUT USER"`
					TAddAreaEnd                      string `json:"[t]add_area_end"`
					TAddAreaStart                    string `json:"[t]add_area_start"`
					TBack                            string `json:"[t]back"`
					TBacklerNavCapcomFightersNetwork string `json:"[t]backler__nav__capcom_fighters_network"`
					TBacklerNavCfn                   string `json:"[t]backler__nav__cfn"`
					TBacklerNavEvent                 string `json:"[t]backler__nav__event"`
					TBacklerNavFightersclub          string `json:"[t]backler__nav__fightersclub"`
					TBacklerNavFighterslist          string `json:"[t]backler__nav__fighterslist"`
					TBacklerNavInfo                  string `json:"[t]backler__nav__info"`
					TBacklerNavInformation           string `json:"[t]backler__nav__information"`
					TBacklerNavMore                  string `json:"[t]backler__nav__more"`
					TBacklerNavOther                 string `json:"[t]backler__nav__other"`
					TBacklerNavProfile               string `json:"[t]backler__nav__profile"`
					TBacklerNavRanking               string `json:"[t]backler__nav__ranking"`
					TBacklerNavReward                string `json:"[t]backler__nav__reward"`
					TBattleDraw                      string `json:"[t]battle_draw"`
					TBattleLose                      string `json:"[t]battle_lose"`
					TBattleWin                       string `json:"[t]battle_win"`
					TBrowserTimezone                 string `json:"[t]browser_timezone"`
					TCancel                          string `json:"[t]cancel"`
					TChange                          string `json:"[t]change"`
					TCHARACTER                       string `json:"[t]CHARACTER"`
					TClose                           string `json:"[t]close"`
					TClubPointsUnit                  string `json:"[t]club_points_unit"`
					TConfirm                         string `json:"[t]confirm"`
					TDeleteAreaEnd                   string `json:"[t]delete_area_end"`
					TDeleteAreaStart                 string `json:"[t]delete_area_start"`
					TErrorLogin                      string `json:"[t]error_login"`
					TFAVORITECHARACTER               string `json:"[t]FAVORITE CHARACTER"`
					TFooterSupportFaq                string `json:"[t]footer_support_faq"`
					TFooterSupportManual             string `json:"[t]footer_support_manual"`
					TFooterSupportTerms              string `json:"[t]footer_support_terms"`
					TGamemodeBattleHub               string `json:"[t]gamemode-battle_hub"`
					TGamemodeFightingGround          string `json:"[t]gamemode-fighting_ground"`
					TGamemodeWorldTour               string `json:"[t]gamemode-world_tour"`
					THeaderNavEsports                string `json:"[t]header__nav__esports"`
					THeaderNavEsportsCpt             string `json:"[t]header__nav__esports__cpt"`
					THeaderNavEsportsSfl             string `json:"[t]header__nav__esports__sfl"`
					THeaderNavLogin                  string `json:"[t]header__nav__login"`
					THeaderNavLogout                 string `json:"[t]header__nav__logout"`
					THeaderNavOfficial               string `json:"[t]header__nav__official"`
					THeaderNavOfficialCharacter      string `json:"[t]header__nav__official__character"`
					THeaderNavOfficialColumn         string `json:"[t]header__nav__official__column"`
					THeaderNavOfficialGame           string `json:"[t]header__nav__official__game"`
					THeaderNavOfficialGameMode       string `json:"[t]header__nav__official__game__mode"`
					THeaderNavOfficialGameModeBh     string `json:"[t]header__nav__official__game__mode__bh"`
					THeaderNavOfficialGameModeFg     string `json:"[t]header__nav__official__game__mode__fg"`
					THeaderNavOfficialGameModeWt     string `json:"[t]header__nav__official__game__mode__wt"`
					THeaderNavOfficialNews           string `json:"[t]header__nav__official__news"`
					THeaderNavOfficialNewsAll        string `json:"[t]header__nav__official__news__all"`
					THeaderNavOfficialNewsCampaign   string `json:"[t]header__nav__official__news__campaign"`
					THeaderNavOfficialNewsEsports    string `json:"[t]header__nav__official__news__esports"`
					THeaderNavOfficialNewsEvent      string `json:"[t]header__nav__official__news__event"`
					THeaderNavOfficialNewsGame       string `json:"[t]header__nav__official__news__game"`
					THeaderNavOfficialNewsGoods      string `json:"[t]header__nav__official__news__goods"`
					THeaderNavOfficialProduct        string `json:"[t]header__nav__official__product"`
					THeaderNavOfficialProductDlc     string `json:"[t]header__nav__official__product__dlc"`
					THeaderNavOfficialProductMain    string `json:"[t]header__nav__official__product__main"`
					THeaderNavOfficialProductSeason  string `json:"[t]header__nav__official__product__season"`
					THeaderNavPointHistory           string `json:"[t]header__nav__point__history"`
					THeaderNavProfile                string `json:"[t]header__nav__profile"`
					THeaderNavSeries                 string `json:"[t]header__nav__series"`
					THeaderNavSettingBase            string `json:"[t]header__nav__setting__base"`
					THeaderNavSettingProfile         string `json:"[t]header__nav__setting__profile"`
					THeaderNavSettings               string `json:"[t]header__nav__settings"`
					THeaderNavSupport                string `json:"[t]header__nav__support"`
					THeaderNavSupportFaq             string `json:"[t]header__nav__support__faq"`
					THeaderNavSupportManual          string `json:"[t]header__nav__support__manual"`
					TInfoCaution                     string `json:"[t]info__caution"`
					TInfoListGenre                   string `json:"[t]info__list__genre"`
					TInfoListGenreTitle              string `json:"[t]info__list__genre__title"`
					TInfoListPlatformTitle           string `json:"[t]info__list__platform__title"`
					TInfoListRatingTitle             string `json:"[t]info__list__rating__title"`
					TInfoListRatingTitleSchedule     string `json:"[t]info__list__rating__title__schedule"`
					TInfoListRelease                 string `json:"[t]info__list__release"`
					TInfoListReleaseTitle            string `json:"[t]info__list__release__title"`
					TInfoListTitle                   string `json:"[t]info__list__title"`
					TIngameFooter                    string `json:"[t]ingame__footer"`
					TIngameFooter00                  string `json:"[t]ingame__footer00"`
					TIngameFooter01                  string `json:"[t]ingame__footer01"`
					TLangAll                         string `json:"[t]lang-all"`
					TLangAr                          string `json:"[t]lang-ar"`
					TLangDe                          string `json:"[t]lang-de"`
					TLangEnUs                        string `json:"[t]lang-en-us"`
					TLangEsEs                        string `json:"[t]lang-es-es"`
					TLangFr                          string `json:"[t]lang-fr"`
					TLangIt                          string `json:"[t]lang-it"`
					TLangJaJp                        string `json:"[t]lang-ja-jp"`
					TLangKoKr                        string `json:"[t]lang-ko-kr"`
					TLangPl                          string `json:"[t]lang-pl"`
					TLangPtBr                        string `json:"[t]lang-pt-br"`
					TLangRu                          string `json:"[t]lang-ru"`
					TLangSelect                      string `json:"[t]lang-select"`
					TLangZhHans                      string `json:"[t]lang-zh-hans"`
					TLangZhHant                      string `json:"[t]lang-zh-hant"`
					TLeaguePointUnit                 string `json:"[t]league_point_unit"`
					TLinksBbc                        string `json:"[t]links_bbc"`
					TLinksSeries                     string `json:"[t]links_series"`
					TLinksSf6                        string `json:"[t]links_sf6"`
					TListCrossplayActive             string `json:"[t]list_crossplay_active"`
					TListCrossplayDisactive          string `json:"[t]list_crossplay_disactive"`
					TLogin                           string `json:"[t]login"`
					TMaintenanceLink                 string `json:"[t]maintenance_link"`
					TMaintenanceText                 string `json:"[t]maintenance_text"`
					TMaintenanceTitle                string `json:"[t]maintenance_title"`
					TMore                            string `json:"[t]more"`
					TNext                            string `json:"[t]next"`
					TNo                              string `json:"[t]no"`
					TNoData                          string `json:"[t]no_data"`
					TNotRegisteredRegister           string `json:"[t]not_registered_register"`
					TNotRegisteredToWellcome         string `json:"[t]not_registered_to_wellcome"`
					TPagetop                         string `json:"[t]pagetop"`
					TPLAYINGTIME                     string `json:"[t]PLAYING TIME"`
					TPraise30DayDisplay              string `json:"[t]praise_30day_display"`
					TPraiseExit                      string `json:"[t]praise_exit"`
					TPraisePointHistory              string `json:"[t]praise_point_history"`
					TPraisePointUnit                 string `json:"[t]praise_point_unit"`
					TPraiseTitleBattleHub            string `json:"[t]praise_title_battle_hub"`
					TPraiseTitleDeparture            string `json:"[t]praise_title_departure"`
					TPraiseTitleEveryday             string `json:"[t]praise_title_everyday"`
					TPraiseTitleFightingGround       string `json:"[t]praise_title_fighting_ground"`
					TPraiseTitleTweet                string `json:"[t]praise_title_tweet"`
					TPraiseTitleView                 string `json:"[t]praise_title_view"`
					TPraiseTitleWorldTogether        string `json:"[t]praise_title_world_together"`
					TPraiseTitleWorldTour            string `json:"[t]praise_title_world_tour"`
					TQuickNavCharacter               string `json:"[t]quick__nav__character"`
					TQuickNavHeadline                string `json:"[t]quick__nav__headline"`
					TSearchCharacterSelect           string `json:"[t]search_character_select"`
					TSearchHomeArea                  string `json:"[t]search_home_area"`
					TSeparator                       string `json:"[t]separator"`
					TSkip                            string `json:"[t]skip"`
					TSNS                             string `json:"[t]SNS"`
					TSnsDc                           string `json:"[t]sns__dc"`
					TSnsFb                           string `json:"[t]sns__fb"`
					TSnsIg                           string `json:"[t]sns__ig"`
					TSnsTw                           string `json:"[t]sns__tw"`
					TSnsYt                           string `json:"[t]sns__yt"`
					TTitle                           string `json:"[t]title"`
					TYes                             string `json:"[t]yes"`
					TYYYYMMDD                        string `json:"[t]YYYY/MM/DD"`
					TYYYYMMDDHHmm                    string `json:"[t]YYYY/MM/DD HHmm"`
					T                                string `json:"[t]シリーズ歴"`
				} `json:"common"`
				Error struct {
					T404Text               string `json:"[t]404_text"`
					T404Title              string `json:"[t]404_title"`
					T500Text               string `json:"[t]500_text"`
					T500Title              string `json:"[t]500_title"`
					TDescription           string `json:"[t]description"`
					TErrorExpiredText      string `json:"[t]error_expired_text"`
					TErrorLinkFaq          string `json:"[t]error_link_faq"`
					TErrorLinkTop          string `json:"[t]error_link_top"`
					TErrorNoplaydataText   string `json:"[t]error_noplaydata_text"`
					TErrorNoplaydataTitle  string `json:"[t]error_noplaydata_title"`
					TErrorSystemText       string `json:"[t]error_system_text"`
					TLoginLogoutBtn        string `json:"[t]login_logout_btn"`
					TLoginReloginErrorText string `json:"[t]login_relogin_error_text"`
					TTitle                 string `json:"[t]title"`
				} `json:"error"`
				Fighter struct {
					TAddAreaEnd                 string `json:"[t]add_area_end"`
					TAddAreaStart               string `json:"[t]add_area_start"`
					TDeleteAreaEnd              string `json:"[t]delete_area_end"`
					TDeleteAreaStart            string `json:"[t]delete_area_start"`
					TDescription                string `json:"[t]description"`
					TFighterNavApplication      string `json:"[t]fighter__nav__application"`
					TFighterNavBlock            string `json:"[t]fighter__nav__block"`
					TFighterNavFollow           string `json:"[t]fighter__nav__follow"`
					TFighterNavFriend           string `json:"[t]fighter__nav__friend"`
					TFighterNavLive             string `json:"[t]fighter__nav__live"`
					TFighterNavMvp              string `json:"[t]fighter__nav__mvp"`
					TFighterNavReceivedrequest  string `json:"[t]fighter__nav__receivedrequest"`
					TFighterNavRequest          string `json:"[t]fighter__nav__request"`
					TFighterNavSearch           string `json:"[t]fighter__nav__search"`
					TFighterNavWinning          string `json:"[t]fighter__nav__winning"`
					TFighterArea                string `json:"[t]fighter_area"`
					TFighterError               string `json:"[t]fighter_error"`
					TFighterLogin               string `json:"[t]fighter_login"`
					TFighterNotBlock            string `json:"[t]fighter_not_block"`
					TFighterNotConditions       string `json:"[t]fighter_not_conditions"`
					TFighterNotFollow           string `json:"[t]fighter_not_follow"`
					TFighterNotFriend           string `json:"[t]fighter_not_friend"`
					TFighterNotRequest          string `json:"[t]fighter_not_request"`
					TFighterNotUser             string `json:"[t]fighter_not_user"`
					TFighterSort                string `json:"[t]fighter_sort"`
					TFighterSortDateApplication string `json:"[t]fighter_sort_date_application"`
					TFighterSortDateRegister    string `json:"[t]fighter_sort_date_register"`
					TFighterSortGamemode        string `json:"[t]fighter_sort_gamemode"`
					TFighterSortLastplay        string `json:"[t]fighter_sort_lastplay"`
					TFighterSortLeague          string `json:"[t]fighter_sort_league"`
					TFighterSortWinning         string `json:"[t]fighter_sort_winning"`
					TFighterTimezone            string `json:"[t]fighter_timezone"`
					TFighterWinning0            string `json:"[t]fighter_winning_0"`
					TFighterWinning1            string `json:"[t]fighter_winning_1"`
					TNotRegisteredEx1           string `json:"[t]not_registered_ex1"`
					TNotRegisteredEx2           string `json:"[t]not_registered_ex2"`
					TNotRegisteredEx3           string `json:"[t]not_registered_ex3"`
					TNotRegisteredEx4           string `json:"[t]not_registered_ex4"`
					TNotRegisteredEx5           string `json:"[t]not_registered_ex5"`
					TNotRegisteredEx6           string `json:"[t]not_registered_ex6"`
					TNotRegisteredLead          string `json:"[t]not_registered_lead"`
					TSearch4Word                string `json:"[t]search_4word"`
					TSearch9Word                string `json:"[t]search_9word"`
					TSearchAbility              string `json:"[t]search_ability"`
					TSearchBack                 string `json:"[t]search_back"`
					TSearchCharacter            string `json:"[t]search_character"`
					TSearchCommentTag           string `json:"[t]search_comment_tag"`
					TSearchCrossplay            string `json:"[t]search_crossplay"`
					TSearchDetail               string `json:"[t]search_detail"`
					TSearchEasy                 string `json:"[t]search_easy"`
					TSearchHomecategory         string `json:"[t]search_homecategory"`
					TSearchHomeID               string `json:"[t]search_homeId"`
					TSearchID                   string `json:"[t]search_id"`
					TSearchLastPlay             string `json:"[t]search_last_play"`
					TSearchLeagueMax            string `json:"[t]search_league_max"`
					TSearchLeagueMin            string `json:"[t]search_league_min"`
					TSearchPlayername           string `json:"[t]search_playername"`
					TSearchPlaytimezone         string `json:"[t]search_playtimezone"`
					TSearchTagBlank             string `json:"[t]search_tag_blank"`
					TSearchTagCharacter         string `json:"[t]search_tag_character"`
					TSearchTagInput             string `json:"[t]search_tag_input"`
					TSearchTagProficiency       string `json:"[t]search_tag_proficiency"`
					TSearchTagTimezone          string `json:"[t]search_tag_timezone"`
					TSearchTagTitle             string `json:"[t]search_tag_title"`
					TSearchUsercode             string `json:"[t]search_usercode"`
					TSeparator                  string `json:"[t]separator"`
					TTitle                      string `json:"[t]title"`
				} `json:"fighter"`
				Terms struct {
					TTermsDate   string `json:"[t]terms_date"`
					TTermsRead   string `json:"[t]terms_read"`
					TTermsText11 string `json:"[t]terms_text1-1"`
					TTermsText12 string `json:"[t]terms_text1-2"`
					TTermsText21 string `json:"[t]terms_text2-1"`
					TTermsText22 string `json:"[t]terms_text2-2"`
					TTermsText31 string `json:"[t]terms_text3-1"`
					TTermsText32 string `json:"[t]terms_text3-2"`
					TTermsText41 string `json:"[t]terms_text4-1"`
					TTermsText51 string `json:"[t]terms_text5-1"`
					TTermsText61 string `json:"[t]terms_text6-1"`
					TTermsText71 string `json:"[t]terms_text7-1"`
					TTermsText81 string `json:"[t]terms_text8-1"`
					TTermsTitle1 string `json:"[t]terms_title1"`
					TTermsTitle2 string `json:"[t]terms_title2"`
					TTermsTitle3 string `json:"[t]terms_title3"`
					TTermsTitle4 string `json:"[t]terms_title4"`
					TTermsTitle5 string `json:"[t]terms_title5"`
					TTermsTitle6 string `json:"[t]terms_title6"`
					TTermsTitle7 string `json:"[t]terms_title7"`
					TTermsTitle8 string `json:"[t]terms_title8"`
				} `json:"terms"`
			} `json:"__namespaces"`
			Common struct {
				AppEnv    string `json:"appEnv"`
				IsError   bool   `json:"isError"`
				LoginUser struct {
					FighterID  string `json:"fighterId"`
					Flg        bool   `json:"flg"`
					PlatformID int    `json:"platformId"`
					RegionID   int    `json:"regionId"`
					ShortID    int64  `json:"shortId"`
				} `json:"loginUser"`
				StatusCode int `json:"statusCode"`
			} `json:"common"`
			FighterBannerList []FighterBanner `json:"fighter_banner_list"`
			Page              int             `json:"page"`
			SearchParams      struct {
				FighterID string `json:"fighter_id"`
			} `json:"search_params"`
		} `json:"pageProps"`
	} `json:"props"`
	Query struct {
		FighterID string `json:"fighter_id"`
	} `json:"query"`
	ScriptLoader []any `json:"scriptLoader"`
}

func MapFighterBannerToPlayer(fb *FighterBanner) model.Player {
	return model.Player{
		DisplayName:       fb.PersonalInfo.FighterID,
		Code:              fmt.Sprint(fb.PersonalInfo.ShortID),
		FavoriteCharacter: fb.FavoriteCharacterName,
	}
}

var COUNTRIES = []string{
	`Afghanistan`,
	`Aland Islands`,
	`Albania`,
	`Algeria`,
	`American Samoa`,
	`Andorra`,
	`Angola`,
	`Anguilla`,
	`Antarctica`,
	`Argentina`,
	`Armenia`,
	`Aruba`,
	`Australia`,
	`Austria`,
	`Azerbaijan`,
	`Bahamas`,
	`Bahrain`,
	`Bangladesh`,
	`Barbados`,
	`Belarus`,
	`Belgium`,
	`Belize`,
	`Benin`,
	`Bermuda`,
	`Bhutan`,
	`Bolivia`,
	`Bonaire, Sint Eustatius and Saba`,
	`Bosnia and Herzegovina`,
	`Botswana`,
	`Bouvet Island`,
	`Brazil`,
	`British Indian Ocean Territory`,
	`British Virgin Islands`,
	`Brunei Darussalam`,
	`Bulgaria`,
	`Burkina Faso`,
	`Burundi`,
	`Cambodia`,
	`Cameroon`,
	`Canada`,
	`Cabo Verde`,
	`Cayman Islands`,
	`Central African Republic`,
	`Chad`,
	`Chile`,
	`China`,
	`Christmas Island`,
	`Cocos (Keeling) Islands`,
	`Colombia`,
	`Comoros`,
	`Congo`,
	`Congo, the Democratic Republic of the`,
	`Cook Islands`,
	`Costa Rica`,
	`Cote D'ivoire (Ivory Coast)`,
	`Croatia`,
	`Cuba`,
	`Curacao`,
	`Cyprus`,
	`Czech Republic`,
	`Denmark`,
	`Djibouti`,
	`Dominica`,
	`Dominican Republic`,
	`Ecuador`,
	`Egypt`,
	`El Salvador`,
	`Equatorial Guinea`,
	`Eritrea`,
	`Estonia`,
	`Ethiopia`,
	`Falkland Islands (Malvinas)`,
	`Faroe Islands`,
	`Fiji`,
	`Finland`,
	`France`,
	`French Guiana`,
	`French Polynesia`,
	`French Southern Territories`,
	`Gabon`,
	`Gambia`,
	`Georgia`,
	`Germany`,
	`Ghana`,
	`Gibraltar`,
	`Greece`,
	`Greenland`,
	`Grenada`,
	`Guadeloupe`,
	`Guam`,
	`Guatemala`,
	`Guernsey`,
	`Guinea`,
	`Guinea-Bissau`,
	`Guyana`,
	`Haiti`,
	`Honduras`,
	`Hong Kong`,
	`Hungary`,
	`Iceland`,
	`India`,
	`Indonesia`,
	`Iraq`,
	`Ireland`,
	`Islamic Republic of Iran`,
	`Isle of Man`,
	`Israel`,
	`Italy`,
	`Jamaica`,
	`日本`,
	`Jersey`,
	`Jordan`,
	`Kazakhstan`,
	`Kenya`,
	`Kiribati`,
	`Korea, Democratic People's Republic of`,
	`Korea, Republic of`,
	`Kuwait`,
	`Kyrgyzstan`,
	`Laos`,
	`Latvia`,
	`Lebanon`,
	`Lesotho`,
	`Liberia`,
	`Libya`,
	`Liechtenstein`,
	`Lithuania`,
	`Luxembourg`,
	`Macau`,
	`Macedonia, The Former Yugoslav Republic of`,
	`Madagascar`,
	`Malawi`,
	`Malaysia`,
	`Maldives`,
	`Mali`,
	`Malta`,
	`Marshall Islands`,
	`Martinique`,
	`Mauritania`,
	`Mauritius`,
	`Mayotte`,
	`Mexico`,
	`Micronesia`,
	`Moldova, Republic of`,
	`Monaco`,
	`Mongolia`,
	`Montserrat`,
	`Montenegro`,
	`Morocco`,
	`Mozambique`,
	`Myanmar`,
	`Namibia`,
	`Nauru`,
	`Nepal`,
	`Netherlands`,
	`New Caledonia`,
	`New Zealand`,
	`Nicaragua`,
	`Niger`,
	`Nigeria`,
	`Niue`,
	`Norfolk Island`,
	`Northern Mariana Islands`,
	`Norway`,
	`Oman`,
	`Pakistan`,
	`Palau`,
	`Palestinian Territory, Occupied`,
	`Panama`,
	`Papua New Guinea`,
	`Paraguay`,
	`Peru`,
	`Philippines`,
	`Pitcairn`,
	`Poland`,
	`Portugal`,
	`Puerto Rico`,
	`Qatar`,
	`Reunion`,
	`Romania`,
	`Russian Federation`,
	`Rwanda`,
	`Saint Barthelemy`,
	`Saint Lucia`,
	`Saint Martin (French part)`,
	`Samoa`,
	`San Marino`,
	`Saudi Arabia`,
	`Senegal`,
	`Serbia`,
	`Seychelles`,
	`Sierra Leone`,
	`Singapore`,
	`Sint Maarten (Dutch part)`,
	`Slovakia`,
	`Slovenia`,
	`Solomon Islands`,
	`Somalia`,
	`South Africa`,
	`South Georgia and the South Sandwich Islands`,
	`South Sudan`,
	`Spain`,
	`Sri Lanka`,
	`St. Helena`,
	`St. Kitts and Nevis`,
	`Sudan`,
	`Suriname`,
	`Swaziland`,
	`Sweden`,
	`Switzerland`,
	`Syrian Arab Republic`,
	`Taiwan`,
	`Tajikistan`,
	`Tanzania, United Republic of`,
	`Thailand`,
	`Timor-Leste`,
	`Togo`,
	`Tokelau`,
	`Tonga`,
	`Tunisia`,
	`Turkey`,
	`Turkmenistan`,
	`Tuvalu`,
	`Uganda`,
	`Ukraine`,
	`United Arab Emirates`,
	`United Kingdom`,
	`United States`,
	`United States Minor Outlying`,
	`United States Virgin Islands`,
	`Uruguay`,
	`Uzbekistan`,
	`Vanuatu`,
	`Vatican City State (Holy See)`,
	`Venezuela`,
	`Viet Nam`,
	`Western Sahara`,
	`Yemen`,
	`Zambia`,
	`Zimbabwe`,
}
