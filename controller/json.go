package controller

type TwitterLoginUrlJson struct {
	LoginUrl string `json:"loginUrl"`
}

type TwitterUser struct {
	ID          int64  `json:"id"`
	ScreenName  string `json:"screenName"`
	DisplayName string `json:"displayName"`
	ImageUrl    string `json:"imageUrl"`
	Biography   string `json:"biography"`
}

type LovePoint struct {
	LovePoint int `json:"lovePoint"`
}

type MatchResult struct {
	MatchSuccess bool `json:"matchSuccess"`
}
