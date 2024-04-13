package tinybot

import (
	"container/list"
	"fmt"
	"github.com/gorilla/websocket"
)

type MockDriver struct {
	Dialed bool //链接状态
	Event  list.List
	Resp   map[string]APIResponse
	Error  map[string]error
}

func (m *MockDriver) Dial() {
	m.Dialed = true
}

func NewMockDriver() *MockDriver {
	return &MockDriver{
		Dialed: false,
		Event:  list.List{},
		Resp:   make(map[string]APIResponse),
		Error:  make(map[string]error),
	}
}

func (m *MockDriver) IsDialed() {
	if !m.Dialed {
		panic("MockDriver.RetrieveEvent() called before Dial()")
	}
}

func (m *MockDriver) CallApi(req APIRequest) (APIResponse, error) {
	m.IsDialed()

	if resp, ok := m.Resp[req.Echo]; !ok {
		return APIResponse{}, fmt.Errorf("[MockDriver] cannot produce response body")
	} else if err, ok1 := m.Error[req.Echo]; ok1 {
		return resp, err
	} else {
		return resp, nil
	}
}

func (m *MockDriver) RetrieveEvent() (Event, error) {
	m.IsDialed()

	if !(m.Event.Len() > 0) {
		err := &websocket.CloseError{Code: websocket.CloseNormalClosure}
		return Event{}, err
	}

	return m.Event.Remove(m.Event.Front()).(Event), nil
}

func (m *MockDriver) SendEvent(event Event) {
	m.Event.PushBack(event)
}

func (m *MockDriver) SetResponse(echo string, resp APIResponse, err error) {
	if err != nil {
		m.Error[echo] = err
	}

	m.Resp[echo] = resp
}
