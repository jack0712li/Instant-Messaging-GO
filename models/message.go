package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FormId   string //sender
	TargetId string //receiver
	Type     string //message type group single or broadcast
	Media    int    //message media type text image voice
	Content  string //message content
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}