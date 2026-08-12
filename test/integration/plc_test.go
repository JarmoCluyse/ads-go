//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jarmocluyse/ads-go/pkg/ads"
	adssymbol "github.com/jarmocluyse/ads-go/pkg/ads/ads-symbol"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const plcPort = 852

var (
	loadIntegrationEnvOnce sync.Once
	loadIntegrationEnvErr  error
)

func integrationEnvCandidates() []string {
	seen := map[string]struct{}{}
	candidates := make([]string, 0, 6)

	add := func(path string) {
		if path == "" {
			return
		}
		cleanPath := filepath.Clean(path)
		if _, exists := seen[cleanPath]; exists {
			return
		}
		seen[cleanPath] = struct{}{}
		candidates = append(candidates, cleanPath)
	}

	add(os.Getenv("ADS_TEST_ENV_FILE"))
	add(".env.integration")
	add(".env")

	if _, currentFile, _, ok := runtime.Caller(0); ok {
		integrationDir := filepath.Dir(currentFile)
		repoRoot := filepath.Clean(filepath.Join(integrationDir, "..", ".."))

		add(filepath.Join(integrationDir, ".env.integration"))
		add(filepath.Join(integrationDir, ".env"))
		add(filepath.Join(repoRoot, ".env.integration"))
		add(filepath.Join(repoRoot, ".env"))
	}

	return candidates
}

func loadIntegrationEnv() {
	loadIntegrationEnvOnce.Do(func() {
		for _, candidate := range integrationEnvCandidates() {
			info, err := os.Stat(candidate)
			if err != nil || info.IsDir() {
				continue
			}
			loadIntegrationEnvErr = godotenv.Load(candidate)
			return
		}
	})
}

// newClient creates an ADS client from environment variables and registers cleanup.
// ADS_TARGET_NET_ID is required. ADS_ROUTER_HOST defaults to 127.0.0.1.
// The tests also try to load a .env.integration or .env file from the repo root
// or the test/integration directory before checking the environment.
func newClient(t *testing.T) *ads.Client {
	t.Helper()
	loadIntegrationEnv()
	require.NoError(t, loadIntegrationEnvErr, "load integration env file")

	targetNetID := os.Getenv("ADS_TARGET_NET_ID")
	require.NotEmpty(t, targetNetID, "ADS_TARGET_NET_ID env var must be set (e.g. 192.168.1.5.1.1), or provided via .env.integration/.env or ADS_TEST_ENV_FILE")

	settings := ads.ClientSettings{
		TargetNetID: targetNetID,
		RouterHost:  os.Getenv("ADS_ROUTER_HOST"), // empty → LoadDefaults fills 127.0.0.1
	}

	client := ads.NewClient(settings, nil)
	require.NoError(t, client.Connect(), "connect to ADS router")

	t.Cleanup(func() { _ = client.Disconnect() })
	return client
}

// -----------------------------------------------------------------------
// Expected value computation — pure function, no I/O.
// Mirrors the deterministic formulas in example3/src/Programs/MAIN.TcPOU.
// -----------------------------------------------------------------------

type snapshotExpected struct {
	Seed  uint32
	Bool  bool
	Sint  int8
	Usint uint8
	Byte_ uint8
	Int   int16
	Uint  uint16
	Word  uint16
	Dint  int32
	Udint uint32
	Dword uint32
	Lint  int64
	Ulint uint64
	Lword uint64
	Real  float32
	Lreal float64
	// Date/time raw values (wire encoding: uint32 or uint64)
	Time_ uint32 // ms
	Tod   uint32 // ms since midnight (seed % 86_400_000)
	Date  uint32 // seconds since 1970-01-01 (raw)
	Dt    uint32 // seconds since 1970-01-01 (raw)
	// Ltime  uint64 — not supported in TC 4024
	// Ltod   uint64 — not supported in TC 4024
	// Ldate  uint64 — not supported in TC 4024
	// Ldt    uint64 — not supported in TC 4024
	String string
}

func expectedValues(seed uint32) snapshotExpected {
	return snapshotExpected{
		Seed:   seed,
		Bool:   seed%2 == 0,
		Sint:   int8(seed),
		Usint:  uint8(seed),
		Byte_:  uint8(seed),
		Int:    int16(seed),
		Uint:   uint16(seed),
		Word:   uint16(seed),
		Dint:   int32(seed),
		Udint:  seed,
		Dword:  seed,
		Lint:   int64(seed),
		Ulint:  uint64(seed),
		Lword:  uint64(seed),
		Real:   float32(seed),
		Lreal:  float64(seed),
		Time_:  seed,
		Tod:    seed % 86_400_000,
		Date:   seed,
		Dt:     seed,
		String: fmt.Sprintf("S=%d", seed),
	}
}

