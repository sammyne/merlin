package merlin_test

import (
	"fmt"
	"strings"

	"github.com/sammyne/merlin"
)

func ExampleRand_Read() {
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

	rng := strings.NewReader("5c8557cb56ea03216b0e4723b7a6ac6cbc0df94fc947af23031c38ddde4ae672")

	r, err := merlin.NewRand(transcript, rng, witness)
	if err != nil {
		panic(fmt.Sprintf("failed to new rng: %v", err))
	}

	var buf [32]byte
	if _, err := r.Read(buf[:]); err != nil {
		panic(fmt.Sprintf("fail to read random data: %v", err))
	}
	fmt.Printf("%x\n", buf[:])

	// Output:
	// daf9e3c3479ec98a47a1dd63ac9f063b73ac2b3134ceafa0adef8413c56fc4ed
}
