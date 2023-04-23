package sniffer

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Config struct {
	ServiceName string
	Port        int32
}

type Sniffer struct {
	interfaceName string
	errChan       chan error

	data map[int][]string
}

func New(interfaceName string) *Sniffer {
	return &Sniffer{
		interfaceName: interfaceName,
		errChan:       make(chan error),
		data:          make(map[int][]string),
	}
}

func (s *Sniffer) Run() error {
	// Check an interface listening
	_, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	log.Printf("[OK] Check listening interface with name `%s`\n", s.interfaceName)

	return nil
}

func (s *Sniffer) AddConfig(config Config) error {
	if handle, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever); err != nil {
		return err
	} else if err = handle.SetBPFFilter(fmt.Sprintf("tcp and port %d", config.Port)); err != nil {
		return err
	} else {
		log.Printf(
			"START LISTEN service with name `%s` and port `%d`\n",
			config.ServiceName, config.Port)

		go s.process(context.Background(), config, handle)
	}

	return nil
}

func (s *Sniffer) Listen() error {
	return <-s.errChan
}

func (s *Sniffer) process(ctx context.Context, cfg Config, handle *pcap.Handle) {
	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeTCP:
				tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

				payload := string(tcpPacket.Payload)

				// log.WithField("service", cfg.ServiceName).Info(
				// 	tcpPacket.FIN,
				// 	tcpPacket.RST,
				// 	tcpPacket.Ack,
				// 	tcpPacket.Seq,
				// 	tcpPacket.DstPort.String(),
				// )

				log.Println(packet.Metadata().Timestamp)

				// we wan't handle KEEP-ALIVE
				if len(payload) <= 1 {
					continue
				}

				if _, exists := s.data[int(tcpPacket.Ack)]; !exists {
					s.data[int(tcpPacket.Ack)] = []string{payload}

					log.Println(len(s.data))
					continue
				}

				s.data[int(tcpPacket.Ack)] = append(s.data[int(tcpPacket.Ack)], payload)

				log.Println(s.data)

				// log.Println(string(payload))

			}
		}
		// appLayer := packet.ApplicationLayer()
		// if appLayer != nil {
		// 	// fmt.Println(string(appLayer.Payload()))

		// 	log.Infoln(string(appLayer.LayerContents()))
		// }
	}

	return
}
