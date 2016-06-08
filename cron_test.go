package gron

import "testing"

// Test that invoking stop() before start() silently returns,
// without blocking the stop channel
func TestStopWithoutStart(t *testing.T) {
	cron := New()
	cron.Stop()
}
