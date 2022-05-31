package message

import (
	"encoding/json"
	"github.com/google/uuid"
	"trust/common"
)

type MessageType int

const (
	Accident MessageType = iota
	Communication
	Multimedia
)

type Message struct {
	Id      int         `json:"message_id"`
	Type    MessageType `json:"message_type"`
	Content string      `json:"message_content"`
}

func (m *Message) ToJsonBytes() []byte {
	messageAsBytes, err := json.Marshal(m)
	if err != nil {
		return common.EMPTY_JSON_BYTE
	}
	return messageAsBytes
}

func NewMessageId() string {
	return uuid.NewString()
}