// -----------------------------------------------------------------------
// Assertion helpers
// -----------------------------------------------------------------------

func assertSnapshot(t *testing.T, raw any, seed uint32) {
	t.Helper()

	m, ok := raw.(map[string]any)
	require.Truef(t, ok, "snapshot: expected map[string]any, got %T", raw)

	exp := expectedValues(seed)

	assert.Equal(t, exp.Seed, m["nSeed"], "nSeed")
	assert.Equal(t, exp.Bool, m["bBoolVar"], "bBoolVar")
	assert.Equal(t, exp.Sint, m["nSintVar"], "nSintVar")
	assert.Equal(t, exp.Usint, m["nUsintVar"], "nUsintVar")
	assert.Equal(t, exp.Byte_, m["nByteVar"], "nByteVar")
	assert.Equal(t, exp.Int, m["nIntVar"], "nIntVar")
	assert.Equal(t, exp.Uint, m["nUintVar"], "nUintVar")
	assert.Equal(t, exp.Word, m["nWordVar"], "nWordVar")
	assert.Equal(t, exp.Dint, m["nDintVar"], "nDintVar")
	assert.Equal(t, exp.Udint, m["nUdintVar"], "nUdintVar")
	assert.Equal(t, exp.Dword, m["nDwordVar"], "nDwordVar")
	assert.Equal(t, exp.Lint, m["nLintVar"], "nLintVar")
	assert.Equal(t, exp.Ulint, m["nUlintVar"], "nUlintVar")
	assert.Equal(t, exp.Lword, m["nLwordVar"], "nLwordVar")
	// float32 is only exact for seeds ≤ 2^24-1 = 16_777_215
	if seed <= 16_777_215 {
		assert.Equal(t, exp.Real, m["fRealVar"], "fRealVar")
	}
	assert.Equal(t, exp.Lreal, m["fLrealVar"], "fLrealVar")
	assert.Equal(t, exp.Time_, m["tTimeVar"], "tTimeVar")
	assert.Equal(t, exp.Tod, m["tdTimeOfDayVar"], "tdTimeOfDayVar")
	assert.Equal(t, exp.Date, m["dDateVar"], "dDateVar")
	assert.Equal(t, exp.Dt, m["dtDateTimeVar"], "dtDateTimeVar")
	// ltime_var / ltod_var / ldate_var / ldt_var not tested
	assert.Equal(t, exp.String, m["sStringVar"], "sStringVar")
}

func assertIntArray(t *testing.T, raw any, seed uint32) {
	t.Helper()
	name := "aIntArray"
	arr, ok := raw.([]any)
	require.Truef(t, ok, "%s: expected []any, got %T", name, raw)
	require.Len(t, arr, 10, "%s length", name)
	for i, v := range arr {
		assert.Equal(t, int16(seed+uint32(i)), v, "%s	[%d]", name, i)
	}
}

func assertDintArray(t *testing.T, raw any, seed uint32) {
	t.Helper()
	name := "aDintArray"
	arr, ok := raw.([]any)
	require.Truef(t, ok, "%s: expected []any, got %T", name, raw)
	require.Len(t, arr, 10, "%s length", name)
	for i, v := range arr {
		assert.Equal(t, int32(seed+uint32(i)), v, "%s[%d]", name, i)
	}
}

func assert2DArray(t *testing.T, raw any, seed uint32) {
	t.Helper()
	name := "aIntArray2d"
	outer, ok := raw.([]any)
	require.Truef(t, ok, "%s: expected []any, got %T", name, raw)
	require.Len(t, outer, 3, "%s outer dimension", name)
	for i, row := range outer {
		inner, ok := row.([]any)
		require.Truef(t, ok, "%s[%d]: expected []any, got %T", name, i, row)
		require.Len(t, inner, 3, "%s[%d] inner dimension", name, i)
		for j, v := range inner {
			assert.Equal(t, int16(seed+uint32(i)*3+uint32(j)), v, "%s[%d][%d]", name, i, j)
		}
	}
}

