package server

import "encoding/json"

type ParserJson struct{}

func (p ParserJson) ParseEvent(payload []byte) (*EventJson, error) {
	var event *EventJson
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}
	return event, nil
}
