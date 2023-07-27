package nostr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// TODO: Split this into DecodeMessageFromRelay and DecodeMessageFromClient.
func DecodeMessage(msg []byte) Message {

	// Extract message label (EVENT, REQ, CLOSE, EOSE, NOTICE) from byte slice.
	firstComma := bytes.Index(msg, []byte{','})
	if firstComma == -1 {
		return nil
	}
	label := msg[0:firstComma]

	var v Message
	switch {
	case bytes.Contains(label, []byte("EVENT")):
		v = &MessageEvent{}
	case bytes.Contains(label, []byte("REQ")):
		v = &MessageReq{}
	case bytes.Contains(label, []byte("OK")):
		v = &MessageOk{}
	case bytes.Contains(label, []byte("EOSE")):
		v = &MessageEose{}
	default:
		log.Fatalln("cannot decode message")
	}

	if err := v.UnmarshalJSON(msg); err != nil {
		return nil
	}

	return v
}

type MessageType string

type Message interface {

	// Return the message type.
	Type() MessageType

	// Implement json.Unmarshaler interface
	UnmarshalJSON([]byte) error

	// Implement json.Marshaler interface
	MarshalJSON() ([]byte, error)
}

// NIP-20 - ["OK", <event_id>, <true|false>, <message>]

type MessageOk struct {
	EventId string
	Ok      bool
	Message string
}

func (s MessageOk) GetEventId() string {
	return strings.Trim(s.EventId, "\"")
}

func (s MessageOk) Type() MessageType {
	return MessageType("OK")
}

func (s *MessageOk) UnmarshalJSON(data []byte) error {

	var tmp []json.RawMessage

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln("unable to unmarshal result msg from relay")
	}

	s.EventId = string(tmp[1])

	s.Ok = false
	if string(tmp[2]) == "true" {
		s.Ok = true
	}

	s.Message = string(tmp[3])

	return nil
}

func (s MessageOk) MarshalJSON() ([]byte, error) {

	msg := append([]byte(nil), []byte(`["OK",`)...)

	if len(s.EventId) != 0 {
		msg = append(msg, []byte(s.EventId+`,`)...)
	}

	if s.Ok {
		msg = append(msg, []byte("true")...)
	} else {
		msg = append(msg, []byte("false")...)
	}

	msg = append(msg, []byte(s.Message+`]`)...)

	return msg, nil
}

// NIP-01 - ["EVENT", <subscription_id>, <event JSON as defined above>]

type MessageEvent struct {
	SubscriptionId string
	Event
}

// Jesus christ, the Subscription map didnt want to map unless trimmed. Mother fucker. Fuck trimming string always. jesus christ
func (s MessageEvent) GetSubId() string {
	return strings.Trim(s.SubscriptionId, "\"")
}

func (s MessageEvent) Type() MessageType {
	return MessageType("EVENT")
}

func (s *MessageEvent) UnmarshalJSON(data []byte) error {

	var tmp []json.RawMessage

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln("unable to unmarshal EVENT msg")
	}

	// Is len(2) then it is an event from client-to-relay
	// Otherwise its an event from relay-to-client
	switch len(tmp) {
	case 2:
		return json.Unmarshal(tmp[1], &s.Event)
	case 3:
		s.SubscriptionId = string(tmp[1])
		return json.Unmarshal(tmp[2], &s.Event)
	default:
		return fmt.Errorf("failed to decode EVENT message")
	}
}

func (s MessageEvent) MarshalJSON() ([]byte, error) {

	msg := append([]byte(nil), []byte(`["EVENT",`)...)

	if len(s.SubscriptionId) != 0 {
		msg = append(msg, []byte(s.SubscriptionId+`,`)...)
	}

	// Marshal the signed event to a slice of bytes ready for transmission.
	bytes, err := json.Marshal(s.Event)
	if err != nil {
		log.Fatalln("unable to marchal incoming event")
	}

	msg = append(msg, bytes...)

	msg = append(msg, []byte(`]`)...)

	return msg, nil
}

// NIP-01 - ["REQ", <subscription_id>, <filters JSON>...]

type MessageReq struct {
	SubscriptionId string
	Filters
}

func (s MessageReq) GetSubId() string {
	return strings.Trim(s.SubscriptionId, "\"")
}

func (s MessageReq) Type() MessageType {
	return MessageType("REQ")
}

func (s *MessageReq) UnmarshalJSON(data []byte) error {

	var tmp []json.RawMessage

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln("unable to unmarshal REQ msg")
	}

	s.SubscriptionId = string(tmp[1])

	f := 0
	s.Filters = make(Filters, len(tmp)-2)
	for i := 2; i < len(tmp); i++ {
		err := json.Unmarshal(tmp[i], &s.Filters[f])
		if err != nil {
			return fmt.Errorf("filter %d err: %w", f, err)
		}
		f += 1
	}

	return nil
}

func (s MessageReq) MarshalJSON() ([]byte, error) {

	msg := []byte(nil)

	// Open message array.
	msg = append(msg, []byte(`[`)...)

	// Add message label
	msg = append(msg, []byte(`"REQ",`)...)

	// Add subscription ID between string braces.
	msg = append(msg, []byte(`"`+s.GetSubId()+`",`)...)

	for i, v := range s.Filters {

		// Encode the individual filter.
		bytes, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		// Add filter data to json list.
		msg = append(msg, bytes...)

		// Add delimiter to next item, except the last one
		if i != len(s.Filters)-1 {
			msg = append(msg, []byte(`,`)...)
		}
	}

	// Close the entire message.
	msg = append(msg, []byte(`]`)...)

	return msg, nil
}

// NIP-01 - ["EOSE", <subscription_id>]

type MessageEose struct {
	SubscriptionId string
}

func (s MessageEose) GetSubId() string {
	return strings.Trim(s.SubscriptionId, "\"")
}

func (s MessageEose) Type() MessageType {
	return MessageType("EOSE")
}

func (s *MessageEose) UnmarshalJSON(data []byte) error {

	var tmp []json.RawMessage

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln("unable to unmarshal REQ msg")
	}

	s.SubscriptionId = string(tmp[1])

	return nil
}

func (s MessageEose) MarshalJSON() ([]byte, error) {

	msg := []byte(nil)

	// Open message array.
	msg = append(msg, []byte(`[`)...)

	// Add message label
	msg = append(msg, []byte(`"REQ",`)...)

	// Add subscription ID between string braces.
	msg = append(msg, []byte(`"`+s.SubscriptionId+`"`)...)

	// Close the entire message.
	msg = append(msg, []byte(`]`)...)

	return msg, nil
}
