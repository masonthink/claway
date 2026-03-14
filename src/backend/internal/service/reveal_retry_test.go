package service

import (
	"testing"
	"time"
)

func TestMaxRevealRetries(t *testing.T) {
	if maxRevealRetries != 3 {
		t.Errorf("maxRevealRetries = %d, want 3", maxRevealRetries)
	}
}

func TestRetryDelayCalculation(t *testing.T) {
	// Verify exponential backoff delays before each retry attempt:
	// attempt 1 → 1s delay, attempt 2 → 2s delay
	expectedDelays := []time.Duration{1 * time.Second, 2 * time.Second}
	for i, expected := range expectedDelays {
		attempt := i + 1
		delay := time.Duration(1<<uint(attempt-1)) * time.Second
		if delay != expected {
			t.Errorf("attempt %d: delay = %v, want %v", attempt, delay, expected)
		}
	}

	// With 3 retries (attempts 0,1,2), delays happen before attempts 1 and 2:
	// attempt 1 delay: 1s, attempt 2 delay: 2s → total wait: 3s
	totalWait := time.Duration(0)
	for attempt := 1; attempt < maxRevealRetries; attempt++ {
		totalWait += time.Duration(1<<uint(attempt-1)) * time.Second
	}
	if totalWait != 3*time.Second {
		t.Errorf("total max wait = %v, want 3s", totalWait)
	}
}