// assertLinear reads every scalar field individually from basePath and checks
// against the deterministic expectation for the given seed.
// Works for flat FB vars (base = "Main.fbTypeTest") and struct linear access
// (base = "Main.fbTypeTest.stStructVar" or "Main.fbWriteTest.stStructVar").
func assertLinear(t *testing.T, client *ads.Client, base string, seed uint32) {
	t.Helper()
	exp := expectedValues(seed)

	read := func(field string) any {
		t.Helper()
		got, err := client.ReadValue(plcPort, base+"."+field)
		require.NoError(t, err, "ReadValue %s.%s", base, field)
		return got
	}

	assert.Equal(t, exp.Seed, read("nSeed"), "nSeed")
	assert.Equal(t, exp.Bool, read("bBoolVar"), "bBoolVar")
	assert.Equal(t, exp.Sint, read("nSintVar"), "nSintVar")
	assert.Equal(t, exp.Usint, read("nUsintVar"), "nUsintVar")
	assert.Equal(t, exp.Byte_, read("nByteVar"), "nByteVar")
	assert.Equal(t, exp.Int, read("nIntVar"), "nIntVar")
	assert.Equal(t, exp.Uint, read("nUintVar"), "nUintVar")
	assert.Equal(t, exp.Word, read("nWordVar"), "nWordVar")
	assert.Equal(t, exp.Dint, read("nDintVar"), "nDintVar")
	assert.Equal(t, exp.Udint, read("nUdintVar"), "nUdintVar")
	assert.Equal(t, exp.Dword, read("nDwordVar"), "nDwordVar")
	assert.Equal(t, exp.Lint, read("nLintVar"), "nLintVar")
	assert.Equal(t, exp.Ulint, read("nUlintVar"), "nUlintVar")
	assert.Equal(t, exp.Lword, read("nLwordVar"), "nLwordVar")
	if seed <= 16_777_215 {
		assert.Equal(t, exp.Real, read("fRealVar"), "fRealVar")
	}
	assert.Equal(t, exp.Lreal, read("fLrealVar"), "fLrealVar")
	assert.Equal(t, exp.Time_, read("tTimeVar"), "tTimeVar")
	assert.Equal(t, exp.Tod, read("tdTimeOfDayVar"), "tdTimeOfDayVar")
	assert.Equal(t, exp.Date, read("dDateVar"), "dDateVar")
	assert.Equal(t, exp.Dt, read("dtDateTimeVar"), "dtDateTimeVar")
	// assert.Equal(t, exp.Time_, read("tLTimeVar"), "tLTimeVar")
	// assert.Equal(t, exp.Tod, read("tdLTimeOfDayVar"), "tdLTimeOfDayVar")
	// assert.Equal(t, exp.Date, read("dLDateVar"), "dLDateVar")
	// assert.Equal(t, exp.Dt, read("dtLDateTimeVar"), "dtLDateTimeVar")
	assert.Equal(t, exp.String, read("sStringVar"), "sStringVar")
}

// -----------------------------------------------------------------------
// Static seed test — boundary & overflow cases
// Reads struct_var structurally (whole struct in one call) for all 14 seeds.
// -----------------------------------------------------------------------

var staticSeedCases = []struct {
	name string
	seed uint32
}{
	{"zero", 0},
	{"one", 1},
	{"sint_max", 127},
	{"sint_overflow", 128},
	{"byte_max", 255},
	{"byte_overflow", 256},
	{"int_max", 32_767},
	{"int_sign_flip", 32_768},
	{"uint_max", 65_535},
	{"uint_overflow", 65_536},
	{"real_exact_max", 16_777_215},
	{"dint_max", 2_147_483_647},
	{"dint_sign_boundary", 2_147_483_648},
	{"udint_max", 4_294_967_295},
}

func TestStaticSeed(t *testing.T) {
	client := newClient(t)

	const fb = "Main.fbTypeTest"

	for _, tc := range staticSeedCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Write the driving seed over ADS
			require.NoError(t, client.WriteValue(plcPort, fb+".nSeed", tc.seed), "WriteValue seed")

			// Wait for at least 5 PLC cycles (50 ms @ 10 ms cycle)
			time.Sleep(50 * time.Millisecond)

			// --- struct_var (structural: whole struct in one read) ---
			snapshotRaw, err := client.ReadValue(plcPort, fb+".stStructVar")
			require.NoError(t, err, "ReadValue "+fb+".stStructVar")
			assertSnapshot(t, snapshotRaw, tc.seed)

			// --- 1-D int array ---
			intArrRaw, err := client.ReadValue(plcPort, fb+".aIntArray")
			require.NoError(t, err, "ReadValue "+fb+".aIntArray")
			assertIntArray(t, intArrRaw, tc.seed)

			// --- 1-D dint array ---
			dintArrRaw, err := client.ReadValue(plcPort, fb+".aDintArray")
			require.NoError(t, err, "ReadValue "+fb+".aDintArray")
			assertDintArray(t, dintArrRaw, tc.seed)

			// --- 2-D int array ---
			arr2DRaw, err := client.ReadValue(plcPort, fb+".aIntArray2d")
			require.NoError(t, err, "ReadValue "+fb+".aIntArray2d")
			assert2DArray(t, arr2DRaw, tc.seed)
		})
	}
}

// -----------------------------------------------------------------------
// Read access-path coverage — verifies all three read patterns for one seed:
//   1. Flat FB vars    — Main.fbTypeTest.<var>           (individual scalar reads)
//   2. Struct structural — Main.fbTypeTest.struct_var    (whole struct in one call)
//   3. Struct linear   — Main.fbTypeTest.struct_var.<var> (individual field reads)
// -----------------------------------------------------------------------

