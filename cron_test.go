package gron

// Test that invoking stop() before start() silently returns,
// without blocking the stop channel
func TestStopWithoutStart() {
	cron := New()
	cron.Stop()
}
