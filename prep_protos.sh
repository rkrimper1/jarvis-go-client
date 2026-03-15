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
#    Module:         prep_protos.sh
#    Description:    Strips out features of the Jarvis protos that do not work in GO
#    Assumptions:    The script assumes tagger and validator are on 
#                    a line by themselves.  
#                    The script assumes you've cloned the jarvis repo alongside
#                    this repo. That's where it gets the .protos.
#


SOURCE_DIR=../jarvis/proto/
PROTO_DIR=proto
PROTO_REF_DIR=google

# Copy Javis proto files to the proto directory. This will be the basis for the Go client.
cp -a $SOURCE_DIR. $PROTO_DIR

cp -a $PROTO_REF_DIR $PROTO_DIR/$PROTO_REF_DIR
# Move into the proto directory and use the stripper to remove the 
# excess proto file features that do not work in GO.
cd proto && protoc --stripper_out=. --stripper_opt=paths=source_relative **/*.proto

# The google proto references are not copied over, so we need to refresh them.
rm -rf $PROTO_REF_DIR
# The proto files are now ready to be used in the go client. 
# They have been stripped of features that do not work in GO, 
# and the google proto references have been refreshed.