func TestReadAllAccessPaths(t *testing.T) {
	client := newClient(t)

	const fb = "Main.fbTypeTest"
	const seed = uint32(100)
	require.NoError(t, client.WriteValue(plcPort, fb+".nSeed", seed), "write seed")
	time.Sleep(50 * time.Millisecond)

	t.Run("flat_fb_vars", func(t *testing.T) {
		assertLinear(t, client, "Main.fbTypeTest", seed)
	})

	t.Run("struct_structural", func(t *testing.T) {
		raw, err := client.ReadValue(plcPort, fb+".stStructVar")
		require.NoError(t, err, "ReadValue "+fb+".stStructVar")
		assertSnapshot(t, raw, seed)
	})

	t.Run("struct_linear", func(t *testing.T) {
		assertLinear(t, client, fb+".stStructVar", seed)
	})
}

// -----------------------------------------------------------------------
// Subscription test — auto-increment mode + on-change notifications
// -----------------------------------------------------------------------

func TestSubscription(t *testing.T) {
	client := newClient(t)
	const fb = "Main.fbTypeTest"
	// Put PLC into a known state: seed=0
	require.NoError(t, client.WriteValue(plcPort, fb+".nSeed", uint32(0)), "reset seed")
	time.Sleep(50 * time.Millisecond)

	// Subscribe to the snapshot struct (on-change, checked every PLC cycle)
	notifCh := make(chan ads.SubscriptionData, 32)
	sub, err := client.SubscribeValue(
		plcPort,
		fb+".stStructVar",
		func(data ads.SubscriptionData) {
			select {
			case notifCh <- data:
			default: // drop if buffer full; test will fail on timeout
			}
		},
		ads.SubscriptionSettings{
			CycleTime:    10 * time.Millisecond,
			SendOnChange: true,
		},
	)
	require.NoError(t, err, "SubscribeValue "+fb+".stStructVar")
	defer func() { _ = client.Unsubscribe(sub) }()

	// Drain the mandatory initial notification sent on subscribe
	select {
	case <-notifCh:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for initial subscription notification")
	}

	// Drive seed directly from Go and wait until the notification reflects the written value.
	// Stale notifications (from before the write) are drained until gotSeed == nextSeed.
	const wantNotifs = 10

	for i := 0; i < wantNotifs; i++ {
		nextSeed := uint32(i + 1)
		require.NoError(t, client.WriteValue(plcPort, fb+".nSeed", nextSeed), "write seed %d", nextSeed)

		deadline := time.After(5 * time.Second)
		for {
			select {
			case data := <-notifCh:
				m, ok := data.Value.(map[string]any)
				require.Truef(t, ok, "notification %d: value should be map[string]any, got %T", i, data.Value)

				gotSeed, ok := m["nSeed"].(uint32)
				require.Truef(t, ok, "notification %d: seed should be uint32", i)

				if gotSeed < nextSeed {
					// stale notification — keep draining
					continue
				}

				// The PLC updates stStructVar.nSeed before the derived fields, so a
				// transitional notification can contain the next seed while the rest
				// of the struct still reflects the previous cycle. The string field is
				// assigned last, so once it matches the target seed the snapshot is
				// stable for full-struct assertions.
				if gotString, ok := m["sStringVar"].(string); !ok || gotString != expectedValues(gotSeed).String {
					continue
				}

				require.Equalf(t, nextSeed, gotSeed, "notification %d: unexpected seed", i)
				assertSnapshot(t, data.Value, gotSeed)

			case <-deadline:
				t.Fatalf("timeout waiting for subscription notification %d/%d", i+1, wantNotifs)
			}
			break
		}
	}
}

// -----------------------------------------------------------------------
// Write coverage — exercises WriteValue for every PLC type via:
//   1. Flat MAIN-level vars  — Main.fbWriteTest.<type>Var  (atomic scalar writes)
//   2. Struct linear fields  — Main.fbWriteTest.stStructVar.<field>
//   3. Struct structural     — Main.fbWriteTest.stStructVar (whole struct at once) [SKIP: not yet implemented]
//   4. Arrays                — Main.fbWriteTest.aIntArray / aDintArray / aIntArray2d
// -----------------------------------------------------------------------

type writeCase struct {
	path  string
	write any
	check func(t *testing.T, got any)
}

func eqCheck(want any) func(t *testing.T, got any) {
	return func(t *testing.T, got any) { t.Helper(); assert.Equal(t, want, got) }
}

func deltaCheck(want, delta float64) func(t *testing.T, got any) {
	return func(t *testing.T, got any) { t.Helper(); assert.InDelta(t, want, got, delta) }
}

