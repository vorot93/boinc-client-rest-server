package main

import "encoding/xml"

type AuthNonceRequest struct {
	Hash string `xml:"nonce_hash"`
}

type AccountManagerRPCRequest struct {
	UseConfigFile *struct{} `xml:"use_config_file"`
	URL           *string   `xml:"url" json:"url"`
	Name          *string   `xml:"name" json:"name"`
	Password      *string   `xml:"password" json:"password"`
}

type BoincRequest struct {
	XMLName                   xml.Name                  `xml:"boinc_gui_rpc_request"`
	Auth1                     *struct{}                 `xml:"auth1"`
	Auth2                     *AuthNonceRequest         `xml:"auth2"`
	DoQuit                    *struct{}                 `xml:"quit"`
	Projects                  *struct{}                 `xml:"get_all_projects_list"`
	AccountManagerRPCRequest  *AccountManagerRPCRequest `xml:"acct_mgr_rpc_request"`
	AccountManagerRPCPoll     *struct{}                 `xml:"acct_mgr_rpc_poll"`
	AccountManagerInfoRequest *struct{}                 `xml:"acct_mgr_info"`
	MessagesFromN             *int                      `xml:"get_messages>seqno"`
}

type ProjectInfo struct {
	Name         string   `xml:"name" json:"name"`
	Summary      string   `xml:"summary" json:"summary"`
	URL          string   `xml:"url" json:"url"`
	GeneralArea  string   `xml:"general_area" json:"general_area"`
	SpecificArea string   `xml:"specific_area" json:"specific_area"`
	Description  string   `xml:"description" json:"description"`
	Location     string   `xml:"home" json:"home"`
	Platforms    []string `xml:"platforms>name" json:"platforms"`
	Image        string   `xml:"image" json:"image"`
}

type AccountManagerInfo struct {
	URL                   string    `xml:"acct_mgr_url" json:"url"`
	Name                  string    `xml:"acct_mgr_name" json:"name"`
	HaveCredentialsStruct *struct{} `xml:"have_credentials" json:"-"`
	HaveCredentials       bool      `xml:"-" json:"have_credentials"`
	CookieRequiredStruct  *struct{} `xml:"cookie_required" json:"-"`
	CookieRequired        bool      `xml:"-" json:"cookie_required"`
	CookieFailureURL      string    `xml:"cookie_failure_url" json:"cookie_failure_url"`
}

type AccountManagerRPCReply struct {
	Status int `xml:"error_num" json:"status"`
}

type Message struct {
	ProjectName string `xml:"project" json:"project"`
	Priority    int    `xml:"pri" json:"pri"`
	SeqNum      int    `xml:"seqno" json:"seqno"`
	Body        string `xml:"body" json:"body"`
	Timestamp   int64  `xml:"time" json:"time"`
}

type BoincReply struct {
	XMLName                xml.Name                `xml:"boinc_gui_rpc_reply"`
	Error                  *string                 `xml:"error"`
	BadRequest             *struct{}               `xml:"bad_request"`
	Authorized             *struct{}               `xml:"authorized"`
	Unauthorized           *struct{}               `xml:"unauthorized"`
	Nonce                  *string                 `xml:"nonce"`
	Projects               []ProjectInfo           `xml:"projects>project"`
	AccountManagerInfo     *AccountManagerInfo     `xml:"acct_mgr_info"`
	AccountManagerRPCReply *AccountManagerRPCReply `xml:"acct_mgr_rpc_reply"`
	Messages               []Message               `xml:"msgs>msg"`
	Scratchpad             interface{}             `xml:"scratchpad"`
}
