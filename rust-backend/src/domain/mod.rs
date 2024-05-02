pub struct Service {
    pub name: String,
    pub address: String,
}

#[derive(Debug, Clone)]
pub enum PacketDirection {
    IN,
    OUT,
}