func runWriteCases(t *testing.T, client *ads.Client, cases []writeCase) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			require.NoError(t, client.WriteValue(plcPort, tc.path, tc.write), "WriteValue %s", tc.path)
			time.Sleep(20 * time.Millisecond)
			got, err := client.ReadValue(plcPort, tc.path)
			require.NoError(t, err, "ReadValue %s", tc.path)
			tc.check(t, got)
		})
	}
}

func TestWriteAllTypes(t *testing.T) {
	client := newClient(t)

	const fb = "Main.fbWriteTest"

	// ------------------------------------------------------------------
	// 1. Flat MAIN-level vars — PLC never touches these
	// ------------------------------------------------------------------
	t.Run("flat_main_vars", func(t *testing.T) {
		runWriteCases(t, client, []writeCase{
			{fb + ".bBoolVar", true, eqCheck(true)},
			{fb + ".bBoolVar", false, eqCheck(false)},
			{fb + ".nSintVar", int8(-99), eqCheck(int8(-99))},
			{fb + ".nUsintVar", uint8(200), eqCheck(uint8(200))},
			{fb + ".nByteVar", uint8(0xFF), eqCheck(uint8(0xFF))},
			{fb + ".nIntVar", int16(-32000), eqCheck(int16(-32000))},
			{fb + ".nUintVar", uint16(65000), eqCheck(uint16(65000))},
			{fb + ".nWordVar", uint16(0xABCD), eqCheck(uint16(0xABCD))},
			{fb + ".nDintVar", int32(-2_000_000), eqCheck(int32(-2_000_000))},
			{fb + ".nUdintVar", uint32(3_000_000), eqCheck(uint32(3_000_000))},
			{fb + ".nDwordVar", uint32(0xDEADBEEF), eqCheck(uint32(0xDEADBEEF))},
			{fb + ".nLintVar", int64(-9_000_000_000), eqCheck(int64(-9_000_000_000))},
			{fb + ".nUlintVar", uint64(18_000_000_000), eqCheck(uint64(18_000_000_000))},
			{fb + ".nLwordVar", uint64(0xCAFEBABEDEADBEEF), eqCheck(uint64(0xCAFEBABEDEADBEEF))},
			{fb + ".fRealVar", float32(3.14), deltaCheck(float64(float32(3.14)), 1e-4)},
			{fb + ".fLrealVar", float64(2.718281828), deltaCheck(2.718281828, 1e-9)},
			{fb + ".tTimeVar", uint32(5000), eqCheck(uint32(5000))},
			{fb + ".tdTimeOfDayVar", uint32(3_600_000), eqCheck(uint32(3_600_000))},
			{fb + ".dDateVar", uint32(1_000_000), eqCheck(uint32(1_000_000))},
			{fb + ".dtDateTimeVar", uint32(1_700_000_000), eqCheck(uint32(1_700_000_000))},
			// {fb + ".tLTimeVar", uint64(5000), eqCheck(uint64(5000))},
			// {fb + ".tdLTimeOfDayVar", uint64(3_600_000), eqCheck(uint64(3_600_000))},
			// {fb + ".dLDateVar", uint64(1_000_000), eqCheck(uint64(1_000_000))},
			// {fb + ".dtLDateTimeVar", uint64(1_700_000_000), eqCheck(uint64(1_700_000_000))},
			{fb + ".sStringVar", "hello ADS", eqCheck("hello ADS")},
		})
	})

	// ------------------------------------------------------------------
	// 2. Struct fields — linear access, one field at a time
	// ------------------------------------------------------------------
	t.Run("struct_linear", func(t *testing.T) {
		const s = fb + ".stStructVar"
		runWriteCases(t, client, []writeCase{
			{s + ".nSeed", uint32(77), eqCheck(uint32(77))},
			{s + ".bBoolVar", true, eqCheck(true)},
			{s + ".bBoolVar", false, eqCheck(false)},
			{s + ".nSintVar", int8(-50), eqCheck(int8(-50))},
			{s + ".nUsintVar", uint8(150), eqCheck(uint8(150))},
			{s + ".nByteVar", uint8(0xAB), eqCheck(uint8(0xAB))},
			{s + ".nIntVar", int16(-1000), eqCheck(int16(-1000))},
			{s + ".nUintVar", uint16(1000), eqCheck(uint16(1000))},
			{s + ".nWordVar", uint16(0x1234), eqCheck(uint16(0x1234))},
			{s + ".nDintVar", int32(-100_000), eqCheck(int32(-100_000))},
			{s + ".nUdintVar", uint32(100_000), eqCheck(uint32(100_000))},
			{s + ".nDwordVar", uint32(0x12345678), eqCheck(uint32(0x12345678))},
			{s + ".nLintVar", int64(-1_000_000_000), eqCheck(int64(-1_000_000_000))},
			{s + ".nUlintVar", uint64(1_000_000_000), eqCheck(uint64(1_000_000_000))},
			{s + ".nLwordVar", uint64(0xDEADBEEF12345678), eqCheck(uint64(0xDEADBEEF12345678))},
			{s + ".fRealVar", float32(1.5), deltaCheck(float64(float32(1.5)), 1e-4)},
			{s + ".fLrealVar", float64(3.141592653589793), deltaCheck(3.141592653589793, 1e-9)},
			{s + ".tTimeVar", uint32(2000), eqCheck(uint32(2000))},
			{s + ".tdTimeOfDayVar", uint32(7_200_000), eqCheck(uint32(7_200_000))},
			{s + ".dDateVar", uint32(500_000), eqCheck(uint32(500_000))},
			{s + ".dtDateTimeVar", uint32(1_600_000_000), eqCheck(uint32(1_600_000_000))},
			// {m + ".tLTimeVar", uint64(2000), eqCheck(uint64(2000))},
			// {m + ".tdLTimeOfDayVar", uint64(7_200_000), eqCheck(uint64(7_200_000))},
			// {m + ".dLDateVar", uint64(500_000), eqCheck(uint64(500_000))},
			// {m + ".dtLDateTimeVar", uint64(1_600_000_000), eqCheck(uint64(1_600_000_000))},
			{s + ".sStringVar", "write test", eqCheck("write test")},
		})
	})

	// ------------------------------------------------------------------
	// 3. Struct structural write — write the whole struct in one call
	//    Serialize already handles map[string]any for structs.
	// ------------------------------------------------------------------
	t.Run("struct_structural", func(t *testing.T) {
		const seed = uint32(77)

		exp := expectedValues(seed)
		m := map[string]any{
			"nSeed":             exp.Seed,
			"bAutoMode":         false,
			"nAutoTickInterval": uint32(50),
			"bBoolVar":          exp.Bool,
			"nSintVar":          exp.Sint,
			"nUsintVar":         exp.Usint,
			"nByteVar":          exp.Byte_,
			"nIntVar":           exp.Int,
			"nUintVar":          exp.Uint,
			"nWordVar":          exp.Word,
			"nDintVar":          exp.Dint,
			"nUdintVar":         exp.Udint,
			"nDwordVar":         exp.Dword,
			"nLintVar":          exp.Lint,
			"nUlintVar":         exp.Ulint,
			"nLwordVar":         exp.Lword,
			"fRealVar":          exp.Real,
			"fLrealVar":         exp.Lreal,
			"tTimeVar":          exp.Time_,
			"tdTimeOfDayVar":    exp.Tod,
			"dDateVar":          exp.Date,
			"dtDateTimeVar":     exp.Dt,
			"sStringVar":        exp.String,
		}
		require.NoError(t, client.WriteValue(plcPort, fb+".stStructVar", m))
		time.Sleep(20 * time.Millisecond)
		got, err := client.ReadValue(plcPort, fb+".stStructVar")
		require.NoError(t, err)
		assertSnapshot(t, got, seed)
	})

	// ------------------------------------------------------------------
	// 4. FB control vars — public vars on the Deterministic FB instance
	// ------------------------------------------------------------------
	t.Run("fb_control_vars", func(t *testing.T) {
		const controlFB = "Main.fbTypeTest"
		runWriteCases(t, client, []writeCase{
			{controlFB + ".nSeed", uint32(999), eqCheck(uint32(999))},
			{controlFB + ".bAutoMode", true, eqCheck(true)},
			{controlFB + ".bAutoMode", false, eqCheck(false)},
			{controlFB + ".nAutoTickInterval", uint32(123), eqCheck(uint32(123))},
		})
		// Restore stable state
		require.NoError(t, client.WriteValue(plcPort, controlFB+".bAutoMode", false), "restore auto_mode")
		require.NoError(t, client.WriteValue(plcPort, controlFB+".nSeed", uint32(0)), "restore seed")
		require.NoError(t, client.WriteValue(plcPort, controlFB+".nAutoTickInterval", uint32(50)), "restore tick interval")
	})

	// ------------------------------------------------------------------
	// 5. Array write + read
	// ------------------------------------------------------------------
	t.Run("int_array", func(t *testing.T) {
		writeVal := make([]int16, 10)
		for i := range writeVal {
			writeVal[i] = int16(i * 10)
		}
		require.NoError(t, client.WriteValue(plcPort, fb+".aIntArray", writeVal))
		time.Sleep(20 * time.Millisecond)
		got, err := client.ReadValue(plcPort, fb+".aIntArray")
		require.NoError(t, err)
		arr, ok := got.([]any)
		require.True(t, ok, "expected []any got %T", got)
		require.Len(t, arr, 10)
		for i, v := range arr {
			assert.Equal(t, int16(i*10), v, "aIntArray[%d]", i)
		}
	})

	t.Run("dint_array", func(t *testing.T) {
		writeVal := make([]int32, 10)
		for i := range writeVal {
			writeVal[i] = int32(i * 1000)
		}
		require.NoError(t, client.WriteValue(plcPort, fb+".aDintArray", writeVal))
		time.Sleep(20 * time.Millisecond)
		got, err := client.ReadValue(plcPort, fb+".aDintArray")
		require.NoError(t, err)
		arr, ok := got.([]any)
		require.True(t, ok, "expected []any got %T", got)
		require.Len(t, arr, 10)
		for i, v := range arr {
			assert.Equal(t, int32(i*1000), v, "aDintArray[%d]", i)
		}
	})

	t.Run("2d_array", func(t *testing.T) {
		writeVal := [][]int16{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}
		require.NoError(t, client.WriteValue(plcPort, fb+".aIntArray2d", writeVal))
		time.Sleep(20 * time.Millisecond)
		got, err := client.ReadValue(plcPort, fb+".aIntArray2d")
		require.NoError(t, err)
		outer, ok := got.([]any)
		require.True(t, ok, "expected []any got %T", got)
		require.Len(t, outer, 3)
		for i, row := range outer {
			inner, ok := row.([]any)
			require.True(t, ok)
			require.Len(t, inner, 3)
			for j, v := range inner {
				assert.Equal(t, writeVal[i][j], v, "aIntArray2d[%d][%d]", i, j)
			}
		}
	})
}

