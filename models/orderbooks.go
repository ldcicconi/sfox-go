package models

type Orderbook struct {
	Book
	MarketMaking    Book   `json:"market_making"`
	LastUpdatedMs   int64  `json:"lastupdated"`
	LastPublishedMs int64  `json:"lastpublished"`
	Pair            string `json:"pair"`
}

func (ob *Orderbook) Mid() float64 {
	if len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		return 0.0
	}
	return (ob.Bids[0].Price() + ob.Asks[0].Price()) / 2.0
}

type Book struct {
	Bids []Quote `json:"bids"` // sorted best to worst, always
	Asks []Quote `json:"asks"` // <- this too
}

type Quote []interface{}

func (q *Quote) Price() float64    { return (*q)[0].(float64) }
func (q *Quote) Quantity() float64 { return (*q)[1].(float64) }
func (q *Quote) Exchange() string  { return (*q)[2].(string) }
