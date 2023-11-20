package client

import (
	"reflect"
	"testing"

	"github.com/LiU-SeeGoals/controller/internal/action"
)

type MockConnection struct {
    SentMessages [][]byte
}

func (m *MockConnection) Write(b []byte) (n int, err error) {
    m.SentMessages = append(m.SentMessages, b)
    return len(b), nil
}

func (m *MockConnection) Close() error {
    return nil
}

func TestSendAction(t *testing.T) {
    testCases := []struct {
        actions         []action.Action
        expectedMessages [][]byte
    }{
        {
            []action.Action{&action.Stop{Id: 1}},
            [][]byte{{0x03, 0x00, 0x01}}, // Replace with the actual expected bytes for Stop{Id: 1}
        },
        {
            []action.Action{&action.Stop{Id: 2}},
            [][]byte{{0x03, 0x00, 0x02}}, // Replace with the actual expected bytes for Stop{Id: 2}
        },
        // Add more test cases with different actions and their expected byte representations
    }

    for _, tc := range testCases {
        mockConn := &MockConnection{}
        client := NewBaseStationClient("localhost:8080")
        client.connection = mockConn

        client.SendActions(tc.actions)

        if len(mockConn.SentMessages) != len(tc.expectedMessages) {
            t.Errorf("Expected %d messages to be sent, got %d", len(tc.expectedMessages), len(mockConn.SentMessages))
            continue
        }

        for i, msg := range mockConn.SentMessages {
            if !reflect.DeepEqual(msg, tc.expectedMessages[i]) {
                t.Errorf("Sent message #%d = %v, want %v", i+1, msg, tc.expectedMessages[i])
            }
        }
    }
}
