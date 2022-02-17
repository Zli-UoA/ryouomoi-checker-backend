package controller

type TwitterLoginUrlJson struct {
	LoginUrl string `json:"loginUrl"`
}

type TwitterUser struct {
	ID          int64  `db:"id"`
	ScreenNme   string `db:"screenNme"`
	DisplayName string `db:"displayName"`
	ImageUrl    string `db:"imageUrl"`
	Biography   string `db:"biography"`
}
