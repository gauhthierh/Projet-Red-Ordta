// Command mapgen creates the static 3D version of the RED world map.
// It is an authoring tool: the game only needs the generated OBJ/MTL files.
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type point struct{ x, y, z float64 }

type triangleFace struct{ a, b, c point }

type objWriter struct {
	file            *os.File
	nextVertex      int
	nextNormal      int
	nextUV          int
	currentObject   string
	currentMaterial string
	partNumber      int
	partOpen        bool
	groupByMaterial bool
	facesByMaterial map[string][]triangleFace
	materialOrder   []string
}

type landmark struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Radius      float64 `json:"radius"`
	Description string  `json:"description"`
}

// collision describes a blocking footprint on the XY plane. Rectangle
// rotations are expressed in degrees to keep the generated JSON readable.
type collision struct {
	ID              string  `json:"id"`
	Category        string  `json:"category"`
	Shape           string  `json:"shape"`
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
	Width           float64 `json:"width,omitempty"`
	Height          float64 `json:"height,omitempty"`
	Radius          float64 `json:"radius,omitempty"`
	RotationDegrees float64 `json:"rotation_degrees,omitempty"`
}

type layout struct {
	Name       string      `json:"name"`
	WorldSize  float64     `json:"world_size"`
	NorthAxis  string      `json:"north_axis"`
	UpAxis     string      `json:"up_axis"`
	ImageWidth int         `json:"image_width"`
	Mapping    string      `json:"image_to_world"`
	Spawn      [3]float64  `json:"spawn"`
	Landmarks  []landmark  `json:"landmarks"`
	Collisions []collision `json:"collisions"`
	Bridges    []collision `json:"bridges"`
}

var materials = map[string][3]float64{
	"grass":       {0.22, 0.43, 0.16},
	"dark_grass":  {0.10, 0.25, 0.10},
	"path":        {0.58, 0.42, 0.25},
	"stone":       {0.45, 0.47, 0.49},
	"light_stone": {0.68, 0.68, 0.64},
	"rock":        {0.30, 0.32, 0.34},
	"wood":        {0.31, 0.17, 0.08},
	"roof_red":    {0.52, 0.14, 0.08},
	"roof_blue":   {0.10, 0.25, 0.43},
	"wall":        {0.64, 0.59, 0.48},
	"water":       {0.05, 0.42, 0.66},
	"magic":       {0.10, 0.88, 0.94},
	"sand":        {0.66, 0.50, 0.29},
	"field":       {0.53, 0.39, 0.12},
	"crop":        {0.90, 0.66, 0.14},
	"marsh":       {0.15, 0.31, 0.22},
	"leaf":        {0.08, 0.30, 0.11},
	"leaf_light":  {0.18, 0.48, 0.18},
	"metal":       {0.55, 0.59, 0.62},
	"forge":       {1.00, 0.28, 0.04},
	"black":       {0.03, 0.03, 0.03},
	"white":       {0.88, 0.88, 0.82},
	"window":      {0.20, 0.68, 0.92},
	"cloth_red":   {0.76, 0.10, 0.08},
	"cloth_blue":  {0.08, 0.32, 0.72},
	"hay":         {0.82, 0.61, 0.16},
	"flower":      {0.88, 0.20, 0.55},
	"snow":        {0.88, 0.93, 0.96},
}

var materialTextures = map[string]string{
	"leaf":        "textures/foliage_detailed.png",
	"leaf_light":  "textures/foliage_detailed.png",
	"grass":       "textures/grass_detailed.png",
	"dark_grass":  "textures/grass_detailed.png",
	"field":       "textures/grass_detailed.png",
	"marsh":       "textures/marsh_detailed.png",
	"crop":        "textures/crop_detailed.png",
	"path":        "textures/path_detailed.png",
	"sand":        "textures/path_detailed.png",
	"hay":         "textures/path_detailed.png",
	"stone":       "textures/stone_detailed.png",
	"light_stone": "textures/stone_detailed.png",
	"rock":        "textures/stone_detailed.png",
	"wall":        "textures/stone_detailed.png",
	"snow":        "textures/stone_detailed.png",
	"wood":        "textures/wood_detailed.png",
	"water":       "textures/water_detailed.png",
	"magic":       "textures/water_detailed.png",
	"roof_red":    "textures/roof_detailed.png",
	"roof_blue":   "textures/roof_detailed.png",
	"metal":       "textures/metal_detailed.png",
	"cloth_red":   "textures/cloth_detailed.png",
	"cloth_blue":  "textures/cloth_detailed.png",
}

func backgroundTrees() []point {
	landmarkClearings := []struct{ x, y, radius float64 }{
		{0, 0, 100}, {-150, 150, 48}, {-195, 0, 44}, {170, -165, 56},
		{0, -200, 50}, {-170, -165, 62}, {0, 230, 48}, {175, 165, 64}, {210, 15, 58},
	}

	points := make([]point, 0, 110)
	for i := 0; i < 150; i++ {
		a := float64(i) * 2.399963229728653
		radius := 205 + float64((i*19)%58)
		x, y := math.Cos(a)*radius, math.Sin(a)*radius

		// Keep the northern river corridor and the principal outer gates open.
		if y > 95 && math.Abs(x) < 25 {
			continue
		}
		if math.Abs(x) < 11 || math.Abs(y) < 11 {
			continue
		}

		insideClearing := false
		for _, clearing := range landmarkClearings {
			if math.Hypot(x-clearing.x, y-clearing.y) < clearing.radius {
				insideClearing = true
				break
			}
		}
		if !insideClearing {
			points = append(points, point{x: x, y: y, z: .65 + float64(i%4)*.09})
		}
	}
	return points
}

func main() {
	output := filepath.Clean(filepath.Join("..", "assets", "maps", "red_world"))
	if len(os.Args) > 1 {
		output = os.Args[1]
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		panic(err)
	}

	writeMaterials(filepath.Join(output, "red_world_map_3d.mtl"))
	writeWorld(filepath.Join(output, "red_world_map_3d.obj"))
	writeCollision(filepath.Join(output, "red_world_collision.obj"))
	writeLayout(filepath.Join(output, "red_world_layout.json"))
}

func newOBJ(path, materialFile string, groupByMaterial bool) *objWriter {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(f, "# Static Z-up world map generated by src/tools/mapgen")
	fmt.Fprintln(f, "mtllib", materialFile)
	fmt.Fprintln(f, "s off")
	return &objWriter{
		file: f, nextVertex: 1, nextNormal: 1, nextUV: 1,
		groupByMaterial: groupByMaterial,
		facesByMaterial: make(map[string][]triangleFace),
	}
}

func (w *objWriter) close() error {
	if w.groupByMaterial {
		for _, materialName := range w.materialOrder {
			w.currentObject = "world"
			w.currentMaterial = materialName
			w.partNumber = 0
			w.partOpen = false
			for _, face := range w.facesByMaterial[materialName] {
				w.writeTriangle(face.a, face.b, face.c)
			}
		}
	}
	return w.file.Close()
}

