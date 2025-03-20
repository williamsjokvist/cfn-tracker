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
	export class FGCTrackerError {
	    localizationKey: ErrorLocalizationKey;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new FGCTrackerError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.localizationKey = source["localizationKey"];
	        this.message = source["message"];
	    }
	}
	export class GUIConfig {
	    locale: string;
	    theme: ThemeName;
	    sidebar: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GUIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.locale = source["locale"];
	        this.theme = source["theme"];
	        this.sidebar = source["sidebar"];
	    }
	}
	export class Localization {
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
	
	    static createFrom(source: any = {}) {
	        return new Localization(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appVersion = source["appVersion"];
	        this.source = source["source"];
	        this.about = source["about"];
	        this.settings = source["settings"];
	        this.tracking = source["tracking"];
	        this.history = source["history"];
	        this.language = source["language"];
	        this.changelog = source["changelog"];
	        this.startTracking = source["startTracking"];
	        this.cfnName = source["cfnName"];
	        this.start = source["start"];
	        this.opponent = source["opponent"];
	        this.character = source["character"];
	        this.lpGain = source["lpGain"];
	        this.mrGain = source["mrGain"];
	        this.deleteLog = source["deleteLog"];
	        this.goBack = source["goBack"];
	        this.loading = source["loading"];
	        this.wins = source["wins"];
	        this.losses = source["losses"];
	        this.winRate = source["winRate"];
	        this.stop = source["stop"];
	        this.files = source["files"];
	        this.enterCfnName = source["enterCfnName"];
	        this.result = source["result"];
	        this.replayId = source["replayId"];
	        this.copy = source["copy"];
	        this.time = source["time"];
	        this.winStreak = source["winStreak"];
	        this.newVersionAvailable = source["newVersionAvailable"];
	        this.pause = source["pause"];
	        this.unpause = source["unpause"];
	        this.statistics = source["statistics"];
	        this.date = source["date"];
	        this.minimize = source["minimize"];
	        this.restoreSession = source["restoreSession"];
	        this.exportLog = source["exportLog"];
	        this.league = source["league"];
	        this.pickGame = source["pickGame"];
	        this.follow = source["follow"];
	        this.continueStep = source["continueStep"];
	        this.placement = source["placement"];
	        this.cfnError = source["cfnError"];
	        this.output = source["output"];
	        this.theme = source["theme"];
	        this.defaultTheme = source["defaultTheme"];
	        this.pickTheme = source["pickTheme"];
	        this.usingBrowserSource = source["usingBrowserSource"];
	        this.copyBrowserSourceLink = source["copyBrowserSourceLink"];
	        this.browserSourceDescription = source["browserSourceDescription"];
	        this.displayStats = source["displayStats"];
	        this.importFiles = source["importFiles"];
	        this.obsCustomize = source["obsCustomize"];
	        this.lastMatch = source["lastMatch"];
	        this.statsWillBeDisplayed = source["statsWillBeDisplayed"];
	        this.selectTheme = source["selectTheme"];
	        this.started = source["started"];
	        this.user = source["user"];
	        this.matchesWon = source["matchesWon"];
	        this.matchesLost = source["matchesLost"];
	        this.sessions = source["sessions"];
	        this.refresh = source["refresh"];
	        this.cooldown = source["cooldown"];
	        this.T8_BEGINNER = source["T8_BEGINNER"];
	        this.T8_1ST_DAN = source["T8_1ST_DAN"];
	        this.T8_2ND_DAN = source["T8_2ND_DAN"];
	        this.T8_FIGHTER = source["T8_FIGHTER"];
	        this.T8_STRATEGIST = source["T8_STRATEGIST"];
	        this.T8_COMBATANT = source["T8_COMBATANT"];
	        this.T8_BRAWLER = source["T8_BRAWLER"];
	        this.T8_RANGER = source["T8_RANGER"];
	        this.T8_CAVALRY = source["T8_CAVALRY"];
	        this.T8_WARRIOR = source["T8_WARRIOR"];
	        this.T8_ASSAILANT = source["T8_ASSAILANT"];
	        this.T8_DOMINATOR = source["T8_DOMINATOR"];
	        this.T8_VANQUISHER = source["T8_VANQUISHER"];
	        this.T8_DESTROYER = source["T8_DESTROYER"];
	        this.T8_ELIMINATOR = source["T8_ELIMINATOR"];
	        this.T8_GARYU = source["T8_GARYU"];
	        this.T8_SHINRYU = source["T8_SHINRYU"];
	        this.T8_TENRYU = source["T8_TENRYU"];
	        this.T8_MIGHTY_RULER = source["T8_MIGHTY_RULER"];
	        this.T8_FLAME_RULER = source["T8_FLAME_RULER"];
	        this.T8_BATTLE_RULER = source["T8_BATTLE_RULER"];
	        this.T8_FUJIN = source["T8_FUJIN"];
	        this.T8_RAIJIN = source["T8_RAIJIN"];
	        this.T8_KISHIN = source["T8_KISHIN"];
	        this.T8_BUSHIN = source["T8_BUSHIN"];
	        this.T8_TEKKEN_KING = source["T8_TEKKEN_KING"];
	        this.T8_TEKKEN_EMPEROR = source["T8_TEKKEN_EMPEROR"];
	        this.T8_TEKKEN_GOD = source["T8_TEKKEN_GOD"];
	        this.T8_TEKKEN_GOD_SUPREME = source["T8_TEKKEN_GOD_SUPREME"];
	        this.T8_GOD_OF_DESTRUCTION = source["T8_GOD_OF_DESTRUCTION"];
	        this.errUnknown = source["errUnknown"];
	        this.errSelectGame = source["errSelectGame"];
	        this.errAuth = source["errAuth"];
	        this.errGetLatestSession = source["errGetLatestSession"];
	        this.errGetUser = source["errGetUser"];
	        this.errGetMatches = source["errGetMatches"];
	        this.errSaveLocale = source["errSaveLocale"];
	        this.errCheckForUpdate = source["errCheckForUpdate"];
	        this.errGetGuiConfig = source["errGetGuiConfig"];
	        this.errSaveTheme = source["errSaveTheme"];
	        this.errSaveUser = source["errSaveUser"];
	        this.errSaveSidebar = source["errSaveSidebar"];
	        this.errGetSessions = source["errGetSessions"];
	        this.errGetTranslations = source["errGetTranslations"];
	        this.errGetSessionStatistics = source["errGetSessionStatistics"];
	        this.errCreateSession = source["errCreateSession"];
	        this.errOpenResultsDirectory = source["errOpenResultsDirectory"];
	        this.errReadThemeCSS = source["errReadThemeCSS"];
	    }
	}
	export class Match {
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
	
	    static createFrom(source: any = {}) {
	        return new Match(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.userName = source["userName"];
	        this.sessionId = source["sessionId"];
	        this.replayId = source["replayId"];
	        this.character = source["character"];
	        this.lp = source["lp"];
	        this.lpGain = source["lpGain"];
	        this.mr = source["mr"];
	        this.mrGain = source["mrGain"];
	        this.opponent = source["opponent"];
	        this.opponentCharacter = source["opponentCharacter"];
	        this.opponentLp = source["opponentLp"];
	        this.opponentMr = source["opponentMr"];
	        this.opponentLeague = source["opponentLeague"];
	        this.victory = source["victory"];
	        this.date = source["date"];
	        this.time = source["time"];
	        this.winStreak = source["winStreak"];
	        this.wins = source["wins"];
	        this.losses = source["losses"];
	        this.winRate = source["winRate"];
	    }
	}
	export class Session {
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
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.userName = source["userName"];
	        this.createdAt = source["createdAt"];
	        this.lp = source["lp"];
	        this.mr = source["mr"];
	        this.matches = this.convertValues(source["matches"], Match);
	        this.matchesWon = source["matchesWon"];
	        this.matchesLost = source["matchesLost"];
	        this.endingLp = source["endingLp"];
	        this.endingMr = source["endingMr"];
	        this.startingLp = source["startingLp"];
	        this.startingMr = source["startingMr"];
	        this.lpGain = source["lpGain"];
	        this.mrGain = source["mrGain"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SessionMonth {
	    Date: string;
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionMonth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Date = source["Date"];
	        this.Count = source["Count"];
	    }
	}
	export class SessionsStatistics {
	    Months: SessionMonth[];
	
	    static createFrom(source: any = {}) {
	        return new SessionsStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Months = this.convertValues(source["Months"], SessionMonth);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Theme {
	    name: string;
	    css: string;
	
	    static createFrom(source: any = {}) {
	        return new Theme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.css = source["css"];
	    }
	}
	export class User {
	    id: number;
	    displayName: string;
	    code: string;
	    LP: number;
	    MR: number;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.displayName = source["displayName"];
	        this.code = source["code"];
	        this.LP = source["LP"];
	        this.MR = source["MR"];
	    }
	}

}

