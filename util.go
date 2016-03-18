package main

import "strings"

// ParseHostPort splits hostname into host and port parts.
func ParseHostPort(server string) (host string, port string, is6 bool) {
	var Nsemicolons = strings.Count(server, ":")

	if Nsemicolons == 0 {
		host = server
		port = ""
	} else if Nsemicolons == 1 && !strings.ContainsAny(server, "[]") {
		hostArray := strings.Split(server, ":")
		host = hostArray[0]
		port = hostArray[1]
	} else if Nsemicolons >= 2 && Nsemicolons <= 7+1 {
		if strings.ContainsAny(server, "[]") {
			leftIndex := strings.Index(server, "[")
			rightIndex := strings.Index(server, "]")
			if leftIndex < rightIndex && leftIndex != -1 && rightIndex != -1 {
				host = server[leftIndex+1 : rightIndex]

				rightSplit := strings.Split(server, "]")
				if len(rightSplit) == 2 {
					scSplit := strings.Split(rightSplit[1], ":")
					if len(scSplit) == 2 {
						port = scSplit[1]
					}
				}
			}
		} else {
			host = server
			port = ""
		}

		is6 = true
	}

	return host, port, is6
}
