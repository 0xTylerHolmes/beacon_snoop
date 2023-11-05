package beacon_snoop

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

var beaconEndpoint = "http://10.0.20.5:5052"
var listenAddr = "127.0.0.1:4000"

func Test_Snooper(t *testing.T) {
	config := SnooperConfig{
		remote:     beaconEndpoint,
		listenAddr: listenAddr,
	}
	snooper, err := NewSnooper(config)
	require.NoError(t, err)
	snooper.Snoop(context.Background())
}
