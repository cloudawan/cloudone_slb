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

package restapi

import (
	"encoding/json"
	"github.com/cloudawan/cloudone_slb/control"
	"github.com/cloudawan/cloudone_utility/slb"
	"github.com/emicklei/go-restful"
)

func registerWebServiceSLB() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/slb")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.PUT("/").To(putSLB).
		Doc("Modify the SLB configuration").
		Do(returns200, returns400, returns422, returns500).
		Reads(slb.Command{}))
}

func putSLB(request *restful.Request, response *restful.Response) {
	command := slb.Command{}
	err := request.ReadEntity(&command)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Read body failure"
		jsonMap["ErrorMessage"] = err.Error()
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	err = control.ConfigureSLB(command)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Configure the SLB failure"
		jsonMap["ErrorMessage"] = err.Error()
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(422, string(errorMessageByteSlice))
		return
	}
}
