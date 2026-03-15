// Proprietary and Confidential
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
// This source code is the property of:
// 
// Robert Krimper (c) 2026
// 
// https://www.krimper.com
// 
// Author:         Robert Krimper, https://www.linkedin.com/in/robert-krimper
// Modified by:    
// Module:         main.go
// Description:    example gRPC client for Jarvis API with authentication and TLS support
// 
// 

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"os"
	"time"

	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	
	securitypb "github.com/rkrimper1/jarvis/api/pb/security"
	commonpb "github.com/rkrimper1/jarvis/api/pb/common"
	businesspb "github.com/rkrimper1/jarvis/api/pb/business"
)

const (
	FATAL = "FATAL "
	INFO  = "INFO "
	ERROR = "ERROR "
	GRPC_SERVER_ADDRESS = "localhost:50051"
	AUDIENCE = "Tony-Stark"
)

type Client struct {
	conn 		*grpc.ClientConn
	authRes 	*securitypb.AuthenticateResponse
	err 		error
	appEnv 		string
	credFile 	string
	timeout 	time.Duration
	address 	string
	audience 	string
	ctx 		context.Context
}

func main() {
	log.Println("Starting Main")
	client, err := NewClient()
	if err != nil {
		log.Fatalf(FATAL+"Error creating client: %v", err)
	}
	ctx := context.Background()
	client.setConn(ctx)
	if client.err != nil {
		log.Fatalf(FATAL+"Error creating gRPC connection: %v", client.err)
	}
	defer client.conn.Close()

	client.grpcConnLocalAuthCredentials()
	if client.err != nil {
		log.Fatalf(FATAL+"Error fetching auth token: %v", client.err)
	}

	resp := client.grpcCallBusinsessOpsMessage()
	log.Printf(INFO+"Response from BusinessOpsService.SendMessage: %s\n", toJSONString(resp))	

	log.Println("Completing Main")
}

func NewClient() (*Client, error) {
	a := os.Getenv("APP_ENV")
	log.Printf("APP_ENV: %s \n", a)
	c := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	log.Printf("GOOGLE_APPLICATION_CREDENTIALS: %s \n", c)
	t :=  time.Second*20
	log.Printf("Timeout: %s \n", t)

	return &Client{appEnv: a, credFile: c, timeout: t, address: GRPC_SERVER_ADDRESS, audience: AUDIENCE}, nil
}

// setConn - connect to the grpc-server
func (c *Client) setConn(ctx context.Context) {
	var opts []grpc.DialOption

	switch c.appEnv {
	case "local":
		// The gRPC Server is running locally
		c.conn, c.err = grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock())
		c.ctx = ctx
	case "local-tls":
		// The gRPC Server is running locally with TLS
		opts = c.grpcConnLocalTlsCredentials(ctx)
		c.conn, c.err = grpc.Dial(c.address, opts...)
	case "cloud":	
		// The gRPC Server is running in the cloud with TLS
		opts = c.grpcConnCloudTlsCredentials(ctx)
		c.conn, c.err = grpc.Dial(c.address, opts...)
	}

	if c.conn != nil {
		log.Printf(INFO+"Connected to jarvis API server at: %s audience %s\n", c.address, c.audience)
	}
}

// grpcConnLocalAuthCredentials - manage Local TLS connection token with insecure
func (c *Client) grpcConnLocalAuthCredentials() {
	sClient := securitypb.NewSecurityServiceClient(c.conn)
	authReq := &securitypb.AuthenticateRequest{
		Meta: &commonpb.RequestMeta{
			RequestId: "auth-001",
		},
		SubjectId: "tony-stark",
		Method: securitypb.AuthMethod_AUTH_METHOD_TOKEN,
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	c.authRes, c.err = sClient.Authenticate(ctx, authReq)
	if c.err != nil {
		log.Fatalf(FATAL+"Error authenticating: %v", c.err)
		return
	}

	// Append the access token to the context for subsequent gRPC calls
	c.ctx = grpcMetadata.AppendToOutgoingContext(c.ctx, "authorization", "Bearer "+ c.authRes.AccessToken)

	log.Printf(INFO+"Authenticated with token: %s\n", c.authRes.AccessToken)
}

// grpcCallBusinsessOpsMessage - call the business service to send a message
func (c *Client) grpcCallBusinsessOpsMessage() *businesspb.SendMessageResponse {
	bClient := businesspb.NewBusinessOpsServiceClient(c.conn)
	msgReq := &businesspb.SendMessageRequest{
		Meta: &commonpb.RequestMeta{
			RequestId: "msg-001",
		},
		Recipients: []string{"pepper-potts"},
		Channel: businesspb.MessageChannel_MESSAGE_CHANNEL_SECURE,
		Subject: "Urgent: Board meeting rescheduled",
		Body: "Please move the Q4 review to 1400 tomorrow.",
		Encrypt: true,
	}

	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	msgRes, err := bClient.SendMessage(ctx, msgReq)
	if err != nil {
		log.Fatalf(FATAL+"Error sending message: %v", err)
		return nil
	}

	log.Printf(INFO+"Message sent with ID: %s\n", msgRes.MessageId)
	return msgRes
}

// grpcConnLocalCredentials - manage Local TLS connection token with insecure
func (c *Client) grpcConnLocalTlsCredentials(ctx context.Context) ([]grpc.DialOption) {

	tokenSource, err := idtoken.NewTokenSource(ctx, c.audience, idtoken.WithCredentialsFile(c.credFile))
	if err != nil {
		log.Fatalf(FATAL+"idtoken.NewTokenSource: %v", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf(FATAL+"TokenSource.Token: %v", err)
	}

	c.ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	return opts
}

// grpcConnCloudTlsCredentials - manage Cloud TLS connection token and certs
func (c *Client) grpcConnCloudTlsCredentials(ctx context.Context) ([]grpc.DialOption) {

	tokenSource, err := idtoken.NewTokenSource(ctx, c.audience, idtoken.WithCredentialsFile(c.credFile))
	if err != nil {
		log.Fatalf(FATAL+"idtoken.NewTokenSource: %v", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf(FATAL+"TokenSource.Token: %v", err)
	}

	c.ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf(FATAL+"Creating Cert Pool %v", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(creds))
	return opts
}

// toJSONString - creates a printable JSON version of input structure v
func toJSONString(v interface{}) string {
	b, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}