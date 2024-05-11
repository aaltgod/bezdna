use std::fmt::Display;

use regex::bytes;

#[derive(Debug, Clone)]
pub struct Service {
    pub id: u64,
    pub name: String,
    pub port: u16,
    pub flag_regexp: bytes::Regex,
}

#[derive(Debug, Clone)]
pub struct Stream {
    pub id: u64,
    pub service_port: u16,
}

#[derive(Debug, Clone)]
pub struct Packet {
    pub id: u64,
    pub direction: PacketDirection,
    pub payload: String,
    pub stream_id: u64,
    pub at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone)]
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
