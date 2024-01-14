package domain

import "time"

type Stream struct {
	ID          int64         `json:"id"`
	ServiceName string        `json:"service_name"`
	ServicePort int32         `json:"service_port"`
	Text        *string       `json:"text"`
	FlagRegexp  string        `json:"flag_regexp"`
	FlagAction  FlagDirection `json:"flag_action"`
	Flags       []Flag
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
}

type StreamWithService struct {
	ServiceName string
	ServicePort int32
	Payload     string
}

type Fake struct {
	UserId    int64  `json:"userId"`
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
