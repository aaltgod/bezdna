mod sniffer;
use sniffer::sniffer::Sniffer;

#[macro_use]
extern crate log;

fn main() {
    env_logger::init();

    let sniffer = Sniffer::new("lo");

    sniffer.run().unwrap()
}
