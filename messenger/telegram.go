package messenger

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func NewTelegramMessenger(ChatID int64, TelegramToken string, Version string) *TelegramMessenger {
	var telegram TelegramMessenger
	telegram.ChatID = ChatID
	telegram.TelegramToken = TelegramToken
	telegram.updateID = 0
	telegram.isFirstRun = true
	telegram.Version = Version
	return &telegram
}

type TelegramMessenger struct {
	ChatID        int64
	TelegramToken string
	updateID      int64
	isFirstRun    bool
	Version       string
}

// send msg to telegram
func (t *TelegramMessenger) Send(msg string) error {

	params := "chat_id=" + url.QueryEscape(fmt.Sprint(t.ChatID)) + "&" + "text=" + url.QueryEscape(msg)
	path := fmt.Sprintf("https://api.telegram.org/bot%s/SendMessage?%s", t.TelegramToken, params)

	_, err := http.Get(path)
	if err != nil {
		return err
	}
	return nil
}

// get last admin command from telegram
func (t *TelegramMessenger) Recive() ([]string, error) {

	// create Getupdates api string
	path := fmt.Sprintf("https://api.telegram.org/bot%s/Getupdates?limit=100&offset=%d", t.TelegramToken, t.updateID)

	// get update from telegram
	response, err := http.Get(path)

	if err != nil {
		return nil, err

	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("statusCode: %d", response.StatusCode)
	}

	// read response body
	defer response.Body.Close()
	var getUpdate telegramGetUpdateResponse
	json.NewDecoder(response.Body).Decode(&getUpdate)

	if len(getUpdate.Result) < 1 {

		// if first run and getUpdate page is empty, set isFirstRun to False.
		if t.isFirstRun {
			t.isFirstRun = false
		}

		return nil, errors.New("new message not found")
	}

	// get last admin command
	adminCommand, err := t.lastAdminmessage(getUpdate)

	if err != nil {
		return adminCommand, err
	}

	// ignore first command (old command)
	if t.isFirstRun {
		t.isFirstRun = false
		return nil, errors.New("ignore first command (old messages)")
	}

	return adminCommand, err
}

// extract last admin message and remove previous messages
func (t *TelegramMessenger) lastAdminmessage(getUpdate telegramGetUpdateResponse) ([]string, error) {

	var adminCommanFound bool
	var alreadyExecuted bool
	var thisResult result

	// extract last admin command
	for i := len(getUpdate.Result) - 1; i >= 0; i-- {
		thisResult = getUpdate.Result[i]

		if thisResult.Message.Text == "" {
			continue
		}

		if t.ChatID == thisResult.Message.Chat.Id {
			//Ignore the previous executed messages
			if t.updateID == thisResult.UpdateID {
				alreadyExecuted = true
				break
			}
			adminCommanFound = true
			break
		}
	}

	// clean previous messages
	if !(adminCommanFound || alreadyExecuted) {
		// not any admin command in Result page
		t.updateID = getUpdate.Result[0].UpdateID // remove all Result page messages
	} else if len(getUpdate.Result) >= 100 && getUpdate.Result[0].UpdateID == thisResult.UpdateID {
		// if response page is full && if first response on response page is last admin message
		t.updateID = getUpdate.Result[0].UpdateID // remove all Result page messages
	} else {
		t.updateID = thisResult.UpdateID // remove previous messages
	}

	if !adminCommanFound {
		return nil, errors.New("new admin command not found")
	}
	return strings.Split(thisResult.Message.Text, " "), nil
}

type telegramGetUpdateResponse struct {
	Result []result
}

type result struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Text string
		Chat struct {
			Id int64
		}
	}
}
