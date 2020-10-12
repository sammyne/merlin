package merlin

import "github.com/sammyne/strobe"

type Transcript struct {
	strobe strobe.Strobe
}

func (t *Transcript) AppendMessage(label, msg []byte) {
	panic("not implemented")
}

func (t *Transcript) ChallengeBytes(label, out []byte) {
	panic("not implemented")
}

func New() *Transcript {
	panic("not implemented")
}
