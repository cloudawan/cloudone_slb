// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package control

import (
	"bytes"
	"strconv"
)

const (
	FrontEndTypeHTTP  = "http"
	FrontEndTypeHTTPS = "https"
)

type DomainName struct {
	FrontEndName string
	FrontEndType string
	FrontEndPort int
	BackEndMap   map[string]BackEnd // Key is port
}

type BackEnd struct {
	Name       string
	DomainName string
	Balance    string
	Port       int
	HostSlice  []string
}

func (domainName DomainName) convertToHAProxyConfiguration() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("\n")
	// frontend
	buffer.WriteString("frontend ")
	buffer.WriteString(domainName.FrontEndName)
	buffer.WriteString("\n")
	// mode
	buffer.WriteString("    mode http\n")
	// bind
	buffer.WriteString("    bind 0.0.0.0:")
	buffer.WriteString(strconv.Itoa(domainName.FrontEndPort))
	buffer.WriteString("\n")
	// acl
	for _, backEnd := range domainName.BackEndMap {
		buffer.WriteString("    acl is_")
		buffer.WriteString(backEnd.Name)
		buffer.WriteString(" hdr_dom(host) -i ")
		buffer.WriteString(backEnd.DomainName)
		buffer.WriteString("\n")
	}
	// use_backend
	for _, backEnd := range domainName.BackEndMap {
		buffer.WriteString("    use_backend ")
		buffer.WriteString(backEnd.Name)
		buffer.WriteString(" if is_")
		buffer.WriteString(backEnd.Name)
		buffer.WriteString("\n")
	}
	buffer.WriteString("\n")
	// backend
	for _, backEnd := range domainName.BackEndMap {
		// backend
		buffer.WriteString("backend ")
		buffer.WriteString(backEnd.Name)
		buffer.WriteString("\n")
		// mode
		buffer.WriteString("    mode http\n")
		// balance
		buffer.WriteString("    balance ")
		buffer.WriteString(backEnd.Balance)
		buffer.WriteString("\n")
		// server
		for i, host := range backEnd.HostSlice {
			buffer.WriteString("    server server")
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(" ")
			buffer.WriteString(host)
			buffer.WriteString(":")
			buffer.WriteString(strconv.Itoa(backEnd.Port))
			buffer.WriteString("\n")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("\n")

	return buffer.String()
}

// Example
/*
frontend http
        bind 0.0.0.0:80
        acl is_site1 hdr(host) -i site1.com
        acl is_site2 hdr(host) -i site2.com
        use_backend site1 if is_site1
        use_backend site2 if is_site2

backend site1
        balance roundrobin
        server s1 192.168.0.52:31472
        server s2 192.168.0.53:31472
        server s3 192.168.0.101:31472

backend site2
        balance roundrobin
        server s1 192.168.0.52:30476
        server s2 192.168.0.53:30476
        server s3 192.168.0.101:30476
*/
