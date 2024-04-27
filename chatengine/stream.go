package chatengine

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type Stream struct {
	ch      chan *ChatResponse
	body    io.ReadCloser
	sc      *bufio.Scanner
	respCnt int
}

func NewStream(ch chan *ChatResponse, body io.ReadCloser) *Stream {
	return &Stream{
		ch:   ch,
		body: body,
		sc:   bufio.NewScanner(body),
	}
}

func (s *Stream) Recv() {
	for {
		if !s.sc.Scan() {
			err := s.sc.Err()
			if err == nil && s.respCnt == 0 {
				err = ErrNoResponse
			}
			if err != nil {
				s.ch <- &ChatResponse{Err: err}
			}
			s.close()
			return
		}

		text := s.sc.Text()
		text = strings.TrimPrefix(text, "data: ")
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		if text == "[DONE]" {
			s.close()
			return
		}

		var resp ChatResponse
		err := json.Unmarshal([]byte(text), &resp)
		if err != nil {
			s.ch <- &ChatResponse{Err: err}
			s.close()
			return
		}

		s.respCnt++
		s.ch <- &resp
	}
}

func (s *Stream) close() {
	close(s.ch)
	s.body.Close()
}
