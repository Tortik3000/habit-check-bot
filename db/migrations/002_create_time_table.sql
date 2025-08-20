-- +goose Up
CREATE TABLE time_done
(
    name TEXT NOT NULL,
    chatId int NOT NULL ,
    FOREIGN KEY (name, chatId) REFERENCES habits(name, chatId) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now() NOT NULL 
);

-- +goose Down
DROP TABLE time_done;