ALTER TABLE users ADD COLUMN "game" TEXT NOT NULL DEFAULT '';
-- Users saved before this have no game. Their matches show it, since only Tekken leagues
-- start with T8_. Without matches, fall back to the code: SF6 codes are all digits.
UPDATE users SET game = CASE
  WHEN EXISTS (SELECT 1 FROM matches WHERE user_id = users.code AND opponent_league GLOB 'T8_*') THEN 't8'
  WHEN EXISTS (SELECT 1 FROM matches WHERE user_id = users.code) THEN 'sf6'
  WHEN code GLOB '*[^0-9]*' THEN 't8'
  ELSE 'sf6'
END;
