package tokenprice_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/omni-network/omni/lib/tokens"
	"github.com/omni-network/omni/monitor/xfeemngr/ticker"
	"github.com/omni-network/omni/monitor/xfeemngr/tokenprice"

	"github.com/stretchr/testify/require"
)

func TestBufferStream(t *testing.T) {
	t.Parallel()

	initial := map[tokens.Token]float64{
		tokens.OMNI: randPrice(),
		tokens.ETH:  randPrice(),
	}

	pricer := tokens.NewMockPricer(initial)

	thresh := 0.1
	tick := ticker.NewMock()
	ctx := context.Background()

	b := tokenprice.NewBuffer(pricer, []tokens.Token{tokens.OMNI, tokens.ETH}, thresh, tick)

	b.Stream(ctx)

	// tick once
	tick.Tick()

	// buffered price should be initial live
	for token, price := range initial {
		require.InEpsilon(t, price, b.Price(token), 0.001, "initial")
	}

	// 10 steps
	buffed := make(map[tokens.Token]float64)
	for i := 0; i < 10; i++ {
		for token := range initial {
			buffed[token] = b.Price(token)
			pricer.SetPrice(token, randPrice())
		}

		tick.Tick()

		live, err := pricer.Price(ctx, tokens.OMNI, tokens.ETH)
		require.NoError(t, err)

		for token, price := range live {
			if inThreshold(price, buffed[token], thresh) {
				require.InEpsilon(t, buffed[token], b.Price(token), 0.001, "should not update")
			} else {
				require.InEpsilon(t, price, b.Price(token), 0.001, "should update")
			}
		}
	}
}

// randPrice generates a random, reasonable token price.
func randPrice() float64 {
	return float64(rand.Intn(5000)) + rand.Float64()
}

// inThreshold returns true if a greater or less than b by pct.
func inThreshold(a, b, pct float64) bool {
	gt := a > b+(b*pct)
	lt := a < b-(b*pct)

	return !gt && !lt
}
