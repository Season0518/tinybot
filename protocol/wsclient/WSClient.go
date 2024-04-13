package wsclient

import (
	"core/pkg/tinybot"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// WSClient Todo: 迁移所有API操作到此处，增加消息的异步获取
type WSClient struct {
	lock      sync.Mutex
	seqMap    sync.Map
	Conn      *websocket.Conn
	Addr      string
	Token     string
	EventChan chan tinybot.Event
}

func NewClient(addr, token string) *WSClient {
	return &WSClient{
		Conn:      nil,
		lock:      sync.Mutex{},
		seqMap:    sync.Map{},
		Addr:      addr,
		Token:     token,
		EventChan: make(chan tinybot.Event, 128), //最多存放128个事件
	}
}

func (w *WSClient) CallApi(req tinybot.APIRequest) (tinybot.APIResponse, error) {
	req.Echo = uuid.NewV5(uuid.Must(uuid.NewV4()), "send_group_msg").String()

	ch := make(chan tinybot.APIResponse, 1)
	w.seqMap.Store(req.Echo, ch)

	w.lock.Lock()
	err := w.Conn.WriteJSON(req)
	w.lock.Unlock()

	if err != nil {
		return tinybot.APIResponse{}, fmt.Errorf("向WebSocket服务发送消息时发生异常: %v", err)
	}

	select {
	case <-time.After(10 * time.Second):
		return tinybot.APIResponse{}, fmt.Errorf("获取事件响应超时: %v", os.ErrDeadlineExceeded)
	case rsp, ok := <-ch:
		if !ok {
			return tinybot.APIResponse{}, fmt.Errorf("在读取响应时发生异常: %v", io.ErrClosedPipe)
		}
		return rsp, nil
	}
}

func (w *WSClient) Dial() {
	var err error
	header := http.Header{}
	header.Add("Authorization", "Bearer "+w.Token)
	for attempt := 0; attempt < 3; attempt++ {
		w.Conn, _, err = websocket.DefaultDialer.Dial(w.Addr, header)
		if err == nil {
			log.Info().Msg("[driver] 成功与Websocket服务器建立链接")
			return
		}
		log.Warn().Msg("[driver] 无法与Websocket服务器建立链接，正在重试...")
		time.Sleep(3 * time.Second)
	}
	log.Panic().Msg("[driver] 无法和Websocket服务器建立连接，请检查服务器地址和Token是否正确")
}

func (w *WSClient) RetrieveEvent() (tinybot.Event, error) {
	t, data, err := w.Conn.ReadMessage()
	//如果是WebSocket链接造成的错误，则不对data进行额外处理。
	if err != nil {
		return tinybot.Event{}, err
	}

	if t != websocket.TextMessage {
		return tinybot.Event{}, nil
	}

	rsp := gjson.ParseBytes(data)

	// 处理回调事件，这里暂时不对回调事件进行处理，默认请求均为成功。
	if rsp.Get("echo").Exists() {
		if ch, ok := w.seqMap.LoadAndDelete(rsp.Get("echo").String()); ok {
			ch.(chan tinybot.APIResponse) <- func() (result tinybot.APIResponse) {
				_ = json.Unmarshal(data, &result)
				return
			}()
			close(ch.(chan tinybot.APIResponse))
			return tinybot.Event{}, nil
		}
	}

	result := tinybot.Event{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		if gjson.GetBytes(data, "post_type").String() == "message" && !gjson.GetBytes(data, "message").IsArray() {
			return tinybot.Event{}, fmt.Errorf("在反序列化Event时失败，可能是服务端未启用Array上报")
		} else {
			return tinybot.Event{}, fmt.Errorf("[driver] 在反序列化Event时出现不可预料的错误")
		}
	}

	return result, nil
}

func (w *WSClient) ExceptionHandler() {

}
