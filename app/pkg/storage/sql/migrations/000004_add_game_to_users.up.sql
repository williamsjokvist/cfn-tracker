ALTER TABLE users ADD COLUMN "game" TEXT NOT NULL DEFAULT '';
-- Users saved before this have no game, but SF6 user codes are all digits and Tekken IDs aren't.
UPDATE users SET game = CASE WHEN code GLOB '*[^0-9]*' THEN 't8' ELSE 'sf6' END;
