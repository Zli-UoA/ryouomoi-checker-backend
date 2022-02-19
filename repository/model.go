package repository

import "time"

type TwitterUser struct {
	TwitterID         int64  `db:"twitter_id"`
	ScreenName        string `db:"screen_name"`
	DisplayName       string `db:"display_name"`
	ProfileImageUrl   string `db:"profile_image_url"`
	Biography         string `db:"biography"`
	AccessToken       string `db:"access_token"`
	AccessTokenSecret string `db:"access_token_secret"`
}

type UserLovePoint struct {
	ID          int64        `db:"id"`
	User        *TwitterUser `db:"user"`
	LoverUserID int64        `db:"lover_user_id"`
	LovePoint   int          `db:"love_point"`
}

type Couple struct {
	ID        int64        `db:"id"`
	User1     *TwitterUser `db:"user1"`
	User2     *TwitterUser `db:"user2"`
	CreatedAt time.Time    `db:"created_at"`
	BrokenAt  *time.Time   `db:"broken_at"`
}

type ChatRoom struct {
	ID        int64     `db:"id"`
	Couple    *Couple   `db:"couple"`
	CreatedAt time.Time `db:"created_at"`
}

type Chats struct {
	ID        int64        `db:"id"`
	ChatRoom  *ChatRoom    `db:"chat_room"`
	User      *TwitterUser `db:"user"`
	Message   string       `db:"message"`
	CreatedAt time.Time    `db:"created_at"`
}

type CoupleBrokeReport struct {
	ID            int64        `db:"id"`
	Couple        *Couple      `db:"couple"`
	User          *TwitterUser `db:"user"`
	BrokeReasonID int64        `db:"broke_reason_id"`
	AllowShare    bool         `db:"allow_share"`
}
