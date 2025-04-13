package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type Client struct {
	_client *http.Client
}

type retweetUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type retweetUsersWrapper struct {
	Data []retweetUser `json:"data"`
}

func New(key, secret string) (*Client, error) {
	client, err := twitterClient(key, secret)
	if err != nil {
		return nil, err
	}

	return &Client{
		_client: client,
	}, nil
}

func twitterClient(key, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.x.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{}
	return config.Client(context.Background(), &token), nil
}

func (c *Client) Retweeters(tweetId string) ([]string, error) {
	url := fmt.Sprintf("https://api.x.com/2/tweets/%s/retweeted_by", tweetId)
	res, err := c._client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var retweets retweetUsersWrapper
	err = json.Unmarshal(body, &retweets)
	if err != nil {
		return nil, err
	}

	var usernames = make([]string, 0, len(retweets.Data))
	for _, retweet := range retweets.Data {
		usernames = append(usernames, retweet.Username)
	}

	return usernames, nil
}
