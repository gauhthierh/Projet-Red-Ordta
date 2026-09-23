package main

import (
	"fmt"
	"math"
)

// Une liste commune sert au dessin et aux collisions. Chaque croisement
// avec un chemin ouvre une porte assez large pour le joueur et les poteaux.
func fenceSections(name string, x, y, sx, sy, gateWidth float64, gateSide string) [][4]float64 {
	sides := [][4]float64{{x - sx/2, y - sy/2, x + sx/2, y - sy/2}, {x - sx/2, y + sy/2, x + sx/2, y + sy/2}, {x - sx/2, y - sy/2, x - sx/2, y + sy/2}, {x + sx/2, y - sy/2, x + sx/2, y + sy/2}}
	names := []string{"south", "north", "west", "east"}
	var result [][4]float64
	for index, s := range sides {
		length := math.Hypot(s[2]-s[0], s[3]-s[1])
		steps := int(math.Ceil(length / .5))
		start := -1
		pos := func(i int) (float64, float64) {
			t := float64(i) / float64(steps)
			return s[0] + (s[2]-s[0])*t, s[1] + (s[3]-s[1])*t
		}
		for i := 0; i <= steps; i++ {
			blocked := i == steps
			if !blocked {
				px, py := pos(i)
				qx, qy := pos(i + 1)
				blocked = onRoad(px, py, 5.2) || onRoad(qx, qy, 5.2)
				if names[index] == gateSide && math.Abs((float64(i)+.5)*length/float64(steps)-length/2) < gateWidth/2+.5 {
					blocked = true
				}
			}
			if !blocked && start < 0 {
				start = i
			}
			if blocked && start >= 0 {
				ax, ay := pos(start)
				bx, by := pos(i)
				if math.Hypot(bx-ax, by-ay) > 1 {
					result = append(result, [4]float64{ax, ay, bx, by})
				}
				start = -1
			}
		}
	}
	return result
}

func boundarySegments() [][4]float64 {
	var result [][4]float64
	for i := 0; i < 128; i++ {
		a, b := float64(i)*2*math.Pi/128, float64(i+1)*2*math.Pi/128
		result = append(result, [4]float64{286 * math.Cos(a), 286 * math.Sin(a), 286 * math.Cos(b), 286 * math.Sin(b)})
	}
	return result
}

func (w *objWriter) continuousBoundary() {
	// Soubassement continu : plus de passage ni de vue sous les montagnes.
	for i, s := range boundarySegments() {
		w.object(fmt.Sprintf("boundary_cliff_%03d", i))
		w.material("rock")
		// Un léger recouvrement ferme également les coins entre les volumes.
		dx, dy := s[2]-s[0], s[3]-s[1]
		length := math.Hypot(dx, dy)
		w.segment(s[0]-dx/length, s[1]-dy/length, s[2]+dx/length, s[3]+dy/length, 16, -2, 14)
	}
	for i := 0; i < 128; i++ {
		a := float64(i) * 2 * math.Pi / 128
		x, y := 289*math.Cos(a), 289*math.Sin(a)
		w.object(fmt.Sprintf("boundary_peak_%03d", i))
		w.material("rock")
		height := 15 + 5*math.Sin(a*7) + 3*math.Cos(a*11)
		w.boulder(x, y, 8, 10, height)
		w.material("snow")
		w.boulder(x, y, 8+height*.8, 4, height*.22)
	}
}
