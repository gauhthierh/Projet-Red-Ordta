package monde

import "testing"

func TestPositionSurCarte(t *testing.T) {
	for _, cas := range []struct{ x, y, px, py float32 }{
		{0, 0, 500, 500}, {-300, 300, 0, 0}, {300, -300, 1000, 1000}, {0, 150, 500, 250},
	} {
		x, y := positionSurCarte(cas.x, cas.y, 600, 1000)
		if x != cas.px || y != cas.py {
			t.Fatalf("position (%v,%v) : obtenu (%v,%v)", cas.x, cas.y, x, y)
		}
	}
}
