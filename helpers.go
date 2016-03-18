package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
)

func MakeNonce() string {
	randBuf := make([]byte, 32)
	rand.Read(randBuf)
	return base64.URLEncoding.EncodeToString(randBuf)
}

func MakeNonceMD5(nonce, pass string) string {
	hex := fmt.Sprintf("%x", md5.Sum([]byte(nonce+pass)))
	return hex
}

func ExpandEmptyTags(xmlString string, tags []string) string {
	for _, tag := range tags {
		xmlString = strings.Replace(strings.Replace(xmlString, fmt.Sprintf("<%s />", tag), fmt.Sprintf("<%s></%s>", tag, tag), -1), fmt.Sprintf("<%s/>", tag), fmt.Sprintf("<%s></%s>", tag, tag), -1)
	}

	return xmlString
}

func CollapseEmptyTags(xmlString string, tags []string) string {
	for _, tag := range tags {
		xmlString = strings.Replace(xmlString, fmt.Sprintf("<%s></%s>", tag, tag), fmt.Sprintf("<%s/>", tag), -1)
	}

	return xmlString
}

func IsHostAllowed(conn net.Conn, hosts []string) bool {
	var remoteAddr, _, _ = ParseHostPort(conn.RemoteAddr().String())

	var isAllowed bool
	for _, v := range hosts {
		if v == remoteAddr {
			isAllowed = true
			break
		}
	}

	return isAllowed
}
