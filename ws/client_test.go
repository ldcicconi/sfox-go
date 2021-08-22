package ws_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ldcicconi/sfox-go/ws"
)

func TestClient(t *testing.T) {
	c, err := ws.NewClient(false)
	if err != nil {
		t.Fatal("could not create client")
	}
	if err := c.Start(); err != nil {
		t.Fatalf("could not start client: %s", err)
	}
	time.Sleep(time.Second * 1)
	ch, err := c.SubscribeToOrderbook("btcusd", ws.OrderbookType_NPR)
	if err != nil {
		t.Fatalf("could not sub to orderbook: %s", err)
	}
	for ob := range ch {
		fmt.Printf("%+v\n", ob)
	}
}
