package sniffer

import (
	"fmt"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Config struct {
	ServiceName string
	Port        uint16
}

type Sniffer struct {
	dbRepository db.Repository

	interfaceName string
}

func New(interfaceName string, dbRepository db.Repository) *Sniffer {
	return &Sniffer{
		dbRepository:  dbRepository,
		interfaceName: interfaceName,
	}
}

// Run runs handling of services which already exist in db
func (s *Sniffer) Run() error {
	// Check an interface listening
	_, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return errors.Wrap(err, WrapOpenLive)
	}

	services, err := s.dbRepository.GetServices()
	if err != nil {
		return errors.Wrap(err, WrapGetServices)
	}

	for _, service := range services {
		if err = s.AddConfig(Config{
			ServiceName: service.Name,
			Port:        service.Port,
		}); err != nil {
			return errors.Wrap(err, WrapAddConfig)
		}
	}

	return nil
}

func (s *Sniffer) AddConfig(config Config) error {
	if handle, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever); err != nil {
		return errors.Wrap(err, WrapOpenLive)
	} else if err = handle.SetBPFFilter(fmt.Sprintf("tcp and port %d", config.Port)); err != nil {
		return errors.Wrap(err, WrapSetBPFFilter)
	} else {
		log.Infof(
			"START LISTEN service with name `%s` and port `%d`\n",
			config.ServiceName, config.Port)

		go s.process(context.Background(), config, handle)
	}

	return nil
}

func (s *Sniffer) process(ctx context.Context, config Config, handle *pcap.Handle) {
	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeTCP:
				tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
				payload := string(tcpPacket.Payload)

				log.Infoln(config, payload)

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
						Name: config.ServiceName,
						Port: config.Port,
					},
				); err != nil {
					log.Errorln(errors.Wrap(err, WrapInsertStreamByService))
				}
			}
		}
	}
}
