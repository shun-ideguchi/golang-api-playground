package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const BASE_URL = "https://apis.bankcode-jp.com/v3"

type Client struct {
	apiKey      string
	httpClient  *http.Client
	base        *url.URL
	rateLimiter *rate.Limiter
}

type option func(*Client) error

type GetParameter struct {
	Fields []string
}

var endpoint *url.URL

// init はベースURLをパースします
func init() {
	u, err := url.Parse(BASE_URL)
	if err != nil {
		panic(err)
	}

	endpoint = u
}

func NewClient(options ...option) (*Client, error) {
	n := rate.Every(3 * time.Second)
	l := rate.NewLimiter(n, 1)

	c := &Client{
		httpClient:  http.DefaultClient,
		base:        endpoint,
		rateLimiter: l,
	}

	for _, optionFn := range options {
		if err := optionFn(c); err != nil {
			return nil, fmt.Errorf("initialize client: %w", err)
		}
	}

	return c, nil
}

func (c *Client) Call(ctx context.Context, req *http.Request, f func(resp io.ReadCloser) error) (err error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to api: %w", err)
	}
	defer func() {
		defer resp.Body.Close()
		if _, dErr := io.Copy(io.Discard, resp.Body); dErr != nil {
			err = dErr
		}
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http status: %s", resp.Status)
	}

	return f(resp.Body)
}

// getRequest はGET通信のリクエストを作成します。
func (c *Client) GetRequest(ctx context.Context, u *url.URL, param *GetParameter) (*http.Request, error) {
	query := u.Query()
	query.Add("apiKey", c.apiKey)
	if len(param.Fields) > 0 {
		query.Add("fields", strings.Join(param.Fields, ","))
	}
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("generate http get request: %w", err)
	}

	return req, nil
}

// WithEndpoint はデフォルトのエンドポイントを変更します
func WithEndpoint(endpoint string) option {
	return func(c *Client) error {
		if endpoint == "" {
			return errors.New("URL is empty")
		}
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		c.base = u
		return nil
	}
}

// WithApiKey はクライアントにAPIキーをセットします
func WithApiKey(key string) option {
	return func(c *Client) error {
		if key == "" {
			return fmt.Errorf("API key is empty")
		}

		c.apiKey = key
		return nil
	}
}
