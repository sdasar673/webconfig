/**
* Copyright 2021 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
* SPDX-License-Identifier: Apache-2.0
*/
package http

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-akka/configuration"
	log "github.com/sirupsen/logrus"
	owcommon "github.com/rdkcentral/webconfig/common"
	"github.com/rdkcentral/webconfig/util"
)

const (
	upstreamHostDefault        = "http://localhost:1234"
	upstreamUrlTemplateDefault = "/api/v1/%v/xpcupstream"
)

type UpstreamConnector struct {
	*HttpClient
	host                string
	serviceName         string
	upstreamUrlTemplate string
}

func NewUpstreamConnector(conf *configuration.Config, tlsConfig *tls.Config) *UpstreamConnector {
	serviceName := "upstream"
	confKey := fmt.Sprintf("webconfig.%v.host", serviceName)
	host := conf.GetString(confKey, upstreamHostDefault)
	confKey = fmt.Sprintf("webconfig.%v.url_template", serviceName)
	upstreamUrlTemplate := conf.GetString(confKey, upstreamUrlTemplateDefault)

	return &UpstreamConnector{
		HttpClient:          NewHttpClient(conf, serviceName, tlsConfig),
		host:                host,
		serviceName:         serviceName,
		upstreamUrlTemplate: upstreamUrlTemplate,
	}
}

func (c *UpstreamConnector) UpstreamHost() string {
	return c.host
}

func (c *UpstreamConnector) SetUpstreamHost(host string) {
	c.host = host
}

func (c *UpstreamConnector) ServiceName() string {
	return c.serviceName
}

func (c *UpstreamConnector) PostUpstream(mac string, inHeader http.Header, bbytes []byte, fields log.Fields) ([]byte, http.Header, error) {
	url := c.UpstreamHost() + fmt.Sprintf(c.upstreamUrlTemplate, mac)
	headerMap := util.HeaderToMap(inHeader)
	rbytes, header, err := c.DoWithRetries("POST", url, headerMap, bbytes, fields, c.ServiceName())
	if err != nil {
		return rbytes, header, owcommon.NewError(err)
	}
	return rbytes, header, nil
}
