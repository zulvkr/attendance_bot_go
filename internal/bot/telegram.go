package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TelegramAPI handles all Telegram Bot API interactions
type TelegramAPI struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// Update represents a Telegram update
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

// Message represents a Telegram message
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text,omitempty"`
	Date      int64  `json:"date"`
}

// User represents a Telegram user
type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// GetUpdatesResponse represents the response from getUpdates
type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// SendMessageResponse represents the response from sendMessage
type SendMessageResponse struct {
	OK     bool    `json:"ok"`
	Result Message `json:"result"`
}

// NewTelegramAPI creates a new Telegram API client
func NewTelegramAPI(token string) *TelegramAPI {
	return &TelegramAPI{
		token:   token,
		baseURL: "https://api.telegram.org/bot" + token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUpdates retrieves updates from Telegram
func (api *TelegramAPI) GetUpdates(offset int64, timeout int) ([]Update, error) {
	params := url.Values{}
	if offset > 0 {
		params.Set("offset", strconv.FormatInt(offset, 10))
	}
	if timeout > 0 {
		params.Set("timeout", strconv.Itoa(timeout))
	}

	url := api.baseURL + "/getUpdates"
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	resp, err := api.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response GetUpdatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("telegram API error: %s", string(body))
	}

	return response.Result, nil
}

// SendMessage sends a message to a chat
func (api *TelegramAPI) SendMessage(chatID int64, text string) error {
	return api.SendMessageWithOptions(chatID, text, nil)
}

// SendMessageOptions contains optional parameters for sending messages
type SendMessageOptions struct {
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageID      int64  `json:"reply_to_message_id,omitempty"`
}

// SendMessageWithOptions sends a message with additional options
func (api *TelegramAPI) SendMessageWithOptions(chatID int64, text string, options *SendMessageOptions) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	if options != nil {
		if options.ParseMode != "" {
			payload["parse_mode"] = options.ParseMode
		}
		if options.DisableWebPagePreview {
			payload["disable_web_page_preview"] = true
		}
		if options.DisableNotification {
			payload["disable_notification"] = true
		}
		if options.ReplyToMessageID > 0 {
			payload["reply_to_message_id"] = options.ReplyToMessageID
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := api.httpClient.Post(
		api.baseURL+"/sendMessage",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.OK {
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

// SendDocument sends a document to a chat
func (api *TelegramAPI) SendDocument(chatID int64, document io.Reader, filename string) error {
	// This is a simplified implementation
	// In a full implementation, you'd use multipart/form-data
	return fmt.Errorf("sendDocument not implemented yet")
}

// GetMe returns basic information about the bot
func (api *TelegramAPI) GetMe() (*User, error) {
	resp, err := api.httpClient.Get(api.baseURL + "/getMe")
	if err != nil {
		return nil, fmt.Errorf("failed to get bot info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response struct {
		OK     bool `json:"ok"`
		Result User `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("telegram API error: %s", string(body))
	}

	return &response.Result, nil
}
