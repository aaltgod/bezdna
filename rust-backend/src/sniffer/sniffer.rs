use pnet::packet::ip::{IpNextHeaderProtocol, IpNextHeaderProtocols};
use pnet::packet::tcp::{TcpFlags, TcpPacket};
use pnet::{
    datalink::{Channel::Ethernet, EtherType},
    packet::{
        ethernet::{EtherTypes, EthernetPacket},
        ipv4::{Ipv4, Ipv4Packet},
        Packet,
    },
};
use std::io::Error;
use std::sync::Arc;

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
struct ConnectionId {
    source_ip: std::net::IpAddr,
    source_port: u16,
    destination_ip: std::net::IpAddr,
    destination_port: u16,
}

pub struct Sniffer {
    interface_name: String,
    connection_ids: Arc<ConnectionId>,
}

impl Sniffer {
    pub fn new(interface_name: &str) -> Self {
        Sniffer {
            interface_name: interface_name.to_string(),
            connection_ids: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    pub fn run(&self) -> Result<(), Error> {
        let interfaces = pnet::datalink::interfaces()
            .into_iter()
            .filter(|interface| interface.name.eq(self.interface_name.as_str()))
            .next()
            .unwrap();

        let (_tx, mut rx) = match pnet::datalink::channel(&interfaces, Default::default()) {
            Ok(Ethernet(tx, rx)) => (tx, rx),
            Ok(_) => panic!("UNHANDLED type"),
            Err(e) => panic!("{e}"),
        };

        loop {
            match rx.next() {
                Ok(packet) => self.handle_packet(EthernetPacket::new(packet).unwrap()),
                Err(e) => panic!("{e}"),
            }
        }
    }

    fn handle_packet(&self, packet: EthernetPacket) {
        match packet.get_ethertype() {
            EtherTypes::Ipv4 => {
                let ipv4_packet = Ipv4Packet::new(packet.payload()).unwrap();

                match ipv4_packet.get_next_level_protocol() {
                    IpNextHeaderProtocols::Tcp => {
                        let tcp_packet = TcpPacket::new(ipv4_packet.payload()).unwrap();
                        let source_ip = packet.get_source();
                        let source_port = tcp_packet.get_source();
                        let destination_ip = packet.get_destination();
                        let destination_port = tcp_packet.get_destination();

                        let conn_id = ConnectionId {
                            source_ip,
                            source_port,
                            destination_ip,
                            destination_port,
                        };

                        let mut connections = connections.lock().unwrap();

                        // Проверяем, существует ли уже соединение с таким идентификатором
                        if !connections.contains_key(&conn_id) {
                            // Если нет, добавляем новое соединение в HashMap
                            connections.insert(conn_id, (tx.clone(), rx.clone()));
                        };

                        let tcp_flag = match tcp_packet.get_flags() {
                            TcpFlags::FIN => {
                                warn!("FIN");
                            }
                            TcpFlags::ACK => {
                                warn!("ACK")
                            }
                            TcpFlags::PSH => {
                                warn!("PSH")
                            }
                            TcpFlags::CWR => {
                                warn!("CWR")
                            }
                            TcpFlags::ECE => {
                                warn!("ECE")
                            }
                            TcpFlags::SYN => {
                                warn!("SYN")
                            }
                            TcpFlags::URG => {
                                warn!("URG")
                            }
                            TcpFlags::RST => {
                                warn!("RST")
                            }
                            _ => {
                                warn!("UNKNOWN")
                            }
                        };

                        warn!(
                            "{:?} {} {} {} {} {}",
                            String::from_utf8(tcp_packet.payload().to_vec())
                                .unwrap_or_default()
                                .to_string(),
                            tcp_packet.get_destination(),
                            tcp_packet.get_source(),
                            tcp_packet.payload().len(),
                            tcp_packet.get_sequence(),
                            tcp_packet.get_flags(),
                        );
                    }
                    _ => {}
                }
            }
            _ => {}
        }
    }
}
