package ws

import (
	"encoding/json"
	"strings"
	"sync"

	gorilla "github.com/gorilla/websocket"
	"github.com/ldcicconi/sfox-go/models"
	"github.com/sirupsen/logrus"
)

const (
	prodUrl = "wss://ws.sfox.com/ws"
	betaUrl = "wss://ws.staging.sfox.com/ws"
)

type Client struct {
	url      string
	conn     *gorilla.Conn
	log      *logrus.Logger
	obSubsMu *sync.RWMutex
	obSubs   map[string][]chan models.Orderbook
}

func NewClient(beta bool) (c *Client, err error) {
	url := prodUrl
	if beta {
		url = betaUrl
	}

	return &Client{
		url:      url,
		log:      logrus.New(),
		obSubsMu: &sync.RWMutex{},
		obSubs:   make(map[string][]chan models.Orderbook),
	}, nil
}

func (c *Client) Start() error {
	conn, _, err := gorilla.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}

	c.log.WithField("url", c.url).Info("connected to SFOX WS server")
	c.conn = conn

	go func() {
		for {
			_, msgBytes, err := c.conn.ReadMessage()
			if err != nil {
				c.log.WithField("error", err.Error()).Error("could not read from WS connection -- exiting")
				return
			}

			var msg message
			err = json.Unmarshal(msgBytes, &msg)
			if err != nil {
				c.log.WithFields(logrus.Fields{
					"error":     err.Error(),
					"msgString": string(msgBytes),
				}).Error("could not parse msg from WS")
				continue
			}

			if msg.Action != "" {
				// switch msg.Type {
				// case "success":
				// case "error":
				// default:
				// }
				c.log.WithFields(logrus.Fields{
					"body": string(msg.Payload),
				}).Info("action message receieved")
				continue
			}

			if strings.HasPrefix(msg.Recipient, "orderbook") {
				c.handleOrderbookMsg(msg.Recipient, msg.Payload)
			} else {
				c.log.WithFields(logrus.Fields{
					"recipient": msg.Recipient,
					"body":      string(msg.Payload),
				}).Warning("unexpected message receieved")
			}
		}
	}()

	return nil
}

func (c *Client) Stop() {
	if err := c.conn.Close(); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("issue closing WS connection")
	}
}

type message struct {
	Type      string          `json:"type"`
	Sequence  int64           `json:"sequence"`
	Timestamp int64           `json:"timestamp"`
	Recipient string          `json:"recipient"`
	Action    string          `json:"action"`
	Payload   json.RawMessage `json:"payload"`
}
