use pnet::{datalink::Channel::Ethernet, packet::ethernet::EthernetPacket};

#[macro_use]
extern crate log;

fn main() {
    let interfaces = pnet::datalink::interfaces()
        .into_iter()
        .filter(|interface| interface.name.eq("wg"))
        .next()
        .unwrap();

    let (tx, mut rx) = match pnet::datalink::channel(&interfaces, Default::default()) {
        Ok(Ethernet(tx, rx)) => (tx, rx),
        Ok(_) => panic!("UNHANDLED type"),
        Err(e) => panic!("{e}"),
    };

    loop {
        match rx.next() {
            Ok(packet) => {
                warn!("{:?}", EthernetPacket::new(packet).unwrap())
            }
            Err(e) => panic!("{e}"),
        }
    }

    return;
}
