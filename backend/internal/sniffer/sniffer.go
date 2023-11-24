package sniffer

import (
	"fmt"
	"sync"
	"time"

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
	Port        int32
}

type Sniffer struct {
	dbRepository db.Repository

	interfaceName string

	mu                      *sync.Mutex
	tcpStreamInfoByPortBind map[portBind]tcpStreamInfo

	bufferCh chan tcpStreamInfo
}

type portBind struct {
	src layers.TCPPort
	dst layers.TCPPort
}

type tcpStreamInfo struct {
	text      string
	completed bool
	startedAt time.Time
	updatedAt time.Time

	Config
}

func New(interfaceName string, dbRepository db.Repository) *Sniffer {
	return &Sniffer{
		dbRepository:            dbRepository,
		interfaceName:           interfaceName,
		mu:                      &sync.Mutex{},
		tcpStreamInfoByPortBind: make(map[portBind]tcpStreamInfo),
		bufferCh:                make(chan tcpStreamInfo),
	}
}

// Run runs handling of services which already exist in db
func (s *Sniffer) Run(ctx context.Context) error {
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

	go s.manageTCPStream(ctx)
	go s.manageBuffer(ctx)

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

func (s *Sniffer) checkTCPConn(_ context.Context, config Config, tcpPacket *layers.TCP) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	timeNow := time.Now()

	info, exist := s.tcpStreamInfoByPortBind[portBind{
		src: tcpPacket.SrcPort,
		dst: tcpPacket.DstPort,
	}]
	if !exist {
		info, exist = s.tcpStreamInfoByPortBind[portBind{
			src: tcpPacket.DstPort,
			dst: tcpPacket.SrcPort,
		}]
		if exist {
			info.updatedAt = timeNow
			info.text = info.text + string(tcpPacket.Payload)

			s.tcpStreamInfoByPortBind[portBind{
				src: tcpPacket.DstPort,
				dst: tcpPacket.SrcPort,
			}] = info
		} else {
			s.tcpStreamInfoByPortBind[portBind{
				src: tcpPacket.SrcPort,
				dst: tcpPacket.DstPort,
			}] = tcpStreamInfo{
				text:      string(tcpPacket.Payload),
				Config:    config,
				startedAt: timeNow,
				updatedAt: timeNow,
			}
		}

		return false
	}

	if tcpPacket.FIN || tcpPacket.RST {
		info.completed = true

		log.Infoln(s.tcpStreamInfoByPortBind)
	}

	info.updatedAt = timeNow
	info.text = info.text + string(tcpPacket.Payload)

	s.tcpStreamInfoByPortBind[portBind{
		src: tcpPacket.SrcPort,
		dst: tcpPacket.DstPort,
	}] = info

	return true
}

func (s *Sniffer) manageBuffer(ctx context.Context) {
	var (
		batchSize         = 100
		timeout           = 10 * time.Second
		ticker            = time.NewTicker(timeout)
		endedStreamBuffer = make([]tcpStreamInfo, 0)
	)

	for {
		select {
		case <-ctx.Done():
			return

		case res := <-s.bufferCh:
			endedStreamBuffer = append(endedStreamBuffer, res)

		case <-ticker.C:
			var (
				endedStreamAmount    = len(endedStreamBuffer)
				endedStreamsToInsert = make([]domain.Stream, 0, batchSize)
			)

			log.Println("BUFFER: ", endedStreamBuffer)

			if endedStreamAmount != 0 {
				var i int

				for ; i < endedStreamAmount && i < batchSize; i++ {
					endedStreamsToInsert = append(endedStreamsToInsert, domain.Stream{
						ServiceName: endedStreamBuffer[i].ServiceName,
						ServicePort: endedStreamBuffer[i].Port,
						Text:        &endedStreamBuffer[i].text,
						StartedAt:   endedStreamBuffer[i].startedAt,
						EndedAt:     endedStreamBuffer[i].updatedAt,
					})
				}

				log.Println("TO INSERT ", endedStreamsToInsert)

				if err := s.dbRepository.InsertStreams(endedStreamsToInsert); err != nil {
					log.Errorln(errors.Wrap(err, WrapInsertStreamByService))

					break
				}

				endedStreamBuffer = endedStreamBuffer[i:]
			}

			ticker.Reset(timeout)
		}
	}
}

func (s *Sniffer) manageTCPStream(ctx context.Context) {
	timeout := 10 * time.Second
	ticker := time.NewTicker(timeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()

			for portBind, tcpStreamInfo := range s.tcpStreamInfoByPortBind {
				timeToCompare := time.Now().Add(-timeout)

				if tcpStreamInfo.completed || tcpStreamInfo.updatedAt.Before(timeToCompare) || tcpStreamInfo.updatedAt.Equal(timeToCompare) {
					delete(s.tcpStreamInfoByPortBind, portBind)

					s.bufferCh <- tcpStreamInfo

					log.Println("DELETE", tcpStreamInfo)
				}
			}

			s.mu.Unlock()

			ticker.Reset(timeout)
		}
	}
}

func (s *Sniffer) process(ctx context.Context, config Config, handle *pcap.Handle) {
	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeTCP:
				tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
				payload := string(tcpPacket.Payload)

				log.Infoln(config, payload)

				log.Infoln(tcpPacket.SrcPort, tcpPacket.DstPort)

				if !s.checkTCPConn(ctx, config, tcpPacket) {
					break
				}
			}
		}
	}
}
