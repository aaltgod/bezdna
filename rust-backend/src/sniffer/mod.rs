use pnet::packet::ip::IpNextHeaderProtocols;
use pnet::packet::tcp::{TcpFlags, TcpPacket};
use pnet::{
    datalink::Channel::Ethernet,
    packet::{
        ethernet::{EtherTypes, EthernetPacket},
        ipv4::Ipv4Packet,
        Packet,
    },
};
use sqlx::{Pool, Postgres};
use std::collections::HashMap;
use std::io::Error;
use std::net::Ipv4Addr;
use std::sync::{Arc, Mutex};
use std::thread;
use tokio::select;
use tokio::sync::mpsc::Receiver;

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
struct ConnectionId {
    source_ip: Ipv4Addr,
    source_port: u16,
    destination_ip: Ipv4Addr,
    destination_port: u16,
}

pub struct Sniffer {
    db: Pool<Postgres>,
    interface_name: String,
    connection_ids: Arc<Mutex<HashMap<ConnectionId, bool>>>,
    ports_to_sniff: Arc<Mutex<HashMap<u16, bool>>>,
}

impl Sniffer {
    pub fn new(db: Pool<Postgres>, interface_name: &str) -> Self {
        let connection_id_map: HashMap<ConnectionId, bool> = HashMap::new();
        let ports_to_sniff_map: HashMap<u16, bool> = HashMap::new();

        Sniffer {
            db,
            interface_name: interface_name.to_string(),
            connection_ids: Arc::new(Mutex::new(connection_id_map)),
            ports_to_sniff: Arc::new(Mutex::new(ports_to_sniff_map)),
        }
    }

    pub async fn run(self, mut ports_to_wath_rx: Receiver<i32>) -> Result<(), Error> {
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

        futures_util::future::join_all(vec![
            tokio::spawn(async move {
                loop {
                    select! {
                        Some(msg) = ports_to_wath_rx.recv() => {
                            tracing::info!("{}",msg);
                        },
                    }
                }
            }),
            tokio::spawn(async move {
                loop {
                    match rx.next() {
                        Ok(packet) => self.handle_packet(EthernetPacket::new(packet).expect("OK")),
                        Err(e) => error!("{}", e),
                    }
                }
            }),
        ])
        .await;

        Ok(())
    }

    fn handle_packet(&self, packet: EthernetPacket) {
        match packet.get_ethertype() {
            EtherTypes::Ipv4 => {
                let ipv4_packet = Ipv4Packet::new(packet.payload()).unwrap();

                match ipv4_packet.get_next_level_protocol() {
                    IpNextHeaderProtocols::Tcp => {
                        let tcp_packet = TcpPacket::new(ipv4_packet.payload()).unwrap();
                        let source_ip = ipv4_packet.get_source();
                        let source_port = tcp_packet.get_source();
                        let destination_ip = ipv4_packet.get_destination();
                        let destination_port = tcp_packet.get_destination();

                        if self
                            .ports_to_sniff
                            .lock()
                            .unwrap()
                            .get(&destination_port)
                            .is_none()
                        {
                            return;
                        }

                        let conn_id = ConnectionId {
                            source_ip,
                            source_port,
                            destination_ip,
                            destination_port,
                        };

                        let mut connections = self.connection_ids.lock().unwrap();

                        // Проверяем, существует ли уже соединение с таким идентификатором
                        if !connections.contains_key(&conn_id) {
                            // Если нет, добавляем новое соединение в HashMap
                            warn!("INSERTED: {:?}", conn_id);
                            connections.insert(conn_id, true);
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

                        warn!("{:?} {}", connections, connections.len());
                    }
                    _ => {}
                }
            }
            _ => {}
        }
    }
}

pub enum Event {
    AddService,
    DeleteService(i32),
}

unsafe impl Send for Event {}
unsafe impl Sync for Event {}
