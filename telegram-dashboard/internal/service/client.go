package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"yt-solutions-telegram-dashboard/internal/core"
	"yt-solutions-telegram-dashboard/internal/resources"
)

const (
	authUserMethod       = "authUser"
	getAvalabilityMethod = "getAvailability"
	processOrderMethod   = "processOrder"
	getKeyInfoMethod     = "getKeyInfo"
	unbindAddressMethod  = "unbindAddress"

	getRenewalInfoMethod = "getRenewalInfo"
	processRenewalMethod = "processRenewal"
	chargeBalanceMethod  = "chargeBalance"

	sendMessageMethod    = "sendMessage"
	editMessageMethod    = "editMessageText"
	deleteMessageMethod  = "deleteMessage"
	sendDocumentMethod   = "sendDocument"
	answerCallbackMethod = "answerCallbackQuery"
)

type Client interface {
	AuthUser(userId int, username string) (*core.User, error)
	GetAvailability() (int, int, int, error)
	ProcessOrder(userId int) (bool, *core.Key, error)
	GetKeyInfo(userId int) (*core.Key, error)
	UnbindAddress(userId int) error

	GetRenewalInfo() (int, error)
	ProcessRenewal(userId int) (bool, error)
	ChargeBalance(userId, amount int) error

	SendMessage(text string, userId int, keyboard interface{}) (int, error)
	EditMessage(text string, userId, messageId int, keyboard interface{}) error
	DeleteMessage(usedId, messageId int) error
	SendFile(file *os.File, chatId int) error
	AnswerCallback(text, callbackId string) error

	CreateInvoice(userUrl string) (string, string, error)
	CheckFiatPayment(invoiceId string) (int, string, error)
	CheckCryptoPayment(txId string) (*core.TransactionInfo, error)
}

type client struct {
	TelegramHost string
	TelegramPath string
	BackendHost  string
	BackendPath  string
}

func NewClient(tgHost, token, backendHost string) Client {
	return &client{
		TelegramHost: tgHost,
		TelegramPath: "bot" + token,
		BackendHost:  backendHost,
		BackendPath:  "api",
	}
}

func (c client) AuthUser(userId int, username string) (*core.User, error) {
	var buf bytes.Buffer

	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, authUserMethod),
	}

	body := core.User{
		ID:   userId,
		Name: username,
	}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	data, err := c.doRequest(http.MethodPost, u, &buf)
	if err != nil {
		return nil, err
	}

	var user core.User

	if err = json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c client) GetAvailability() (int, int, int, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, getAvalabilityMethod),
	}

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	var available resources.AvailabilityResponse

	if err = json.Unmarshal(data, &available); err != nil {
		return 0, 0, 0, err
	}

	return available.Free, available.Lifetime, available.Price, nil
}

func (c client) ProcessOrder(userId int) (bool, *core.Key, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, processOrderMethod),
	}

	q := url.Values{}
	q.Add("id", strconv.Itoa(userId))

	u.RawQuery = q.Encode()

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, nil, err
	}

	var resp resources.ProcessOrderResponse

	if err = json.Unmarshal(data, &resp); err != nil {
		return false, nil, err
	}

	return resp.Success, &resp.Key, nil
}

func (c client) GetKeyInfo(userId int) (*core.Key, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, getKeyInfoMethod),
	}

	q := url.Values{}
	q.Add("by", "ownerId")
	q.Add("owner", strconv.Itoa(userId))

	u.RawQuery = q.Encode()

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	var key core.Key

	if err = json.Unmarshal(data, &key); err != nil {
		return nil, err
	}

	return &key, nil
}

func (c client) UnbindAddress(userId int) error {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, unbindAddressMethod),
	}

	q := url.Values{}
	q.Add("id", strconv.Itoa(userId))

	u.RawQuery = q.Encode()

	if _, err := c.doRequest(http.MethodGet, u, nil); err != nil {
		return err
	}

	return nil
}

func (c client) GetRenewalInfo() (int, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, getRenewalInfoMethod),
	}

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return 0, err
	}

	var renewInfo resources.RenewalInfoResponse

	if err = json.Unmarshal(data, &renewInfo); err != nil {
		return 0, err
	}

	return renewInfo.Renewal, nil
}

func (c client) ProcessRenewal(userId int) (bool, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, processRenewalMethod),
	}

	q := url.Values{}
	q.Add("owner", strconv.Itoa(userId))

	u.RawQuery = q.Encode()

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, err
	}

	var resp resources.ProcessRenewalResponse

	if err = json.Unmarshal(data, &resp); err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (c client) ChargeBalance(userId, amount int) error {
	u := &url.URL{
		Scheme: "http",
		Host:   c.BackendHost,
		Path:   path.Join(c.BackendPath, chargeBalanceMethod),
	}

	q := url.Values{}
	q.Add("EXTRA", strconv.Itoa(userId))
	q.Add("AMOUNT", strconv.Itoa(amount))

	u.RawQuery = q.Encode()

	if _, err := c.doRequest(http.MethodGet, u, nil); err != nil {
		return err
	}

	return nil
}

