package models

type CookieInfo struct {
	ID          string `json:"id"`
	ViewsCount  int    `json:"views_count"`
	Subscribes  int    `json:"subs_count"`
	Channels    int    `json:"channels_count"`
	VideosCount int    `json:"videos_count"`
	Geo         string `json:"geo"`
	Monetize    bool   `json:"monetize"`
	Brand       bool   `json:"brand"`
	RegDate     string `json:"reg_date"`
	Verified    bool   `json:"verified"`
	Received    string
	Path        string
	TxtPath     string
	SavePath    string
}
