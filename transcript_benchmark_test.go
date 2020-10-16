package merlin_test

import (
	"testing"

	"github.com/sammyne/merlin"
)

func BenchmarkTranscript_ChallengeBytes(b *testing.B) {
	const (
		label          = "hello"
		messageLabel   = "world"
		message        = "nice to meet you"
		challengeLabel = "nice to meet you, too"
	)

	transcript := merlin.NewTranscript([]byte(label))
	transcript.AppendMessage([]byte(messageLabel), []byte(message))

	var dummy error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf [8096]byte
		dummy = transcript.ChallengeBytes([]byte(challengeLabel), buf[:])
	}

	_ = dummy
}
