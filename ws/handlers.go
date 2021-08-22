package ws

import (
	"encoding/json"

	"github.com/ldcicconi/sfox-go/models"
	"github.com/sirupsen/logrus"
)

func (c *Client) handleOrderbookMsg(feed string, payload json.RawMessage) {
	if !c.subbedToFeed(feed) {
		c.log.WithField("feed", feed).Warning("orderbook received for feed without active subscribers")
		return
	}

	var ob models.Orderbook
	if err := json.Unmarshal(payload, &ob); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"feed":  feed,
		}).Error("could not unmarshal orderbook")
		return
	}

	c.obSubsMu.RLock()
	for _, sub := range c.obSubs[feed] {
		select {
		case sub <- ob:
		default:
		}
	}
	c.obSubsMu.RUnlock()
}
