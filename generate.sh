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
#    Module:         generate.sh
#    Description:    Generates the api/pb files from the proto files, these are needed for the go gRPC client
#

API_DIR=api/pb

# set the base api directory, this is where the generated go files will be placed
mkdir $API_DIR

# this generates the buf.lock file that pins the googleapis dependency
# buf generate will be able to resolve google/api/annotations.proto
buf dep update

# use buf to generate the library files for the go client
buf generate

# create the mod.go file 
go mod init github.com/rkrimper1/jarvis
go mod tidy