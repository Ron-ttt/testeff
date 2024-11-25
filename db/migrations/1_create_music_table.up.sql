BEGIN;
CREATE TABLE IF NOT EXISTS music(
    id serial primary key, 
    author text, 
    song text,
    releaseData TIMESTAMP,
    textSong text
    songLink text);

COMMIT;