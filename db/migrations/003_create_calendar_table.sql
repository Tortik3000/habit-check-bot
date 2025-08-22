-- +goose Up
CREATE TABLE calendar
(
    date DATE DEFAULT CURRENT_DATE NOT NULL,
    chatId int NOT NULL,
    mark bool NOT NULL default false,
    PRIMARY KEY (date, chatId)
);

-- +goose Down
DROP TABLE page;