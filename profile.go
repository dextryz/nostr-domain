package nostr

import (
	"encoding/json"
	"fmt"
	"log"
)

type Profile struct {
	Name    string `json:"name,omitempty"`
	About   string `json:"about,omitempty"`
	Picture string `json:"picture,omitempty"`
}

func (s Profile) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatalln("Unable to convert event to string")
	}
	return string(bytes)
}

func ParseMetadata(e Event) (*Profile, error) {

	if e.Kind != KindSetMetadata {
		return nil, fmt.Errorf("event %s is kind %d, not 0", e.Id, e.Kind)
	}

	var profile Profile
	err := json.Unmarshal([]byte(e.Content), &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata from event %s: %w", e.Id, err)
	}

	return &profile, nil
}
