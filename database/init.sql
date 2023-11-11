DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id bigint PRIMARY	KEY,
  name varchar(255) NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);