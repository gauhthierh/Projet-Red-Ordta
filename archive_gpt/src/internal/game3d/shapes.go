package game3d

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Palette and modelling helpers keep all procedural art in Go.
var (
	wood    = color(.28, .15, .08)
	plaster = color(.82, .72, .53)
	roof    = color(.43, .16, .10)
	stone   = color(.36, .39, .42)
	steel   = color(.58, .66, .72)
	gold    = color(.92, .62, .12)
	cloth   = color(.12, .29, .40)
	skin    = color(.83, .58, .38)
	dark    = color(.045, .035, .025)
)

func box(parent *core.Node, x, y, z, sx, sy, sz float32, c *math32.Color) *graphic.Mesh {
	m := graphic.NewMesh(geometry.NewBox(sx, sy, sz), surfaceMaterial(c))
	m.SetPosition(x, y, z)
	parent.Add(m)
	return m
}
func sphere(parent *core.Node, x, y, z, r float32, c *math32.Color) *graphic.Mesh {
	m := graphic.NewMesh(geometry.NewSphere(float64(r), 8, 6), surfaceMaterial(c))
	m.SetPosition(x, y, z)
	parent.Add(m)
	return m
}
func cylinder(parent *core.Node, x, y, z, r, h float32, c *math32.Color) *graphic.Mesh {
	m := graphic.NewMesh(geometry.NewCylinder(float64(r), float64(h), 10, 1, true, true), surfaceMaterial(c))
	m.SetPosition(x, y, z)
	parent.Add(m)
	return m
}
func cone(parent *core.Node, x, y, z, r, h float32, c *math32.Color) *graphic.Mesh {
	m := graphic.NewMesh(geometry.NewCone(float64(r), float64(h), 8, 1, true), surfaceMaterial(c))
	m.SetPosition(x, y, z)
	parent.Add(m)
	return m
}

// shadow is a lightweight contact shadow; it does not pretend to be a full
// shadow-map solution, but keeps characters and props grounded visually.
func shadow(parent *core.Node, x, z, width, depth float32) *graphic.Mesh {
	m := material.NewStandard(color(.02, .04, .025))
	m.SetTransparent(true)
	m.SetOpacity(.27)
	m.SetDepthMask(false)
	mesh := graphic.NewMesh(geometry.NewPlane(width, depth), m)
	mesh.SetRotationX(-math32.Pi / 2)
	mesh.SetPosition(x, -.015, z)
	parent.Add(mesh)
	return mesh
}
