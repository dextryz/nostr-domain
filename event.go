package nostr

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type Status string

const (
	StatusOK   Status = "OK"
	StatusFail Status = "FAIL"
)

const (
	KindSetMetadata     uint32 = 0
	KindTextNote        uint32 = 1
	KindRecommendServer uint32 = 2
)

type Event struct {
	Id        string    `json:"id"`
	PubKey    string    `json:"pubkey"`
	CreatedAt Timestamp `json:"created_at"`
	Kind      uint32    `json:"kind"`
	Tags      Tags      `json:"tags"`
	Content   string    `json:"content"`
	Sig       string    `json:"sig"`
}

// To obtain the event id, we sha256 the serialized event.
func (s Event) GetId() string {
	h := sha256.Sum256(s.Serialize())
	return hex.EncodeToString(h[:])
}

func (s Event) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatalln("Unable to convert event to string")
	}
	return string(bytes)
}

// The serialization is done over the UTF-8 JSON-serialized string (with no white space or line breaks).
func (s Event) Serialize() []byte {

	out := make([]byte, 0)

	out = append(out, []byte(
		fmt.Sprintf(
			"[0,\"%s\",%d,%d,",
			s.PubKey,
			s.CreatedAt,
			s.Kind,
		))...)

	// Add encoded tags.
	out = s.Tags.Encode(out)
	out = append(out, ',')

	// Add encoded user content.
	//out = append(out, []byte(s.Content)...)
	out = EscapeString(out, s.Content)
	out = append(out, ']')

	return out
}

func (s *Event) Sign(key string) error {

	bytes, err := hex.DecodeString(key)
	if err != nil {
		log.Fatalf("unable to decode secret: %v", err)
		return fmt.Errorf("Sign called with invalid private key '%s': %w", key, err)
	}

	if s.Tags == nil {
		s.Tags = make(Tags, 0)
	}

	sk, pk := btcec.PrivKeyFromBytes(bytes)
	pkBytes := pk.SerializeCompressed()
	s.PubKey = hex.EncodeToString(pkBytes[1:])

	h := sha256.Sum256(s.Serialize())
	sig, err := schnorr.Sign(sk, h[:])
	if err != nil {
		return err
	}

	s.Id = hex.EncodeToString(h[:])
	s.Sig = hex.EncodeToString(sig.Serialize())

	//log.Printf("\n[\033[32m*\033[0m] Client")
	//log.Printf("  Event signed with SK: %s", key[:10])

	return nil
}

// Escaping strings for JSON encoding according to RFC8259.
// Also encloses result in quotation marks "".
func EscapeString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			// quotation mark
			dst = append(dst, []byte{'\\', '"'}...)
		case c == '\\':
			// reverse solidus
			dst = append(dst, []byte{'\\', '\\'}...)
		case c >= 0x20:
			// default, rest below are control chars
			dst = append(dst, c)
		case c == 0x08:
			dst = append(dst, []byte{'\\', 'b'}...)
		case c < 0x09:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '0', '0' + c}...)
		case c == 0x09:
			dst = append(dst, []byte{'\\', 't'}...)
		case c == 0x0a:
			dst = append(dst, []byte{'\\', 'n'}...)
		case c == 0x0c:
			dst = append(dst, []byte{'\\', 'f'}...)
		case c == 0x0d:
			dst = append(dst, []byte{'\\', 'r'}...)
		case c < 0x10:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '0', 0x57 + c}...)
		case c < 0x1a:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '1', 0x20 + c}...)
		case c < 0x20:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '1', 0x47 + c}...)
		}
	}
	dst = append(dst, '"')
	return dst
}
