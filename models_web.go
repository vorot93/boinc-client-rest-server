package main

import (
	"net/url"
	"time"
)

type jsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

type PostRequest struct {
	BOINCAddr *string `json:"host"`
	Password  *string `json:"pwd"`
}

type boincPacket struct {
	DT   time.Time
	Addr url.URL
	Data []byte
}