// -----------------------------------------------------------------------
// Struct packing — the same four fields under pack_mode 0/2/4/8.
//
// Field layout:  bBoolMember1 BOOL | nDwordMember DWORD | bBoolMember2 BOOL | nLwordMember LWORD
// This layout causes the DWORD and LWORD to land at different byte offsets
// under each pack mode, so the serializer must use subItem.Offset (as
// reported by the PLC's ADS type info) rather than packing fields naively.
//
// Expected offsets per Beckhoff docs:
//   pack_mode 0:  bBoolMember1=0  nDwordMember=1  bBoolMember2=5  nLwordMember=6   (total 14)
//   pack_mode 2:  bBoolMember1=0  nDwordMember=2  bBoolMember2=6  nLwordMember=8   (total 16)
//   pack_mode 4:  bBoolMember1=0  nDwordMember=4  bBoolMember2=8  nLwordMember=12  (total 20)
//   pack_mode 8:  bBoolMember1=0  nDwordMember=4  bBoolMember2=8  nLwordMember=16  (total 24)
//
// Two-phase strategy to avoid circular validation:
//
//   Phase 1 — plc_read: PLC fills structs from seed; Go reads structurally
//     and asserts against independently computed expected values.
//     If Serialize had a bug but Deserialize was symmetric, this would catch it.
//
//   Phase 2 — write_then_linear_read: Go writes structs structurally
//     (WriteValue with map[string]any), then reads each field individually
//     via its linear ADS path. Structural write is validated by linear reads —
//     different code paths, no circular dependency.
// -----------------------------------------------------------------------

