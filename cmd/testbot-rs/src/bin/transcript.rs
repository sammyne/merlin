use std::fs;

use rand::{Rng, RngCore};
use serde::{Serialize, Serializer};

struct Bytes(Vec<u8>);

#[derive(Serialize)]
struct LabeledChallenge {
    pub label: Bytes,
    pub body: Bytes,
}

#[derive(Serialize)]
struct LabeledMsg {
    pub label: Bytes,
    pub body: Bytes,
}

#[derive(Serialize)]
struct Transcript {
    pub label: Bytes,
    pub message_label: Bytes,
    pub challenge_label: Bytes,
    pub message_challenge_pairs: Vec<MessageChallengePair>,
}

#[derive(Serialize)]
struct MessageChallengePair {
    pub message: Bytes,
    pub challenge: Bytes,
}

impl Serialize for Bytes {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(&base64::encode(&self.0))
    }
}

impl Bytes {
    pub fn rand() -> Self {
        let mut r = rand::thread_rng();

        let mut data = {
            let ell = (r.gen::<u8>() as usize) + 1;
            vec![0u8; ell]
        };

        r.fill_bytes(&mut data);

        Self(data)
    }
}

fn main() {
    const LABEL: &'static str = "hello";
    const MESSAGE_LABEL: &'static str = "hello/message";
    const CHALLENGE_LABEL: &'static str = "hello/challenge";

    let ell = (rand::random::<u8>() as usize) + 1;
    let mut message_challenge_pairs = Vec::with_capacity(ell);
    while message_challenge_pairs.len() < message_challenge_pairs.capacity() {
        let message = Bytes::rand();

        let mut t = merlin::Transcript::new(LABEL.as_bytes());
        t.append_message(MESSAGE_LABEL.as_bytes(), &message.0);

        let mut challenge = Bytes::rand();

        t.challenge_bytes(CHALLENGE_LABEL.as_bytes(), challenge.0.as_mut_slice());

        let c = MessageChallengePair { message, challenge };

        message_challenge_pairs.push(c);
    }

    {
        // add a empty message
        let mut t = merlin::Transcript::new(LABEL.as_bytes());
        t.append_message(MESSAGE_LABEL.as_bytes(), &[]);

        let mut challenge = Bytes::rand();

        t.challenge_bytes(CHALLENGE_LABEL.as_bytes(), challenge.0.as_mut_slice());

        let c = MessageChallengePair { message: Bytes(vec![]), challenge };

        message_challenge_pairs.push(c);
    }

    let test_vector = Transcript {
        label: Bytes(LABEL.as_bytes().to_vec()),
        message_label: Bytes(MESSAGE_LABEL.as_bytes().to_vec()),
        challenge_label: Bytes(CHALLENGE_LABEL.as_bytes().to_vec()),
        message_challenge_pairs,
    };

    let out = serde_json::to_vec(&test_vector).unwrap();
    fs::write("transcripts.json", out).unwrap();
}