func (w *objWriter) object(name string) {
	w.currentObject = name
	w.partNumber = 0
	w.partOpen = false
}

func (w *objWriter) material(name string) {
	if name == w.currentMaterial && !w.partOpen {
		return
	}
	w.currentMaterial = name
	w.partOpen = false
}

// beginPart guarantees that each generated OBJ object uses exactly one
// material. The material name is embedded in the object name so the Go loader
// can attach repeatable textures without duplicating special-case code.
func (w *objWriter) beginPart() {
	if w.partOpen {
		return
	}
	name := w.currentObject
	if name == "" {
		name = "unnamed"
	}
	materialName := w.currentMaterial
	if materialName == "" {
		materialName = "white"
	}
	fmt.Fprintf(w.file, "o %s__material_%s__part_%02d\n", name, materialName, w.partNumber)
	fmt.Fprintln(w.file, "usemtl", materialName)
	w.partNumber++
	w.partOpen = true
}

func (w *objWriter) triangle(a, b, c point) {
	if w.groupByMaterial {
		if _, exists := w.facesByMaterial[w.currentMaterial]; !exists {
			w.materialOrder = append(w.materialOrder, w.currentMaterial)
		}
		w.facesByMaterial[w.currentMaterial] = append(w.facesByMaterial[w.currentMaterial], triangleFace{a, b, c})
		return
	}
	w.writeTriangle(a, b, c)
}

func (w *objWriter) writeTriangle(a, b, c point) {
	w.beginPart()
	u := point{b.x - a.x, b.y - a.y, b.z - a.z}
	v := point{c.x - a.x, c.y - a.y, c.z - a.z}
	n := point{u.y*v.z - u.z*v.y, u.z*v.x - u.x*v.z, u.x*v.y - u.y*v.x}
	l := math.Sqrt(n.x*n.x + n.y*n.y + n.z*n.z)
	if l != 0 {
		n.x, n.y, n.z = n.x/l, n.y/l, n.z/l
	}
	fmt.Fprintf(w.file, "v %.4f %.4f %.4f\nv %.4f %.4f %.4f\nv %.4f %.4f %.4f\n", a.x, a.y, a.z, b.x, b.y, b.z, c.x, c.y, c.z)
	fmt.Fprintf(w.file, "vn %.5f %.5f %.5f\n", n.x, n.y, n.z)
	ua, va := textureCoordinates(a, n)
	ub, vb := textureCoordinates(b, n)
	uc, vc := textureCoordinates(c, n)
	fmt.Fprintf(w.file, "vt %.5f %.5f\nvt %.5f %.5f\nvt %.5f %.5f\n", ua, va, ub, vb, uc, vc)
	fmt.Fprintf(
		w.file,
		"f %d/%d/%d %d/%d/%d %d/%d/%d\n",
		w.nextVertex, w.nextUV, w.nextNormal,
		w.nextVertex+1, w.nextUV+1, w.nextNormal,
		w.nextVertex+2, w.nextUV+2, w.nextNormal,
	)
	w.nextVertex += 3
	w.nextNormal++
	w.nextUV += 3
}

// textureCoordinates projects each face on its most suitable plane.
// One texture tile covers twelve world units. Values greater than one are
// repeated by the material configured in src/3D/monde/carte.go.
func textureCoordinates(p, normal point) (float64, float64) {
	const tileSize = 12.0
	ax, ay, az := math.Abs(normal.x), math.Abs(normal.y), math.Abs(normal.z)
	switch {
	case az >= ax && az >= ay:
		return p.x / tileSize, p.y / tileSize
	case ay >= ax:
		return p.x / tileSize, p.z / tileSize
	default:
		return p.y / tileSize, p.z / tileSize
	}
}

func (w *objWriter) quad(a, b, c, d point) {
	w.triangle(a, b, c)
	w.triangle(a, c, d)
}

func (w *objWriter) box(x, y, z, sx, sy, sz float64) {
	x0, x1 := x-sx/2, x+sx/2
	y0, y1 := y-sy/2, y+sy/2
	z0, z1 := z, z+sz
	a, b, c, d := point{x0, y0, z0}, point{x1, y0, z0}, point{x1, y1, z0}, point{x0, y1, z0}
	e, f, g, h := point{x0, y0, z1}, point{x1, y0, z1}, point{x1, y1, z1}, point{x0, y1, z1}
	w.quad(a, d, c, b)
	w.quad(e, f, g, h)
	w.quad(a, b, f, e)
	w.quad(b, c, g, f)
	w.quad(c, d, h, g)
	w.quad(d, a, e, h)
}

func (w *objWriter) segment(x1, y1, x2, y2, width, z, height float64) {
	dx, dy := x2-x1, y2-y1
	l := math.Hypot(dx, dy)
	px, py := dy/l*width/2, -dx/l*width/2
	a := point{x1 + px, y1 + py, z}
	b := point{x2 + px, y2 + py, z}
	c := point{x2 - px, y2 - py, z}
	d := point{x1 - px, y1 - py, z}
	e, f, g, h := a, b, c, d
	e.z, f.z, g.z, h.z = z+height, z+height, z+height, z+height
	w.quad(e, f, g, h)
	if height > 0 {
		w.quad(a, b, f, e)
		w.quad(b, c, g, f)
		w.quad(c, d, h, g)
		w.quad(d, a, e, h)
	}
}

func (w *objWriter) cylinder(x, y, z, radius, height float64, sides int) {
	top := point{x, y, z + height}
	bottom := point{x, y, z}
	for i := 0; i < sides; i++ {
		a0 := float64(i) * 2 * math.Pi / float64(sides)
		a1 := float64(i+1) * 2 * math.Pi / float64(sides)
		b0 := point{x + math.Cos(a0)*radius, y + math.Sin(a0)*radius, z}
		b1 := point{x + math.Cos(a1)*radius, y + math.Sin(a1)*radius, z}
		t0, t1 := b0, b1
		t0.z, t1.z = z+height, z+height
		w.quad(b0, b1, t1, t0)
		w.triangle(top, t0, t1)
		w.triangle(bottom, b1, b0)
	}
}

func (w *objWriter) cone(x, y, z, radius, height float64, sides int) {
	top := point{x, y, z + height}
	center := point{x, y, z}
	for i := 0; i < sides; i++ {
		a0 := float64(i) * 2 * math.Pi / float64(sides)
		a1 := float64(i+1) * 2 * math.Pi / float64(sides)
		p0 := point{x + math.Cos(a0)*radius, y + math.Sin(a0)*radius, z}
		p1 := point{x + math.Cos(a1)*radius, y + math.Sin(a1)*radius, z}
		w.triangle(p0, p1, top)
		w.triangle(center, p1, p0)
	}
}

func (w *objWriter) gableRoof(x, y, z, sx, sy, height float64) {
	x0, x1 := x-sx/2, x+sx/2
	y0, y1 := y-sy/2, y+sy/2
	r0, r1 := point{x, y0, z + height}, point{x, y1, z + height}
	a, b := point{x0, y0, z}, point{x1, y0, z}
	c, d := point{x1, y1, z}, point{x0, y1, z}
	w.quad(a, r0, r1, d)
	w.quad(b, c, r1, r0)
	w.quad(a, d, c, b)
	w.triangle(a, b, r0)
	w.triangle(d, r1, c)
}