func (c client) SendMessage(text string, userId int, keyboard interface{}) (int, error) {
	var buf bytes.Buffer

	u := &url.URL{
		Scheme: "https",
		Host:   c.TelegramHost,
		Path:   path.Join(c.TelegramPath, sendMessageMethod),
	}

	msg := resources.SendMessageRequest{
		ChatId:   userId,
		Text:     text,
		Keyboard: keyboard,
	}

	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return 0, err
	}

	data, err := c.doRequest(http.MethodPost, u, &buf)
	if err != nil {
		return 0, err
	}

	var res resources.SendMessageResult

	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res.Result.ID, nil
}

func (c client) EditMessage(text string, userId, messageId int, keyboard interface{}) error {
	var buf bytes.Buffer

	u := &url.URL{
		Scheme: "https",
		Host:   c.TelegramHost,
		Path:   path.Join(c.TelegramPath, editMessageMethod),
	}

	body := resources.EditMessageRequest{
		ChatId:    userId,
		MessageId: messageId,
		Text:      text,
		Keyboard:  keyboard,
	}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	if _, err := c.doRequest(http.MethodPost, u, &buf); err != nil {
		return err
	}

	return nil
}

func (c client) DeleteMessage(userId, messageId int) error {
	var buf bytes.Buffer

	u := &url.URL{
		Scheme: "https",
		Host:   c.TelegramHost,
		Path:   path.Join(c.TelegramPath, deleteMessageMethod),
	}

	body := resources.DeleteMessageRequest{
		ChatId:    userId,
		MessageId: messageId,
	}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	if _, err := c.doRequest(http.MethodPost, u, &buf); err != nil {
		return err
	}

	return nil
}

func (c client) SendFile(file *os.File, chatId int) error {
	var buf bytes.Buffer

	strChatId := strconv.Itoa(chatId)
	strReader := strings.NewReader(strChatId)

	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile("document", "yt_solutions.exe")
	if err != nil {
		return err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
	}

	strWriter, err := writer.CreateFormField("chat_id")
	if err != nil {
		return err
	}

	if _, err = io.Copy(strWriter, strReader); err != nil {
		return err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   c.TelegramHost,
		Path:   path.Join(c.TelegramPath, sendDocumentMethod),
	}

	fmt.Println("DO REQ")

	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	fmt.Println(string(b))

	return nil
}

func (c client) AnswerCallback(text, callbackId string) error {
	var buf bytes.Buffer

	u := &url.URL{
		Scheme: "https",
		Host:   c.TelegramHost,
		Path:   path.Join(c.TelegramPath, answerCallbackMethod),
	}

	body := resources.AnswerCallbackRequest{
		ID:        callbackId,
		Text:      text,
		ShowAlert: false,
		CacheTime: 1,
	}

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}

	if _, err := c.doRequest(http.MethodPost, u, &buf); err != nil {
		return err
	}

	return nil
}

func (c client) CreateInvoice(userUrl string) (string, string, error) {
	u, err := url.Parse(userUrl)
	if err != nil {
		return "", "", err
	}

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", "", err
	}

	var resp resources.CreateInvoiceResponse

	if err = json.Unmarshal(data, &resp); err != nil {
		return "", "", err
	}

	return resp.ID, resp.Url, nil
}

func (c client) CheckFiatPayment(invoiceId string) (int, string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.crystalpay.ru",
		Path:   "v1",
	}

	q := url.Values{}
	q.Add("s", "17592c19723a2f86dc96e22b6d28f18897c0ea35")
	q.Add("n", "youtubemarket")
	q.Add("o", "invoice-check")
	q.Add("i", invoiceId)

	u.RawQuery = q.Encode()

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return 0, "", nil
	}

	var resp resources.InvoiceStatusResponse

	if err = json.Unmarshal(data, &resp); err != nil {
		return 0, "", nil
	}

	return resp.Amount, resp.State, nil
}

func (c client) CheckCryptoPayment(txId string) (*core.TransactionInfo, error) {

	u := &url.URL{
		Scheme: "https",
		Host:   "apilist.tronscanapi.com",
		Path:   "api/transaction-info",
	}

	q := url.Values{}
	q.Add("hash", txId)

	u.RawQuery = q.Encode()

	data, err := c.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	var resp core.TransactionInfo

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c client) doRequest(httpMethod string, u *url.URL, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	//TODO:
	if u.Host == c.BackendHost {
		req.Header.Add("Authorization", "")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
