export namespace model {
	
	export enum ThemeName {
	    DEFAULT = "default",
	    ENTH = "enth",
	    TEKKEN = "tekken",
	}
	export enum GameType {
	    STREET_FIGHTER_6 = "sf6",
	    TEKKEN_8 = "t8",
	}
	export enum ErrorLocalizationKey {
	    errUnknown = "errUnknown",
	    errSelectGame = "errSelectGame",
	    errAuth = "errAuth",
	    errGetLatestSession = "errGetLatestSession",
	    errGetUser = "errGetUser",
	    errGetMatches = "errGetMatches",
	    errSaveLocale = "errSaveLocale",
	    errCheckForUpdate = "errCheckForUpdate",
	    errGetGuiConfig = "errGetGuiConfig",
	    errSaveTheme = "errSaveTheme",
	    errSaveUser = "errSaveUser",
	    errSaveSidebar = "errSaveSidebar",
	    errGetSessions = "errGetSessions",
	    errGetTranslations = "errGetTranslations",
	    errGetSessionStatistics = "errGetSessionStatistics",
	    errCreateSession = "errCreateSession",
	    errOpenResultsDirectory = "errOpenResultsDirectory",
	    errReadThemeCSS = "errReadThemeCSS",
	}
	export interface FGCTrackerError {
	    localizationKey: ErrorLocalizationKey;
	    message: string;
	}
	export interface GUIConfig {
	    locale: string;
	    theme: ThemeName;
	    sidebar: boolean;
	}
	export interface Localization {
	    appVersion: string;
	    source: string;
	    about: string;
	    settings: string;
	    tracking: string;
	    history: string;
	    language: string;
	    changelog: string;
	    startTracking: string;
	    cfnName: string;
	    tekkenId: string;
	    start: string;
	    opponent: string;
	    character: string;
	    lpGain: string;
	    mrGain: string;
	    deleteLog: string;
	    goBack: string;
	    loading: string;
	    wins: string;
	    losses: string;
	    winRate: string;
	    stop: string;
	    files: string;
	    enterCfnName: string;
	    enterTekkenId: string;
	    result: string;
	    replayId: string;
	    copy: string;
	    time: string;
	    winStreak: string;
	    newVersionAvailable: string;
	    pause: string;
	    unpause: string;
	    statistics: string;
	    date: string;
	    minimize: string;
	    restoreSession: string;
	    exportLog: string;
	    league: string;
	    pickGame: string;
	    follow: string;
	    continueStep: string;
	    placement: string;
	    cfnError: string;
	    output: string;
	    theme: string;
	    defaultTheme: string;
	    pickTheme: string;
	    usingBrowserSource: string;
	    copyBrowserSourceLink: string;
	    browserSourceDescription: string;
	    displayStats: string;
	    importFiles: string;
	    obsCustomize: string;
	    lastMatch: string;
	    statsWillBeDisplayed: string;
	    selectTheme: string;
	    started: string;
	    user: string;
	    matchesWon: string;
	    matchesLost: string;
	    sessions: string;
	    refresh: string;
	    cooldown: string;
	    T8_BEGINNER: string;
	    T8_1ST_DAN: string;
	    T8_2ND_DAN: string;
	    T8_FIGHTER: string;
	    T8_STRATEGIST: string;
	    T8_COMBATANT: string;
	    T8_BRAWLER: string;
	    T8_RANGER: string;
	    T8_CAVALRY: string;
	    T8_WARRIOR: string;
	    T8_ASSAILANT: string;
	    T8_DOMINATOR: string;
	    T8_VANQUISHER: string;
	    T8_DESTROYER: string;
	    T8_ELIMINATOR: string;
	    T8_GARYU: string;
	    T8_SHINRYU: string;
	    T8_TENRYU: string;
	    T8_MIGHTY_RULER: string;
	    T8_FLAME_RULER: string;
	    T8_BATTLE_RULER: string;
	    T8_FUJIN: string;
	    T8_RAIJIN: string;
	    T8_KISHIN: string;
	    T8_BUSHIN: string;
	    T8_TEKKEN_KING: string;
	    T8_TEKKEN_EMPEROR: string;
	    T8_TEKKEN_GOD: string;
	    T8_TEKKEN_GOD_SUPREME: string;
	    T8_GOD_OF_DESTRUCTION: string;
	    errUnknown: string;
	    errSelectGame: string;
	    errAuth: string;
	    errGetLatestSession: string;
	    errGetUser: string;
	    errGetMatches: string;
	    errSaveLocale: string;
	    errCheckForUpdate: string;
	    errGetGuiConfig: string;
	    errSaveTheme: string;
	    errSaveUser: string;
	    errSaveSidebar: string;
	    errGetSessions: string;
	    errGetTranslations: string;
	    errGetSessionStatistics: string;
	    errCreateSession: string;
	    errOpenResultsDirectory: string;
	    errReadThemeCSS: string;
	}
	export interface Match {
	    userId: string;
	    userName: string;
	    sessionId: number;
	    replayId: string;
	    character: string;
	    lp: number;
	    lpGain: number;
	    mr: number;
	    mrGain: number;
	    opponent: string;
	    opponentCharacter: string;
	    opponentLp: number;
	    opponentMr: number;
	    opponentLeague: string;
	    victory: boolean;
	    date: string;
	    time: string;
	    winStreak: number;
	    wins: number;
	    losses: number;
	    winRate: number;
	}
	export interface Session {
	    id: number;
	    userId: string;
	    userName: string;
	    createdAt: string;
	    lp: number;
	    mr: number;
	    matches: Match[];
	    matchesWon: number;
	    matchesLost: number;
	    endingLp: number;
	    endingMr: number;
	    startingLp: number;
	    startingMr: number;
	    lpGain: number;
	    mrGain: number;
	}
	export interface SessionMonth {
	    Date: string;
	    Count: number;
	}
	export interface SessionsStatistics {
	    Months: SessionMonth[];
	}
	export interface Theme {
	    name: string;
	    css: string;
	}
	export interface User {
	    id: number;
	    displayName: string;
	    code: string;
	    LP: number;
	    MR: number;
	}

}

