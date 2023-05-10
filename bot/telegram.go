package bot

import (
	"bytes"
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

func (bot *Bot) SendMessage(chatID int, message string, buttons []string) error {
	var botMessage = Message{
		ChatID: chatID,
		Text:   message,
		ReplyMarkup: ReplyKeyboardMarkup{
			Keyboard:       makeButtons(buttons),
			ResizeKeyboard: true,
		},
	}

	_, err := bot.doMessageRequest(botMessage, sendMessageMethod)
	if err != nil {
		return errors.Wrap("unable to send message: %w", err)
	}

	return nil
}

func makeButtons(buttonsText []string) [][]KeyboardButton {
	buttons := make([]KeyboardButton, len(buttonsText))
	for i := 0; i < len(buttonsText); i++ {
		buttons[i] = makeButton(buttonsText[i])
	}
	return [][]KeyboardButton{buttons}
}

func makeButton(button string) KeyboardButton {
	return KeyboardButton{Text: button}
}

func (bot *Bot) doRequest(query url.Values, method string) (content []byte, err error) {
	defer func() {
		err = errors.WrapIfError("unable to do request: %w", err)
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

func (bot *Bot) doMessageRequest(message Message, method string) (content []byte, err error) {
	defer func() {
		err = errors.WrapIfError("unable to do request: %w", err)
	}()

	apiUrl := url.URL{
		Scheme: "https",
		Host:   bot.host,
		Path:   path.Join(bot.basePath, method),
	}

	buf, err := json.Marshal(message)
	if err != nil {
		return nil, errors.Wrap("unable to encode message: %w", err)
	}

	response, err := http.Post(apiUrl.String(), "application/json", bytes.NewBuffer(buf))
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
