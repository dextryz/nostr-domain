package nostr

type Tag []string

func (tag Tag) Key() string {
	if len(tag) > 0 {
		return tag[0]
	}
	return ""
}

func (tag Tag) Value() string {
	if len(tag) > 1 {
		return tag[1]
	}
	return ""
}

func (tag Tag) Encode(out []byte) []byte {
	out = append(out, '[')
	for i, s := range tag {
		if i > 0 {
			out = append(out, ',')
		}
		// TODO: Escape string as describes by RFC8259
		//out = append(out, []byte(s)...)
		out = EscapeString(out, s)
	}
	out = append(out, ']')
	return out
}

type Tags []Tag

// Encoding appends the JSON encoded byte of Tags as [][]string to out.
func (tags Tags) Encode(out []byte) []byte {
	out = append(out, '[')
	for i, tag := range tags {
		if i > 0 {
			out = append(out, ',')
		}
		out = tag.Encode(out)
	}
	out = append(out, ']')
	return out
}
