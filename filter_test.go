package nostr

import (
	"encoding/json"
	"testing"
)

func TestUnitUnmarshal(t *testing.T) {

	cases := []struct {
		name        string
		raw         string
		wantIds     int
		wantAuthors []string
	}{
		{
			name:        "Unmarshal a single filter",
			raw:         `{"ids":["a","b"],"authors":["alice","bob"],"kinds":[1]}`,
			wantIds:     2,
			wantAuthors: []string{"alice", "bob"},
		},
	}

	for _, c := range cases {

		t.Logf("Test: %s", c.name)

		var filter Filter

		err := json.Unmarshal([]byte(c.raw), &filter)
		ok(t, err)
		equals(t, c.wantIds, len(filter.Ids))
		equals(t, c.wantAuthors[0], "alice")

	}
}

func TestUnitMarshal(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "Marshal a single filter",
			want: `{"ids":["a","b"],"authors":["alice","bob"],"kinds":[1]}`,
		},
	}

	for _, c := range cases {

		t.Logf("Test: %s", c.name)

		filter, err := json.Marshal(Filter{
			Ids:     []string{"a", "b"},
			Authors: []string{"alice", "bob"},
			Kinds:   []uint32{1},
		})

		ok(t, err)
		equals(t, c.want, string(filter))
	}
}
