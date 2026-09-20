package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

type atlas struct {
	File  string
	Names [16]string
}

type manifest struct {
	Style   string            `json:"style"`
	Atlases []string          `json:"atlases"`
	Icons   map[string]string `json:"icons"`
}

var atlases = []atlas{
	{
		File: "items_equipment.png",
		Names: [16]string{
			"health_potion", "poison_potion", "mana_potion", "fireball_spellbook",
			"wolf_fur", "troll_hide", "boar_leather", "raven_feather",
			"gold_coins", "inventory_backpack", "adventurer_hat", "adventurer_tunic",
			"adventurer_boots", "short_sword", "anvil", "blacksmith_hammer",
		},
	},
	{
		File: "characters_spells.png",
		Names: [16]string{
			"human", "elf", "dwarf", "training_goblin",
			"punch", "fireball", "heal", "poison",
			"level_up", "mana_crystal", "experience_scroll", "initiative",
			"victory", "defeat", "merchant", "blacksmith",
		},
	},
	{
		File: "ui_world.png",
		Names: [16]string{
			"main_menu", "character_creation", "character_info", "inventory_chest",
			"merchant_stall", "forge", "crafting", "combat",
			"village", "forge_furnace", "combat_arena", "mystery_characters",
			"health", "mana", "experience", "gold_pouch",
		},
	},
}

type uiIcon struct {
	Name   string
	Source string
}

var uiIcons = []uiIcon{
	{"human", "human"}, {"elf", "elf"}, {"dwarf", "dwarf"}, {"training_goblin", "training_goblin"},
	{"health_potion", "health_potion"}, {"poison_potion", "poison_potion"}, {"mana_potion", "mana_potion"}, {"fireball_spellbook", "fireball_spellbook"},
	{"wolf_fur", "wolf_fur"}, {"troll_hide", "troll_hide"}, {"boar_leather", "boar_leather"}, {"raven_feather", "raven_feather"},
	{"gold_coin", "gold_coins"}, {"inventory_upgrade", "inventory_backpack"}, {"adventurer_hat", "adventurer_hat"}, {"adventurer_tunic", "adventurer_tunic"},
	{"adventurer_boots", "adventurer_boots"}, {"basic_attack", "short_sword"}, {"punch", "punch"}, {"fireball", "fireball"},
	{"heal", "heal"}, {"poison", "poison"}, {"level_up", "level_up"}, {"mana", "mana"},
	{"experience", "experience"}, {"initiative", "initiative"}, {"victory", "victory"}, {"defeat", "defeat"},
	{"main_menu", "main_menu"}, {"character_creation", "character_creation"}, {"character_info", "character_info"}, {"inventory", "inventory_chest"},
	{"merchant", "merchant"}, {"blacksmith", "blacksmith"}, {"crafting", "crafting"}, {"combat", "combat"},
	{"village", "village"}, {"arena", "combat_arena"}, {"who_are_they", "mystery_characters"},
}

func main() {
	root := filepath.Join("assets", "pixel_art")
	outDir := filepath.Join(root, "icons")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	m := manifest{
		Style: "Detailed chunky fantasy pixel art, warm top-left lighting, transparent background",
		Icons: make(map[string]string),
	}

	for _, spec := range atlases {
		path := filepath.Join(root, "atlases", spec.File)
		img := mustOpenPNG(path)
		m.Atlases = append(m.Atlases, filepath.ToSlash(filepath.Join("atlases", spec.File)))

		bounds := img.Bounds()
		// The artwork rows are not exactly aligned with quarter-height cells.
		// Use the empty gutters rather than clipping objects at an arbitrary grid.
		rows := [5]int{0, 340, 630, 920, bounds.Dy()}
		if spec.File == "ui_world.png" {
			rows = [5]int{0, 345, 640, 950, bounds.Dy()}
		}
		for index, name := range spec.Names {
			col, row := index%4, index/4
			cell := image.Rect(
				bounds.Min.X+col*bounds.Dx()/4,
				bounds.Min.Y+rows[row],
				bounds.Min.X+(col+1)*bounds.Dx()/4,
				bounds.Min.Y+rows[row+1],
			)
			icon := extract(img, cell)
			filename := name + ".png"
			mustWritePNG(filepath.Join(outDir, filename), icon)
			m.Icons[name] = filepath.ToSlash(filepath.Join("icons", filename))
		}
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), data, 0o644); err != nil {
		panic(err)
	}
	generateUIAssets(filepath.Dir(root), outDir)
	fmt.Printf("Generated %d source icons and refreshed %d UI icons\n", len(m.Icons), len(uiIcons))
}

