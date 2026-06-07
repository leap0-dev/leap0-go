// Package sse implements a Server-Sent Events stream reader.
package sse

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// Stream reads Server-Sent Events from an HTTP response body.
type Stream struct {
	reader *bufio.Reader
	closer io.Closer
	done   bool
}

func New(body io.ReadCloser) *Stream {
	return &Stream{reader: bufio.NewReader(body), closer: body}
}

func (s *Stream) Next() (json.RawMessage, error) {
	if s.done {
		return nil, io.EOF
	}
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			s.done = true
			return nil, io.EOF
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				s.done = true
				return nil, io.EOF
			}
			return json.RawMessage(data), nil
		}
	}
}

func (s *Stream) Close() error {
	s.done = true
	return s.closer.Close()
}
