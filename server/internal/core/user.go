package core

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	HasKey  bool    `json:"hasKey"`
}

type UserPhotoResult struct {
	Success bool       `json:"ok"`
	Result  UserPhotos `json:"result"`
}

type UserPhotos struct {
	Total  int       `json:"total_count"`
	Photos [][]Photo `json:"photos"`
}

type Photo struct {
	FileId string `json:"file_id"`
}

type PhotoInfoResult struct {
	Success bool      `json:"ok"`
	Result  PhotoInfo `json:"result"`
}

type PhotoInfo struct {
	Path string `json:"file_path"`
}