// packExpected computes the deterministic values for any pack struct from seed.
// Formula matches StructTests.TcPOU exactly.
func packExpected(seed uint32) (b1 bool, d1 uint32, b2 bool, lw1 uint64) {
	return seed%2 == 0, seed, seed%3 == 0, uint64(seed) * 3
}

var structTestFieldNames = struct{ b1, d1, b2, lw1 string }{
	"bBoolMember1", "nDwordMember", "bBoolMember2", "nLwordMember",
}

func assertPackStruct(t *testing.T, got any, seed uint32) {
	t.Helper()
	m, ok := got.(map[string]any)
	require.True(t, ok, "expected map[string]any got %T", got)
	b1, d1, b2, lw1 := packExpected(seed)
	assert.Equal(t, b1, m[structTestFieldNames.b1], structTestFieldNames.b1)
	assert.Equal(t, d1, m[structTestFieldNames.d1], structTestFieldNames.d1)
	assert.Equal(t, b2, m[structTestFieldNames.b2], structTestFieldNames.b2)
	assert.Equal(t, lw1, m[structTestFieldNames.lw1], structTestFieldNames.lw1)
}

func TestStructPacking(t *testing.T) {
	client := newClient(t)

	fbName := "Main.fbStructTest"

	packPaths := []struct{ name, plcRead, writeTarget string }{
		{"pack_none", fbName + ".stPackNone", fbName + ".stPackNoneWrite"},
		{"pack_two", fbName + ".stPackTwo", fbName + ".stPackTwoWrite"},
		{"pack_four", fbName + ".stPackFour", fbName + ".stPackFourWrite"},
		{"pack_eight", fbName + ".stPackEight", fbName + ".stPackEightWrite"},
	}

	// ------------------------------------------------------------------
	// Phase 1: PLC-driven read.
	// Go writes a seed, PLC fills all 4 structs deterministically.
	// Go reads structurally and asserts against independently computed values.
	// ------------------------------------------------------------------
	t.Run("plc_read", func(t *testing.T) {
		for _, seed := range []uint32{0, 1, 2, 3, 6, 100, 0xFFFFFFFF} {
			seed := seed
			t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
				require.NoError(t, client.WriteValue(plcPort, fbName+".nSeed", seed))
				time.Sleep(20 * time.Millisecond)
				for _, p := range packPaths {
					p := p
					t.Run(p.name, func(t *testing.T) {
						got, err := client.ReadValue(plcPort, p.plcRead)
						require.NoError(t, err, "ReadValue %s", p.plcRead)
						assertPackStruct(t, got, seed)
					})
				}
			})
		}
	})

	// ------------------------------------------------------------------
	// Phase 2: Write structurally, read back each field linearly.
	// Structural Serialize and linear Deserialize are independent code paths.
	// ------------------------------------------------------------------
	t.Run("write_then_linear_read", func(t *testing.T) {
		const seed = uint32(42)
		b1, d1, b2, lw1 := packExpected(seed)

		val := map[string]any{structTestFieldNames.b1: b1, structTestFieldNames.d1: d1, structTestFieldNames.b2: b2, structTestFieldNames.lw1: lw1}

		for _, p := range packPaths {
			p := p
			t.Run(p.name, func(t *testing.T) {
				require.NoError(t, client.WriteValue(plcPort, p.writeTarget, val), "WriteValue %s", p.writeTarget)
				time.Sleep(20 * time.Millisecond)

				gotB1, err := client.ReadValue(plcPort, p.writeTarget+"."+structTestFieldNames.b1)
				require.NoError(t, err)
				assert.Equal(t, b1, gotB1, structTestFieldNames.b1)

				gotD1, err := client.ReadValue(plcPort, p.writeTarget+"."+structTestFieldNames.d1)
				require.NoError(t, err)
				assert.Equal(t, d1, gotD1, structTestFieldNames.d1)

				gotB2, err := client.ReadValue(plcPort, p.writeTarget+"."+structTestFieldNames.b2)
				require.NoError(t, err)
				assert.Equal(t, b2, gotB2, structTestFieldNames.b2)

				gotLw1, err := client.ReadValue(plcPort, p.writeTarget+"."+structTestFieldNames.lw1)
				require.NoError(t, err)
				assert.Equal(t, lw1, gotLw1, structTestFieldNames.lw1)
			})
		}
	})
}
func TestSymbolAttributes(t *testing.T) {
	client := newClient(t)

	sym, err := client.GetSymbol(plcPort, "Main.fbAttributeTest")
	require.NoError(t, err, "GetSymbol Main.fbAttributeTest")
	require.NotNil(t, sym)

	// Log diagnostic info to help debug attribute parsing in integration environment.
	t.Logf("Symbol Flags: 0x%X, TypeGUID: %s, AttributesCount: %d", uint32(sym.Flags), sym.TypeGUID, len(sym.Attributes))
	if len(sym.Attributes) == 0 {
		t.Logf("Attributes slice is empty; full symbol: %+v", sym)
		// Fallback: try uploading the entire symbol table and search for the symbol
		t.Log("Attempting UploadSymbols fallback to locate attributes")
		all, uerr := client.UploadSymbols(plcPort)
		if uerr != nil {
			t.Logf("UploadSymbols error: %v", uerr)
		} else {
			var found *adssymbol.AdsSymbol
			for i := range all {
				if all[i].Name == "Main.fbAttributeTest" {
					found = &all[i]
					break
				}
			}
			if found != nil {
				t.Logf("Found symbol in UploadSymbols: Flags=0x%X AttributesCount=%d", uint32(found.Flags), len(found.Attributes))
				for i, a := range found.Attributes {
					t.Logf("Upload Attr %d: name=%q value=%q", i, a.Name, a.Value)
				}
				// replace sym with found for subsequent assertions
				sym = found
			} else {
				t.Log("Symbol not found in upload results")
			}
		}
	} else {
		for i, a := range sym.Attributes {
			t.Logf("Attribute %d: name=%q value=%q", i, a.Name, a.Value)
		}
	}

	// Expect one attribute: a flag-only attribute.
	require.Len(t, sym.Attributes, 1, "expected 1 attribute on Main.fbAttributeTest")

	var foundCustom bool
	for _, a := range sym.Attributes {
		switch a.Name {
		case "some_made_up_var_key":
			foundCustom = true
			assert.Equal(t, "some_made_up_var_value", a.Value, "some_made_up_var_key value")
		}
	}

	require.True(t, foundCustom, "missing attribute: some_made_up_var_key")
}
