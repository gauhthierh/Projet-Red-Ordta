package main

import (
	"bytes"
	"fmt"
	"html"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type iconSpec struct {
	name        string
	label       string
	description string
	shape       string
	accent      string
	accentRGBA  color.RGBA
}

func generateIcons(c *catalog) error {
	icons := []iconSpec{
		{"human", "HU", "Classe Humain", "character", "#3b82f6", color.RGBA{59, 130, 246, 255}},
		{"elf", "EL", "Classe Elfe", "character", "#49a942", color.RGBA{73, 169, 66, 255}},
		{"dwarf", "NA", "Classe Nain", "character", "#d26932", color.RGBA{210, 105, 50, 255}},
		{"training_goblin", "GO", "Gobelin d'entrainement", "character", "#67b23f", color.RGBA{103, 178, 63, 255}},
		{"health_potion", "PV", "Potion de vie", "potion", "#e53935", color.RGBA{229, 57, 53, 255}},
		{"poison_potion", "PO", "Potion de poison", "potion", "#8e44ad", color.RGBA{142, 68, 173, 255}},
		{"mana_potion", "PM", "Potion de mana", "potion", "#2196f3", color.RGBA{33, 150, 243, 255}},
		{"fireball_spellbook", "LF", "Livre de sort Boule de Feu", "book", "#f4511e", color.RGBA{244, 81, 30, 255}},
		{"wolf_fur", "FL", "Fourrure de loup", "material", "#8d8277", color.RGBA{141, 130, 119, 255}},
		{"troll_hide", "PT", "Peau de troll", "material", "#597c4d", color.RGBA{89, 124, 77, 255}},
		{"boar_leather", "CS", "Cuir de sanglier", "material", "#9b5d2e", color.RGBA{155, 93, 46, 255}},
		{"raven_feather", "PC", "Plume de corbeau", "feather", "#596275", color.RGBA{89, 98, 117, 255}},
		{"gold_coin", "OR", "Piece d'or", "coin", "#f2b705", color.RGBA{242, 183, 5, 255}},
		{"inventory_upgrade", "+10", "Augmentation d'inventaire", "bag", "#b77935", color.RGBA{183, 121, 53, 255}},
		{"adventurer_hat", "CH", "Chapeau de l'aventurier", "equipment", "#a16233", color.RGBA{161, 98, 51, 255}},
		{"adventurer_tunic", "TU", "Tunique de l'aventurier", "equipment", "#3578c7", color.RGBA{53, 120, 199, 255}},
		{"adventurer_boots", "BO", "Bottes de l'aventurier", "equipment", "#79502d", color.RGBA{121, 80, 45, 255}},
		{"punch", "CP", "Sort Coup de poing", "spell", "#e5e7eb", color.RGBA{229, 231, 235, 255}},
		{"fireball", "BF", "Sort Boule de Feu", "spell", "#ff6b1a", color.RGBA{255, 107, 26, 255}},
		{"heal", "SO", "Effet de soin", "spell", "#2ecc71", color.RGBA{46, 204, 113, 255}},
		{"poison", "PS", "Effet de poison", "spell", "#a855f7", color.RGBA{168, 85, 247, 255}},
		{"basic_attack", "AT", "Attaque basique", "spell", "#cbd5e1", color.RGBA{203, 213, 225, 255}},
		{"level_up", "LV", "Montee de niveau", "spell", "#facc15", color.RGBA{250, 204, 21, 255}},
		{"mana", "MN", "Ressource de mana", "stat", "#38bdf8", color.RGBA{56, 189, 248, 255}},
		{"experience", "XP", "Experience", "stat", "#eab308", color.RGBA{234, 179, 8, 255}},
		{"initiative", "IN", "Initiative", "stat", "#f97316", color.RGBA{249, 115, 22, 255}},
		{"main_menu", "MM", "Menu principal", "ui", "#60a5fa", color.RGBA{96, 165, 250, 255}},
		{"character_creation", "CC", "Creation de personnage", "ui", "#34d399", color.RGBA{52, 211, 153, 255}},
		{"character_info", "IF", "Informations du personnage", "ui", "#22d3ee", color.RGBA{34, 211, 238, 255}},
		{"inventory", "IV", "Inventaire", "bag", "#d6a25e", color.RGBA{214, 162, 94, 255}},
		{"merchant", "MA", "Marchand", "building", "#fb923c", color.RGBA{251, 146, 60, 255}},
		{"blacksmith", "FO", "Forgeron", "building", "#94a3b8", color.RGBA{148, 163, 184, 255}},
		{"crafting", "FA", "Fabrication", "building", "#f59e0b", color.RGBA{245, 158, 11, 255}},
		{"combat", "CO", "Combat", "spell", "#ef4444", color.RGBA{239, 68, 68, 255}},
		{"victory", "VI", "Victoire", "stat", "#facc15", color.RGBA{250, 204, 21, 255}},
		{"defeat", "DE", "Defaite", "stat", "#64748b", color.RGBA{100, 116, 139, 255}},
		{"who_are_they", "QA", "Qui sont-ils ?", "ui", "#c084fc", color.RGBA{192, 132, 252, 255}},
		{"village", "VL", "Village", "building", "#65a30d", color.RGBA{101, 163, 13, 255}},
		{"arena", "AR", "Arene", "building", "#dc6040", color.RGBA{220, 96, 64, 255}},
	}

	for _, icon := range icons {
		pngData, err := renderIconPNG(icon)
		if err != nil {
			return err
		}
		if err := c.write("ui/icons/png/"+icon.name+".png", "icon-png", icon.description, pngData); err != nil {
			return err
		}
		svgData := renderIconSVG(icon)
		if err := c.write("ui/icons/svg/"+icon.name+".svg", "icon-svg", icon.description+" (editable)", []byte(svgData)); err != nil {
			return err
		}
	}
	preview, err := renderIconSheet(icons)
	if err != nil {
		return err
	}
	if err := c.write("previews/icon_sheet.png", "preview-png", "Planche de controle des icones", preview); err != nil {
		return err
	}
	return nil
}

func renderIconSheet(icons []iconSpec) ([]byte, error) {
	const columns = 5
	rows := (len(icons) + columns - 1) / columns
	canvas := image.NewRGBA(image.Rect(0, 0, columns*256, rows*256))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.RGBA{8, 13, 22, 255}}, image.Point{}, draw.Src)
	for index, icon := range icons {
		data, err := renderIconPNG(icon)
		if err != nil {
			return nil, err
		}
		decoded, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		x := (index % columns) * 256
		y := (index / columns) * 256
		draw.Draw(canvas, image.Rect(x, y, x+256, y+256), decoded, image.Point{}, draw.Src)
	}
	var output bytes.Buffer
	if err := png.Encode(&output, canvas); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func renderIconPNG(spec iconSpec) ([]byte, error) {
	canvas := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.RGBA{14, 22, 36, 255}}, image.Point{}, draw.Src)
	drawRect(canvas, 12, 12, 244, 244, color.RGBA{40, 55, 74, 255})
	drawRect(canvas, 18, 18, 238, 238, color.RGBA{20, 31, 48, 255})
	drawPNGShape(canvas, spec.shape, spec.accentRGBA)

	label := spec.label
	width := font.MeasureString(basicfont.Face7x13, label).Ceil()
	drawer := font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(color.RGBA{245, 247, 250, 255}),
		Face: basicfont.Face7x13,
		Dot:  fixed.P((256-width)/2, 226),
	}
	drawer.DrawString(label)

	var output bytes.Buffer
	if err := png.Encode(&output, canvas); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func drawPNGShape(canvas *image.RGBA, shape string, accent color.RGBA) {
	switch shape {
	case "character":
		drawCircle(canvas, 128, 82, 38, accent)
		drawPolygon(canvas, []image.Point{{64, 182}, {82, 116}, {174, 116}, {192, 182}}, accent)
	case "potion":
		drawRect(canvas, 105, 48, 151, 78, color.RGBA{202, 213, 225, 255})
		drawPolygon(canvas, []image.Point{{92, 78}, {164, 78}, {190, 174}, {166, 196}, {90, 196}, {66, 174}}, accent)
	case "book":
		drawRect(canvas, 66, 54, 190, 190, accent)
		drawRect(canvas, 76, 54, 91, 190, color.RGBA{250, 204, 21, 255})
		drawDiamond(canvas, 140, 122, 28, color.RGBA{255, 218, 92, 255})
	case "feather":
		drawPolygon(canvas, []image.Point{{78, 177}, {105, 76}, {178, 48}, {154, 125}}, accent)
		drawLine(canvas, 81, 188, 165, 62, 6, color.RGBA{220, 226, 235, 255})
	case "coin":
		drawCircle(canvas, 128, 124, 70, accent)
		drawCircle(canvas, 128, 124, 48, color.RGBA{255, 218, 80, 255})
	case "bag":
		drawRect(canvas, 72, 89, 184, 190, accent)
		drawRect(canvas, 91, 59, 165, 101, color.RGBA{112, 72, 34, 255})
		drawRect(canvas, 111, 126, 145, 148, color.RGBA{250, 204, 21, 255})
	case "building":
		drawRect(canvas, 76, 102, 180, 190, accent)
		drawPolygon(canvas, []image.Point{{58, 106}, {128, 48}, {198, 106}}, color.RGBA{190, 83, 51, 255})
		drawRect(canvas, 112, 141, 144, 190, color.RGBA{67, 43, 27, 255})
	case "equipment":
		drawPolygon(canvas, []image.Point{{128, 46}, {190, 78}, {174, 166}, {128, 199}, {82, 166}, {66, 78}}, accent)
	case "spell":
		drawDiamond(canvas, 128, 124, 72, accent)
		drawDiamond(canvas, 128, 124, 36, color.RGBA{255, 245, 210, 255})
	case "material":
		drawPolygon(canvas, []image.Point{{61, 91}, {94, 51}, {137, 71}, {174, 55}, {198, 105}, {175, 181}, {128, 196}, {73, 174}}, accent)
	case "stat":
		drawDiamond(canvas, 128, 116, 72, accent)
		drawCircle(canvas, 128, 116, 24, color.RGBA{255, 250, 225, 255})
	default:
		drawRect(canvas, 58, 62, 198, 184, accent)
		drawRect(canvas, 74, 82, 182, 164, color.RGBA{24, 36, 54, 255})
	}
}

