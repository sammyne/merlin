package merlin

import (
	"fmt"
	"io"

	"github.com/sammyne/strobe"
)

type Rand struct {
	strobe *strobe.Strobe
}

type Witness struct {
	Label []byte
	Body  []byte
}

func (r *Rand) Read(buf []byte) (n int, err error) {
	var bLen [4]byte
	byteOrder.PutUint32(bLen[:], uint32(len(buf)))

	_ = r.strobe.AD(bLen[:], &strobe.Options{})

	out, err := r.strobe.PRF(len(buf))
	if err != nil {
		return 0, err
	}

	copy(buf, out)

	return len(buf), nil
}

// rekey
func (r *Rand) rekey(label, witness []byte) {
	var wLen [4]byte
	byteOrder.PutUint32(wLen[:], uint32(len(witness)))

	_ = r.strobe.AD(label, &strobe.Options{Meta: true})
	_ = r.strobe.AD(wLen[:], &strobe.Options{Meta: true, Streaming: true})
	_ = r.strobe.AD(witness, &strobe.Options{})
}

func (r *Rand) finalize(rng io.Reader) error {
	var entropy [32]byte
	if _, err := io.ReadFull(rng, entropy[:]); err != nil {
		return fmt.Errorf("not enough entropy: %w(%v)", ErrLackingEntropy, err)
	}

	const label = "rng"
	_ = r.strobe.AD([]byte(label), &strobe.Options{Meta: true})

	_ = r.strobe.KEY(entropy[:], false)

	return nil
}

// NewRand
// @note: witnesses are optional
func NewRand(t *Transcript, r io.Reader, witnesses ...Witness) (*Rand, error) {
	out := &Rand{strobe: t.strobe.Clone()}

	for _, w := range witnesses {
		out.rekey(w.Label, w.Body)
	}

	if err := out.finalize(r); err != nil {
		return nil, fmt.Errorf("fail to finalize: %w", err)
	}

	return out, nil
}
