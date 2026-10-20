export namespace model {
	
	export enum ErrorLocalizationKey {
	    errAuth = "errAuth",
	    errAuthManualLoginFailed = "errAuthManualLoginFailed",
	    errBackupData = "errBackupData",
	    errCheckForUpdate = "errCheckForUpdate",
	    errCreateSession = "errCreateSession",
	    errExport = "errExport",
	    errGetDataLocation = "errGetDataLocation",
	    errGetGuiConfig = "errGetGuiConfig",
	    errGetLatestSession = "errGetLatestSession",
	    errGetMatches = "errGetMatches",
	    errGetSessionStatistics = "errGetSessionStatistics",
	    errGetSessions = "errGetSessions",
	    errGetTranslations = "errGetTranslations",
	    errGetUser = "errGetUser",
	    errImportBackup = "errImportBackup",
	    errInvalidBackup = "errInvalidBackup",
	    errInvalidBrowserSourcePort = "errInvalidBrowserSourcePort",
	    errOpenDataDirectory = "errOpenDataDirectory",
	    errOpenResultsDirectory = "errOpenResultsDirectory",
	    errReadThemeCSS = "errReadThemeCSS",
	    errReportIssue = "errReportIssue",
	    errSaveBrowserSourcePort = "errSaveBrowserSourcePort",
	    errSaveDataDirectory = "errSaveDataDirectory",
	    errSaveLocale = "errSaveLocale",
	    errSaveLogFile = "errSaveLogFile",
	    errSaveSidebar = "errSaveSidebar",
	    errSaveStreamOverlay = "errSaveStreamOverlay",
	    errSaveTheme = "errSaveTheme",
	    errSaveUser = "errSaveUser",
	    errSelectGame = "errSelectGame",
	    errTrackingInProgress = "errTrackingInProgress",
	    errUnknown = "errUnknown",
	}
	export enum GameType {
	    STREET_FIGHTER_6 = "sf6",
	    TEKKEN_8 = "t8",
	}
	export enum ThemeName {
	    DEFAULT = "default",
	    DUSK = "dusk",
	    EMBER = "ember",
	    ENTH = "enth",
	    JADE = "jade",
	    SAKURA = "sakura",
	    TEKKEN = "tekken",
	    TIDE = "tide",
	}
	export interface CharacterStatistics {
	    character: string;
	    wins: number;
	    losses: number;
	    lp: number;
	    mr: number;
	}
	export interface DataLocation {
	    dir: string;
	    defaultDir: string;
	    restartRequired: boolean;
	}
	export interface FGCTrackerError {
	    localizationKey: ErrorLocalizationKey;
	    message: string;
	}
	export interface GUIConfig {
	    locale: string;
	    theme: ThemeName;
	    sidebar: boolean;
	    logFile: boolean;
	    browserSourcePort: number;
	    lastGame: GameType;
	    lastUsers: Record<string, string>;
	}
	export interface HistoryMatch {
	    date: string;
	    time: string;
	    character: string;
	    victory: boolean;
	    lp: number;
	    mr: number;
	}
	export interface IssueReport {
	    title: string;
	    description: string;
	    contact: string;
	}
	export interface Localization {
	    appVersion: string;
	    source: string;
	    about: string;
	    settings: string;
	    tracking: string;
	    sf6Comparison: string;
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
	    loading: string;
	    wins: string;
	    losses: string;
	    winRate: string;
	    stop: string;
	    files: string;
	    enterCfnName: string;
	    enterTekkenId: string;
	    comparisonMetric: string;
	    comparisonTopAvg: string;
	    comparisonCheckingLogin: string;
	    comparisonFetchingBattleLog: string;
	    comparisonFetchingStats: string;
	    comparisonFetchingRanking: string;
	    comparisonFetchingTopPlayer: string;
	    comparisonComparing: string;
	    comparisonInterrupted: string;
	    result: string;
	    replayId: string;
	    copy: string;
	    time: string;
	    winStreak: string;
	    newVersionAvailable: string;
	    pause: string;
	    unpause: string;
	    statistics: string;
	    recentSessions: string;
	    statOverTime: string;
	    netWins: string;
	    notEnoughMatches: string;
	    allTime: string;
	    compareWithTopPlayers: string;
	    noStatistics: string;
	    date: string;
	    minimize: string;
	    writeLogFile: string;
	    browserSourcePort: string;
	    reportIssue: string;
	    reportIssueDescription: string;
	    issueTitle: string;
	    issueDescription: string;
	    issueContact: string;
	    send: string;
	    issueReported: string;
	    viewIssue: string;
	    restoreSession: string;
	    exportLog: string;
	    export: string;
	    league: string;
	    pickGame: string;
	    stopTrackingTitle: string;
	    stopTrackingDescription: string;
	    stopComparisonTitle: string;
	    stopComparisonDescription: string;
	    follow: string;
	    continueStep: string;
	    placement: string;
	    cfnError: string;
	    output: string;
	    theme: string;
	    themeDefault: string;
	    themeEnth: string;
	    themeTekken: string;
	    themeJade: string;
	    themeEmber: string;
	    themeTide: string;
	    themeDusk: string;
	    themeSakura: string;
	    browserSourceThemeBlades: string;
	    browserSourceThemeJaeger: string;
	    browserSourceThemeNord: string;
	    browserSourceThemePills: string;
	    browserSourceThemeComic: string;
	    browserSourceThemeTerminal: string;
	    browserSourceThemeVersus: string;
	    browserSourceThemeDeco: string;
	    browserSourceThemeBroadcast: string;
	    browserSourceThemeNeon: string;
	    browserSourceThemeRetro: string;
	    browserSourceThemeNotebook: string;
	    browserSourceThemeGlass: string;
	    browserSourceThemeMinimal: string;
	    browserSourceThemeMinimalDark: string;
	    browserSourceThemeTicker: string;
	    browserSourceThemeTextOnly: string;
	    browserSourceThemeSticker: string;
	    browserSourceThemeGlitch: string;
	    browserSourceThemeChips: string;
	    browserSourceThemeMostWanted: string;
	    browserSourceThemeHearts: string;
	    browserSourceThemeBanner: string;
	    browserSourceThemeMostWantedStrip: string;
	    browserSourceThemeFlashPortal: string;
	    browserSourceThemeLink: string;
	    browserSourceThemeBlackout: string;
	    browserSourceThemeSilver: string;
	    browserSourceThemeTop8: string;
	    browserSourceThemePivot: string;
	    browserSourceThemeVelvet: string;
	    browserSourceThemeVapor: string;
	    browserSourceThemeZebra: string;
	    browserSourceThemeCrimson: string;
	    browserSourceThemeWebcore: string;
	    browserSourceThemeUltra: string;
	    browserSourceThemeProwess: string;
	    createYourOwnTheme: string;
	    userId: string;
	    userName: string;
	    victory: string;
	    lp: string;
	    mr: string;
	    opponentCharacter: string;
	    opponentLp: string;
	    opponentMr: string;
	    opponentLeague: string;
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
	    themeLayoutRows: string;
	    themeLayoutSideBySide: string;
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
	    T8_GOD_OF_DESTRUCTION_1: string;
	    T8_GOD_OF_DESTRUCTION_2: string;
	    T8_GOD_OF_DESTRUCTION_3: string;
	    T8_GOD_OF_DESTRUCTION_4: string;
	    T8_GOD_OF_DESTRUCTION_5: string;
	    T8_GOD_OF_DESTRUCTION_6: string;
	    T8_GOD_OF_DESTRUCTION_7: string;
	    T8_GOD_OF_DESTRUCTION_INFINITY: string;
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
	    errSaveLogFile: string;
	    errGetSessions: string;
	    errGetTranslations: string;
	    errGetSessionStatistics: string;
	    errCreateSession: string;
	    errOpenResultsDirectory: string;
	    errReadThemeCSS: string;
	    errInvalidBrowserSourcePort: string;
	    errSaveBrowserSourcePort: string;
	    errReportIssue: string;
	    errExport: string;
	    errTrackingInProgress: string;
	    errGetDataLocation: string;
	    errSaveDataDirectory: string;
	    errOpenDataDirectory: string;
	    errBackupData: string;
	    errImportBackup: string;
	    errInvalidBackup: string;
	    errSaveStreamOverlay: string;
	    retrying: string;
	    retryingDetail: string;
	    reconnected: string;
	    errNetwork: string;
	    errThrottled: string;
	    errStructureChanged: string;
	    authNeedRelogin: string;
	    authNeedReloginHint: string;
	    authCheckingSession: string;
	    authRestartingBrowser: string;
	    authVerifyingLogin: string;
	    authInterrupted: string;
	    errAuthManualLoginFailed: string;
	    saveData: string;
	    changeFolder: string;
	    openFolder: string;
	    resetToDefault: string;
	    backUp: string;
	    importBackup: string;
	    backupSaved: string;
	    restartToApply: string;
	    folderHasSaveData: string;
	    folderHasSaveDataDescription: string;
	    useExisting: string;
	    replace: string;
	    importBackupDescription: string;
	    import: string;
	    generalSettings: string;
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
	export interface SF6BattleStatsMetric {
	    key: string;
	    name: string;
	    unit: string;
	    kind: string;
	    polarity: number;
	    current: number;
	    topAvg: number;
	    topMin: number;
	    topMax: number;
	}
	export interface SF6BattleStatsComparisonReport {
	    userCode: string;
	    userName: string;
	    characterName: string;
	    characterToolName: string;
	    topN: number;
	    metrics: SF6BattleStatsMetric[];
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
	export interface StreamOverlayConfig {
	    theme: string;
	    lang: string;
	    stats: string[];
	}
	export interface Theme {
	    name: string;
	    css: string;
	}
	export interface User {
	    id: number;
	    displayName: string;
	    code: string;
	    game: GameType;
	    LP: number;
	    MR: number;
	}
	export interface UserStatistics {
	    mainCharacter: string;
	    sessions: number;
	    wins: number;
	    losses: number;
	    characters: CharacterStatistics[];
	    history: HistoryMatch[];
	}

}

