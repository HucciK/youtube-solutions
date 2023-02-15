package client

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	"yt-solutions-soft/internal/models"

	"golang.org/x/net/publicsuffix"
)

const (
	googleLoginServive    = "https://accounts.google.com/ServiceLogin?service=youtube&uilel=3&passive=true&continue=https%3A%2F%2Fwww.youtube.com%2Fsignin%3Faction_handle_signin%3Dtrue%26app%3Ddesktop%26hl%3Den%26next%3Dhttps%253A%252F%252Fwww.youtube.com%252F&hl=en%22"
	youtubeDatasyncUrl    = "https://www.youtube.com/getDatasyncIdsEndpoint"
	youtubeAccountUrl     = "https://www.youtube.com/getAccountSwitcherEndpoint"
	youtubeChannelInfoUrl = "https://www.youtube.com/youtubei/v1/browse?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false"
)

type Client struct {
	Client   *http.Client
	GoogAuth string
}

func NewClient(proxy *url.URL) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: 4 * time.Second,

			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		},
		GoogAuth: "0",
	}
}

func (c *Client) CheckCookies(cookies []*http.Cookie, cookieGeo, cookiesPath, cookiesTxtPath string, checkControl *models.CheckedAmount, updatesChan chan models.Update, resultChan chan *models.CookieInfo) {

	if err := c.setCookieJar(); err != nil {
		return
	}

	data, err := c.startAuthProcess(cookies)
	if err != nil {
		return
	}

	if len(data) < 4 {
		return
	}
	data = data[4:]

	var account YoutubeAccountData

	if err := json.Unmarshal(data, &account); err != nil {
		return
	}

	if len(account.Data.Actions) == 0 {
		return
	}

	for s := 0; s < len(account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections); s++ {

		if len(account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents) == 0 {
			continue
		}

		channels := 0
		for content := 0; content < len(account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents); content++ {
			var channelInfo YoutubeChannelData

			if account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents[content].AccountItem.AccountName.SimpleText == "" {
				continue
			}

			channelId := ""
			for token := 0; token < len(account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents[content].AccountItem.ServiceEndpoint.SelectActiveIdentityEndpoint.SupportedTokens); token++ {
				if account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents[content].AccountItem.ServiceEndpoint.SelectActiveIdentityEndpoint.SupportedTokens[token].OfflineCacheKeyToken.ClientCacheKey == "" {
					continue
				}

				channelId = "UC" + account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections[s].AccountSectionListRenderer.Contents[0].AccountItemSectionRenderer.Contents[content].AccountItem.ServiceEndpoint.SelectActiveIdentityEndpoint.SupportedTokens[token].OfflineCacheKeyToken.ClientCacheKey
			}
			channels++

			data, err := c.channelInformation(channelId)
			if err != nil {
				continue
			}

			if err := json.Unmarshal(data, &channelInfo); err != nil {
				continue
			}

			if len(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs) < 1 {
				continue
			}

			index := 2
			if len(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs[len(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs)-index].TabRenderer.RendererContent.SectionListRenderer.ListRendererContents) == 0 {
				index = 1
			}

			info := &models.CookieInfo{}
			info.ID = channelId
			info.Monetize = false
			info.ViewsCount, _ = c.parseViews(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs[len(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs)-index].TabRenderer.RendererContent.SectionListRenderer.ListRendererContents[0].ItemSectionRenderer.SectionContents[0].ChannelAboutFullMetadataRenderer.ViewCountText.SimpleText)
			info.Subscribes, _ = c.convertSubs(channelInfo.Header.C4TabbedHeaderRenderer.SubscriberCountText.SimpleText)
			//info.Channels = len(account.Data.Actions[0].GetMultiPageMenuAction.Menu.MultiPageMenuRenderer.Sections)
			info.VideosCount, _ = strconv.Atoi(channelInfo.Header.C4TabbedHeaderRenderer.VideosCountText.Runs[0].Text)
			info.RegDate = channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs[len(channelInfo.Contents.TwoColumnBrowseResultsRenderer.Tabs)-index].TabRenderer.RendererContent.SectionListRenderer.ListRendererContents[0].ItemSectionRenderer.SectionContents[0].ChannelAboutFullMetadataRenderer.JoinedDateText.Runs[1].Text
			info.Geo = cookieGeo
			info.Path = cookiesPath
			info.TxtPath = cookiesTxtPath

			if channels >= 2 {
				info.Brand = true
			}

			upd := models.Update{Type: models.FoundChannelType, Data: *info}
			updatesChan <- upd
			checkControl.Valid++

			upd = models.Update{Type: models.CheckStatusType, Data: *checkControl}
			updatesChan <- upd

			resultChan <- info
		}
	}
}

