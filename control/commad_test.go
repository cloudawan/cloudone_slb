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

/*
import (
	"fmt"
	"github.com/cloudawan/cloudone_utility/slb"
	"testing"
	"time"
)

func TestConfigureSLB(t *testing.T) {
	nodeHostSlice := make([]string, 0)
	nodeHostSlice = append(nodeHostSlice, "192.168.0.52")
	nodeHostSlice = append(nodeHostSlice, "192.168.0.53")
	nodeHostSlice = append(nodeHostSlice, "192.168.0.101")

	kubernetesServiceHTTPSlice := make([]slb.KubernetesServiceHTTP, 0)
	kubernetesServiceHTTPSlice = append(
		kubernetesServiceHTTPSlice,
		slb.KubernetesServiceHTTP{
			"default",
			"test",
			80,
			31472,
		})
	kubernetesServiceHTTPSlice = append(
		kubernetesServiceHTTPSlice,
		slb.KubernetesServiceHTTP{
			"default",
			"dev",
			80,
			30476,
		})

	command := slb.Command{
		time.Now(),
		nodeHostSlice,
		kubernetesServiceHTTPSlice,
	}

	fmt.Println(ConfigureSLB(command))
}
*/
