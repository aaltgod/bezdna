use std::collections::HashMap;
use std::io::Error;
use std::ops::Sub;
use std::sync::{Arc, Mutex};
use std::time;

use pnet::{
    datalink::Channel::Ethernet,
    packet::{
        ethernet::{EthernetPacket, EtherTypes},
        ipv4::Ipv4Packet,
        Packet,
    },
};
use pnet::datalink::{bpf, ChannelType, Config, EtherType};
use pnet::packet::ethernet::EtherTypes::Ipv4;
use pnet::packet::ip::IpNextHeaderProtocols;
use pnet::packet::PrimitiveValues;
use pnet::packet::tcp::{TcpFlags, TcpPacket};
use sqlx::{Pool, Postgres};
use tokio::select;
use tokio::sync::mpsc::Receiver;

use crate::domain::PacketDirection;

lazy_static! {
    static ref PORTS_TO_SNIFF: Mutex<HashMap<u16, bool>> = Mutex::new(HashMap::new());
}

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
struct PortPair {
    src: u16,
    dst: u16,
}

#[derive(Debug, Clone)]
struct TcpPacketInfo {
    payload: String,
    packet_direction: PacketDirection,
    completed: bool,
    at: chrono::DateTime<chrono::Utc>,
}

pub struct Sniffer {
    db: Pool<Postgres>,
    interface_name: String,
    tcp_packet_info_by_port_pair: Arc<Mutex<HashMap<PortPair, Vec<TcpPacketInfo>>>>,
    tcp_stream_ttl: chrono::Duration,
    max_stream_ttl: chrono::Duration,
}

impl Sniffer {
    pub fn new(
        db: Pool<Postgres>,
        interface_name: &str,
        tcp_stream_ttl: chrono::Duration,
        max_stream_ttl: chrono::Duration,
    ) -> Self {
        Sniffer {
            db,
            interface_name: interface_name.to_string(),
            tcp_packet_info_by_port_pair: Arc::new(Mutex::new(HashMap::new())),
            tcp_stream_ttl,
            max_stream_ttl,
        }
    }

    pub async fn run(self, mut ports_to_watch_rx: Receiver<u16>) -> Result<(), Error> {
        let arc_self = Arc::new(self);
        let self_handler = Arc::clone(&arc_self);
        let self_manager = Arc::clone(&arc_self);

        let interface = pnet::datalink::interfaces()
            .into_iter()
            .find(|interface| interface.name.eq(self_handler.interface_name.as_str()))
            .unwrap();

        let (_tx, mut rx) = match bpf::channel(
            &interface,
            bpf::Config {
                write_buffer_size: 4096,
                read_buffer_size: 4096,
                read_timeout: None,
                write_timeout: None,
                bpf_fd_attempts: 1000,
            },
        ) {
            Ok(Ethernet(tx, rx)) => (tx, rx),
            Ok(_) => panic!("INVALID type"),
            Err(e) => panic!("{e}"),
        };

        futures_util::future::join_all(vec![
            tokio::spawn(async move {
                loop {
                    select! {
                        Some(msg) = ports_to_watch_rx.recv() => {
                            tracing::info!("{}",msg);

                            let mut m = PORTS_TO_SNIFF.lock().unwrap();
                            m.insert(msg, true);
                        },
                    }
                }
            }),
            tokio::spawn(async move {
                loop {
                    match rx.next() {
                        Ok(packet) => {
                            self_handler.handle_eth_packet(EthernetPacket::new(packet).expect("OK"))
                        }
                        Err(e) => error!("{}", e),
                    }
                }
            }),
            tokio::spawn(async move { self_manager.manage_tcp_streams().await }),
        ])
        .await;

        Ok(())
    }

    async fn manage_tcp_streams(&self) {
        loop {
            tokio::time::sleep(time::Duration::from_secs(5)).await;

            let mut tcp_packet_info_by_port_pair =
                self.tcp_packet_info_by_port_pair.lock().unwrap();

            let mut port_pairs_to_remove: Vec<PortPair> =
                Vec::with_capacity(tcp_packet_info_by_port_pair.len());

            for (port_pair, tcp_packet_info) in tcp_packet_info_by_port_pair.iter() {
                if !tcp_packet_info.is_empty() {
                    let time_now = chrono::Local::now().to_utc();

                    if time_now
                        .sub(tcp_packet_info.last().unwrap().at)
                        .gt(&self.tcp_stream_ttl)
                        || time_now
                            .sub(tcp_packet_info.first().unwrap().at)
                            .gt(&self.max_stream_ttl)
                    {
                        warn!("DONE: {} {:?}", tcp_packet_info.len(), port_pair);

                        port_pairs_to_remove.push(*port_pair);
                    }
                }
            }

            warn!("port_pairs_to_remove LEN: {}", port_pairs_to_remove.len());

            for pair in port_pairs_to_remove {
                tcp_packet_info_by_port_pair.remove(&pair);
            }
        }
    }

    fn handle_eth_packet(&self, packet: EthernetPacket) {
        if packet.get_ethertype() == Ipv4
            // TODO: почему-то на loopback'е не возвращает тип протокола.
            || packet.get_ethertype().to_primitive_values().0 == 0
        {
            let ipv4_packet = Ipv4Packet::new(packet.payload()).unwrap();

            if ipv4_packet.get_next_level_protocol() == IpNextHeaderProtocols::Tcp {
                let tcp_packet = TcpPacket::new(ipv4_packet.payload()).unwrap();

                self.handle_tcp_packet(tcp_packet);
            }
        }
    }

    fn handle_tcp_packet(&self, packet: TcpPacket) {
        let source_port = packet.get_source();
        let destination_port = packet.get_destination();

        let (port_pair, packet_direction): (PortPair, Option<PacketDirection>) = {
            let ports = PORTS_TO_SNIFF.lock().unwrap();

            if ports.contains_key(&destination_port) {
                (
                    PortPair {
                        src: source_port,
                        dst: destination_port,
                    },
                    Some(PacketDirection::IN),
                )
            } else if ports.contains_key(&source_port) {
                (
                    PortPair {
                        src: destination_port,
                        dst: source_port,
                    },
                    Some(PacketDirection::OUT),
                )
            } else {
                return;
            }
        };

        if packet_direction.is_none() {
            return;
        }

        let mut tcp_packet_info_by_port_pair = self.tcp_packet_info_by_port_pair.lock().unwrap();

        let time_now = chrono::Local::now().to_utc();

        let payload = packet.payload();
        let mut info = TcpPacketInfo {
            payload: String::from_utf8(payload.to_vec()).unwrap(),
            packet_direction: packet_direction.unwrap(),
            completed: false,
            at: time_now,
        };

        tcp_packet_info_by_port_pair
            .entry(port_pair)
            .and_modify(|i| i.push(info.clone()))
            .or_insert(vec![info.clone()]);

        warn!(
            "tcp_packet_info_by_port_pair LEN: {}",
            tcp_packet_info_by_port_pair.len()
        );

        if packet.get_flags().eq(&TcpFlags::FIN) || packet.get_flags().eq(&TcpFlags::RST) {
            tcp_packet_info_by_port_pair
                .entry(port_pair)
                .and_modify(|i| {
                    info.completed = true;
                    i.push(info)
                });
        };
    }
}
