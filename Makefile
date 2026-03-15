#!/bin/sh
#    Proprietary and Confidential
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
#    This source code is the property of:
#
#     Robert Krimper (c) 2026
#
#     https://www.krimper.com
#
#    Author:         Robert Krimper, https://www.linkedin.com/in/robert-krimper
#    Modified by:    
#    Module:         Makefile
#    Description:    Makefile for the Go client example for Jarvis API
#
#

build: clean prep generate create_client

create_client:
	./create_client.sh

prep:
	./prep_protos.sh

generate:
	./generate.sh

clean:
	rm -rf proto api/* buf.lock go.mod go.sum jarvis-go-client jarvis-client-version_number-go.zip