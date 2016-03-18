package main

import (
	"net/url"
)

type Client struct {
	conn        *BOINCConn
	numAttempts int
}

func (c *Client) Connect(boincAddr url.URL, rpcKey *string) error {
	conn, err := Connect(boincAddr, rpcKey, []string{"bad_request", "authorized", "unauthorized", "have_credentials", "cookie_required"}, []string{})
	if err == nil {
		c.conn = conn
	}
	return err
}

func (c *Client) GetProjects() ([]ProjectInfo, error) {
	var rsp, err = c.SendReceive(BoincRequest{Projects: &struct{}{}})

	var data []ProjectInfo
	if rsp != nil {
		data = rsp.Projects
	}

	return data, err
}

func (c *Client) GetMessages() ([]Message, error) {
	var reqN = 0
	var rsp, err = c.SendReceive(BoincRequest{MessagesFromN: &reqN})

	var data []Message
	if rsp != nil {
		data = rsp.Messages
	}

	return data, err
}

func (c *Client) GetAcctMgrInfo() (*AccountManagerInfo, error) {
	var rsp, err = c.SendReceive(BoincRequest{AccountManagerInfoRequest: &struct{}{}})

	var data *AccountManagerInfo
	if rsp != nil {
		data = rsp.AccountManagerInfo
		if data.HaveCredentialsStruct != nil {
			data.HaveCredentials = true
		}
		if data.CookieRequiredStruct != nil {
			data.CookieRequired = true
		}
	}

	return data, err
}

func (c *Client) PollAcctMgrRPC() (*AccountManagerRPCReply, error) {
	var rsp, err = c.SendReceive(BoincRequest{AccountManagerRPCPoll: &struct{}{}})

	var data *AccountManagerRPCReply
	if rsp != nil {
		data = rsp.AccountManagerRPCReply
	}

	return data, err
}

func (c *Client) SendReceiveSingle(request BoincRequest) (*BoincReply, error) {
	sendErr := c.conn.Send(request)
	if sendErr != nil {
		return nil, sendErr
	}

	response, respErr := c.conn.ReceiveOne()
	if respErr != nil {
		return nil, respErr
	}

	return response, nil
}

func (c *Client) SendReceive(request BoincRequest) (*BoincReply, error) {
	var response *BoincReply
	var err error
	for i := 0; i < c.numAttempts; i++ {
		response, err = c.SendReceiveSingle(request)
		if err == nil {
			break
		}
	}

	return response, err
}

func (c *Client) Disconnect() {
	c.conn.Close()
}

func MakeClient() *Client {
	var c = Client{numAttempts: 3}

	return &c
}
