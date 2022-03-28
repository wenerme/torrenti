package magnet

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	multihash "github.com/multiformats/go-multihash"
)

type Hash struct {
	Name   string // default to sha1
	Digest []byte
}

func (ih Hash) IsZero() bool {
	return len(ih.Digest) == 0
}

func (ih Hash) String() string {
	if ih.Name == "" || ih.Name == "sha1" {
		return urnBtihPrefix + hex.EncodeToString(ih.Digest)
	}
	return urnBtmhPrefix + hex.EncodeToString(ih.Digest)
}

func (ih Hash) Magent() Magnet {
	return Magnet{
		Hash: ih,
	}
}

func (ih Hash) HexHash() string {
	return hex.EncodeToString(ih.Digest)
}

// Magnet link components.
type Magnet struct {
	Hash        Hash       // From "xt"
	Trackers    []string   // From "tr"
	DisplayName string     // From "dn"
	Peers       []string   // From "x.pe"
	Params      url.Values // All other values, such as "as", "xs", etc
}

const (
	urnBtihPrefix = "urn:btih:"
	urnBtmhPrefix = "urn:btmh:"
)

func (m Magnet) String() string {
	vs := make(url.Values, len(m.Params)+len(m.Trackers)+2)
	for k, v := range m.Params {
		vs[k] = append([]string(nil), v...)
	}
	for _, tr := range m.Trackers {
		vs.Add("tr", tr)
	}
	for _, tr := range m.Peers {
		vs.Add("x.pe", tr)
	}
	if m.DisplayName != "" {
		vs.Add("dn", m.DisplayName)
	}

	// Transmission and Deluge both expect "urn:btih:" to be unescaped.
	// Deluge wants it to be at the start of the magnet link.
	// The InfoHash field is expected to be BitTorrent in this implementation.
	u := url.URL{
		Scheme:   "magnet",
		RawQuery: "xt=" + m.Hash.String(),
	}
	if len(vs) != 0 {
		u.RawQuery += "&" + vs.Encode()
	}
	return u.String()
}

// Parse parses Magnet-formatted URIs into a Magnet instance.
func Parse(uri string) (m Magnet, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		err = fmt.Errorf("error parsing uri: %s", err)
		return
	} else if u.Scheme != "magnet" {
		err = fmt.Errorf("unexpected scheme %q", u.Scheme)
		return
	}

	q := u.Query()
	xt := q.Get("xt")
	if m.Hash, err = ParseHash(q.Get("xt")); err != nil {
		err = fmt.Errorf("error parsing infohash %q: %s", xt, err)
		return
	}
	dropFirst(q, "xt")

	m.DisplayName = q.Get("dn")
	dropFirst(q, "dn")

	m.Trackers = q["tr"]
	delete(q, "tr")

	m.Peers = q["x.pe"]
	delete(q, "x.pe")

	if len(q) == 0 {
		q = nil
	}

	m.Params = q
	return
}

func ParseHash(raw string) (ih Hash, err error) {
	var n int
	encoded := raw
	switch {
	case strings.HasPrefix(encoded, urnBtihPrefix):
		encoded = encoded[len(urnBtihPrefix):]
	case strings.HasPrefix(encoded, urnBtmhPrefix):
		encoded = encoded[len(urnBtmhPrefix):]
		var mh *multihash.DecodedMultihash
		ih.Digest, err = hex.DecodeString(encoded)
		if err != nil {
			err = fmt.Errorf("error hex decoding hash: %s", err)
			return
		}

		mh, err = multihash.Decode(ih.Digest)
		if err != nil {
			err = fmt.Errorf("error multihash decoding xt: %s", err)
			return
		}
		ih.Name = mh.Name
		ih.Digest = mh.Digest
		return
	}
	// info hash
	switch len(encoded) {
	case 40:
		ih.Digest = make([]byte, 20)
		n, err = hex.Decode(ih.Digest[:], []byte(encoded))
	case 32:
		ih.Digest = make([]byte, 20)
		n, err = base32.StdEncoding.Decode(ih.Digest[:], []byte(encoded))
	default:
		err = fmt.Errorf("unhandled xt parameter encoding (encoded length %d)", len(encoded))
		return
	}
	if err != nil {
		err = fmt.Errorf("error decoding xt: %s", err)
	} else if n != 20 {
		err = fmt.Errorf("invalid length '%d' of the decoded bytes", n)
	}
	return
}

func dropFirst(vs url.Values, key string) {
	sl := vs[key]
	switch len(sl) {
	case 0, 1:
		vs.Del(key)
	default:
		vs[key] = sl[1:]
	}
}
