package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"dodevops-api/common/util"
	"dodevops-api/pkg/db"
	appLog "dodevops-api/pkg/log"

	"github.com/gorilla/websocket"
)

type wsMsg struct {
	Type  string `json:"type"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
	Close bool   `json:"close"`
}

type wsConn struct {
	sync.RWMutex
	Ws *websocket.Conn
}

func (c *wsConn) WriteMessage(messageType int, data []byte) error {
	c.Lock()
	err := c.Ws.WriteMessage(messageType, data)
	c.Unlock()
	return err
}

func NewWsConn(conn *websocket.Conn) *wsConn {
	return &wsConn{Ws: conn}
}

type RecordData struct {
	Event string  `json:"event"`
	Time  float64 `json:"time"`
	Data  []byte  `json:"data"`
}

type Meta struct {
	TERM      string
	Width     int
	Height    int
	UserName  string
	ConnectId string
	HostId    uint
	HostName  string
}

type WebSocketStream struct {
	sync.RWMutex
	Terminal    *Terminal
	Conn        *wsConn
	messageType int
	recorder    []*RecordData
	CreatedAt   util.HTime
	UpdatedAt   util.HTime
	Meta        Meta
	written     bool
	closed      bool
}

func NewWebSocketSteam(terminal *Terminal, connection *wsConn, meta Meta) *WebSocketStream {
	return &WebSocketStream{
		Terminal:    terminal,
		Conn:        connection,
		messageType: websocket.BinaryMessage,
		CreatedAt: util.HTime{
			Time: time.Now(),
		},
		UpdatedAt: util.HTime{
			Time: time.Now(),
		},
		recorder: make([]*RecordData, 0),
		Meta:     meta,
	}
}

func (r *WebSocketStream) Read(p []byte) (n int, err error) {
	if r.closed {
		return 0, io.EOF
	}

	t, message, err := r.Conn.Ws.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			r.closed = true
		}
		return 0, err
	}

	var msgObj wsMsg
	if t == websocket.TextMessage {
		if err := json.Unmarshal(message, &msgObj); err == nil {
			switch msgObj.Type {
			case "resizePty":
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					r.Meta.Width = msgObj.Cols
					r.Meta.Height = msgObj.Rows
					if err := r.Terminal.session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						appLog.Log().Error(fmt.Sprintf("ssh pty change windows size failed: %v", err))
					}
				}
				return 0, nil
			case "closePty":
				if msgObj.Close {
					r.closed = true
					if err := r.Terminal.Close(); err != nil {
						appLog.Log().Error(fmt.Sprintf("Close pty failed: %v", err))
					}
				}
				return 0, nil
			}
		}
	}

	r.Lock()
	defer r.Unlock()

	r.UpdatedAt = util.HTime{Time: time.Now()}
	r.messageType = t
	n = len(message)
	copy(p, message)
	return
}

func (r *WebSocketStream) Write(p []byte) (n int, err error) {
	if r.closed {
		return 0, io.EOF
	}

	n = len(p)
	if n == 0 {
		return 0, nil
	}

	var msgObj wsMsg
	if json.Unmarshal(p, &msgObj) == nil {
		switch msgObj.Type {
		case "resizePty", "closePty":
			return n, nil
		}
	}

	r.Lock()
	defer r.Unlock()

	data := make([]byte, len(p))
	copy(data, p)
	r.recorder = append(r.recorder, &RecordData{
		Time:  time.Since(r.CreatedAt.Time).Seconds(),
		Event: "o",
		Data:  data,
	})

	if r.Conn != nil {
		_ = r.Conn.Ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err = r.Conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				r.closed = true
			}
			return 0, err
		}
	}

	r.UpdatedAt = util.HTime{Time: time.Now()}
	return n, nil
}

func (r *WebSocketStream) Write2Log() error {
	r.Lock()
	defer r.Unlock()

	if r.written {
		return nil
	}
	recorders := r.recorder
	if len(recorders) != 0 {
		b := new(bytes.Buffer)
		meta := castV2Header{
			Width:     r.Meta.Width,
			Height:    r.Meta.Height,
			Timestamp: time.Now().Unix(),
			Title:     r.Meta.ConnectId,
			Env: &map[string]string{
				"SHELL": "/bin/bash", "TERM": r.Meta.TERM,
			},
		}
		cast, buffer := newCastV2(meta, b)
		for _, v := range recorders {
			cast.Record(v.Time, v.Data, v.Event)
		}
		compressData := zlibCompress(buffer.Bytes())
		if len(compressData) > 320 {
			record := sshRecord{
				ConnectID:   r.Meta.ConnectId,
				HostName:    r.Meta.HostName,
				UserName:    r.Meta.UserName,
				Records:     compressData,
				ConnectTime: r.CreatedAt,
				LogoutTime: util.HTime{
					Time: time.Now(),
				},
				HostId: r.Meta.HostId,
			}
			if db.Db == nil {
				return fmt.Errorf("database is not initialized")
			}
			if err := db.Db.Create(&record).Error; err != nil {
				return err
			}
		}
	}
	r.written = true
	return nil
}
