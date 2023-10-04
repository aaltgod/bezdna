use pnet::{
    datalink::{Channel::Ethernet, EtherType},
    packet::{
        ethernet::{EtherTypes, EthernetPacket},
        ipv4::{Ipv4, Ipv4Packet},
        Packet,
    },
};
use std::io::Error;

pub struct Sniffer {
    interface_name: String,
}

impl Sniffer {
    pub fn new(interface_name: &str) -> Self {
        Sniffer {
            interface_name: interface_name.to_string(),
        }
    }

    pub fn run(&self) -> Result<(), Error> {
        let interfaces = pnet::datalink::interfaces()
            .into_iter()
            .filter(|interface| interface.name.eq("wg"))
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
        Ok(())
    }

    fn handle_packet(&self, packet: EthernetPacket) {
        match packet.get_ethertype() {
            EtherTypes::Ipv4 => {
                let ipv4_packet = Ipv4Packet::new(packet.payload()).unwrap();

                warn!("{:?}", ipv4_packet);
            }
            _ => warn!(
                "{:?}",
                String::from_utf8(packet.payload().to_vec())
                    .unwrap_or_default()
                    .to_string()
            ),
        }
    }
}