func (c *Client) startAuthProcess(cookies []*http.Cookie) ([]byte, error) {

	u, err := url.Parse(googleLoginServive)
	if err != nil {
		return nil, err
	}

	c.Client.Jar.SetCookies(u, cookies)
	res, err := c.Client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	u, err = url.Parse(youtubeDatasyncUrl)
	if err != nil {
		return nil, err
	}

	if len(c.Client.Jar.Cookies(u)) < 10 {
		return c.userCookiesAuthRequest(cookies)
	}

	return c.googleCookiesAuthRequest(c.Client.Jar.Cookies(u))
}

func (c *Client) googleCookiesAuthRequest(cookies []*http.Cookie) ([]byte, error) {

	sapisidHash, err := c.generateAuthToken(cookies)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", youtubeAccountUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", sapisidHash)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("x-goog-authuser", c.GoogAuth)
	req.Header.Add("origin", "https://www.youtube.com")
	req.Header.Add("referer", "https://www.youtube.com/")
	req.Header.Add("x-origin", "https://www.youtube.com")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (c *Client) userCookiesAuthRequest(cookies []*http.Cookie) ([]byte, error) {
	u, err := url.Parse(youtubeDatasyncUrl)
	if err != nil {
		return nil, err
	}

	c.Client.Jar.SetCookies(u, cookies)

	sapisidHash, err := c.generateAuthToken(c.Client.Jar.Cookies(u))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", youtubeAccountUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", sapisidHash)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("x-goog-authuser", c.GoogAuth)
	req.Header.Add("origin", "https://www.youtube.com")
	req.Header.Add("referer", "https://www.youtube.com/")
	req.Header.Add("x-origin", "https://www.youtube.com")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (c *Client) generateAuthToken(cookies []*http.Cookie) (string, error) {
	var sapisid string
	for _, cookie := range cookies {
		if cookie.Name == "SAPISID" || cookie.Name == "__Secure-1PAPISID" || cookie.Name == "__Secure-3PAPISID" {
			sapisid = cookie.Value
		}
	}

	if sapisid == "" {
		return "", errors.New("sapisid not found")
	}

	t := time.Now().Unix()
	auth_str := fmt.Sprintf("%v %s %s", t, sapisid, "https://www.youtube.com")

	hash := sha1.New()
	hash.Write([]byte(auth_str))

	return fmt.Sprintf("SAPISIDHASH %v_%x", t, hash.Sum(nil)), nil
}

func (c *Client) channelInformation(channel_id string) ([]byte, error) {
	u, err := url.Parse(youtubeChannelInfoUrl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	r.BrowseID = channel_id

	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		return nil, err
	}

	sapisidHash, err := c.generateAuthToken(c.Client.Jar.Cookies(u))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", youtubeChannelInfoUrl, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", sapisidHash)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("x-goog-authuser", c.GoogAuth)
	req.Header.Add("origin", "https://www.youtube.com")
	req.Header.Add("referer", "https://www.youtube.com/")
	req.Header.Add("x-origin", "https://www.youtube.com")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) googAuthToInt() int {
	gInt, _ := strconv.Atoi(c.GoogAuth)
	return gInt
}

func (c *Client) increaseGoogAuth() {
	intGoog, err := strconv.Atoi(c.GoogAuth)
	if err != nil {
		// handle later
	}
	intGoog++

	c.GoogAuth = strconv.Itoa(intGoog)
}

func (c *Client) convertSubs(input string) (int, error) {
	var result float64

	if input == "" {
		return 0, nil
	}

	subsStr := strings.Split(input, " ")[0]
	multiplier := subsStr[len(subsStr)-1]

	if string(multiplier) == "K" || string(multiplier) == "M" {
		subs, err := strconv.ParseFloat(strings.ReplaceAll(subsStr, string(multiplier), ""), 64)
		if err != nil {
			return 0, err
		}

		switch string(multiplier) {
		case "K":
			result = subs * 1000
		case "M":
			result = subs * 1000000
		default:
			result = subs
		}

		return int(result), nil
	}

	if subsStr == "No" {
		return 0, nil
	}

	return strconv.Atoi(subsStr)
}

func (c *Client) parseViews(input string) (int, error) {
	if input == "" {
		return 0, nil
	}

	viewsStr := strings.Split(input, " ")[0]

	return strconv.Atoi(strings.ReplaceAll(viewsStr, ",", ""))
}

func (c *Client) setCookieJar() error {

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}

	c.Client.Jar = jar

	return nil
}
