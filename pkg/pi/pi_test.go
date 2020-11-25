package pi

import "testing"

func TestPiParallel(t *testing.T) {
	t.Logf("%f", Parallel(1*1000*1000))
}

func TestPiSingle(t *testing.T) {
	t.Logf("%e", Single(1*1000*1000))
	t.Logf("%e", Single(2*1000*1000))
	t.Logf("%e", Single(4*1000*1000))
	t.Logf("%f", Single(8*1000*1000))
	t.Logf("%f", Single(16*1000*1000))
	t.Logf("%f", Single(32*1000*1000))
	t.Logf("%f", Single(64*1000*1000))
	t.Logf("%f", Single(128*1000*1000))
}
