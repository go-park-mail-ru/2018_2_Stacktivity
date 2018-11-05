CREATE TABLE IF NOT EXISTS users
(
  uid      SERIAL       NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  username VARCHAR(30)  NOT NULL,
  email    VARCHAR(30)  NOT NULL,
  pass     VARCHAR(120) NOT NULL,
  avatar   VARCHAR(120) NOT NULL,
  score    INTEGER DEFAULT 0,
  level    INTEGER DEFAULT 0
);

CREATE UNIQUE INDEX users_uid_uindex
  ON users (uid);

CREATE UNIQUE INDEX users_username_uindex
  ON users (username);

CREATE UNIQUE INDEX users_email_uindex
  ON users (email);
