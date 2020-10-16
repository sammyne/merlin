package merlin_test

import (
	"fmt"

	"github.com/sammyne/merlin"
)

func ExampleTranscript_ChallengeBytes() {
	const (
		label          = "hello"
		messageLabel   = "world"
		message        = "nice to meet you"
		challengeLabel = "nice to meet you, too"
	)

	transcript := merlin.NewTranscript([]byte(label))
	transcript.AppendMessage([]byte(messageLabel), []byte(message))

	var buf [32]byte
	if err := transcript.ChallengeBytes([]byte(challengeLabel), buf[:]); err != nil {
		panic(fmt.Sprintf("fail to read challenge bytes: %v", err))
	}

	fmt.Printf("%x\n", buf[:])

	// Output:
	// 92544e2811dec627ae330e93512e12199e45d798b2ae9e2e794548882cc55938
}
