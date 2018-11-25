CREATE TABLE IF NOT EXISTS chat_message
(
  id         SERIAL                                 NOT NULL
    CONSTRAINT chat_message_pkey
    PRIMARY KEY,
  user_id    INTEGER
    CONSTRAINT chat_message_user_uid_fk
    REFERENCES "user",
  chat_id    INTEGER
    CONSTRAINT chat_message_chat_id_fk
    REFERENCES chat,
  chat_text  VARCHAR(512),
  created    TIME DEFAULT now(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS chat_message_id_uindex
  ON chat_message (id);