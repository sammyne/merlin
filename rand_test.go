package merlin_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/sammyne/merlin"
)

func TestRand_Read(t *testing.T) {
	testVector := mustReadTestVector4Rand(t)

	//out, _ := json.MarshalIndent(testVector, "", " ")
	//fmt.Printf("%s\n", out)

	for i, c := range testVector {
		transcript := merlin.NewTranscript(c.TranscriptLabel)
		transcript.AppendMessage(c.TranscriptMessageLabel, c.TranscriptMessage)

		witnesses := make([]merlin.Witness, len(c.Witnesses))
		for i, w := range c.Witnesses {
			witnesses[i].Label, witnesses[i].Body = append([]byte{}, c.WitnessLabel...), w
		}

		r, err := merlin.NewRand(transcript, bytes.NewReader(c.Rand), witnesses...)
		if err != nil {
			t.Fatalf("#%d failed to new rng: %v", i, err)
		}

		var got [32]byte
		if _, err := r.Read(got[:]); err != nil {
			t.Fatalf("#%d failed to read rand: %v", i, err)
		}
		if !bytes.Equal(c.Expect, got[:]) {
			t.Fatalf("#%d failed: expect %x, got %x", i, c.Expect, got[:])
		}
	}
}

type TestCase struct {
	TranscriptLabel        []byte `json:"transcript_label"`
	TranscriptMessageLabel []byte `json:"transcript_message_label"`
	TranscriptMessage      []byte `json:"transcript_message"`

	WitnessLabel []byte   `json:"witness_label"`
	Witnesses    [][]byte `json:"witnesses"`

	Rand []byte `json:"rand"`

	Expect []byte `json:"expect"`
}

func mustReadTestVector4Rand(t *testing.T) []TestCase {
	const src = "testdata/rand.json"

	data, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatalf("fail to read %q: %v", src, err)
	}

	var testVector []TestCase
	if err := json.Unmarshal(data, &testVector); err != nil {
		t.Fatalf("fail to unmarshal test vector: %v", err)
	}

	return testVector
}
