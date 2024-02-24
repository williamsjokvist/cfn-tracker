CREATE TABLE sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id TEXT NOT NULL,
	lp INTEGER DEFAULT 0 NOT NULL,
	mr INTEGER DEFAULT 0 NOT NULL,
	created_at TEXT DEFAULT (DATETIME('NOW')) NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(code)
)

CREATE TABLE matches (
	user_id TEXT,
	session_id INTEGER,
	character TEXT NOT NULL,
	lp INTEGER,
	lp_gain INTEGER,
	mr INTEGER,
	mr_gain INTEGER,
	opponent TEXT,
	opponent_character TEXT,
	opponent_lp TEXT,
	opponent_mr INTEGER,
	opponent_league TEXT,
	victory BOOLEAN,
	wins INTEGER,
	losses INTEGER,
	win_streak INTEGER,
	win_rate INTEGER,
	date TEXT,
	time TEXT,
	PRIMARY KEY (session_id, date, time),
	FOREIGN KEY (session_id) REFERENCES sessions(id)
	FOREIGN KEY (user_id) REFERENCES users(code)
)

CREATE TABLE users (
	code TEXT NOT NULL,
	display_name TEXT NOT NULL,
	PRIMARY KEY (code)
)
