use std::fs;

use rand::{CryptoRng, RngCore};
use serde::Serialize;

use testbot::Bytes;

struct BufferedRand {
    buf: Bytes,
}

#[derive(Serialize)]
struct TestCase {
    transcript_label: Bytes,
    transcript_message_label: Bytes,
    transcript_message: Bytes,

    witness_label: Bytes,
    witnesses: Vec<Bytes>,
    rand: Bytes,

    expect: Bytes,
}

impl CryptoRng for BufferedRand {}

impl RngCore for BufferedRand {
    fn next_u32(&mut self) -> u32 {
        let mut out = [0u8; 4];
        self.fill_bytes(&mut out[..]);

        u32::from_be_bytes(out)
    }

    fn next_u64(&mut self) -> u64 {
        let mut out = [0u8; 8];
        self.fill_bytes(&mut out[..]);

        u64::from_be_bytes(out)
    }

    fn fill_bytes(&mut self, dest: &mut [u8]) {
        rand::thread_rng().fill_bytes(dest);

        self.buf.0.extend_from_slice(dest);
    }

    fn try_fill_bytes(&mut self, dest: &mut [u8]) -> Result<(), rand::Error> {
        Ok(self.fill_bytes(dest))
    }
}

fn main() {
    const TRANSCRIPT_LABEL: &'static str = "transcript";
    const TRANSCRIPT_MESSAGE_LABEL: &'static str = "transcript/message";
    const RAND_WITNESS_LABEL: &'static str = "rand/witness";

    let new_case = |transcript_message, witnesses, rand, expect| -> TestCase {
        TestCase {
            transcript_label: Bytes(TRANSCRIPT_LABEL.as_bytes().to_vec()),
            transcript_message_label: Bytes(TRANSCRIPT_MESSAGE_LABEL.as_bytes().to_vec()),
            transcript_message,

            witness_label: Bytes(RAND_WITNESS_LABEL.as_bytes().to_vec()),
            witnesses,
            rand,

            expect: Bytes(expect),
        }
    };

    let ell = (rand::random::<u8>() as usize) + 32;
    let mut test_vector = Vec::with_capacity(ell);
    while test_vector.len() < test_vector.capacity() {
        let message = Bytes::rand();
        let mut witnesses = Vec::with_capacity(rand::random::<u8>() as usize % 16);
        while witnesses.len() < witnesses.capacity() {
            witnesses.push(Bytes::rand());
        }

        let mut rng = BufferedRand {
            buf: Bytes(Vec::new()),
        };

        let mut t = merlin::Transcript::new(TRANSCRIPT_LABEL.as_bytes());
        t.append_message(TRANSCRIPT_MESSAGE_LABEL.as_bytes(), &message.0);

        let mut builder = t.build_rng();
        for w in &witnesses {
            builder = builder.rekey_with_witness_bytes(RAND_WITNESS_LABEL.as_bytes(), &w.0);
        }

        let mut expect = vec![0u8; 32];
        builder.finalize(&mut rng).fill_bytes(expect.as_mut_slice());

        let c = new_case(message, witnesses, rng.buf, expect);
        test_vector.push(c);
    }

    let out = serde_json::to_vec(&test_vector).unwrap();
    fs::write("rand.json", out).unwrap();
}