func (w *objWriter) fenceRect(name string, x, y, sx, sy, gateWidth float64, gateSide string) {
	w.object(name)
	w.material("wood")
	for _, section := range fenceSections(name, x, y, sx, sy, gateWidth, gateSide) {
		w.segment(section[0], section[1], section[2], section[3], .28, .8, .35)
		w.segment(section[0], section[1], section[2], section[3], .28, 1.7, .25)
		length := math.Hypot(section[2]-section[0], section[3]-section[1])
		steps := int(math.Ceil(length / 4))
		for i := 0; i <= steps; i++ {
			t := float64(i) / float64(steps)
			w.cylinder(section[0]+t*(section[2]-section[0]), section[1]+t*(section[3]-section[1]), 0, .32, 2.4, 8)
		}
	}
}

func (w *objWriter) lamp(name string, x, y float64) {
	w.object(name)
	w.material("metal")
	w.cylinder(x, y, 0, .22, 4.2, 7)
	// La lanterne possède un boîtier métallique texturé. Le matériau lumineux
	// reste limité aux quatre vitres afin d'éviter un cube orange uniforme.
	w.box(x, y, 4.15, 1.05, 1.05, 1.15)
	w.cone(x, y, 5.30, .82, .48, 4)
	w.material("forge")
	w.box(x, y-.54, 4.40, .66, .05, .62)
	w.box(x, y+.49, 4.40, .66, .05, .62)
	w.box(x-.54, y, 4.40, .05, .66, .62)
	w.box(x+.49, y, 4.40, .05, .66, .62)
}

func (w *objWriter) crystal(name string, x, y, scale float64) {
	w.object(name)
	w.material("magic")
	w.cone(x, y, .15, 1.2*scale, 5.5*scale, 6)
}

func (w *objWriter) bush(name string, x, y, scale float64) {
	w.object(name)
	w.material("leaf")
	w.cone(x-.7*scale, y, .1, 1.8*scale, 2.2*scale, 8)
	w.material("leaf_light")
	w.cone(x+.7*scale, y+.25*scale, .1, 1.55*scale, 1.9*scale, 8)
}

func (w *objWriter) house(name string, x, y, sx, sy, height float64, roof string) {
	w.object(name)
	w.material("stone")
	w.box(x, y, 0, sx+1, sy+1, .55)
	w.material("wall")
	w.box(x, y, .55, sx, sy, height-.55)
	w.material(roof)
	w.gableRoof(x, y, height, sx+2, sy+2, 4)
	w.material("wood")
	w.box(x, y-sy/2-.18, .10, 2.2, .35, 3.5)
	w.material("window")
	for _, wx := range []float64{x - sx*.27, x + sx*.27} {
		w.box(wx, y-sy/2-.20, height*.48, 2.1, .25, 2.1)
		w.box(wx, y+sy/2+.02, height*.48, 2.1, .25, 2.1)
	}
	w.material("stone")
	w.box(x+sx*.30, y+sy*.12, height, 1.7, 1.7, 5)
	w.houseDetails(x, y, sx, sy, height)
}

func (w *objWriter) tree(name string, x, y, scale float64) {
	if onRoad(x, y, 5+3.9*scale) {
		return
	}
	generatedTrees = append(generatedTrees, circleCollision(name, "tree", x, y, .7*scale+.15))
	w.object(name)
	w.material("wood")
	w.cylinder(x, y, 0, .7*scale, 4*scale, 7)
	w.material("leaf")
	w.cone(x, y, 2.8*scale, 3.9*scale, 6.3*scale, 9)
	w.material("leaf_light")
	w.cone(x, y, 5.8*scale, 3.1*scale, 5.2*scale, 9)
	w.material("leaf")
	w.cone(x, y, 8.5*scale, 2.1*scale, 4*scale, 9)
}

func (w *objWriter) rock(name string, x, y, radius, height float64) {
	w.object(name)
	w.material("rock")
	w.boulder(x, y, 0, radius, height)
}

