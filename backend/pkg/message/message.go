package message

import "time"

type Message struct {
	ClientId   int32
	ClientName string
	Time       time.Time
	Content    string
}