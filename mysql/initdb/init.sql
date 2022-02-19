CREATE TABLE IF NOT EXISTS twitter_users (
    twitter_id BIGINT PRIMARY KEY,
    screen_name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    profile_image_url VARCHAR(255) NOT NULL,
    biography VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    access_token_secret VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_love_points (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL REFERENCES twitter_users(twitter_id),
    lover_user_id BIGINT NOT NULL,
    love_point INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS couples (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id_1 BIGINT NOT NULL REFERENCES twitter_users(twitter_id),
    user_id_2 BIGINT NOT NULL REFERENCES twitter_users(twitter_id),
    created_at TIMESTAMP NOT NULL,
    broken_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS couple_broke_reasons (
    id BIGINT PRIMARY KEY,
    text VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO couple_broke_reasons (id, text) VALUES (1, '価値観の違い');
INSERT INTO couple_broke_reasons (id, text) VALUES (2, '趣味・趣向の違い');
INSERT INTO couple_broke_reasons (id, text) VALUES (3, '冷めた・嫌いになった');
INSERT INTO couple_broke_reasons (id, text) VALUES (4, '他に好きな人ができた');
INSERT INTO couple_broke_reasons (id, text) VALUES (5, '自然消滅');
INSERT INTO couple_broke_reasons (id, text) VALUES (6, 'その他');

CREATE TABLE IF NOT EXISTS couple_broke_reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    couple_id BIGINT NOT NULL REFERENCES couples(id),
    user_id BIGINT NOT NULL REFERENCES twitter_users(twitter_id),
    broke_reason_id BIGINT NOT NULL REFERENCES couple_broke_reasons(id)
);

CREATE TABLE IF NOT EXISTS chat_rooms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    couple_id BIGINT NOT NULL REFERENCES couples(id),
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    chat_room_id BIGINT NOT NULL REFERENCES chat_rooms(id),
    user_id BIGINT NOT NULL REFERENCES twitter_users(twitter_id),
    message VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
