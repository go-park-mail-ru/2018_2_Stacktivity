-- auto-generated definition
CREATE TABLE IF NOT EXISTS "user"
(
  uid      SERIAL       NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  username VARCHAR(30)  NOT NULL,
  email    VARCHAR(30)  NOT NULL,
  pass     VARCHAR(120) NOT NULL,
  avatar   VARCHAR(120) DEFAULT NULL,
  score    INTEGER      DEFAULT 0
);

CREATE UNIQUE INDEX users_uid_uindex
  ON "user" (uid);

CREATE UNIQUE INDEX users_username_uindex
  ON "user" (username);

CREATE UNIQUE INDEX users_email_uindex
  ON "user" (email);

