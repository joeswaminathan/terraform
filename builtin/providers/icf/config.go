package icf

import (
	"github.com/joeswaminathan/icf-sdk-go/icf"
)

type Config struct {
	Credentials icf.Credentials
	EndPoint    string
	Protocol    string
	Root        string
}

// Client configures and returns a fully initialized AWSClient
func (c *Config) Client() (client interface{}) {
	sdkc := &icf.Config{
		Credentials: c.Credentials,
		EndPoint:    c.EndPoint,
		Protocol:    c.Protocol,
		Root:        c.Root,
	}
	client = icf.NewClient(sdkc)

	return
}