func renderIconSVG(spec iconSpec) string {
	shape := svgShape(spec.shape, spec.accent)
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" viewBox="0 0 256 256">
  <rect width="256" height="256" rx="28" fill="#0e1624"/>
  <rect x="12" y="12" width="232" height="232" rx="22" fill="#142030" stroke="#33475f" stroke-width="6"/>
  %s
  <text x="128" y="226" text-anchor="middle" fill="#f5f7fa" font-family="sans-serif" font-size="24" font-weight="700">%s</text>
</svg>
`, shape, html.EscapeString(spec.label))
}

func svgShape(shape, accent string) string {
	switch shape {
	case "character":
		return fmt.Sprintf(`<circle cx="128" cy="82" r="38" fill="%s"/><path d="M64 182L82 116H174L192 182Z" fill="%s"/>`, accent, accent)
	case "potion":
		return fmt.Sprintf(`<rect x="105" y="48" width="46" height="30" fill="#cbd5e1"/><path d="M92 78H164L190 174L166 196H90L66 174Z" fill="%s"/>`, accent)
	case "book":
		return fmt.Sprintf(`<rect x="66" y="54" width="124" height="136" rx="8" fill="%s"/><rect x="76" y="54" width="15" height="136" fill="#facc15"/><path d="M140 94L168 122L140 150L112 122Z" fill="#ffda5c"/>`, accent)
	case "feather":
		return fmt.Sprintf(`<path d="M78 177L105 76L178 48L154 125Z" fill="%s"/><path d="M81 188L165 62" stroke="#dce2eb" stroke-width="7"/>`, accent)
	case "coin":
		return fmt.Sprintf(`<circle cx="128" cy="124" r="70" fill="%s"/><circle cx="128" cy="124" r="48" fill="#ffda50"/>`, accent)
	case "bag":
		return fmt.Sprintf(`<rect x="72" y="89" width="112" height="101" rx="18" fill="%s"/><path d="M91 101V70Q128 42 165 70V101" fill="none" stroke="#704822" stroke-width="14"/><rect x="111" y="126" width="34" height="22" fill="#facc15"/>`, accent)
	case "building":
		return fmt.Sprintf(`<rect x="76" y="102" width="104" height="88" fill="%s"/><path d="M58 106L128 48L198 106Z" fill="#be5333"/><rect x="112" y="141" width="32" height="49" fill="#432b1b"/>`, accent)
	case "equipment":
		return fmt.Sprintf(`<path d="M128 46L190 78L174 166L128 199L82 166L66 78Z" fill="%s" stroke="#f8e7b0" stroke-width="6"/>`, accent)
	case "spell", "stat":
		return fmt.Sprintf(`<path d="M128 44L200 124L128 196L56 124Z" fill="%s"/><circle cx="128" cy="124" r="26" fill="#fff5d2"/>`, accent)
	case "material":
		return fmt.Sprintf(`<path d="M61 91L94 51L137 71L174 55L198 105L175 181L128 196L73 174Z" fill="%s"/>`, accent)
	default:
		return fmt.Sprintf(`<rect x="58" y="62" width="140" height="122" rx="16" fill="%s"/><rect x="74" y="82" width="108" height="82" rx="9" fill="#182436"/>`, accent)
	}
}

func drawRect(canvas *image.RGBA, x0, y0, x1, y1 int, fill color.RGBA) {
	draw.Draw(canvas, image.Rect(x0, y0, x1, y1), &image.Uniform{C: fill}, image.Point{}, draw.Src)
}

func drawCircle(canvas *image.RGBA, cx, cy, radius int, fill color.RGBA) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				canvas.SetRGBA(cx+x, cy+y, fill)
			}
		}
	}
}

func drawDiamond(canvas *image.RGBA, cx, cy, radius int, fill color.RGBA) {
	drawPolygon(canvas, []image.Point{{cx, cy - radius}, {cx + radius, cy}, {cx, cy + radius}, {cx - radius, cy}}, fill)
}

func drawPolygon(canvas *image.RGBA, points []image.Point, fill color.RGBA) {
	minY, maxY := points[0].Y, points[0].Y
	for _, point := range points[1:] {
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	for y := minY; y <= maxY; y++ {
		intersections := make([]int, 0, len(points))
		for index, current := range points {
			next := points[(index+1)%len(points)]
			if (current.Y <= y && next.Y > y) || (next.Y <= y && current.Y > y) {
				x := current.X + (y-current.Y)*(next.X-current.X)/(next.Y-current.Y)
				intersections = append(intersections, x)
			}
		}
		if len(intersections) >= 2 {
			if intersections[0] > intersections[1] {
				intersections[0], intersections[1] = intersections[1], intersections[0]
			}
			for x := intersections[0]; x <= intersections[1]; x++ {
				canvas.SetRGBA(x, y, fill)
			}
		}
	}
}

func drawLine(canvas *image.RGBA, x0, y0, x1, y1, width int, fill color.RGBA) {
	steps := max(abs(x1-x0), abs(y1-y0))
	for step := 0; step <= steps; step++ {
		t := float64(step) / float64(steps)
		x := int(float64(x0) + float64(x1-x0)*t)
		y := int(float64(y0) + float64(y1-y0)*t)
		drawCircle(canvas, x, y, width/2, fill)
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