func generateUIAssets(assetsRoot, sourceDir string) {
	pngDir := filepath.Join(assetsRoot, "ui", "icons", "png")
	svgDir := filepath.Join(assetsRoot, "ui", "icons", "svg")
	previewPath := filepath.Join(assetsRoot, "previews", "icon_sheet.png")
	for _, directory := range []string{pngDir, svgDir, filepath.Dir(previewPath)} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			panic(err)
		}
	}

	const iconSize = 256
	preview := image.NewNRGBA(image.Rect(0, 0, 1280, 2048))
	draw.Draw(preview, preview.Bounds(), &image.Uniform{C: color.NRGBA{R: 8, G: 17, B: 31, A: 255}}, image.Point{}, draw.Src)

	for index, spec := range uiIcons {
		source := mustOpenPNG(filepath.Join(sourceDir, spec.Source+".png"))
		icon := fitTransparent(source, iconSize, 18)
		pngPath := filepath.Join(pngDir, spec.Name+".png")
		mustWritePNG(pngPath, icon)
		mustWriteEmbeddedSVG(filepath.Join(svgDir, spec.Name+".svg"), icon)

		column, row := index%5, index/5
		cell := image.Rect(column*256+12, row*256+12, (column+1)*256-12, (row+1)*256-12)
		draw.Draw(preview, cell, &image.Uniform{C: color.NRGBA{R: 20, G: 34, B: 52, A: 255}}, image.Point{}, draw.Src)
		drawBorder(preview, cell, 4, color.NRGBA{R: 51, G: 70, B: 92, A: 255})
		drawIconCentered(preview, cell, icon, 194)
	}
	mustWritePNG(previewPath, preview)
}

func fitTransparent(src image.Image, size, padding int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	available := size - 2*padding
	b := src.Bounds()
	dw, dh := available, available
	if b.Dx() > b.Dy() {
		dh = b.Dy() * available / b.Dx()
	} else {
		dw = b.Dx() * available / b.Dy()
	}
	target := image.Rect((size-dw)/2, (size-dh)/2, (size+dw)/2, (size+dh)/2)
	drawNearest(dst, target, src)
	return dst
}

func drawIconCentered(dst *image.NRGBA, cell image.Rectangle, src image.Image, maxSize int) {
	b := src.Bounds()
	dw, dh := maxSize, maxSize
	if b.Dx() > b.Dy() {
		dh = b.Dy() * maxSize / b.Dx()
	} else {
		dw = b.Dx() * maxSize / b.Dy()
	}
	x := cell.Min.X + (cell.Dx()-dw)/2
	y := cell.Min.Y + (cell.Dy()-dh)/2
	drawNearest(dst, image.Rect(x, y, x+dw, y+dh), src)
}

func drawNearest(dst *image.NRGBA, target image.Rectangle, src image.Image) {
	sb := src.Bounds()
	for y := target.Min.Y; y < target.Max.Y; y++ {
		sy := sb.Min.Y + (y-target.Min.Y)*sb.Dy()/target.Dy()
		for x := target.Min.X; x < target.Max.X; x++ {
			sx := sb.Min.X + (x-target.Min.X)*sb.Dx()/target.Dx()
			pixel := color.NRGBAModel.Convert(src.At(sx, sy)).(color.NRGBA)
			if pixel.A != 0 {
				dst.SetNRGBA(x, y, pixel)
			}
		}
	}
}

func drawBorder(dst *image.NRGBA, rect image.Rectangle, width int, c color.NRGBA) {
	draw.Draw(dst, image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y+width), &image.Uniform{C: c}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(rect.Min.X, rect.Max.Y-width, rect.Max.X, rect.Max.Y), &image.Uniform{C: c}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+width, rect.Max.Y), &image.Uniform{C: c}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(rect.Max.X-width, rect.Min.Y, rect.Max.X, rect.Max.Y), &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func mustWriteEmbeddedSVG(path string, icon image.Image) {
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, icon); err != nil {
		panic(err)
	}
	encoded := base64.StdEncoding.EncodeToString(pngData.Bytes())
	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" viewBox="0 0 256 256">
  <image width="256" height="256" image-rendering="pixelated" href="data:image/png;base64,%s"/>
