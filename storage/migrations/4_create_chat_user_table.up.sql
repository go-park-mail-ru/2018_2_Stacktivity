CREATE TABLE IF NOT EXISTS chat_user
(
  id      SERIAL  NOT NULL
    CONSTRAINT chat_user_pkey
    PRIMARY KEY,
  chat_id INTEGER NOT NULL
    CONSTRAINT chat_user_chat_id_fk
    REFERENCES chat,
  user_id INTEGER NOT NULL
    CONSTRAINT chat_user_user_uid_fk
    REFERENCES "user"
);

CREATE UNIQUE INDEX IF NOT EXISTS chat_user_id_uindex
  ON chat_user (id);