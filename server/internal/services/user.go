package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"yt-solutions-server/internal/core"
)

type UserRepo interface {
	CreateUser(userId int, username string) error
	GetUserById(userId int) (*core.User, error)
	ChargeBalance(userId, amount int) error
}

type UserService struct {
	UserRepo
	DashboardAuth string
}

func NewUserService(r UserRepo, token string) *UserService {
	return &UserService{
		UserRepo:      r,
		DashboardAuth: newDashboardAuth(token),
	}
}

func (u UserService) AuthUser(userId int, username string) (*core.User, error) {
	if err := u.UserRepo.CreateUser(userId, username); err != nil {
		return u.UserRepo.GetUserById(userId)
	}

	user := &core.User{
		ID:   userId,
		Name: username,
	}

	return user, nil
}

func (u UserService) GetUserInfo(userId int) (*core.User, error) {
	user, err := u.UserRepo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, err
	}

	return user, nil
}

func (u UserService) GetUserPhoto(userId int) ([]byte, error) {
	reqUrl := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   path.Join(u.DashboardAuth, "getUserProfilePhotos"),
	}

	q := url.Values{}
	q.Add("user_id", strconv.Itoa(userId))

	reqUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result core.UserPhotoResult

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Result.Total == 0 {
		return nil, nil
	}

	reqUrl.Path = path.Join(u.DashboardAuth, "getFile")

	q = url.Values{}
	q.Add("file_id", result.Result.Photos[0][1].FileId)

	reqUrl.RawQuery = q.Encode()

	req, err = http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var info core.PhotoInfoResult

	if err = json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	reqUrl.Path = path.Join("file", u.DashboardAuth, info.Result.Path)

	req, err = http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (u UserService) ChargeBalance(userId, amount int) error {
	return u.UserRepo.ChargeBalance(userId, amount)
}

func newDashboardAuth(token string) string {
	return "bot" + token
}
