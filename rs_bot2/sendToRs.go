package rs_bot2

import (
	"context"
	"fmt"

	"github.com/mentalisit/restapi/models"

	"github.com/mentalisit/conf/logger"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client LogicServiceClient
	log    *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	conn, err := grpc.Dial("rs_bot2:50051", grpc.WithInsecure())
	if err != nil {
		log.ErrorErr(err)
		return nil
	}
	fmt.Println("connect to rs grpc ok")
	return &Client{
		conn:   conn,
		client: NewLogicServiceClient(conn),
		log:    log,
	}
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendToRs2(i models.InMessageV2) error {
	in := InMessageV2{
		Text:        i.Text,
		Tip:         i.Tip,
		NameNick:    i.NameNick,
		UserName:    i.Username,
		UserId:      i.UserId,
		NameMention: i.NameMention,
		Messenger:   i.Messenger.ToMap(),
		Config: &CorporationConfigV2{
			Name: i.Config.Uid,
		},
		Options: i.Options,
	}

	_, err := c.client.LogicRs2(context.Background(), &in)
	if err != nil {
		return err
	}
	return nil
}
