package event

import (
	"encoding/json"
	"trust/common"
	"trust/common/entity/message"
)

type Event struct {
	CarId string `json:"car_id"`
	message.Message
}

func (e *Event) ToJsonBytes() []byte {
	eventAsBytes, err := json.Marshal(e)
	if err != nil {
		return common.EMPTY_JSON_BYTE
	}
	return eventAsBytes
}

func FromJSON(data string) (Event, error) {
	eventAsBytes := []byte(data)
	event := Event{}
	err := json.Unmarshal(eventAsBytes, &event)
	return event, err
}
