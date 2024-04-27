package chatengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// API接口文档: https://platform.moonshot.cn/docs/api-reference
type Engine struct {
	cfg        Config
	httpClient *http.Client
}

func NewEngine(cfg Config) *Engine {
	return &Engine{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

// 流式响应
func (e *Engine) StreamRequest(ctx context.Context, req ChatReq) (chan *ChatResponse, error) {
	req.Stream = true
	httpResp, err := e.prepare(ctx, &req)
	if err != nil {
		return nil, err
	}

	code := httpResp.StatusCode
	if code != http.StatusOK {
		defer httpResp.Body.Close()
		bytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}
		log.Println("error: http status code:", code, ", body: "+string(bytes))
		if code == http.StatusTooManyRequests {
			return nil, ErrRateLimit
		}
		return nil, fmt.Errorf("invalid http status %v", code)
	}

	// httpResp.Body will be closed by Stream
	ch := make(chan *ChatResponse)
	stream := NewStream(ch, httpResp.Body)
	go stream.Recv()
	return ch, nil
}

// 非流式响应
func (e *Engine) ChatRequest(ctx context.Context, req ChatReq) (*ChatResponse, error) {
	req.Stream = false
	httpResp, err := e.prepare(ctx, &req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	code := httpResp.StatusCode
	if code != http.StatusOK {
		if code == http.StatusTooManyRequests {
			return nil, ErrRateLimit
		}
		return nil, fmt.Errorf("invalid http status %v", code)
	}

	bytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ChatResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (e *Engine) prepare(ctx context.Context, req *ChatReq) (*http.Response, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := "https://api.moonshot.cn/v1/chat/completions"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+e.cfg.ApiKey)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	return e.httpClient.Do(request)
}
