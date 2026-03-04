package zima

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success int    `json:"success"` // 200
	Message string `json:"message"`
	Data    struct {
		// Token 是个对象，不是字符串
		Token struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresAt    int64  `json:"expires_at"`
		} `json:"token"`

		User interface{} `json:"user"`
	} `json:"data"`
}

type ShareEntry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(ip string) *Client {
	return &Client{
		BaseURL:    fmt.Sprintf("http://%s", ip),
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Login execute
func (c *Client) Login(username, password string) error {
	url := c.BaseURL + "/v1/users/login"

	reqBody, _ := json.Marshal(LoginRequest{Username: username, Password: password})

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("登录请求失败，状态码: %d, 错误：%s", resp.StatusCode, string(bodyBytes))
	}

	var res LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if res.Success != 200 {
		return fmt.Errorf("API 错误: %s", res.Message)
	}

	c.Token = res.Data.Token.AccessToken

	if c.Token == "" {
		return fmt.Errorf("未获取到 Access Token")
	}

	return nil
}

// GetShareList获取共享列表，按照那个文档来
func (c *Client) GetShareList() ([]string, error) {
	url := c.BaseURL + "/v2_1/files/share?status=true"

	req, _ := http.NewRequest("GET", url, nil)

	// 通常 JWT 是放在 Authorization 头的
	// 假如不需要 Bearer 前缀，就用 c.Token；如果报错401，试试 "Bearer " + c.Token
	req.Header.Set("Authorization", c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取列表失败，状态码: %d，错误详情：%s", resp.StatusCode, string(bodyBytes))
	}

	var shares []ShareEntry
	if err := json.NewDecoder(resp.Body).Decode(&shares); err != nil {
		return nil, fmt.Errorf("解析列表JSON失败: %v", err)
	}

	var names []string
	for _, s := range shares {
		if s.Status == "enable" {
			names = append(names, s.Name)
		}
	}
	return names, nil
}
