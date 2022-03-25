package magnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	raw := "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"
	m, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a", m.Hash.String())
	assert.Equal(t, raw, m.String())
	assert.Equal(t, "", m.Hash.Name)
}
