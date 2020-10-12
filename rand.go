package merlin

import (
	"io"

	"github.com/sammyne/strobe"
)

type Rand struct {
}

type RandBuilder struct {
	strobe *strobe.Strobe
}

// Rekey
// @TODO: consider if Rekey is a must
func (b *RandBuilder) Rekey(label, witness []byte) *RandBuilder {
	panic("not implemented")
}

func (b *RandBuilder) Finalize(r io.Reader) *Rand {
	panic("not implemented")
}

func NewRand(t *Transcript, label, witness []byte, r io.Reader) (*Rand, error) {
	panic("not implemented")
}

func NewRandBuilder(t *Transcript) *RandBuilder {
	return &RandBuilder{strobe: t.strobe.Clone()}
}
