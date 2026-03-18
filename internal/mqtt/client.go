package mqtt

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ZhangPengPaul/LightOTA/internal/config"
	"github.com/ZhangPengPaul/LightOTA/internal/model"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Client struct {
	conn *autopaho.ConnectionManager
	cfg  config.MQTTConfig
}

func NewClient(cfg config.MQTTConfig) (*Client, error) {
	brokerURL := fmt.Sprintf("tcp://%s:%d", cfg.Broker, cfg.Port)
	u, err := url.Parse(brokerURL)
	if err != nil {
		return nil, err
	}

	clientCfg := autopaho.ClientConfig{
		BrokerUrls: []*url.URL{u},
	}

	if cfg.Username != "" {
		clientCfg.ConnectUsername = cfg.Username
		clientCfg.ConnectPassword = []byte(cfg.Password)
	}

	clientCfg.ClientID = cfg.ClientID

	ctx := context.Background()
	conn, err := autopaho.NewConnection(ctx, clientCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (c *Client) Disconnect() {
	ctx := context.Background()
	c.conn.Disconnect(ctx)
}

func (c *Client) PublishUpgradeNotification(deviceID string, firmware *model.Firmware) error {
	topic := fmt.Sprintf("ota/%s/upgrade", deviceID)
	payload := fmt.Sprintf(`{"deviceId":"%s","version":"%s","firmwareId":"%s","downloadUrl":"","md5":"%s","fileSize":%d}`,
		deviceID, firmware.Version, firmware.ID, firmware.MD5, firmware.FileSize)

	ctx := context.Background()
	_, err := c.conn.Publish(ctx, &paho.Publish{
		Topic:   topic,
		QoS:     1,
		Payload: []byte(payload),
	})

	return err
}
