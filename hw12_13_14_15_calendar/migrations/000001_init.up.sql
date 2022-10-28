CREATE TABLE IF NOT EXISTS users (
    user_id     SERIAL NOT NULL PRIMARY KEY,
    user_name   VARCHAR(255) NOT NULL,
    user_desc   TEXT
);

CREATE TABLE IF NOT EXISTS events (
    event_id    SERIAL NOT NULL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    start_date  TIMESTAMP NOT NULL,
    end_date    TIMESTAMP NOT NULL,
    event_desc  TEXT,
    user_id     INTEGER,
    notify_date TIMESTAMP NOT NULL
);
