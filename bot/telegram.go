package bot

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegramBot/bot/errors"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Bot struct {
	token string

	client   http.Client
	host     string
	basePath string
}

func NewBot(token string, host string) *Bot {
	return &Bot{
		token:    token,
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{}}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (bot *Bot) Updates(offset, limit int) ([]Update, error) {
	values := url.Values{}
	values.Add("offset", strconv.Itoa(offset))
	values.Add("limit", strconv.Itoa(limit))

	content, err := bot.doRequest(values, getUpdatesMethod)
	if err != nil {
		return nil, err
	}

	var response UpdateResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (bot *Bot) SendMessage(chatID int, message string) error {
	values := url.Values{}
	values.Add("chat_id", strconv.Itoa(chatID))
	values.Add("text", message)

	_, err := bot.doRequest(values, sendMessageMethod)
	if err != nil {
		return errors.Wrap("unable to send message: %w", err)
	}

	return err
}

func (bot *Bot) doRequest(query url.Values, method string) (content []byte, err error) {
	defer func() {
		errors.WrapIfError("unable to do request: %w", err)
	}()

	apiUrl := url.URL{
		Scheme: "https",
		Host:   bot.host,
		Path:   path.Join(bot.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, errors.Wrap("unable to do request: %w", err)
	}

	request.URL.RawQuery = query.Encode()

	response, err := bot.client.Do(request)
	if err != nil {
		return nil, errors.Wrap("unable to do request: %w", err)
	}
	defer response.Body.Close()

	content, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap("unable to read content: %w", err)
	}

	return content, err
}
