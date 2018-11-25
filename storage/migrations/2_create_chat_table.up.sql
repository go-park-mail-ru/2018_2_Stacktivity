CREATE TABLE IF NOT EXISTS chat
(
  id        SERIAL NOT NULL
    CONSTRAINT chat_pkey
    PRIMARY KEY,
  chat_name VARCHAR(30),
  chat_type VARCHAR(30)
);

CREATE UNIQUE INDEX IF NOT EXISTS chat_id_uindex
  ON chat (id);