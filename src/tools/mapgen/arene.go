package main

import (
	"fmt"
	"math"
)

func arenaPoint(r, a, z float64) point { return point{170 + r*math.Cos(a), -165 + r*math.Sin(a), z} }

// Conserve une ouverture sur toute l'épaisseur des gradins pour chaque route.
func arenaSectorClear(a, b float64) bool {
	for r := 33.0; r <= 45; r++ {
		for j := 0; j <= 8; j++ {
			p := arenaPoint(r, a+(b-a)*float64(j)/8, 0)
			if onRoad(p.x, p.y, 6) {
				return false
			}
		}
	}
	return true
}

// Gradins pleins, avec des bords communs : aucun banc suspendu.
func (w *objWriter) arenaSector(inner, outer, a, b, bottom, top float64) {
	p, q := arenaPoint(inner, a, top), arenaPoint(outer, a, top)
	r, s := arenaPoint(outer, b, top), arenaPoint(inner, b, top)
	w.quad(p, q, r, s)
	pts := []point{p, q, r, s}
	for i := 0; i < 4; i++ {
		u, v := pts[i], pts[(i+1)%4]
		w.quad(point{u.x, u.y, bottom}, point{v.x, v.y, bottom}, v, u)
	}
}

var arenaPavilions = [][2]float64{{170, -211}, {216, -165}}
var arenaDummies = [][2]float64{{151, -186}, {155, -190}, {159, -193}}

func (w *objWriter) arenaDetails() {
	w.object("arena_floor")
	w.material("sand")
	w.cylinder(170, -165, .20, 45, .04, 128)
	for i := 0; i < 96; i++ {
		a, b := float64(i)*2*math.Pi/96, float64(i+1)*2*math.Pi/96
		w.material("light_stone")
		w.arenaSector(24.7, 25, a, b, .24, .25)
		if !arenaSectorClear(a, b) {
			continue
		}
		w.object(fmt.Sprintf("arena_stand_%02d", i))
		for row := 0; row < 4; row++ {
			inner := 34 + float64(row)*2
			top := .9 + float64(row)*.65
			w.material("stone")
			w.arenaSector(inner, inner+2, a, b, .24, top)
			w.material("wood")
			w.arenaSector(inner+.45, inner+1.35, a, b, top, top+.16)
		}
		w.material("stone")
		w.arenaSector(42, 43, a, b, .24, 3.65)
		w.material("wood")
		w.arenaSector(41.9, 43.1, a, b, 3.65, 3.83)
	}
	// Pavillons placés hors des accès ; suppression des anciennes tours et tentes.
	for i, p := range arenaPavilions {
		w.object(fmt.Sprintf("arena_pavilion_%d", i))
		w.material("stone")
		w.box(p[0], p[1], .24, 6, 5, .3)
		w.material("wood")
		for _, dx := range []float64{-2.6, 2.6} {
			for _, dy := range []float64{-2.1, 2.1} {
				w.box(p[0]+dx, p[1]+dy, .54, .3, .3, 4)
			}
		}
		w.box(p[0], p[1], 4.4, 6, 5, .25)
		w.material("roof_red")
		w.cone(p[0], p[1], 4.65, 4.8, 2, 4)
		w.material("wood")
		w.box(p[0], p[1], .54, 3, 1.2, .8)
	}
	for i, p := range arenaDummies {
		w.object(fmt.Sprintf("arena_dummy_%d", i))
		w.material("wood")
		w.cylinder(p[0], p[1], .24, .17, 2.2, 8)
		w.box(p[0], p[1], 1.5, 1.6, .2, .2)
		w.material("hay")
		w.cylinder(p[0], p[1], 1, .38, .8, 10)
		w.cylinder(p[0], p[1], 1.9, .25, .4, 10)
	}
}

func arenaDetailCollisions() []collision {
	var result []collision
	for i := 0; i < 96; i++ {
		a, b := float64(i)*2*math.Pi/96, float64(i+1)*2*math.Pi/96
		if !arenaSectorClear(a, b) {
			continue
		}
		mid := (a + b) / 2
		p := arenaPoint(38.5, mid, 0)
		result = append(result, rectangleCollision(fmt.Sprintf("arena_stand_%02d", i), "stand", p.x, p.y, 2*43.1*math.Sin((b-a)/2), 9.3, mid*180/math.Pi+90))
	}
	for i, p := range arenaPavilions {
		result = append(result, rectangleCollision(fmt.Sprintf("arena_pavilion_%d", i), "prop", p[0], p[1], 6, 5, 0))
	}
	for i, p := range arenaDummies {
		result = append(result, circleCollision(fmt.Sprintf("arena_dummy_%d", i), "training", p[0], p[1], .8))
	}
	return result
}
