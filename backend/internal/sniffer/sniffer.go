package sniffer

import (
	"regexp"
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
	ServicePort int32
	FlagRegexp  *regexp.Regexp
}

type Sniffer struct {
	dbRepository db.Repository

	interfaceName string

	configMu            *sync.Mutex
	configByServicePort map[int32]Config

	streamMu                *sync.Mutex
	tcpStreamInfoByPortBind map[portBind]tcpStreamInfo

	bufferCh chan tcpStreamInfo
}

type direction string

const (
	in  direction = "IN"
	out direction = "OUT"
)

func (d direction) String() string {
	return string(d)
}

type portBind struct {
	src layers.TCPPort
	dst layers.TCPPort
}

type flag struct {
	text      string
	direction direction
}

type tcpStreamInfo struct {
	text      string
	flags     []flag
	completed bool
	startedAt time.Time
	updatedAt time.Time

	Config
}

func New(interfaceName string, dbRepository db.Repository) *Sniffer {
	return &Sniffer{
		dbRepository:            dbRepository,
		interfaceName:           interfaceName,
		configMu:                &sync.Mutex{},
		configByServicePort:     make(map[int32]Config),
		streamMu:                &sync.Mutex{},
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
		s.configByServicePort[service.Port] = Config{
			ServiceName: service.Name,
			ServicePort: service.Port,
			// trust database data
			FlagRegexp: regexp.MustCompile(service.FlagRegexp),
		}
	}

	handle, err := pcap.OpenLive(s.interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return errors.Wrap(err, WrapOpenLive)
	}

	if err = handle.SetBPFFilter("tcp"); err != nil {
		return errors.Wrap(err, WrapSetBPFFilter)
	}

	go s.process(ctx, handle)
	go s.manageTCPStream(ctx)
	go s.manageBuffer(ctx)

	return nil
}

func (s *Sniffer) AddConfig(config Config) error {
	s.configMu.Lock()
	defer s.configMu.Unlock()

	s.configByServicePort[config.ServicePort] = config

	return nil
}

func (s *Sniffer) checkTCPConn(_ context.Context, config Config, tcpPacket *layers.TCP, direction direction) bool {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()

	var (
		timeNow = time.Now()
		payload = tcpPacket.Payload
		flags   = func() []flag {
			var (
				regFlags = config.FlagRegexp.FindAllString(string(payload), -1)
				flags    = make([]flag, 0, len(regFlags))
			)

			for _, rf := range regFlags {
				flags = append(flags, flag{
					text:      rf,
					direction: direction,
				})
			}

			return flags
		}()
	)

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
			info.text = info.text + string(payload)
			info.flags = append(info.flags, flags...)

			s.tcpStreamInfoByPortBind[portBind{
				src: tcpPacket.DstPort,
				dst: tcpPacket.SrcPort,
			}] = info
		} else {
			s.tcpStreamInfoByPortBind[portBind{
				src: tcpPacket.SrcPort,
				dst: tcpPacket.DstPort,
			}] = tcpStreamInfo{
				text:      string(payload),
				Config:    config,
				startedAt: timeNow,
				updatedAt: timeNow,
				flags:     flags,
			}
		}

		return false
	}

	if tcpPacket.FIN || tcpPacket.RST {
		info.completed = true
	}

	info.updatedAt = timeNow
	info.text = info.text + string(payload)
	info.flags = append(info.flags, flags...)

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
				flagsSlice           = make([][]flag, 0, batchSize)
			)

			if endedStreamAmount != 0 {
				var (
					i int
				)

				endedStreamBuffer[i].FlagRegexp.FindAllString(endedStreamBuffer[i].text, 1)

				for ; i < endedStreamAmount && i < batchSize; i++ {
					endedStreamsToInsert = append(endedStreamsToInsert, domain.Stream{
						ServiceName: endedStreamBuffer[i].ServiceName,
						ServicePort: endedStreamBuffer[i].ServicePort,
						FlagRegexp:  endedStreamBuffer[i].FlagRegexp.String(),
						Text:        &endedStreamBuffer[i].text,
						StartedAt:   endedStreamBuffer[i].startedAt,
						EndedAt:     endedStreamBuffer[i].updatedAt,
					})

					flagsSlice = append(flagsSlice, endedStreamBuffer[i].flags)
				}

				streamIDs, err := s.dbRepository.InsertStreams(endedStreamsToInsert)
				if err != nil {
					log.Errorln(errors.Wrap(err, WrapInsertStreamByService))

					break
				}

				flagsToInsert := make([]domain.Flag, 0, len(flagsSlice))

				for i, flags := range flagsSlice {
					for _, flag := range flags {
						flagsToInsert = append(flagsToInsert, domain.Flag{
							StreamID:  streamIDs[i],
							Text:      flag.text,
							Direction: domain.FlagDirection(flag.direction),
						})
					}
				}

				if err := s.dbRepository.InsertFlags(flagsToInsert); err != nil {
					log.Errorln(errors.Wrap(err, "dbRepository.InsertFlags"))

					break
				}

				endedStreamBuffer = endedStreamBuffer[i:]
			}

			ticker.Reset(timeout)
		}
	}
}

func (s *Sniffer) manageTCPStream(ctx context.Context) {
	ttl := 10 * time.Second
	ticker := time.NewTicker(ttl)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.streamMu.Lock()

			for portBind, tcpStreamInfo := range s.tcpStreamInfoByPortBind {
				timeToCompare := time.Now().Add(-ttl)

				if tcpStreamInfo.completed || tcpStreamInfo.startedAt.Before(timeToCompare) || tcpStreamInfo.startedAt.Equal(timeToCompare) {
					delete(s.tcpStreamInfoByPortBind, portBind)

					if len(tcpStreamInfo.text) != 0 {
						s.bufferCh <- tcpStreamInfo
					}
				}
			}

			s.streamMu.Unlock()

			ticker.Reset(ttl)
		}
	}
}

func (s *Sniffer) process(ctx context.Context, handle *pcap.Handle) {
	for packet := range gopacket.NewPacketSource(handle, handle.LinkType()).Packets() {
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeTCP:
				tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

				src, dst := tcpPacket.SrcPort, tcpPacket.DstPort

				func() {
					s.configMu.Lock()
					defer s.configMu.Unlock()

					var (
						config    Config
						direction direction
					)

					config, exists := s.configByServicePort[int32(src)]
					if exists {
						direction = out
					} else if config, exists = s.configByServicePort[int32(dst)]; exists {
						direction = in
					} else {
						return
					}

					// log.Infoln("CONFIG PAYLOAD ", src, dst, payload)

					// log.Infoln(tcpPacket.SrcPort, tcpPacket.DstPort)

					if !s.checkTCPConn(ctx, config, tcpPacket, direction) {
						return
					}
				}()
			}
		}
	}
}
