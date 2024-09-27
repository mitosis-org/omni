//nolint:gosec // no need for secure randomneness
package admin

import (
	"context"
	"math/rand"
	"sort"

	"github.com/omni-network/omni/e2e/app"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"
	"github.com/omni-network/omni/lib/netconf"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Test tests all admin commands against an ephemeral network.
func Test(ctx context.Context, def app.Definition) error {
	if !def.Testnet.Network.IsEphemeral() {
		return errors.New("only ephemeral networks")
	}

	log.Info(ctx, "Running contract admin tests.")

	network := app.NetworkFromDef(def)

	if err := testEnsurePortalSpec(ctx, def, network); err != nil {
		return err
	}

	if err := testUpgradePortal(ctx, def, network); err != nil {
		return err
	}

	if err := tesUpgradeFeeOracleV1(ctx, def, network); err != nil {
		return err
	}

	if err := testUpgradeGasStation(ctx, def); err != nil {
		return err
	}

	if err := testUpgradeGasPump(ctx, def, network); err != nil {
		return err
	}

	if err := testUpgradeSlashing(ctx, def); err != nil {
		return err
	}

	if err := testUpgradeStaking(ctx, def); err != nil {
		return err
	}

	if err := testUpgradeBridgeNative(ctx, def); err != nil {
		return err
	}

	if err := testUpgradeBridgeL1(ctx, def); err != nil {
		return err
	}

	log.Info(ctx, "Done.")

	return nil
}

// noCheck always returns nil. Use for upgrade actions, where only check is if upgrade succeeds.
func noCheck(context.Context, app.Definition, netconf.Chain) error { return nil }

// testUpgradePortal tests UpgradePortal command.
func testUpgradePortal(ctx context.Context, def app.Definition, network netconf.Network) error {
	err := forOne(ctx, def, randChain(network), UpgradePortal, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade portal")
	}

	err = forAll(ctx, def, network, UpgradePortal, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade all portals")
	}

	return nil
}

func tesUpgradeFeeOracleV1(ctx context.Context, def app.Definition, network netconf.Network) error {
	err := forOne(ctx, def, randChain(network), UpgradeFeeOracleV1, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade feeoracle")
	}

	err = forAll(ctx, def, network, UpgradeFeeOracleV1, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade all feeoracles")
	}

	return nil
}

func testUpgradeGasStation(ctx context.Context, def app.Definition) error {
	err := UpgradeGasStation(ctx, def)
	if err != nil {
		return errors.Wrap(err, "upgrade gas station")
	}

	return nil
}

func testUpgradeGasPump(ctx context.Context, def app.Definition, network netconf.Network) error {
	// cannot UpgradeGasPump on omni evm
	c := randChain(network)
	for {
		if c.Name != omniEVMName {
			break
		}

		c = randChain(network)
	}

	err := forOne(ctx, def, c, UpgradeGasPump, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade gas pump")
	}

	err = forAll(ctx, def, network, UpgradeGasPump, noCheck)
	if err != nil {
		return errors.Wrap(err, "upgrade all gas pumps")
	}

	return nil
}

func testUpgradeSlashing(ctx context.Context, def app.Definition) error {
	err := UpgradeSlashing(ctx, def)
	if err != nil {
		return errors.Wrap(err, "upgrade slashing")
	}

	return nil
}

func testUpgradeStaking(ctx context.Context, def app.Definition) error {
	err := UpgradeStaking(ctx, def)
	if err != nil {
		return errors.Wrap(err, "upgrade staking")
	}

	return nil
}

func testUpgradeBridgeNative(ctx context.Context, def app.Definition) error {
	err := UpgradeBridgeNative(ctx, def)
	if err != nil {
		return errors.Wrap(err, "upgrade bridge native")
	}

	return nil
}

func testUpgradeBridgeL1(ctx context.Context, def app.Definition) error {
	err := UpgradeBridgeL1(ctx, def)
	if err != nil {
		return errors.Wrap(err, "upgrade bridge l1")
	}

	return nil
}

func testEnsurePortalSpec(ctx context.Context, def app.Definition, network netconf.Network) error {
	expected := randPortalSpec(network)

	ensurePortalSpec := func(ctx context.Context, def app.Definition, cfg Config) error {
		return EnsurePortalSpec(ctx, def, cfg, expected)
	}

	err := forOne(ctx, def, randChain(network), ensurePortalSpec, checkPortalSpec(network, expected))
	if err != nil {
		return errors.Wrap(err, "ensure portal spec")
	}

	// new random expected values
	*expected = *randPortalSpec(network)

	err = forAll(ctx, def, network, ensurePortalSpec, checkPortalSpec(network, expected))
	if err != nil {
		return errors.Wrap(err, "ensure all portal specs")
	}

	return nil
}

// forOne runs an action & check configured for a single chain (Config{Chain: "name"}).
func forOne(
	ctx context.Context,
	def app.Definition,
	chain netconf.Chain,
	action func(context.Context, app.Definition, Config) error,
	check func(context.Context, app.Definition, netconf.Chain) error,
) error {
	if err := action(ctx, def, Config{Chain: chain.Name}); err != nil {
		return errors.Wrap(err, "act", "chain", chain.Name)
	}

	if err := check(ctx, def, chain); err != nil {
		return errors.Wrap(err, "check", "chain", chain.Name)
	}

	return nil
}

// forAll runs an action & check configured for all chains (Config{Chain: ""}).
func forAll(
	ctx context.Context,
	def app.Definition,
	network netconf.Network,
	action func(context.Context, app.Definition, Config) error,
	check func(context.Context, app.Definition, netconf.Chain) error,
) error {
	if err := action(ctx, def, Config{}); err != nil {
		return errors.Wrap(err, "act")
	}

	for _, chain := range network.EVMChains() {
		if err := check(ctx, def, chain); err != nil {
			return errors.Wrap(err, "check", "chain", chain.Name)
		}
	}

	return nil
}

func checkPortalSpec(network netconf.Network, expected *PortalSpec) func(context.Context, app.Definition, netconf.Chain) error {
	return func(ctx context.Context, def app.Definition, chain netconf.Chain) error {
		backend, err := def.Backends().Backend(chain.ID)
		if err != nil {
			return errors.Wrap(err, "backend", "chain", chain.Name)
		}

		live, err := livePortalSpec(ctx, network, chain, backend)
		if err != nil {
			return errors.Wrap(err, "live portal spec", "chain", chain.Name)
		}

		// sort chain IDs
		if len(live.PauseXCallTo) != len(expected.PauseXCallTo) {
			return errors.New("live portal spec mismatch", "chain", chain.Name, "live", live, "expected", *expected)
		}

		if len(live.PauseXSubmitFrom) != len(expected.PauseXSubmitFrom) {
			return errors.New("live portal spec mismatch", "chain", chain.Name, "live", live, "expected", *expected)
		}

		// sort chain IDs, for comparison
		sortUint64(live.PauseXCallTo)
		sortUint64(live.PauseXSubmitFrom)
		sortUint64(expected.PauseXCallTo)
		sortUint64(expected.PauseXSubmitFrom)

		if !cmp.Equal(live, *expected, cmpopts.EquateEmpty()) {
			return errors.New("live portal spec mismatch", "chain", chain.Name, "live", live, "expected", *expected)
		}

		return nil
	}
}

func randPortalSpec(network netconf.Network) *PortalSpec {
	pauseAll := randBool()
	if pauseAll {
		return &PortalSpec{PauseAll: true}
	}

	spec := &PortalSpec{
		PauseXCall:   randBool(),
		PauseXSubmit: randBool(),
	}

	if !spec.PauseXCall {
		spec.PauseXCallTo = randChainIDs(network)
	}

	if !spec.PauseXSubmit {
		spec.PauseXSubmitFrom = randChainIDs(network)
	}

	return spec
}

func sortUint64(ns []uint64) {
	sort.Slice(ns, func(i, j int) bool { return ns[i] < ns[j] })
}

func randChain(network netconf.Network) netconf.Chain {
	chains := network.EVMChains()
	return chains[rand.Intn(len(chains))]
}

func randChains(network netconf.Network) []netconf.Chain {
	chains := network.EVMChains()

	n := rand.Intn(len(chains))
	if n == 0 {
		return nil
	}

	rand.Shuffle(len(chains), func(i, j int) {
		chains[i], chains[j] = chains[j], chains[i]
	})

	return chains[:n]
}

func randChainIDs(network netconf.Network) []uint64 {
	chains := randChains(network)

	chainIDs := make([]uint64, len(chains))
	for i, chain := range chains {
		chainIDs[i] = chain.ID
	}

	return chainIDs
}

func randBool() bool {
	return rand.Intn(2) == 0
}
