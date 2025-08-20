-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE habits
(
    name TEXT NOT NULL,
    chatId int NOT NULL ,
    mark bool NOT NULL default false,
    PRIMARY KEY (name, chatId)
);


-- +goose Down
DROP TABLE page;