func writeWorld(path string) {
	generatedTrees = nil
	w := newOBJ(path, "red_world_map_3d.mtl", true)
	defer w.close()

	// Ground and natural boundaries.
	w.object("terrain_base")
	w.material("grass")
	w.box(0, 0, -2, 600, 600, 2)
	w.continuousBoundary()

	// Rivers and bridges corresponding to the blue channels on the 2D map.
	w.material("water")
	for i, s := range rivers {
		w.object(fmt.Sprintf("river_%02d", i))
		w.segment(s[0], s[1], s[2], s[3], 15, .05, .08)
	}
	// Stone banks make the winding river readable from the gameplay camera.
	for i, p := range [][2]float64{{-2, 112}, {-20, 92}, {-52, 73}, {-82, 43}, {-99, 12}, {-101, -34}, {-94, -71}, {-65, -92}, {-25, -103}, {25, -103}, {69, -96}, {108, -78}, {134, -49}, {169, -41}, {213, -46}, {258, -51}} {
		w.rock(fmt.Sprintf("river_bank_rock_%02d", i), p[0], p[1], 2.2+float64(i%3), 2.5+float64(i%4))
	}
	// Northern waterfall, cliff ledges and the bridge seen on the 2D plan.
	w.object("north_waterfall")
	w.material("water")
	w.segment(4, 122, 4, 170, 11, .4, 7)
	for i, p := range [][2]float64{{-17, 145}, {-12, 170}, {21, 148}, {18, 176}} {
		w.rock(fmt.Sprintf("waterfall_cliff_%02d", i), p[0], p[1], 7+float64(i%2)*2, 14+float64(i%3)*4)
	}

	// Main roads. They use several short sections so their silhouette follows
	// the winding paths of the illustrated 2D map instead of straight chords.
	w.material("path")
	routes := [][][2]float64{
		{{-20, 0}, {-82, 0}, {-125, 0}, {-162, 0}, {-165, -22}, {-195, -22}},
		{{-125, 0}, {-125, 75}, {-112, 122}, {-120, 120}},
		{{0, 25}, {-22, 30}, {-22, 72}, {0, 72}, {0, 82}, {-28, 105}, {-28, 158}, {-48, 195}, {-48, 225}},
		{{82, 0}, {110, 20}, {112, 102}, {147, 118}},
		{{20, 0}, {25, -16}, {70, -16}, {70, 0}, {82, 0}, {135, 4}, {178, 11}, {183, 15}},
		{{0, -125}, {65, -125}, {105, -101}, {139, -137}, {150, -150}},
		{{0, -20}, {10, -35}, {10, -70}, {0, -82}, {0, -125}, {0, -160}},
		{{0, -125}, {-72, -125}, {-115, -145}, {-118, -169}},
		{{-150, 150}, {-108, 181}, {-62, 207}, {-25, 224}, {0, 230}},
		{{0, 230}, {48, 220}, {93, 204}, {136, 187}, {175, 165}},
		{{175, 165}, {193, 126}, {201, 87}, {218, 50}, {210, 15}},
		{{210, 15}, {207, -30}, {194, -75}, {181, -122}, {170, -165}},
		{{170, -165}, {126, -180}, {82, -188}, {40, -183}, {0, -180}},
		{{0, -180}, {-45, -186}, {-90, -179}, {-132, -173}, {-170, -165}},
		{{-170, -165}, {-187, -126}, {-195, -84}, {-188, -42}, {-195, 0}},
		{{-195, 0}, {-193, 44}, {-184, 86}, {-166, 126}, {-150, 150}},
	}
	mapRoutes = routes
	w.connectedRoads(routes)

	// Distinct biome floors make every gameplay zone readable at a glance.
	for _, zone := range []struct {
		name, material string
		x, y, radius   float64
	}{
		{"guild_clearing", "dark_grass", -150, 150, 42},
		{"forge_ground", "stone", -195, 0, 39},
		{"cemetery_ground", "dark_grass", 0, -205, 43},
		{"wolf_forest_ground", "dark_grass", -172, 92, 84},
		{"raven_cliff_ground", "rock", 0, 230, 43},
		{"mana_grove_ground", "dark_grass", 210, 15, 51},
	} {
		w.object(zone.name)
		w.material(zone.material)
		w.cylinder(zone.x, zone.y, .07, zone.radius, .08, 40)
	}

	// Fortified village, central plaza and four gates.
	w.material("stone")
	for i := 0; i < 40; i++ {
		a0 := float64(i) * 2 * math.Pi / 40
		a1 := float64(i+1) * 2 * math.Pi / 40
		// Openings at the four cardinal gates.
		mid := (a0 + a1) / 2
		if math.Abs(math.Sin(mid)) < .08 || math.Abs(math.Cos(mid)) < .08 {
			continue
		}
		w.object(fmt.Sprintf("village_wall_%02d", i))
		w.segment(math.Cos(a0)*82, math.Sin(a0)*82, math.Cos(a1)*82, math.Sin(a1)*82, 3.5, 0, 6)
		for _, a := range []float64{a0, a1} {
			w.object(fmt.Sprintf("village_battlement_%02d_%0.2f", i, a))
			w.box(math.Cos(a)*82, math.Sin(a)*82, 6, 2.3, 2.3, 1.7)
		}
	}
	w.fortifiedGates()
	w.object("central_plaza")
	w.material("light_stone")
	w.cylinder(0, 0, .16, 27, .18, 32)
	w.object("central_fountain")
	w.material("stone")
	w.cylinder(0, 0, .35, 7, 1.2, 18)
	w.material("water")
	w.cylinder(0, 0, 1.56, 5.6, .15, 18)
	w.material("stone")
	w.cylinder(0, 0, 1.7, 1, 6, 12)
	w.material("magic")
	w.cone(0, 0, 7.7, 1.3, 3.4, 8)
	for i := 0; i < 8; i++ {
		a := float64(i) * 2 * math.Pi / 8
		w.lamp(fmt.Sprintf("plaza_lamp_%02d", i), math.Cos(a)*22, math.Sin(a)*22)
	}

	// Village houses and market.
	houses := [][4]float64{{-49, 37, 17, 13}, {-31, 55, 14, 11}, {31, 55, 14, 11}, {50, 34, 17, 13}, {-54, 8, 15, 12}, {56, 8, 15, 12}, {-50, -30, 18, 14}, {-27, -51, 15, 12}, {27, -51, 15, 12}, {50, -29, 18, 14}, {-8, -56, 13, 10}}
	for i, h := range houses {
		roof := "roof_red"
		if i%3 == 0 {
			roof = "roof_blue"
		}
		w.house(fmt.Sprintf("village_house_%02d", i), h[0], h[1], h[2], h[3], 8, roof)
	}
	// The keep and its two towers reproduce the dominant building at the
	// northern side of the central village.
	w.house("village_keep", 0, 52, 25, 17, 13, "roof_blue")
	w.keepDetails()
	for i, x := range []float64{-15, 15} {
		w.object(fmt.Sprintf("village_keep_tower_%02d", i))
		w.material("stone")
		w.cylinder(x, 53, 0, 4.5, 17, 10)
		w.material("roof_blue")
		w.cone(x, 53, 17, 5.5, 5, 10)
	}
	for i, p := range [][2]float64{{35, 14}, {48, 10}, {39, 1}, {55, -2}} {
		w.object(fmt.Sprintf("market_stall_%02d", i))
		w.material("wood")
		w.box(p[0], p[1], .15, 7, 5, 3)
		w.material([]string{"roof_red", "roof_blue"}[i%2])
		w.box(p[0], p[1], 3.15, 8, 6, 1)
		w.material("cloth_red")
		if i%2 == 1 {
			w.material("cloth_blue")
		}
		for stripe := -2; stripe <= 2; stripe++ {
			w.box(p[0]+float64(stripe)*1.25, p[1]-3.1, 4.15, .65, .25, 1.5)
		}
	}
	w.fenceRect("market_fence", 44, 6, 35, 30, 8, "west")

	// Adventurers guild and equipment hall (northwest).
	w.house("adventurers_guild", -150, 150, 32, 24, 12, "roof_blue")
	w.material("sand")
	w.cylinder(-150, 120, .12, 22, .12, 20)
	w.fenceRect("guild_training_fence", -150, 120, 52, 42, 10, "east")
	w.material("wood")
	for i := 0; i < 5; i++ {
		w.box(-166+float64(i)*8, 120, .3, .4, 5, 2.5)
	}
	for i, p := range [][2]float64{{-168, 111}, {-156, 111}, {-144, 111}, {-132, 111}} {
		w.object(fmt.Sprintf("guild_target_%02d", i))
		w.material("wood")
		w.cylinder(p[0], p[1], 0, .45, 3.8, 7)
		w.material("white")
		w.cylinder(p[0], p[1], 3.2, 1.35, .35, 12)
		w.material("cloth_red")
		w.cylinder(p[0], p[1]-.05, 3.25, .55, .40, 12)
	}
	for i := 0; i < 5; i++ {
		w.object(fmt.Sprintf("guild_banner_%02d", i))
		w.material("metal")
		x := -172 + float64(i)*11
		w.cylinder(x, 148, 0, .2, 7, 7)
		w.material("cloth_blue")
		w.box(x+1.5, 148, 4.7, 3, .18, 3.8)
	}

	// Forge and crafting yard (west).
	w.house("blacksmith_forge", -195, 0, 28, 22, 9, "roof_red")
	w.object("forge_chimney")
	w.material("stone")
	w.cylinder(-205, 4, 9, 3.2, 17, 10)
	w.object("forge_fire")
	w.material("forge")
	w.box(-190, -6, .3, 5, 4, 2)
	w.material("metal")
	for i := 0; i < 4; i++ {
		w.object(fmt.Sprintf("forge_anvil_%02d", i))
		w.box(-210+float64(i)*10, -18, .2, 4, 2.5, 2)
	}
	w.fenceRect("forge_yard_fence", -195, -2, 58, 52, 9, "east")
	for i, p := range [][2]float64{{-218, 17}, {-210, 19}, {-182, 19}, {-174, 13}} {
		w.object(fmt.Sprintf("forge_crate_%02d", i))
		w.material("wood")
		w.box(p[0], p[1], .1, 4, 4, 3.5)
		w.material("metal")
		w.box(p[0], p[1], 1.35, 4.2, .3, .28)
	}
	for i := 0; i < 3; i++ {
		w.object(fmt.Sprintf("forge_woodpile_%02d", i))
		w.material("wood")
		w.cylinder(-220+float64(i)*2.2, -10, .2, .75, 5, 8)
	}

	// Training arena and goblin practice zone (southeast).
	w.object("training_arena")
	w.material("sand")
	w.cylinder(170, -165, .10, 45, .12, 36)
	w.material("wood")
	for i := 0; i < 36; i++ {
		a := float64(i) * 2 * math.Pi / 36
		if onRoad(170+math.Cos(a)*47, -165+math.Sin(a)*47, 5.5) {
			continue
		}
		w.cylinder(170+math.Cos(a)*47, -165+math.Sin(a)*47, 0, .55, 4.5, 6)
		if i%6 == 0 {
			w.material("cloth_red")
			w.box(170+math.Cos(a)*47, -165+math.Sin(a)*47, 3.2, 2.8, .22, 3.5)
			w.material("wood")
		}
	}
	for i, p := range [][2]float64{{150, -165}, {170, -185}, {190, -165}} {
		w.object(fmt.Sprintf("training_dummy_%02d", i))
		w.cylinder(p[0], p[1], 0, .7, 5, 7)
		w.box(p[0], p[1], 4, 4.5, .5, .5)
	}
	for i, p := range [][2]float64{{133, -195}, {207, -195}, {133, -135}, {207, -135}} {
		w.object(fmt.Sprintf("arena_watchtower_%02d", i))
		w.material("wood")
		w.box(p[0], p[1], 0, 6, 6, 8)
		w.material("roof_red")
		w.cone(p[0], p[1], 8, 5, 4, 4)
	}
	for i, p := range [][2]float64{{125, -178}, {215, -178}, {128, -151}, {212, -151}} {
		w.object(fmt.Sprintf("arena_tent_%02d", i))
		w.material([]string{"cloth_red", "white"}[i%2])
		w.cone(p[0], p[1], 0, 5, 7, 4)
	}

	// Cemetery and resurrection shrine (south).
	w.house("resurrection_shrine", 0, -180, 24, 18, 11, "roof_blue")
	w.fenceRect("cemetery_fence", 0, -205, 76, 74, 9, "north")
	w.material("stone")
	for row := 0; row < 4; row++ {
		for col := 0; col < 7; col++ {
			x := -30 + float64(col)*10
			y := -205 - float64(row)*9
			w.object(fmt.Sprintf("grave_%d_%d", row, col))
			w.box(x, y, .1, 2.4, .8, 4)
		}
	}
	w.object("resurrection_statue")
	w.material("white")
	w.cylinder(0, -165, .2, 2.3, 8, 10)
	w.cone(0, -165, 8.2, 3.8, 6, 8)
	for i := 0; i < 6; i++ {
		a := float64(i) * 2 * math.Pi / 6
		w.lamp(fmt.Sprintf("cemetery_lamp_%02d", i), math.Cos(a)*33, -205+math.Sin(a)*28)
	}

	// Boar meadow and farms (southwest).
	w.object("boar_meadow")
	w.material("field")
	w.box(-170, -165, .08, 105, 85, .10)
	for row := 0; row < 7; row++ {
		w.material("crop")
		w.segment(-210, -190+float64(row)*6, -165, -190+float64(row)*6, 2.2, .22, .35)
	}
	w.house("farmhouse", -205, -135, 22, 17, 7, "roof_red")
	w.fenceRect("farm_fence", -175, -169, 108, 92, 12, "east")
	w.object("windmill")
	w.material("wall")
	w.cylinder(-145, -135, 0, 6, 17, 10)
	w.material("wood")
	w.box(-145, -141, 13, 1, 13, 1)
	w.box(-151, -135, 13, 13, 1, 1)
	for i, p := range [][2]float64{{-217, -202}, {-210, -200}, {-202, -202}, {-194, -200}, {-217, -192}, {-208, -190}} {
		w.object(fmt.Sprintf("hay_bale_%02d", i))
		w.material("hay")
		w.cylinder(p[0], p[1], .3, 2.4, 3.4, 10)
	}
	for row := 0; row < 5; row++ {
		for col := 0; col < 4; col++ {
			w.object(fmt.Sprintf("orchard_tree_%02d_%02d", row, col))
			w.tree(fmt.Sprintf("orchard_tree_shape_%02d_%02d", row, col), -153+float64(col)*8, -196+float64(row)*8, .55)
		}
	}

	// Wolf forest (west and northwest).
	for i := 0; i < 48; i++ {
		a := float64(i) * 2.399963
		r := 22 + float64((i*17)%95)
		x, y := -172+math.Cos(a)*r, 92+math.Sin(a)*r
		if math.Hypot(x, y) < 100 || math.Hypot(x+150, y-150) < 32 || math.Hypot(x+195, y) < 30 {
			continue
		}
		w.tree(fmt.Sprintf("wolf_forest_tree_%02d", i), x, y, .8+float64(i%4)*.12)
	}
	for i, p := range [][2]float64{{-228, 62}, {-220, 112}, {-202, 145}, {-185, 55}, {-133, 78}, {-118, 118}} {
		w.rock(fmt.Sprintf("wolf_forest_boulder_%02d", i), p[0], p[1], 4+float64(i%3), 6+float64(i%2)*3)
	}
	w.object("wolf_den")
	w.material("black")
	w.box(-224, 92, .2, 9, 4, 7)

	// Raven cliffs and watchtower (north).
	for i := 0; i < 17; i++ {
		a := float64(i) * 2 * math.Pi / 17
		w.rock(fmt.Sprintf("raven_cliff_%02d", i), math.Cos(a)*35, 230+math.Sin(a)*26, 8+float64(i%3)*3, 18+float64(i%5)*3)
	}
	w.object("raven_watchtower")
	w.material("stone")
	w.cylinder(0, 230, 14, 7, 22, 10)
	w.material("roof_blue")
	w.cone(0, 230, 36, 9, 7, 10)
	for i := 0; i < 12; i++ {
		a := float64(i) * 2 * math.Pi / 12
		w.object(fmt.Sprintf("raven_pillar_%02d", i))
		w.material("stone")
		w.box(math.Cos(a)*23, 230+math.Sin(a)*19, 12, 2.2, 2.2, 7+float64(i%3)*2)
	}

	// Troll cave and marsh (northeast).
	w.object("troll_marsh")
	w.material("marsh")
	w.cylinder(175, 155, .08, 52, .10, 24)
	for i := 0; i < 11; i++ {
		a := float64(i) * 2 * math.Pi / 11
		w.rock(fmt.Sprintf("troll_cave_rock_%02d", i), 180+math.Cos(a)*28, 198+math.Sin(a)*17, 8+float64(i%3)*3, 12+float64(i%4)*4)
	}
	w.object("troll_cave_entrance")
	w.material("black")
	w.box(180, 180, .2, 16, 3, 13)
	for i, p := range [][2]float64{{142, 144}, {151, 162}, {166, 132}, {189, 136}, {204, 153}, {211, 173}, {151, 187}} {
		w.object(fmt.Sprintf("marsh_dead_tree_%02d", i))
		w.material("wood")
		w.cylinder(p[0], p[1], 0, .65, 6+float64(i%3), 7)
		w.segment(p[0], p[1], p[0]+4, p[1]+3, .45, 4.8, .5)
	}
	w.object("marsh_boardwalk")
	w.material("wood")
	w.segment(135, 137, 191, 171, 3.8, .45, .45)
	for i := 0; i < 9; i++ {
		w.rock(fmt.Sprintf("marsh_stepping_stone_%02d", i), 145+float64(i)*7, 150+math.Sin(float64(i))*6, 1.8, .8)
	}

	// Magical grove and mana spring (east).
	for i := 0; i < 18; i++ {
		a := float64(i) * 2 * math.Pi / 18
		w.tree(fmt.Sprintf("magic_grove_tree_%02d", i), 210+math.Cos(a)*39, 15+math.Sin(a)*34, .85)
	}
	w.object("mana_spring")
	w.material("water")
	w.cylinder(210, 15, .12, 17, .16, 20)
	w.material("magic")
	w.cone(210, 15, .3, 4, 15, 6)
	for i := 0; i < 6; i++ {
		a := float64(i) * 2 * math.Pi / 6
		w.object(fmt.Sprintf("mana_pillar_%02d", i))
		w.box(210+math.Cos(a)*23, 15+math.Sin(a)*23, .1, 3, 3, 10)
	}
	for i := 0; i < 14; i++ {
		a := float64(i) * 2 * math.Pi / 14
		r := 12 + float64((i%3)*9)
		w.crystal(fmt.Sprintf("mana_crystal_%02d", i), 210+math.Cos(a)*r, 15+math.Sin(a)*r, .55+float64(i%3)*.18)
	}
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		w.object(fmt.Sprintf("mana_ruin_arch_%02d", i))
		w.material("light_stone")
		x, y := 210+math.Cos(a)*31, 15+math.Sin(a)*31
		w.box(x, y, 0, 3.2, 3.2, 9)
		w.box(x+math.Cos(a+math.Pi/2)*6, y+math.Sin(a+math.Pi/2)*6, 0, 3.2, 3.2, 9)
		w.segment(x, y, x+math.Cos(a+math.Pi/2)*6, y+math.Sin(a+math.Pi/2)*6, 3.2, 8, 2)
	}

	// Dense outer vegetation and small roadside details complete the areas
	// between landmarks, as on the 2D reference map.
	for i, p := range backgroundTrees() {
		w.tree(fmt.Sprintf("background_tree_%03d", i), p.x, p.y, p.z)
	}
	for i, p := range [][2]float64{
		{-116, -29}, {-132, 28}, {-111, 65}, {-87, 128}, {-48, 166},
		{46, 181}, {84, 157}, {121, 128}, {145, 76}, {163, 38},
		{143, -3}, {124, -59}, {105, -135}, {63, -151}, {34, -130},
		{-36, -139}, {-75, -148}, {-116, -121}, {-143, -83}, {-157, -39},
		{89, 32}, {67, 101}, {-62, 108}, {-74, -42}, {47, -73},
	} {
		w.bush(fmt.Sprintf("roadside_bush_%02d", i), p[0], p[1], .75+float64(i%3)*.14)
	}
}

