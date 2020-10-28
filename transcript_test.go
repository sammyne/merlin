package merlin_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/sammyne/merlin"
)

func TestTranscript_ChallengeBytes(t *testing.T) {
	testVector := mustReadTestVector4Transcript(t)

	//out, _ := json.MarshalIndent(testVector, "", "  ")
	//fmt.Printf("%s\n", out)

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

func TestTranscript_Clone(t *testing.T) {
	t1 := merlin.NewTranscript([]byte("json"))
	t1.AppendMessage([]byte("hello"), []byte("world"))

	expect, err := json.Marshal(t1)
	if err != nil {
		t.Fatalf("fail marshal transcript: %v", err)
	}

	t2 := t1.Clone()
	got, err := json.Marshal(t2)
	if err != nil {
		t.Fatalf("fail marshal the cloned transcript: %v", err)
	}

	if !bytes.Equal(expect, got) {
		t.Fatalf("invalid clone: expect '%s', got '%s'", expect, got)
	}
}

func TestTranscript_MarshalJSON(t *testing.T) {
	t1 := merlin.NewTranscript([]byte("json"))
	t1.AppendMessage([]byte("hello"), []byte("world"))

	expect, err := json.Marshal(t1)
	if err != nil {
		t.Fatalf("fail marshal transcript: %v", err)
	}

	var t2 merlin.Transcript
	if err := json.Unmarshal(expect, &t2); err != nil {
		t.Fatalf("fail unmarshal transcript: %v", err)
	}

	got, err := json.Marshal(t2)
	if err != nil {
		t.Fatalf("fail marshal the recovered transcript: %v", err)
	}

	if !bytes.Equal(expect, got) {
		t.Fatalf("invalid marshal/unmarshaling: expect '%s', got '%s'", expect, got)
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

func mustReadTestVector4Transcript(t *testing.T) TestVector {
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
