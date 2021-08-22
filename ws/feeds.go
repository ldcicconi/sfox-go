package ws

import (
	"fmt"

	"github.com/ldcicconi/sfox-go/models"
)

type OrderbookType string

const (
	OrderbookType_Smart OrderbookType = "orderbook.sfox."
	OrderbookType_NPR   OrderbookType = "orderbook.net."
)

func (c *Client) SubscribeToOrderbook(pair string, orderbookType OrderbookType) (ch <-chan models.Orderbook, err error) {
	feed := fmt.Sprintf("%s%s", string(orderbookType), pair)
	if c.subbedToFeed(feed) {
		x := make(chan models.Orderbook)
		c.registerObSub(feed, x)
		return x, nil
	}

	if err = c.subToFeeds([]string{feed}); err != nil {
		return ch, err
	}

	x := make(chan models.Orderbook)
	c.registerObSub(feed, x)
	return x, nil
}

func (c *Client) registerObSub(feed string, ch chan models.Orderbook) {
	c.obSubsMu.Lock()
	c.obSubs[feed] = append(c.obSubs[feed], ch)
	c.obSubsMu.Unlock()
}

func (c *Client) subbedToFeed(feed string) (subbed bool) {
	c.obSubsMu.RLock()
	subbed = len(c.obSubs[feed]) > 0
	c.obSubsMu.RUnlock()
	return subbed
}

func (c *Client) subToFeeds(feeds []string) (err error) {
	payload := struct {
		Type  string   `json:"type"`
		Feeds []string `json:"feeds"`
	}{
		Type:  "subscribe",
		Feeds: feeds,
	}

	return c.conn.WriteJSON(&payload)
	// TODO: wait for ack, return error if it is there
}
