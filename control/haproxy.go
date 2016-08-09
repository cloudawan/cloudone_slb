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
	"encoding/json"
	"errors"
	"github.com/cloudawan/cloudone_slb/utility/configuration"
	"github.com/cloudawan/cloudone_utility/logger"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type HAProxy struct {
	serviceName                   string
	configurationFilePath         string
	configurationTemplateFilePath string
	configurationCommandFilePath  string
	command                       Command
	domainNameSuffix              string
}

func CreateHAProxy() (returnedHAProxy *HAProxy, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("CreateHAProxy Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedHAProxy = nil
			returnedError = err.(error)
		}
	}()

	haproxyServiceName, ok := configuration.LocalConfiguration.GetString("haproxyServiceName")
	if ok == false {
		log.Critical("Can't load haproxyServiceName")
		return nil, errors.New("Can't load haproxyServiceName")
	}

	haproxyConfigurationFilePath, ok := configuration.LocalConfiguration.GetString("haproxyConfigurationFilePath")
	if ok == false {
		log.Critical("Can't load haproxyConfigurationFilePath")
		return nil, errors.New("Can't load haproxyConfigurationFilePath")
	}

	haproxyConfigurationTemplateFilePath, ok := configuration.LocalConfiguration.GetString("haproxyConfigurationTemplateFilePath")
	if ok == false {
		log.Critical("Can't load haproxyConfigurationTemplateFilePath")
		return nil, errors.New("Can't load haproxyConfigurationTemplateFilePath")
	}

	haproxyConfigurationCommandFilePath, ok := configuration.LocalConfiguration.GetString("haproxyConfigurationCommandFilePath")
	if ok == false {
		log.Critical("Can't load haproxyConfigurationCommandFilePath")
		return nil, errors.New("Can't load haproxyConfigurationCommandFilePath")
	}

	domainNameSuffix, ok := configuration.LocalConfiguration.GetString("domainNameSuffix")
	if ok == false {
		log.Critical("Can't load domainNameSuffix")
		return nil, errors.New("Can't load domainNameSuffix")
	}

	haproxy := &HAProxy{
		haproxyServiceName,
		haproxyConfigurationFilePath,
		haproxyConfigurationTemplateFilePath,
		haproxyConfigurationCommandFilePath,
		Command{},
		domainNameSuffix,
	}

	err := haproxy.createConfigurationTemplateIfNotExisting()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	/*
		err = haproxy.readCommand()
		if err != nil {
			log.Error(err)
			return nil, err
		}
	*/

	return haproxy, nil
}

func (haproxy *HAProxy) createConfigurationTemplateIfNotExisting() (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("createConfigurationTemplateIfNotExisting Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	if _, err := os.Stat(haproxy.configurationTemplateFilePath); os.IsNotExist(err) {
		// Not exist
		r, err := os.Open(haproxy.configurationFilePath)
		if err != nil {
			log.Error(err)
			return err
		}
		defer r.Close()

		w, err := os.Create(haproxy.configurationTemplateFilePath)
		if err != nil {
			log.Error(err)
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func (haproxy *HAProxy) readConfigurationTemplate() (returnedResult string, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("readConfigurationTemplate Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedResult = ""
			returnedError = err.(error)
		}
	}()

	byteSlice, err := ioutil.ReadFile(haproxy.configurationTemplateFilePath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(byteSlice), nil
}

func (haproxy *HAProxy) writeConfiguration(byteSlice []byte) (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("writeConfiguration Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	err := ioutil.WriteFile(haproxy.configurationFilePath, byteSlice, 0666)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (haproxy *HAProxy) reloadService() (returnedResult string, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("reloadService Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedResult = ""
			returnedError = err.(error)
		}
	}()

	command := exec.Command("service", haproxy.serviceName, "reload")
	outputByte, err := command.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(outputByte))
		return string(outputByte), err
	}

	return string(outputByte), nil
}

func (haproxy *HAProxy) readCommand() (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("readCommand Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	if _, err := os.Stat(haproxy.configurationTemplateFilePath); os.IsNotExist(err) {
		// Not exist
		return nil
	}

	byteSlice, err := ioutil.ReadFile(haproxy.configurationCommandFilePath)
	if err != nil {
		log.Error(err)
		return err
	}

	command := Command{}
	err = json.Unmarshal(byteSlice, &command)
	if err != nil {
		log.Error(err)
		return err
	}

	haproxy.command = command
	return nil
}

func (haproxy *HAProxy) writeCommand() (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("writeCommand Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	byteSlice, err := json.Marshal(haproxy.command)
	if err != nil {
		log.Error(err)
		return err
	}

	err = ioutil.WriteFile(haproxy.configurationCommandFilePath, byteSlice, 0666)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (haproxy *HAProxy) runCommandAndSave() error {
	buffer := bytes.Buffer{}
	// Read template
	templateText, err := haproxy.readConfigurationTemplate()
	if err != nil {
		log.Error("Read configuration template error")
		log.Error(err)
		return err
	}
	buffer.WriteString(templateText)
	// Convert command
	commandText := haproxy.command.convertToHAProxyConfiguration()
	buffer.WriteString(commandText)
	// Wrtie configuration
	err = haproxy.writeConfiguration(buffer.Bytes())
	if err != nil {
		log.Error("Wrtie configuration error")
		log.Error(err)
		return err
	}
	// Reload HAProxy
	result, err := haproxy.reloadService()
	if err != nil {
		log.Error("Reload HAProxy error")
		log.Error(err)
		log.Error(result)
		return err
	}
	// Wrtie Command
	err = haproxy.writeCommand()
	if err != nil {
		log.Error("Wrtie Command error")
		log.Error(err)
		return err
	}

	return nil
}
