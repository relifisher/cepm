package wechat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"cepm-backend/config"
)

const (
	wechatWorkAPIHost = "https://qyapi.weixin.qq.com/cgi-bin"
	accessTokenURL    = wechatWorkAPIHost + "/gettoken?corpid=%s&corpsecret=%s"
	userInfoURL       = wechatWorkAPIHost + "/user/getuserinfo?access_token=%s&code=%s"
	userDetailURL     = wechatWorkAPIHost + "/user/get?access_token=%s&userid=%s"
)

// AccessTokenResponse defines the structure of the access token API response.
type AccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// UserInfoResponse defines the structure of the user info API response (from code).
type UserInfoResponse struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	UserID     string `json:"userid"`
	UserTicket string `json:"user_ticket"`
	OpenID     string `json:"openid"`
}

// UserDetailResponse defines the structure of the user detail API response (from userid).
type UserDetailResponse struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	// Add other fields you might need
}

// WechatClient manages interactions with the WeChat Work API.
type WechatClient struct {
	corpID      string
	corpSecret  string
			agentID    int64
	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

// NewWechatClient creates a new WechatClient instance.
func NewWechatClient(cfg *config.WechatConfig) *WechatClient {
	return &WechatClient{
		corpID:     cfg.CorpID,
		corpSecret: cfg.CorpSecret,
		agentID:    cfg.AgentID,
	}
}

// GetAccessToken retrieves a new access token or returns the cached one if valid.
func (c *WechatClient) GetAccessToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.expiresAt) {
		return c.accessToken, nil // Return cached token if valid
	}

	// Fetch new token
	url := fmt.Sprintf(accessTokenURL, c.corpID, c.corpSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read access token response: %w", err)
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal access token response: %w", err)
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("wechat work API error (%d): %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	c.accessToken = tokenResp.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second) // Subtract 60 seconds for buffer

	return c.accessToken, nil
}

// GetUserInfoByCode gets user basic info (userid) using the code from OAuth.
func (c *WechatClient) GetUserInfoByCode(code string) (*UserInfoResponse, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(userInfoURL, accessToken, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var userInfoResp UserInfoResponse
	if err := json.Unmarshal(body, &userInfoResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info response: %w", err)
	}

	if userInfoResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat work API error (%d): %s", userInfoResp.ErrCode, userInfoResp.ErrMsg)
	}

	return &userInfoResp, nil
}

// GetUserDetail retrieves detailed user information using userid.
func (c *WechatClient) GetUserDetail(userid string) (*UserDetailResponse, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(userDetailURL, accessToken, userid)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user detail: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user detail response: %w", err)
	}

	var userDetailResp UserDetailResponse
	if err := json.Unmarshal(body, &userDetailResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user detail response: %w", err)
	}

	if userDetailResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat work API error (%d): %s", userDetailResp.ErrCode, userDetailResp.ErrMsg)
	}

	return &userDetailResp, nil
}
