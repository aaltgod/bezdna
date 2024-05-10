use std::fmt::Display;

use regex::bytes;

#[derive(Debug, Clone)]
pub struct Service {
    pub name: String,
    pub port: u16,
    pub flag_regexp: bytes::Regex,
}

#[derive(Debug, Clone)]
pub struct Stream {
    pub id: u64,
    pub service_port: u16,
    pub payload: String,
    pub started_at: chrono::DateTime<chrono::Utc>,
    pub ended_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone)]
pub struct Flag {
    pub id: u64,
    pub text: String,
    pub regexp: bytes::Regex,
    pub stream_id: u64,
    pub direction: FlagDirection,
}

#[derive(Debug, Clone)]
pub enum FlagDirection {
    IN,
    OUT,
}

impl Display for FlagDirection {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let str = match self {
            FlagDirection::IN => "IN".to_string(),
            FlagDirection::OUT => "OUT".to_string(),
        };
        write!(f, "{}", str)
    }
}
