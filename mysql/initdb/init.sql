CREATE TABLE IF NOT EXISTS twitter_users (
    twitter_id BIGINT PRIMARY KEY,
    screen_name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    profile_image_url VARCHAR(255) NOT NULL,
    biography VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    access_token_secret VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_lover_points (
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
