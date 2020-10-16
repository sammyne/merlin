package merlin_test

import (
	"crypto/rand"
	"testing"

	"github.com/sammyne/merlin"
)

func BenchmarkRand_Read(b *testing.B) {
	const (
		transcriptLabel        = "hello"
		transcriptMessageLabel = "world"
		transcriptMessage      = "nice to meet you"
		witnessLabel           = "nice to meet you, too"
		witnessBody            = "see you later"
	)

	transcript := merlin.NewTranscript([]byte(transcriptLabel))
	transcript.AppendMessage([]byte(transcriptMessageLabel), []byte(transcriptMessage))

	witness := merlin.Witness{Label: []byte(witnessLabel), Body: []byte(witnessBody)}

	r, err := merlin.NewRand(transcript, rand.Reader, witness)
	if err != nil {
		b.Fatalf("failed to new rng: %v", err)
	}

	var dummy int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf [256]byte
		dummy, _ = r.Read(buf[:])
	}

	_ = dummy
}
