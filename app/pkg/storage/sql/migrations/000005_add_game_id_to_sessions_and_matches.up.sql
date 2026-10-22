ALTER TABLE sessions ADD COLUMN "game_id" TEXT NOT NULL DEFAULT '';
ALTER TABLE matches ADD COLUMN "game_id" TEXT NOT NULL DEFAULT '';
-- Every game tracked before this has a single game ID, the same as its user's game.
-- Fightcade sessions instead store the ID of the emulated game they track.
UPDATE sessions SET game_id = COALESCE((SELECT u.game FROM users u WHERE u.code = sessions.user_id), '');
UPDATE matches SET game_id = COALESCE((SELECT s.game_id FROM sessions s WHERE s.id = matches.session_id), '');
