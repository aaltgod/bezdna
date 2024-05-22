use crate::domain;

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
pub struct PortPair {
    pub src: i16,
    pub dst: i16,
}

#[derive(Debug, Clone)]
pub struct TcpPacketInfo {
    pub payload: String,
    pub packet_direction: domain::PacketDirection,
    pub completed: bool,
    pub at: chrono::DateTime<chrono::Utc>,
}
