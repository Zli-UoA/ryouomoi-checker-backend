package model

import "time"

type UserLovePoint struct {
	ID          int64
	UserID      int64
	LoverUserID int64
	LovePoint   int
}

type Couple struct {
	ID        int64
	User1     *User
	User2     *User
	CreatedAt *time.Time
	BrokenAt  *time.Time
}

type BrokeReport struct {
	ID            int64
	Couple        Couple
	User          TwitterUser
	BrokeReasonID int64 // init.sql参照
}
