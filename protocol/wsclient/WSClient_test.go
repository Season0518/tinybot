package wsclient

import (
	"testing"
	"time"
)

func TestGeneric(t *testing.T) {
	client := NewClient("ws://localhost:1443", "1145141919810")

	//go client.Listen()
	for {
		select {
		case event := <-client.EventChan:
			t.Logf("Received event: %v", event)
		}
	}

}

func TestBusyBlock(t *testing.T) {
	client := NewClient("ws://localhost:1443", "1145141919810")
	client.Dial()

	// Simulate Process is busing now
	t.Logf("Busying")
	time.Sleep(5 * time.Second)
	t.Logf("Free")

	for {
		event, _ := client.RetrieveEvent()
		t.Logf("Received event: %v", event)
	}
}
