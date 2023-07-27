package nostr

import "strings"

type Filter struct {

	// An id here can be a prefix string to an event ID.
	Ids []string `json:"ids,omitempty"`

	// Filter by author public keys.
	Authors []string `json:"authors,omitempty"`

	// Only broadcast messages with kind in list.
	Kinds []uint32 `json:"kinds,omitempty"`

	// Maximum number of events to be returned in the initial query.
	Limit int `json:"limit,omitempty"`
}

type Filters []Filter

func (s Filters) Match(event *Event) bool {
	for _, f := range s {
		if f.Matches(event) {
			return true
		}
	}
	return false
}

func (s Filter) Matches(event *Event) bool {

	if event == nil {
		return false
	}

	if s.Ids != nil && !containsPrefix(s.Ids, event.Id) {
		return false
	}

	if s.Authors != nil && !containsAuthor(s.Authors, event.PubKey) {
		return false
	}

	if s.Kinds != nil && !contains(s.Kinds, event.Kind) {
		return false
	}

	return true
}

func containsAuthor(authors []string, pub string) bool {
	for _, v := range authors {
		if strings.HasPrefix(pub, v) {
			return true
		}
	}
	return false
}

func containsPrefix(prefixlist []string, id string) bool {
	for _, prefix := range prefixlist {
		if strings.HasPrefix(string(id), prefix) {
			return true
		}
	}
	return false
}

func contains(s []uint32, target uint32) bool {
	for _, item := range s {
		if item == target {
			return true
		}
	}
	return false
}
