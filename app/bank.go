package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type (
	Bank struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		HalfWidthKana string `json:"halfWidthKana"`
		FullWidthKana string `json:"fullWidthKana"`
		Hiragana      string `json:"hiragana"`
	}

	Banks struct {
		Data       []*Bank `json:"banks"`
		Size       int     `json:"size"`
		Limit      int     `json:"limit"`
		HasNext    bool    `json:"hasNext"`
		NextCursor string  `json:"nextCursor"`
		HasPrev    bool    `json:"hasPrev"`
		Version    string  `json:"version"`
	}
)

func (c *Client) GetBank(ctx context.Context, code string, param *GetParameter) (*Bank, error) {
	u, err := url.Parse(c.base.String() + "/banks/" + code)
	if err != nil {
		return nil, fmt.Errorf("generate URL: %w", err)
	}
	req, err := c.GetRequest(ctx, u, param)
	if err != nil {
		return nil, fmt.Errorf("generate request: %w", err)
	}

	var res Bank
	err = c.Call(ctx, req, func(resp io.ReadCloser) error {
		if err := json.NewDecoder(resp).Decode(&res); err != nil {
			return fmt.Errorf("decode to response: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}
