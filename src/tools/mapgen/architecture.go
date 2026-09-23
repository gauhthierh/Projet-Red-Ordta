package main

import (
	"fmt"
	"math"
)

// Les huit tours encadrent les quatre passages, au lieu de les boucher.
func gateTowerPositions() []point {
	var positions []point
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		for _, offset := range []float64{-10, 10} {
			positions = append(positions, point{82*math.Cos(a) - offset*math.Sin(a), 82*math.Sin(a) + offset*math.Cos(a), 0})
		}
	}
	return positions
}

func (w *objWriter) fortifiedGates() {
	for i, p := range gateTowerPositions() {
		w.object(fmt.Sprintf("gate_tower_%02d", i))
		w.material("stone")
		w.cylinder(p.x, p.y, 0, 5, 13, 24)
		w.material("light_stone")
		w.cylinder(p.x, p.y, 0, 5.2, .65, 24)
		w.cylinder(p.x, p.y, 9.5, 5.2, .5, 24)
		w.cylinder(p.x, p.y, 12.5, 5.4, .65, 24)
		w.material("roof_blue")
		w.cone(p.x, p.y, 13.15, 6.1, 7, 24)
		w.material("metal")
		w.cylinder(p.x, p.y, 20, .12, 2, 8)
		w.material("cloth_blue")
		w.box(p.x+1, p.y, 20.7, 2, .08, 1)
	}
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		x, y := 82*math.Cos(a), 82*math.Sin(a)
		dx, dy := -math.Sin(a), math.Cos(a)
		w.object(fmt.Sprintf("gate_lintel_%d", i))
		w.material("stone")
		w.segment(x-dx*10, y-dy*10, x+dx*10, y+dy*10, 4, 6, 3)
		w.material("light_stone")
		w.segment(x-dx*10, y-dy*10, x+dx*10, y+dy*10, 4.3, 9, .5)
		for j := -4; j <= 4; j += 2 {
			w.box(x+dx*float64(j), y+dy*float64(j), 9.5, 1.4, 1.4, 1.2)
		}
	}
}

// Encadrements en relief, croisillons, colombages et corniches.
func (w *objWriter) houseDetails(x, y, sx, sy, height float64) {
	w.material("wood")
	for _, side := range []float64{-1, 1} {
		fy := y + side*(sy/2+.24)
		for _, wx := range []float64{x - sx*.27, x + sx*.27} {
			for _, offset := range []float64{-1.15, 1.15} {
				w.box(wx+offset, fy, height*.48-.12, .18, .18, 2.34)
				w.box(wx, fy, height*.48+1.05+offset, 2.5, .18, .18)
			}
			w.box(wx, fy, height*.48, .13, .22, 2.1)
			w.box(wx, fy, height*.48+1, 2.1, .22, .13)
		}
		for _, px := range []float64{x - sx/2 + .12, x, x + sx/2 - .12} {
			w.box(px, fy, .55, .28, .22, height-.55)
		}
		w.box(x, fy, height-.35, sx, .25, .3)
	}
	w.material("light_stone")
	w.box(x, y-sy/2-.4, 0, 3, .65, .22)
	w.material("metal")
	w.box(x+.65, y-sy/2-.38, 1.65, .16, .12, .3)
}

// Volume rocheux à plusieurs niveaux, plutôt qu'une pyramide pointue.
func (w *objWriter) boulder(x, y, z, radius, height float64) {
	const sides = 10
	rings := [][2]float64{{0, .88}, {.30, 1}, {.78, .72}, {1, .38}}
	for level := 0; level < len(rings)-1; level++ {
		for i := 0; i < sides; i++ {
			vertex := func(j, k int) point {
				j %= sides
				a := float64(j) * 2 * math.Pi / sides
				r := radius * rings[k][1] * (1 + .09*math.Sin(float64(j)*2.7))
				return point{x + math.Cos(a)*r, y + math.Sin(a)*r, z + rings[k][0]*height}
			}
			w.quad(vertex(i, level), vertex(i+1, level), vertex(i+1, level+1), vertex(i, level+1))
		}
	}
	for i := 0; i < sides; i++ {
		a, b := float64(i)*2*math.Pi/sides, float64(i+1)*2*math.Pi/sides
		// Le couvercle reprend exactement l'anneau irrégulier des parois.
		ra := radius * .38 * (1 + .09*math.Sin(float64(i)*2.7))
		rb := radius * .38 * (1 + .09*math.Sin(float64((i+1)%sides)*2.7))
		w.triangle(point{x, y, z + height}, point{x + math.Cos(a)*ra, y + math.Sin(a)*ra, z + height}, point{x + math.Cos(b)*rb, y + math.Sin(b)*rb, z + height})
	}
}

func (w *objWriter) keepDetails() {
	w.object("keep_entrance_and_banners")
	w.material("light_stone")
	// Portail en relief et fronton adossés au corps principal.
	w.box(-2.6, 43.25, 0, .8, 1.2, 6)
	w.box(2.6, 43.25, 0, .8, 1.2, 6)
	w.box(0, 43.25, 5.7, 6, 1.2, .8)
	w.material("wood")
	w.box(0, 43.05, .2, 4.2, .3, 5.5)
	w.material("metal")
	for _, z := range []float64{1.4, 3.7} {
		w.box(0, 42.85, z, 4.2, .12, .2)
	}
	w.material("roof_blue")
	w.gableRoof(0, 49, 13, 10, 14, 8)
	w.material("light_stone")
	w.box(0, 41.9, 6.5, 1.7, .25, 3.5)
	w.material("cloth_red")
	for _, x := range []float64{-7, 7} {
		w.box(x, 43.3, 6.5, 1.8, .15, 4.5)
		w.material("metal")
		w.box(x, 43.1, 10.8, 2.3, .2, .18)
		w.material("cloth_red")
	}
}
