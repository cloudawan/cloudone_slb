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
	"errors"
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/cloudawan/cloudone_utility/slb"
	"strconv"
	"time"
)

type Command struct {
	CreatedTime   time.Time
	DomainNameMap map[string]DomainName // Key is FrontEndPort
}

func (command Command) convertToHAProxyConfiguration() string {
	buffer := bytes.Buffer{}
	for _, domainName := range command.DomainNameMap {
		buffer.WriteString(domainName.convertToHAProxyConfiguration())
	}
	return buffer.String()
}

func (command *Command) AddKubernetesServiceHTTP(kubernetesServiceHTTP slb.KubernetesServiceHTTP, nodeHostSlice []string, domainNameSuffix string) (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("AddKubernetesServiceHTTP Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	frontEndPort := kubernetesServiceHTTP.FrontEndPort
	backEndPort := kubernetesServiceHTTP.BackEndPort
	serviceName := kubernetesServiceHTTP.Service
	namespace := kubernetesServiceHTTP.Namespace

	if command.DomainNameMap == nil {
		command.DomainNameMap = make(map[string]DomainName)
	}

	domainName, ok := command.DomainNameMap[strconv.Itoa(frontEndPort)]
	if ok {
		// Exist
		if domainName.FrontEndType != FrontEndTypeHTTP {
			log.Error("FrontEndType is inconsistent to the existing data")
			log.Error("namespace %v serviceName %v frontEndPort %v backEndPort %v nodeHostSlice %v", namespace, serviceName, frontEndPort, backEndPort, nodeHostSlice)
			return errors.New("FrontEndType is inconsistent to the existing data")
		}

		backEnd, ok := domainName.BackEndMap[strconv.Itoa(backEndPort)]
		if ok {
			// Exist
			backEnd.HostSlice = nodeHostSlice
			command.DomainNameMap[strconv.Itoa(frontEndPort)].BackEndMap[strconv.Itoa(backEndPort)] = backEnd
		} else {
			// Not exist
			command.DomainNameMap[strconv.Itoa(frontEndPort)].BackEndMap[strconv.Itoa(backEndPort)] = BackEnd{
				"backend" + strconv.Itoa(backEndPort),
				serviceName + "." + namespace + "." + domainNameSuffix,
				"roundrobin",
				backEndPort,
				nodeHostSlice,
			}
		}
	} else {
		// Not exist
		backEndMap := make(map[string]BackEnd)
		backEndMap[strconv.Itoa(backEndPort)] = BackEnd{
			"backend" + strconv.Itoa(backEndPort),
			serviceName + "." + namespace + "." + domainNameSuffix,
			"roundrobin",
			backEndPort,
			nodeHostSlice,
		}

		command.DomainNameMap[strconv.Itoa(frontEndPort)] = DomainName{
			"frontend" + strconv.Itoa(frontEndPort),
			FrontEndTypeHTTP,
			frontEndPort,
			backEndMap,
		}
	}

	return nil
}

func ConfigureSLB(command slb.Command) error {
	haproxy, err := CreateHAProxy()
	if err != nil {
		log.Error("Create HAProxy error")
		log.Error(err)
		return err
	}

	// Set time
	haproxy.command.CreatedTime = time.Now()

	// Add Kubernetes service http
	for _, kubernetesServiceHTTP := range command.KubernetesServiceHTTPSlice {
		if err := (&haproxy.command).AddKubernetesServiceHTTP(kubernetesServiceHTTP, command.NodeHostSlice, haproxy.domainNameSuffix); err != nil {
			log.Error("Fail to add kuberentes service http to domain name service")
			log.Error("kubernetesServiceHTTP %v nodeHostSlice %v domainNameSuffix %v", kubernetesServiceHTTP, command.NodeHostSlice, haproxy.domainNameSuffix)
			return errors.New("Fail to add kuberentes service http to domain name service")
		}
	}

	// Run
	err = haproxy.runCommandAndSave()
	if err != nil {
		log.Error("Run command and save error")
		log.Error(err)
		return err
	}

	return nil
}
