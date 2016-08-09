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

package configuration

import (
	haproxyRestAPILogger "github.com/cloudawan/cloudone_slb/utility/logger"
	"github.com/cloudawan/cloudone_utility/configuration"
)

var log = haproxyRestAPILogger.GetLogManager().GetLogger("utility")

var configurationContent = `
{
	"certificate": "/etc/cloudone_slb/development_cert.pem",
	"key": "/etc/cloudone_slb/development_key.pem",
	"restapiPort": 8083,
	"haproxyServiceName": "haproxy",
	"haproxyConfigurationFilePath": "/etc/haproxy/haproxy.cfg",
	"haproxyConfigurationTemplateFilePath": "/etc/haproxy/haproxy.cfg.template",
	"haproxyConfigurationCommandFilePath": "/etc/haproxy/haproxy.cfg.command",
	"domainNameSuffix": "cloudawan.com"
}
`

var LocalConfiguration *configuration.Configuration

func init() {
	err := Reload()
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func Reload() error {
	localConfiguration, err := configuration.CreateConfiguration("cloudone_slb", configurationContent)
	if err == nil {
		LocalConfiguration = localConfiguration
	}

	return err
}
