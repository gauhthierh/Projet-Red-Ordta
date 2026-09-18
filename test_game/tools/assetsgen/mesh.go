package main

import (
	"bytes"
	"fmt"
	"math"
)

type vec3 struct {
	x float64
	y float64
	z float64
}

type face struct {
	vertices []int
	material string
}

type mesh struct {
	vertices []vec3
	faces    []face
}

func (m *mesh) vertex(point vec3) int {
	m.vertices = append(m.vertices, point)
	return len(m.vertices)
}

func (m *mesh) polygon(material string, points ...vec3) {
	indices := make([]int, 0, len(points))
	for _, point := range points {
		indices = append(indices, m.vertex(point))
	}
	m.faces = append(m.faces, face{vertices: indices, material: material})
}

func (m *mesh) box(center vec3, size vec3, material string) {
	hx, hy, hz := size.x/2, size.y/2, size.z/2
	vertices := []vec3{
		{center.x - hx, center.y - hy, center.z - hz},
		{center.x + hx, center.y - hy, center.z - hz},
		{center.x + hx, center.y + hy, center.z - hz},
		{center.x - hx, center.y + hy, center.z - hz},
		{center.x - hx, center.y - hy, center.z + hz},
		{center.x + hx, center.y - hy, center.z + hz},
		{center.x + hx, center.y + hy, center.z + hz},
		{center.x - hx, center.y + hy, center.z + hz},
	}
	ids := make([]int, len(vertices))
	for index, vertex := range vertices {
		ids[index] = m.vertex(vertex)
	}
	for _, indices := range [][]int{
		{0, 1, 2, 3}, {5, 4, 7, 6}, {4, 0, 3, 7},
		{1, 5, 6, 2}, {3, 2, 6, 7}, {4, 5, 1, 0},
	} {
		m.faces = append(m.faces, face{vertices: []int{ids[indices[0]], ids[indices[1]], ids[indices[2]], ids[indices[3]]}, material: material})
	}
}

func (m *mesh) cylinder(center vec3, radius, height float64, sides int, material string) {
	bottomCenter := m.vertex(vec3{center.x, center.y - height/2, center.z})
	topCenter := m.vertex(vec3{center.x, center.y + height/2, center.z})
	bottom := make([]int, sides)
	top := make([]int, sides)
	for index := range sides {
		angle := float64(index) * math.Pi * 2 / float64(sides)
		x := center.x + math.Cos(angle)*radius
		z := center.z + math.Sin(angle)*radius
		bottom[index] = m.vertex(vec3{x, center.y - height/2, z})
		top[index] = m.vertex(vec3{x, center.y + height/2, z})
	}
	for index := range sides {
		next := (index + 1) % sides
		m.faces = append(m.faces,
			face{vertices: []int{bottom[index], bottom[next], top[next], top[index]}, material: material},
			face{vertices: []int{bottomCenter, bottom[next], bottom[index]}, material: material},
			face{vertices: []int{topCenter, top[index], top[next]}, material: material},
		)
	}
}

func (m *mesh) cone(center vec3, radius, height float64, sides int, material string) {
	bottomCenter := m.vertex(vec3{center.x, center.y - height/2, center.z})
	tip := m.vertex(vec3{center.x, center.y + height/2, center.z})
	base := make([]int, sides)
	for index := range sides {
		angle := float64(index) * math.Pi * 2 / float64(sides)
		base[index] = m.vertex(vec3{center.x + math.Cos(angle)*radius, center.y - height/2, center.z + math.Sin(angle)*radius})
	}
	for index := range sides {
		next := (index + 1) % sides
		m.faces = append(m.faces,
			face{vertices: []int{base[index], base[next], tip}, material: material},
			face{vertices: []int{bottomCenter, base[next], base[index]}, material: material},
		)
	}
}

func (m *mesh) octahedron(center vec3, radius float64, material string) {
	points := []vec3{
		{center.x, center.y + radius, center.z}, {center.x, center.y - radius, center.z},
		{center.x + radius, center.y, center.z}, {center.x - radius, center.y, center.z},
		{center.x, center.y, center.z + radius}, {center.x, center.y, center.z - radius},
	}
	for _, indices := range [][]int{
		{0, 2, 4}, {0, 4, 3}, {0, 3, 5}, {0, 5, 2},
		{1, 4, 2}, {1, 3, 4}, {1, 5, 3}, {1, 2, 5},
	} {
		m.polygon(material, points[indices[0]], points[indices[1]], points[indices[2]])
	}
}

func (m *mesh) roof(center vec3, width, depth, height float64, material string) {
	hw, hd := width/2, depth/2
	points := []vec3{
		{center.x - hw, center.y, center.z - hd}, {center.x + hw, center.y, center.z - hd},
		{center.x, center.y + height, center.z - hd}, {center.x - hw, center.y, center.z + hd},
		{center.x + hw, center.y, center.z + hd}, {center.x, center.y + height, center.z + hd},
	}
	for _, indices := range [][]int{{0, 1, 2}, {5, 4, 3}, {0, 3, 4, 1}, {1, 4, 5, 2}, {2, 5, 3, 0}} {
		facePoints := make([]vec3, len(indices))
		for index, pointIndex := range indices {
			facePoints[index] = points[pointIndex]
		}
		m.polygon(material, facePoints...)
	}
}

func (m *mesh) obj(name, materialLibrary string) []byte {
	var output bytes.Buffer
	fmt.Fprintf(&output, "# Generated low-poly placeholder\nmtllib %s\no %s\ns off\n", materialLibrary, name)
	for _, vertex := range m.vertices {
		fmt.Fprintf(&output, "v %.5f %.5f %.5f\n", vertex.x, vertex.y, vertex.z)
	}
	currentMaterial := ""
	for _, face := range m.faces {
		if face.material != currentMaterial {
			fmt.Fprintf(&output, "usemtl %s\n", face.material)
			currentMaterial = face.material
		}
		output.WriteString("f")
		for _, index := range face.vertices {
			fmt.Fprintf(&output, " %d", index)
		}
		output.WriteByte('\n')
	}
	return output.Bytes()
}
