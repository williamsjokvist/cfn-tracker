package sf6

type BattleLog struct {
	PageProps struct {
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
					ShortID    int64  `json:"short_id"`
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
		CurrentPage int `json:"current_page"`
		ReplayList  []struct {
			BattleVersion int  `json:"battle_version"`
			IsRegistered  bool `json:"is_registered"`
			Player1Info   struct {
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
						ShortID    int64  `json:"short_id"`
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
			} `json:"player1_info"`
			Player2Info struct {
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
			} `json:"player2_info"`
			ReplayBattleSubType     int    `json:"replay_battle_sub_type"`
			ReplayBattleType        int    `json:"replay_battle_type"`
			ReplayID                string `json:"replay_id"`
			UploadedAt              int    `json:"uploaded_at"`
			Views                   int    `json:"views"`
			ReplayBattleTypeName    string `json:"replay_battle_type_name"`
			ReplayBattleSubTypeName string `json:"replay_battle_sub_type_name"`
		} `json:"replay_list"`
		TotalPage int   `json:"total_page"`
		Sid       int64 `json:"sid"`
		WeekList  []struct {
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
	} `json:"pageProps"`
	NSsp bool `json:"__N_SSP"`
}

type SearchResult struct {
	PageProps struct {
		FighterBannerList []struct {
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
		} `json:"fighter_banner_list"`
		SearchParams struct {
			FighterID string `json:"fighter_id"`
		} `json:"search_params"`
		Page   int `json:"page"`
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
	} `json:"pageProps"`
	NSsp bool `json:"__N_SSP"`
}
