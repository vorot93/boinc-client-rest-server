package main

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type BOINCConn struct {
	addr            url.URL
	rpcKey          *string
	xmlExpandTags   []string
	xmlCollapseTags []string
	conn            *net.TCPConn
	isClosed        bool
	receiveChan     chan boincPacket
}

func (c *BOINCConn) Send(data BoincRequest) error {
	xmlData, xmlErr := xml.Marshal(data)
	if xmlErr != nil {
		return xmlErr
	}

	return c.SendRaw(xmlData)
}

func (c *BOINCConn) SendRaw(data []byte) error {
	sendPkt := boincPacket{DT: time.Now(), Addr: url.URL{Scheme: c.conn.RemoteAddr().Network(), Host: c.conn.RemoteAddr().String()}, Data: append([]byte(CollapseEmptyTags(string(data), c.xmlCollapseTags)), 3)}
	c.sendData(sendPkt)

	return nil
}

func (c *BOINCConn) sendData(pkt boincPacket) bool {
	if c.isClosed {
		return false
	}

	var sendReq = append(pkt.Data, 3)
	c.conn.Write(sendReq)

	return true
}

func (c *BOINCConn) ReceiveOne() (*BoincReply, error) {
	var pkt = <-c.receiveChan

	var recvData = []byte(strings.Replace(ExpandEmptyTags(string(pkt.Data), c.xmlExpandTags), `encoding="ISO-8859-1"`, `encoding="UTF-8"`, -1))
	var response BoincReply

	xmlErr := xml.Unmarshal(recvData, &response)
	if xmlErr != nil {
		return nil, xmlErr
	}

	if &response == nil {
		return nil, errors.New("Empty response.")
	}

	return &response, nil
}

func (c *BOINCConn) ReceiveAll() []*BoincReply {
	var output []*BoincReply
	for i := 0; i < len(c.receiveChan); i++ {
		v, _ := c.ReceiveOne()
		output = append(output, v)
	}

	return output
}

func (c *BOINCConn) receiveLoop() {
	for {
		if c.isClosed {
			close(c.receiveChan)
			return
		}
		buf := make([]byte, 10000000)
		n, err := c.conn.Read(buf)
		if err == nil {
			var pkt boincPacket
			pkt.Data = buf[0:n]
			pkt.DT = time.Now()
			pkt.Addr = c.addr

			c.receiveChan <- pkt
		}
	}
}

func (c *BOINCConn) makeConnection() error {
	if c.conn != nil {
		c.conn.Close()
	}

	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp", c.addr.String())
	if resolveErr != nil {
		return resolveErr
	}

	tcpConnection, dialErr := net.DialTCP("tcp", nil, tcpAddr)
	if dialErr != nil {
		return dialErr
	}

	c.conn = tcpConnection
	return nil
}

func (c *BOINCConn) doAuth() error {
	rpcKey := c.rpcKey

	sendErr := c.Send(BoincRequest{Auth1: &struct{}{}})
	if sendErr != nil {
		return sendErr
	}

	// Nonce
	responseA, _ := c.ReceiveOne()
	nonce := responseA.Nonce

	if nonce == nil {
		return errors.New("Empty nonce.")
	}

	var rpcKeyV string
	if rpcKey != nil {
		rpcKeyV = *rpcKey
	}
	hash := fmt.Sprintf("%x", md5.Sum([]byte(*nonce+rpcKeyV)))

	sendErr2 := c.Send(BoincRequest{Auth2: &AuthNonceRequest{Hash: hash}})
	if sendErr2 != nil {
		return sendErr2
	}

	responseResult, _ := c.ReceiveOne()

	if responseResult.Unauthorized != nil {
		return errors.New("Authentication error.")
	}

	return nil
}

func (c *BOINCConn) Reset() {
	c.Close()
	c.makeConnection()
	c.doAuth()
}
func (c *BOINCConn) Close() {
	c.conn.Close()
	c.isClosed = true
}

func BOINCConnFinalizer(obj *BOINCConn) {
	obj.Close()
}

func Connect(addr url.URL, rpcKey *string, xmlExpandTags []string, xmlCollapseTags []string) (*BOINCConn, error) {
	bConn := BOINCConn{addr: addr, rpcKey: rpcKey, xmlExpandTags: xmlExpandTags, xmlCollapseTags: xmlCollapseTags, receiveChan: make(chan boincPacket)}
	err := bConn.makeConnection()
	if err != nil {
		bConn.Close()
		return nil, err
	}

	go bConn.receiveLoop()

	if rpcKey != nil {
		authErr := bConn.doAuth()
		if authErr != nil {
			bConn.Close()
			return nil, authErr
		}
	}

	return &bConn, nil
}
