package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"go-simple-tg-bot/internal/models"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func (c *Client) Updates(ctx context.Context, offset, limit int) ([]models.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, "getUpdates", q)
	if err != nil {
		return nil, err
	}

	var res models.UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, "sendMessage", q)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendPhotoByURL(ctx context.Context, chatID int, photoURL, caption string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", photoURL)
	if caption != "" {
		q.Add("caption", caption)
	}

	_, err := c.doRequest(ctx, "sendPhoto", q)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if req.Response.StatusCode != 200 {
		return nil, fmt.Errorf("Метод %s вернул код статус %d", method, req.Response.StatusCode)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
