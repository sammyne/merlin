package merlin_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/sammyne/merlin"
)

func TestTranscript_ChallengeBytes(t *testing.T) {
	testVector := mustReadTestVector(t)

	for i, v := range testVector.MessageChallengePairs {
		label := append([]byte{}, testVector.Label...)
		transcript := merlin.NewTranscript(label)

		messageLabel := append([]byte{}, testVector.MessageLabel...)
		transcript.AppendMessage(messageLabel, v.Message)

		challengeLabel := append([]byte{}, testVector.ChallengeLabel...)

		got := make([]byte, len(v.Challenge))
		if err := transcript.ChallengeBytes(challengeLabel, got); err != nil {
			t.Fatalf("#%d fail to generate challenge: %v", i, err)
		}

		if !bytes.Equal(v.Challenge, got) {
			t.Fatalf("#%d invalid challenge: expect %x, got %x", i, v.Challenge, got)
		}
	}
}

type MessageChallengePair struct {
	Message   []byte `json:"message"`
	Challenge []byte `json:"challenge"`
}

type TestVector struct {
	Label                 []byte                 `json:"label"`
	MessageLabel          []byte                 `json:"message_label"`
	ChallengeLabel        []byte                 `json:"challenge_label"`
	MessageChallengePairs []MessageChallengePair `json:"message_challenge_pairs"`
}

func makeCopy(data []byte) []byte {
	return append([]byte{}, data...)
}

func mustReadTestVector(t *testing.T) TestVector {
	const src = "testdata/transcripts.json"

	data, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatalf("fail to read %s: %v", src, err)
	}

	var testVector TestVector
	if err := json.Unmarshal(data, &testVector); err != nil {
		t.Fatalf("fail to unmarshal TestVector: %v", err)
	}

	return testVector
}