</svg>
`, encoded)
	if err := os.WriteFile(path, []byte(svg), 0o644); err != nil {
		panic(err)
	}
}

func mustOpenPNG(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}

func extract(src image.Image, cell image.Rectangle) *image.NRGBA {
	w, h := cell.Dx(), cell.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.Set(x, y, src.At(cell.Min.X+x, cell.Min.Y+y))
		}
	}
	// These corrected source sheets have real alpha. Never flood-fill pale
	// artwork or erase a border: both operations destroyed parts of the sprites.
	removeTinyFragments(out)
	return trimWithPadding(out, 10)
}

func removeTinyFragments(img *image.NRGBA) {
	b := img.Bounds()
	seen := make([]bool, b.Dx()*b.Dy())
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			if seen[y*b.Dx()+x] || img.NRGBAAt(x, y).A < 32 {
				continue
			}
			points := []image.Point{image.Pt(x, y)}
			seen[y*b.Dx()+x] = true
			for i := 0; i < len(points); i++ {
				p := points[i]
				for _, n := range []image.Point{image.Pt(p.X-1, p.Y), image.Pt(p.X+1, p.Y), image.Pt(p.X, p.Y-1), image.Pt(p.X, p.Y+1)} {
					if !n.In(b) || seen[n.Y*b.Dx()+n.X] || img.NRGBAAt(n.X, n.Y).A < 32 {
						continue
					}
					seen[n.Y*b.Dx()+n.X] = true
					points = append(points, n)
				}
			}
			if len(points) < 80 {
				for _, p := range points {
					img.SetNRGBA(p.X, p.Y, color.NRGBA{})
				}
			}
		}
	}
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			if img.NRGBAAt(x, y).A < 32 {
				img.SetNRGBA(x, y, color.NRGBA{})
			}
		}
	}
}

func clearCellBorder(img *image.NRGBA, width int) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if x-b.Min.X < width || b.Max.X-1-x < width || y-b.Min.Y < width || b.Max.Y-1-y < width {
				img.SetNRGBA(x, y, color.NRGBA{})
			}
		}
	}
}

// The image generator may render its transparency preview as a pale checkerboard.
// Only pale neutral pixels connected to a cell edge are removed, preserving white
// highlights enclosed by an icon's dark outline.
func removeConnectedCheckerboard(img *image.NRGBA) {
	b := img.Bounds()
	seen := make([]bool, b.Dx()*b.Dy())
	queue := make([]image.Point, 0, 2*(b.Dx()+b.Dy()))
	push := func(x, y int) {
		if x < 0 || y < 0 || x >= b.Dx() || y >= b.Dy() {
			return
		}
		i := y*b.Dx() + x
		if seen[i] || !isChecker(img.NRGBAAt(x, y)) {
			return
		}
		seen[i] = true
		queue = append(queue, image.Pt(x, y))
	}
	for x := 0; x < b.Dx(); x++ {
		push(x, 0)
		push(x, b.Dy()-1)
	}
	for y := 0; y < b.Dy(); y++ {
		push(0, y)
		push(b.Dx()-1, y)
	}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		img.SetNRGBA(p.X, p.Y, color.NRGBA{})
		push(p.X-1, p.Y)
		push(p.X+1, p.Y)
		push(p.X, p.Y-1)
		push(p.X, p.Y+1)
	}
}

func isChecker(c color.NRGBA) bool {
	min := c.R
	max := c.R
	for _, v := range []uint8{c.G, c.B} {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min >= 185 && max-min <= 22
}

func trimWithPadding(src *image.NRGBA, padding int) *image.NRGBA {
	b := src.Bounds()
	minX, minY, maxX, maxY := b.Max.X, b.Max.Y, b.Min.X-1, b.Min.Y-1
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if src.NRGBAAt(x, y).A == 0 {
				continue
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	if maxX < minX || maxY < minY {
		return image.NewNRGBA(image.Rect(0, 0, 1, 1))
	}
	out := image.NewNRGBA(image.Rect(0, 0, maxX-minX+1+2*padding, maxY-minY+1+2*padding))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			out.SetNRGBA(x-minX+padding, y-minY+padding, src.NRGBAAt(x, y))
		}
	}
	return out
}

func mustWritePNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}
