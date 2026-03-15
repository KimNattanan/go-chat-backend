package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewBucket(2, 1)

	if !tb.Allow() {
		t.Error("first Allow() should succeed")
	}
	if !tb.Allow() {
		t.Error("second Allow() should succeed")
	}
	if tb.Allow() {
		t.Error("third Allow() should fail (bucket empty)")
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	tb := NewBucket(10, 1)

	if !tb.AllowN(5) {
		t.Error("AllowN(5) should succeed")
	}
	if !tb.AllowN(5) {
		t.Error("AllowN(5) again should succeed")
	}
	if tb.AllowN(1) {
		t.Error("AllowN(1) when empty should fail")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewBucket(2, 10)
	tb.AllowN(2)
	if tb.Allow() {
		t.Error("should be empty after AllowN(2)")
	}
	time.Sleep(150 * time.Millisecond)
	if !tb.Allow() {
		t.Error("after 150ms refill, Allow() should succeed")
	}
}

func TestTokenBucket_Tokens(t *testing.T) {
	tb := NewBucket(5, 1)
	if got := tb.Tokens(); got != 5 {
		t.Errorf("initial Tokens() = %v, want 5", got)
	}
	tb.AllowN(3)
	got := tb.Tokens()

	if got < 2 {
		t.Errorf("after AllowN(3), Tokens() = %v, want >= 2", got)
	}
}

func TestTokenBucket_Allow_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		capacity float64
		rate     float64
		take     []float64
		want     []bool
	}{
		{"take_one_by_one", 2, 1, []float64{1, 1, 1}, []bool{true, true, false}},
		{"take_all_at_once", 3, 1, []float64{3, 1}, []bool{true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewBucket(tt.capacity, tt.rate)
			for i, n := range tt.take {
				got := tb.AllowN(n)
				if got != tt.want[i] {
					t.Errorf("AllowN(%v) #%d = %v, want %v", n, i, got, tt.want[i])
				}
			}
		})
	}
}
