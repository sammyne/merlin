use rand::{Rng, RngCore};
use serde::{Serialize, Serializer};

pub struct Bytes(pub Vec<u8>);

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
