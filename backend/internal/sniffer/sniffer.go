package sniffer

import (
	"fmt"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Config struct {
	ServiceName string
	Port        uint16
}

type Sniffer struct {
	dbRepository db.Repository

	serviceName   string
	port          uint16
	interfaceName string
	errChan       chan error
}

func New(interfaceName string, dbRepository db.Repository) *Sniffer {
	return &Sniffer{
		dbRepository:  dbRepository,
		interfaceName: interfaceName,
		errChan:       make(chan error),
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
	s.serviceName = config.ServiceName
	s.port = uint16(config.Port)

	if handle, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever); err != nil {
		return err
	} else if err = handle.SetBPFFilter(fmt.Sprintf("tcp and port %d", s.port)); err != nil {
		return err
	} else {
		log.Printf(
			"START LISTEN service with name `%s` and port `%d`\n",
			s.serviceName, s.port)

		go s.process(context.Background(), handle)
	}

	return nil
}

func (s *Sniffer) Listen() error {
	return <-s.errChan
}

func (s *Sniffer) process(ctx context.Context, handle *pcap.Handle) {
	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeTCP:
				tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

				payload := string(tcpPacket.Payload)

				// we wan't handle KEEP-ALIVE
				if len(payload) <= 1 {
					continue
				}

				if err := s.dbRepository.InsertStreamByService(
					domain.Stream{
						Ack:       uint64(tcpPacket.Ack),
						Timestamp: packet.Metadata().Timestamp,
						Payload:   payload,
					},
					domain.Service{
						Name: s.serviceName,
						Port: s.port,
					},
				); err != nil {
					log.Error(err)
				}
			}
		}
	}

	return
}