func circleCollision(id, category string, x, y, radius float64) collision {
	return collision{ID: id, Category: category, Shape: "circle", X: x, Y: y, Radius: radius}
}

func rectangleCollision(id, category string, x, y, width, height, rotation float64) collision {
	return collision{
		ID: id, Category: category, Shape: "rectangle", X: x, Y: y,
		Width: width, Height: height, RotationDegrees: rotation,
	}
}

func segmentCollision(id, category string, x1, y1, x2, y2, thickness float64) collision {
	dx, dy := x2-x1, y2-y1
	return rectangleCollision(
		id, category, (x1+x2)/2, (y1+y2)/2,
		math.Hypot(dx, dy), thickness, math.Atan2(dy, dx)*180/math.Pi,
	)
}

func appendFenceCollisions(dst []collision, name string, x, y, sx, sy, gateWidth float64, gateSide string) []collision {
	for i, s := range fenceSections(name, x, y, sx, sy, gateWidth, gateSide) {
		dst = append(dst, segmentCollision(fmt.Sprintf("%s_%03d", name, i), "fence", s[0], s[1], s[2], s[3], .8))
	}
	return dst
}

// worldCollisions mirrors the solid footprints emitted by writeWorld. Small
// decorative details (flowers, crops, banners and lamps) intentionally remain
// walkable so movement does not feel cluttered.
func worldCollisions() []collision {
	var result []collision

	// Natural border and river-bank rocks.
	for i, s := range boundarySegments() {
		c := segmentCollision(fmt.Sprintf("boundary_cliff_%03d", i), "mountain", s[0], s[1], s[2], s[3], 16)
		c.Width += 2
		result = append(result, c)
	}
	for i, p := range [][2]float64{{-2, 112}, {-20, 92}, {-52, 73}, {-82, 43}, {-99, 12}, {-101, -34}, {-94, -71}, {-65, -92}, {-25, -103}, {25, -103}, {69, -96}, {108, -78}, {134, -49}, {169, -41}, {213, -46}, {258, -51}} {
		result = append(result, circleCollision(fmt.Sprintf("river_bank_rock_%02d", i), "rock", p[0], p[1], 2+float64(i%3)))
	}
	for i, p := range [][2]float64{{-17, 145}, {-12, 170}, {21, 148}, {18, 176}} {
		result = append(result, circleCollision(fmt.Sprintf("waterfall_cliff_%02d", i), "cliff", p[0], p[1], 7+float64(i%2)*2))
	}
	// River sections are blocking except where a visible bridge creates a
	// deliberate opening. They remain rectangles, like walls and fences.
	for i, s := range rivers {
		result = append(result, segmentCollision(fmt.Sprintf("river_water_%02d", i), "river", s[0], s[1], s[2], s[3], 15))
	}
	result = append(result,
		segmentCollision("north_waterfall_south", "water", 4, 122, 4, 136, 10),
		segmentCollision("north_waterfall_north", "water", 4, 148, 4, 170, 10),
	)

	// Village wall, gate towers, fountain, houses and market.
	for i := 0; i < 40; i++ {
		a0 := float64(i) * 2 * math.Pi / 40
		a1 := float64(i+1) * 2 * math.Pi / 40
		mid := (a0 + a1) / 2
		if math.Abs(math.Sin(mid)) < .08 || math.Abs(math.Cos(mid)) < .08 {
			continue
		}
		result = append(result, segmentCollision(fmt.Sprintf("village_wall_%02d", i), "wall", math.Cos(a0)*82, math.Sin(a0)*82, math.Cos(a1)*82, math.Sin(a1)*82, 4))
	}
	for i, p := range gateTowerPositions() {
		result = append(result, circleCollision(fmt.Sprintf("gate_tower_%02d", i), "tower", p.x, p.y, 5.2))
	}
	result = append(result, circleCollision("central_fountain", "fountain", 0, 0, 7.2))
	for i := 0; i < 8; i++ {
		a := float64(i) * 2 * math.Pi / 8
		result = append(result, circleCollision(
			fmt.Sprintf("plaza_lamp_%02d", i),
			"lamp",
			math.Cos(a)*22,
			math.Sin(a)*22,
			.5,
		))
	}
	houses := [][4]float64{{-49, 37, 17, 13}, {-31, 55, 14, 11}, {31, 55, 14, 11}, {50, 34, 17, 13}, {-54, 8, 15, 12}, {56, 8, 15, 12}, {-50, -30, 18, 14}, {-27, -51, 15, 12}, {27, -51, 15, 12}, {50, -29, 18, 14}, {-8, -56, 13, 10}}
	for i, h := range houses {
		result = append(result, rectangleCollision(fmt.Sprintf("village_house_%02d", i), "building", h[0], h[1], h[2]+1, h[3]+1, 0))
	}
	result = append(result,
		rectangleCollision("village_keep", "building", 0, 52, 26, 18, 0),
		circleCollision("village_keep_tower_00", "tower", -15, 53, 4.7),
		circleCollision("village_keep_tower_01", "tower", 15, 53, 4.7),
	)
	for i, p := range [][2]float64{{35, 14}, {48, 10}, {39, 1}, {55, -2}} {
		result = append(result, rectangleCollision(fmt.Sprintf("market_stall_%02d", i), "stall", p[0], p[1], 8, 6, 0))
	}
	result = appendFenceCollisions(result, "market_fence", 44, 6, 35, 30, 8, "west")

	// Guild and forge.
	result = append(result, rectangleCollision("adventurers_guild", "building", -150, 150, 33, 25, 0))
	result = appendFenceCollisions(result, "guild_training_fence", -150, 120, 52, 42, 10, "east")
	for i, p := range [][2]float64{{-168, 111}, {-156, 111}, {-144, 111}, {-132, 111}} {
		result = append(result, circleCollision(fmt.Sprintf("guild_target_%02d", i), "training", p[0], p[1], 1.4))
	}
	result = append(result,
		rectangleCollision("blacksmith_forge", "building", -195, 0, 29, 23, 0),
		circleCollision("forge_chimney", "chimney", -205, 4, 3.4),
		rectangleCollision("forge_fire", "forge", -190, -6, 5.2, 4.2, 0),
	)
	for i := 0; i < 4; i++ {
		result = append(result, rectangleCollision(fmt.Sprintf("forge_anvil_%02d", i), "prop", -210+float64(i)*10, -18, 4.2, 2.7, 0))
	}
	result = appendFenceCollisions(result, "forge_yard_fence", -195, -2, 58, 52, 9, "east")
	for i, p := range [][2]float64{{-218, 17}, {-210, 19}, {-182, 19}, {-174, 13}} {
		result = append(result, rectangleCollision(fmt.Sprintf("forge_crate_%02d", i), "prop", p[0], p[1], 4.2, 4.2, 0))
	}

	// Arena.
	for i := 0; i < 36; i++ {
		a := float64(i) * 2 * math.Pi / 36
		if onRoad(170+math.Cos(a)*47, -165+math.Sin(a)*47, 5.5) {
			continue
		}
		result = append(result, circleCollision(fmt.Sprintf("arena_post_%02d", i), "fence", 170+math.Cos(a)*47, -165+math.Sin(a)*47, .75))
	}
	for i, p := range [][2]float64{{150, -165}, {170, -185}, {190, -165}} {
		result = append(result, circleCollision(fmt.Sprintf("training_dummy_%02d", i), "training", p[0], p[1], 1.1))
	}
	for i, p := range [][2]float64{{133, -195}, {207, -195}, {133, -135}, {207, -135}} {
		result = append(result, rectangleCollision(fmt.Sprintf("arena_watchtower_%02d", i), "tower", p[0], p[1], 6.2, 6.2, 0))
	}
	for i, p := range [][2]float64{{125, -178}, {215, -178}, {128, -151}, {212, -151}} {
		result = append(result, circleCollision(fmt.Sprintf("arena_tent_%02d", i), "tent", p[0], p[1], 4.7))
	}

	// Cemetery.
	result = append(result,
		rectangleCollision("resurrection_shrine", "building", 0, -180, 25, 19, 0),
		circleCollision("resurrection_statue", "statue", 0, -165, 2.5),
	)
	result = appendFenceCollisions(result, "cemetery_fence", 0, -205, 76, 74, 9, "north")
	for row := 0; row < 4; row++ {
		for col := 0; col < 7; col++ {
			result = append(result, rectangleCollision(fmt.Sprintf("grave_%d_%d", row, col), "grave", -30+float64(col)*10, -205-float64(row)*9, 2.6, 1.2, 0))
		}
	}
	for i := 0; i < 6; i++ {
		a := float64(i) * 2 * math.Pi / 6
		result = append(result, circleCollision(
			fmt.Sprintf("cemetery_lamp_%02d", i),
			"lamp",
			math.Cos(a)*33,
			-205+math.Sin(a)*28,
			.5,
		))
	}

	// Farm and orchard.
	result = append(result,
		rectangleCollision("farmhouse", "building", -205, -135, 23, 18, 0),
		circleCollision("windmill", "building", -145, -135, 6.3),
	)
	result = appendFenceCollisions(result, "farm_fence", -175, -169, 108, 92, 12, "east")
	for i, p := range [][2]float64{{-217, -202}, {-210, -200}, {-202, -202}, {-194, -200}, {-217, -192}, {-208, -190}} {
		result = append(result, circleCollision(fmt.Sprintf("hay_bale_%02d", i), "prop", p[0], p[1], 2.5))
	}
	for row := 0; row < 5; row++ {
		for col := 0; col < 4; col++ {
			result = append(result, circleCollision(fmt.Sprintf("orchard_tree_%02d_%02d", row, col), "tree", -153+float64(col)*8, -196+float64(row)*8, 1.2))
		}
	}

	// Wolf forest.
	for i := 0; i < 48; i++ {
		a := float64(i) * 2.399963
		r := 22 + float64((i*17)%95)
		x, y := -172+math.Cos(a)*r, 92+math.Sin(a)*r
		if math.Hypot(x, y) < 100 || math.Hypot(x+150, y-150) < 32 || math.Hypot(x+195, y) < 30 {
			continue
		}
		result = append(result, circleCollision(fmt.Sprintf("wolf_forest_tree_%02d", i), "tree", x, y, .8+float64(i%4)*.12))
	}
	for i, p := range [][2]float64{{-228, 62}, {-220, 112}, {-202, 145}, {-185, 55}, {-133, 78}, {-118, 118}} {
		result = append(result, circleCollision(fmt.Sprintf("wolf_forest_boulder_%02d", i), "rock", p[0], p[1], 4+float64(i%3)))
	}
	result = append(result, rectangleCollision("wolf_den", "cave", -224, 92, 9.5, 4.5, 0))

	// Raven cliffs.
	for i := 0; i < 17; i++ {
		a := float64(i) * 2 * math.Pi / 17
		result = append(result, circleCollision(fmt.Sprintf("raven_cliff_%02d", i), "cliff", math.Cos(a)*35, 230+math.Sin(a)*26, 8+float64(i%3)*3))
	}
	result = append(result, circleCollision("raven_watchtower", "tower", 0, 230, 7.2))
	for i := 0; i < 12; i++ {
		a := float64(i) * 2 * math.Pi / 12
		result = append(result, rectangleCollision(fmt.Sprintf("raven_pillar_%02d", i), "ruin", math.Cos(a)*23, 230+math.Sin(a)*19, 2.4, 2.4, 0))
	}

	// Troll cave and marsh.
	for i := 0; i < 11; i++ {
		a := float64(i) * 2 * math.Pi / 11
		result = append(result, circleCollision(fmt.Sprintf("troll_cave_rock_%02d", i), "rock", 180+math.Cos(a)*28, 198+math.Sin(a)*17, 8+float64(i%3)*3))
	}
	result = append(result, rectangleCollision("troll_cave_entrance", "cave", 180, 180, 16, 3.5, 0))
	for i, p := range [][2]float64{{142, 144}, {151, 162}, {166, 132}, {189, 136}, {204, 153}, {211, 173}, {151, 187}} {
		result = append(result, circleCollision(fmt.Sprintf("marsh_dead_tree_%02d", i), "tree", p[0], p[1], .9))
	}

	// Mana grove.
	for i := 0; i < 18; i++ {
		a := float64(i) * 2 * math.Pi / 18
		result = append(result, circleCollision(fmt.Sprintf("magic_grove_tree_%02d", i), "tree", 210+math.Cos(a)*39, 15+math.Sin(a)*34, 1.1))
	}
	result = append(result, circleCollision("mana_spring", "water", 210, 15, 17.2))
	for i := 0; i < 6; i++ {
		a := float64(i) * 2 * math.Pi / 6
		result = append(result, rectangleCollision(fmt.Sprintf("mana_pillar_%02d", i), "ruin", 210+math.Cos(a)*23, 15+math.Sin(a)*23, 3.2, 3.2, 0))
	}
	for i := 0; i < 14; i++ {
		a := float64(i) * 2 * math.Pi / 14
		r := 12 + float64((i%3)*9)
		result = append(result, circleCollision(fmt.Sprintf("mana_crystal_%02d", i), "crystal", 210+math.Cos(a)*r, 15+math.Sin(a)*r, 1.2))
	}
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		x, y := 210+math.Cos(a)*31, 15+math.Sin(a)*31
		result = append(result,
			rectangleCollision(fmt.Sprintf("mana_ruin_arch_%02d_a", i), "ruin", x, y, 3.4, 3.4, 0),
			rectangleCollision(fmt.Sprintf("mana_ruin_arch_%02d_b", i), "ruin", x+math.Cos(a+math.Pi/2)*6, y+math.Sin(a+math.Pi/2)*6, 3.4, 3.4, 0),
		)
	}
	for i, p := range backgroundTrees() {
		result = append(result, circleCollision(fmt.Sprintf("background_tree_%03d", i), "tree", p.x, p.y, .85*p.z))
	}

	// Les arbres vivants proviennent directement des maillages effectivement
	// émis, afin de ne plus désynchroniser leur position et leur collision.
	if generatedTrees != nil {
		filtered := result[:0]
		for _, c := range result {
			if c.Category != "tree" || strings.HasPrefix(c.ID, "marsh_dead_tree") {
				filtered = append(filtered, c)
			}
		}
		result = append(filtered, generatedTrees...)
	}
	return result
}

