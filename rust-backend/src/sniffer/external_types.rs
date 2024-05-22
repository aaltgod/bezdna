use std::collections::HashMap;

use regex::bytes;
use tokio::sync::Mutex;

lazy_static! {
    pub static ref PORTS_TO_SNIFF: Mutex<HashMap<i16, bytes::Regex>> = Mutex::new(HashMap::new());
}
