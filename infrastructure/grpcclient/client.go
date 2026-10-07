package grpcclient

import (
	"context"

	"google.golang.org/grpc"

	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Client struct {
	connection   *grpc.ClientConn
	logPublisher loggerinterface.LogPublisher
	target       string
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.connection
}

func (c *Client) Target() string {
	return c.target
}

func (c *Client) Close(ctx context.Context) {
	if c.connection == nil {
		return
	}

	if err := c.connection.Close(); err != nil {
		c.logPublisher.LogError(ctx, err)
	}
}
