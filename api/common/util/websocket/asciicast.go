package websocket

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
)

type castV2Header struct {
	Version      uint               `json:"version"`
	Width        int                `json:"width"`
	Height       int                `json:"height"`
	Timestamp    int64              `json:"timestamp,omitempty"`
	Duration     float64            `json:"duration,omitempty"`
	Title        string             `json:"title,omitempty"`
	Command      string             `json:"command,omitempty"`
	Env          *map[string]string `json:"env,omitempty"`
	outputStream *json.Encoder
}

func newCastV2(meta castV2Header, stream *bytes.Buffer) (*castV2Header, *bytes.Buffer) {
	c := castV2Header{
		Version:      2,
		Width:        meta.Width,
		Height:       meta.Height,
		Title:        meta.Title,
		Timestamp:    meta.Timestamp,
		Duration:     meta.Duration,
		Env:          meta.Env,
		outputStream: json.NewEncoder(stream),
	}
	_ = c.outputStream.Encode(c)
	return &c, stream
}

func (c *castV2Header) Record(t float64, data []byte, event string) {
	out := []interface{}{t, event, string(data)}
	c.Duration = t
	_ = c.outputStream.Encode(out)
}

func zlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, _ = w.Write(src)
	_ = w.Close()
	return in.Bytes()
}
