package merlin

import (
	"encoding/json"

	"github.com/sammyne/strobe"
)

type Transcript struct {
	strobe *strobe.Strobe
}

func (t *Transcript) AppendMessage(label, msg []byte) {
	var msgLen [4]byte
	byteOrder.PutUint32(msgLen[:], uint32(len(msg)))

	_ = t.strobe.AD(label, &strobe.Options{Meta: true})
	_ = t.strobe.AD(msgLen[:], &strobe.Options{Meta: true, Streaming: true})
	_ = t.strobe.AD(msg, &strobe.Options{})
}

func (t *Transcript) ChallengeBytes(label, out []byte) error {
	// TODO: optimize strobe to make the PRF operate on bytes in place
	var outLen [4]byte
	byteOrder.PutUint32(outLen[:], uint32(len(out)))

	_ = t.strobe.AD(label, &strobe.Options{Meta: true})
	_ = t.strobe.AD(outLen[:], &strobe.Options{Meta: true, Streaming: true})
	if err := t.strobe.PRF(out, false); err != nil {
		return err
	}

	return nil
}

// Clone makes a deep copy of t.
func (t *Transcript) Clone() *Transcript {
	return &Transcript{strobe: t.strobe.Clone()}
}

func (t Transcript) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.strobe)
}

func (t *Transcript) UnmarshalJSON(data []byte) error {
	t.strobe = new(strobe.Strobe)
	if err := json.Unmarshal(data, t.strobe); err != nil {
		return err
	}
	return nil
}

func NewTranscript(label []byte) *Transcript {
	strobe, _ := strobe.New(ProtocolLabel, strobe.Bit128)
	transcript := &Transcript{strobe}

	transcript.AppendMessage([]byte("dom-sep"), label)

	return transcript
}
