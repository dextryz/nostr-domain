package nostr

import (
	"encoding/json"
	"testing"
)

func TestUnit_MessageEvent(t *testing.T) {
	cases := []struct {
		name      string
		msg       string
		wantId    int
		wantSubId string
	}{
		{
			name:      "Parse EVENT from relay",
			msg:       `["EVENT","_",{"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"","created_at":1644271588,"kind":1,"tags":null,"content":"ping","sig":""}]`,
			wantId:    64,
			wantSubId: "_",
		},
		{
			name:      "Parse EVENT from client",
			msg:       `["EVENT",{"id":"dc90c95f09947507c1044e8f48bcf6350aa6bff1507dd4acfc755b9239b5c962","pubkey":"","created_at":1644271588,"kind":1,"tags":null,"content":"ping","sig":""}]`,
			wantId:    64,
			wantSubId: "",
		},
	}

	for _, c := range cases {

		t.Logf("Test: %s", c.name)

		var msg MessageEvent

		err := json.Unmarshal([]byte(c.msg), &msg)
		ok(t, err)
		equals(t, c.wantId, len(msg.Id))
		equals(t, c.wantSubId, msg.GetSubId())

		msgJson, err := json.Marshal(msg)
		ok(t, err)
		equals(t, c.msg, string(msgJson))
	}
}

func TestUnit_MessageReq(t *testing.T) {
	cases := []struct {
		name       string
		msg        string
		wantFilter int
		wantIds    int
		wantKinds  int
	}{
		{
			name:       "Parse REQ message",
			msg:        `["REQ","2",{"ids":["a","b"],"authors":["alice"],"kinds":[1]}]`,
			wantFilter: 1,
			wantIds:    2,
			wantKinds:  1,
		},
	}

	for _, c := range cases {

		t.Logf("Test: %s", c.name)

		var msg MessageReq

		err := json.Unmarshal([]byte(c.msg), &msg)
		ok(t, err)
		equals(t, c.wantFilter, len(msg.Filters))
		equals(t, c.wantIds, len(msg.Filters[0].Ids))
		equals(t, c.wantKinds, len(msg.Filters[0].Kinds))

		msgJson, err := json.Marshal(msg)
		ok(t, err)
		equals(t, c.msg, string(msgJson))
	}
}
