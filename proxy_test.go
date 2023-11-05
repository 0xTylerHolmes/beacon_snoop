package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var beaconEndpoint = "http://10.0.20.5:5052"
var listenAddr = ":4040"

func Test_Snooper(t *testing.T) {
	logFile := fmt.Sprintf("/tmp/beacon_snoop_logger_%d.log", time.Now().Unix())
	t.Logf("using temp log: %s", logFile)
	config := SnooperConfig{
		remote:      beaconEndpoint,
		listenAddr:  listenAddr,
		logFilePath: logFile,
		logHeaders:  true,
		logToFile:   true,
	}
	snooper, err := NewSnooper(config)
	require.NoError(t, err)
	snooper.Snoop(context.Background())
}
