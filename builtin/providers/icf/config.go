package icf

import (
	"cto-github.cisco.com/jswamina/icf-sdk-go/src/icf"
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
		EndPoint: c.EndPoint,
		Protocol: c.Protocol,
		Root: c.Root,
	}
	client = icf.NewClient(sdkc)

	return
}