func writeCollision(path string) {
	w := newOBJ(path, "red_world_map_3d.mtl", false)
	defer w.close()
	for _, c := range worldCollisions() {
		w.object(c.ID)
		w.material("stone")
		if c.Shape == "circle" {
			w.cylinder(c.X, c.Y, 0, c.Radius, 2, 10)
			continue
		}
		a := c.RotationDegrees * math.Pi / 180
		dx, dy := math.Cos(a)*c.Width/2, math.Sin(a)*c.Width/2
		w.segment(c.X-dx, c.Y-dy, c.X+dx, c.Y+dy, c.Height, 0, 2)
	}
}

func writeMaterials(path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintln(f, "# Materials for the static RED world map")
	for name, c := range materials {
		fmt.Fprintf(f, "\nnewmtl %s\nKa %.3f %.3f %.3f\nKd %.3f %.3f %.3f\nKs 0.08 0.08 0.08\nNs 8\nillum 2\n", name, c[0]*.45, c[1]*.45, c[2]*.45, c[0], c[1], c[2])
		if texturePath, ok := materialTextures[name]; ok {
			fmt.Fprintln(f, "map_Kd", texturePath)
		}
	}
}

func writeLayout(path string) {
	data := layout{
		Name:       "RED World",
		WorldSize:  600,
		NorthAxis:  "+Y",
		UpAxis:     "+Z",
		ImageWidth: 1254,
		Mapping:    "worldX=(pixelX/1254-0.5)*600; worldY=(0.5-pixelY/1254)*600",
		Spawn:      [3]float64{0, -20, 0.2},
		Landmarks: []landmark{
			{"village", "Village fortifié", 0, 0, 82, "Place centrale, maisons et portes"},
			{"market", "Marché", 43, 8, 24, "Potions, grimoire et améliorations d'inventaire"},
			{"guild", "Guilde des aventuriers", -150, 150, 35, "Création, informations et équipement du personnage"},
			{"forge", "Forge", -195, 0, 34, "Fabrication du chapeau, de la tunique et des bottes"},
			{"arena", "Arène d'entraînement", 170, -165, 48, "Combat contre le gobelin d'entraînement"},
			{"cemetery", "Sanctuaire de résurrection", 0, -180, 42, "Retour du personnage après sa mort"},
			{"boar_fields", "Prés aux sangliers", -170, -165, 55, "Cuir de sanglier et cultures"},
			{"wolf_forest", "Forêt des loups", -172, 92, 72, "Fourrure de loup"},
			{"raven_cliffs", "Falaises aux corbeaux", 0, 230, 38, "Plumes de corbeau"},
			{"troll_marsh", "Marais du troll", 175, 165, 55, "Peau de troll et grotte"},
			{"mana_grove", "Bosquet de mana", 210, 15, 44, "Source magique et contenu optionnel de mana"},
		},
		Collisions: worldCollisions(),
		Bridges:    bridgePassages,
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		panic(err)
	}
}
