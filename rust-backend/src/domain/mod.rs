use std::fmt::Display;

use regex::bytes;

#[derive(Debug, Clone)]
pub struct Service {
    pub id: i64,
    pub name: String,
    pub port: i16,
    pub flag_regexp: bytes::Regex,
}

#[derive(Debug, Clone)]
pub struct Stream {
    pub id: i64,
    pub service_port: i16,
}

#[derive(Debug, Clone)]
pub struct Packet {
    pub id: i64,
    pub direction: PacketDirection,
    pub payload: String,
    pub stream_id: i64,
    pub at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone, Copy)]
pub enum PacketDirection {
    IN,
    OUT,
}

impl Display for PacketDirection {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let str = match self {
            PacketDirection::IN => "IN".to_string(),
            PacketDirection::OUT => "OUT".to_string(),
        };
        write!(f, "{}", str)
    }
}