// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// Amazon EC2 Systems Manager is a collection of capabilities that helps you
// automate management tasks such as collecting system inventory, applying operating
// system (OS) patches, automating the creation of Amazon Machine Images (AMIs),
// and configuring operating systems (OSs) and applications at scale. Systems
// Manager works with managed instances: Amazon EC2 instances and servers or
// virtual machines (VMs) in your on-premises environment that are configured
// for Systems Manager.
//
// This references is intended to be used with the EC2 Systems Manager User
// Guide (Linux (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager.html))
// (Windows (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/systems-manager.html)).
//
// To get started, verify prerequisites and configure managed instances (Linux
// (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/systems-manager-prereqs.html))
// (Windows (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/systems-manager-prereqs.html)).
//The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
type SSM struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "ssm"

// New creates a new instance of the SSM client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a SSM client from just a session.
//     svc := ssm.New(mySession)
//
//     // Create a SSM client with additional configuration
//     svc := ssm.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *SSM {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion string) *SSM {
	svc := &SSM{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2014-11-06",
				JSONVersion:   "1.1",
				TargetPrefix:  "AmazonSSM",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a SSM operation and runs any
// custom request initialization.
func (c *SSM) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
