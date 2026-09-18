package main

import "math"

const materialLibrary = "../../materials/low_poly.mtl"

func generatePalette(c *catalog) error {
	palette := `# Projet RED - palette low-poly
newmtl skin
Kd 0.91 0.66 0.46
newmtl skin_dark
Kd 0.54 0.31 0.18
newmtl hair
Kd 0.16 0.08 0.04
newmtl beard
Kd 0.72 0.28 0.06
newmtl blue
Kd 0.10 0.34 0.82
newmtl green
Kd 0.18 0.62 0.20
newmtl dark_green
Kd 0.08 0.28 0.09
newmtl red
Kd 0.82 0.10 0.10
newmtl purple
Kd 0.44 0.12 0.64
newmtl cyan
Kd 0.08 0.60 0.86
newmtl gold
Kd 0.96 0.66 0.08
newmtl silver
Kd 0.66 0.74 0.82
newmtl iron
Kd 0.25 0.29 0.34
newmtl leather
Kd 0.37 0.17 0.07
newmtl leather_light
Kd 0.61 0.34 0.15
newmtl fur
Kd 0.42 0.36 0.30
newmtl troll
Kd 0.25 0.42 0.22
newmtl raven
Kd 0.04 0.05 0.08
newmtl wood
Kd 0.32 0.16 0.06
newmtl stone
Kd 0.32 0.37 0.42
newmtl grass
Kd 0.16 0.38 0.18
newmtl sand
Kd 0.58 0.47 0.28
newmtl glass
Kd 0.50 0.78 0.86
d 0.65
newmtl fire
Kd 1.00 0.28 0.03
Ke 0.80 0.10 0.00
newmtl magic
Kd 0.15 0.60 1.00
Ke 0.05 0.25 0.80
newmtl heal
Kd 0.20 1.00 0.48
Ke 0.05 0.60 0.18
newmtl poison
Kd 0.52 0.10 0.72
Ke 0.20 0.02 0.30
newmtl white
Kd 0.90 0.92 0.96
newmtl black
Kd 0.02 0.025 0.035
`
	return c.write("materials/low_poly.mtl", "material", "Palette commune de couleurs plates", []byte(palette))
}

type modelSpec struct {
	path        string
	description string
	build       func() *mesh
}

func generateModels(c *catalog) error {
	models := []modelSpec{
		{"models/characters/human.obj", "Personnage humain low-poly", humanModel},
		{"models/characters/elf.obj", "Personnage elfe low-poly", elfModel},
		{"models/characters/dwarf.obj", "Personnage nain low-poly", dwarfModel},
		{"models/characters/training_goblin.obj", "Gobelin d'entrainement low-poly", goblinModel},
		{"models/items/health_potion.obj", "Potion de vie", func() *mesh { return potionModel("red") }},
		{"models/items/poison_potion.obj", "Potion de poison", func() *mesh { return potionModel("poison") }},
		{"models/items/mana_potion.obj", "Potion de mana bonus", func() *mesh { return potionModel("magic") }},
		{"models/items/fireball_spellbook.obj", "Livre de sort Boule de Feu", spellbookModel},
		{"models/items/wolf_fur.obj", "Fourrure de loup", func() *mesh { return hideModel("fur") }},
		{"models/items/troll_hide.obj", "Peau de troll", func() *mesh { return hideModel("troll") }},
		{"models/items/boar_leather.obj", "Cuir de sanglier", func() *mesh { return hideModel("leather_light") }},
		{"models/items/raven_feather.obj", "Plume de corbeau", featherModel},
		{"models/items/gold_coin.obj", "Piece d'or", coinModel},
		{"models/items/inventory_upgrade.obj", "Sac d'augmentation d'inventaire", backpackModel},
		{"models/equipment/adventurer_hat.obj", "Chapeau de l'aventurier", hatModel},
		{"models/equipment/adventurer_tunic.obj", "Tunique de l'aventurier", tunicModel},
		{"models/equipment/adventurer_boots.obj", "Bottes de l'aventurier", bootsModel},
		{"models/environment/merchant_stall.obj", "Etal du marchand", merchantStallModel},
		{"models/environment/forge.obj", "Forge avec four et enclume", forgeModel},
		{"models/environment/combat_arena.obj", "Petite arene de combat", arenaModel},
		{"models/environment/village_square.obj", "Environnement principal de village", villageModel},
		{"models/vfx/punch_impact.obj", "Impact du Coup de poing", func() *mesh { return burstModel("white") }},
		{"models/vfx/basic_attack.obj", "Impact de l'attaque basique", func() *mesh { return burstModel("silver") }},
		{"models/vfx/fireball.obj", "Effet de Boule de Feu", fireballModel},
		{"models/vfx/heal_burst.obj", "Effet de soin", func() *mesh { return auraModel("heal") }},
		{"models/vfx/poison_cloud.obj", "Effet de poison", poisonCloudModel},
		{"models/vfx/level_up.obj", "Effet de montee de niveau", func() *mesh { return auraModel("gold") }},
	}

	for _, model := range models {
		name := model.path
		data := model.build().obj(name, materialLibrary)
		if err := c.write(model.path, "model-obj", model.description, data); err != nil {
			return err
		}
	}
	return nil
}

func humanModel() *mesh {
	m := &mesh{}
	addHumanoid(m, "blue", "skin", 1)
	m.box(vec3{0, 2.43, -0.08}, vec3{0.68, 0.16, 0.58}, "hair")
	return m
}

func elfModel() *mesh {
	m := &mesh{}
	addHumanoid(m, "green", "skin", 1.08)
	m.cone(vec3{-0.53, 2.17, 0}, 0.16, 0.52, 4, "skin")
	m.cone(vec3{0.53, 2.17, 0}, 0.16, 0.52, 4, "skin")
	m.box(vec3{0, 2.47, -0.08}, vec3{0.62, 0.14, 0.54}, "hair")
	return m
}

func dwarfModel() *mesh {
	m := &mesh{}
	addHumanoid(m, "red", "skin", 0.78)
	m.cone(vec3{0, 1.72, 0.28}, 0.44, 0.75, 6, "beard")
	m.cylinder(vec3{0, 2.02, 0}, 0.43, 0.18, 8, "iron")
	return m
}

func goblinModel() *mesh {
	m := &mesh{}
	addHumanoid(m, "dark_green", "green", 0.88)
	m.cone(vec3{-0.55, 1.93, 0}, 0.2, 0.62, 4, "green")
	m.cone(vec3{0.55, 1.93, 0}, 0.2, 0.62, 4, "green")
	m.box(vec3{0, 1.86, 0.38}, vec3{0.34, 0.18, 0.2}, "green")
	return m
}

func addHumanoid(m *mesh, clothing, skin string, scale float64) {
	m.box(vec3{-0.2 * scale, 0.48 * scale, 0}, vec3{0.3 * scale, 0.92 * scale, 0.34 * scale}, "leather")
	m.box(vec3{0.2 * scale, 0.48 * scale, 0}, vec3{0.3 * scale, 0.92 * scale, 0.34 * scale}, "leather")
	m.box(vec3{0, 1.25 * scale, 0}, vec3{0.9 * scale, 0.9 * scale, 0.48 * scale}, clothing)
	m.box(vec3{-0.58 * scale, 1.27 * scale, 0}, vec3{0.25 * scale, 0.86 * scale, 0.28 * scale}, skin)
	m.box(vec3{0.58 * scale, 1.27 * scale, 0}, vec3{0.25 * scale, 0.86 * scale, 0.28 * scale}, skin)
	m.box(vec3{0, 2.02 * scale, 0}, vec3{0.68 * scale, 0.68 * scale, 0.62 * scale}, skin)
}

func potionModel(liquid string) *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, 0.42, 0}, 0.34, 0.64, 8, liquid)
	m.cylinder(vec3{0, 0.82, 0}, 0.17, 0.28, 8, "glass")
	m.cylinder(vec3{0, 1.02, 0}, 0.2, 0.16, 8, "wood")
	return m
}

func spellbookModel() *mesh {
	m := &mesh{}
	m.box(vec3{0, 0.18, 0}, vec3{1.0, 0.24, 1.28}, "red")
	m.box(vec3{-0.48, 0.2, 0}, vec3{0.16, 0.3, 1.34}, "gold")
	m.octahedron(vec3{0.12, 0.34, 0}, 0.22, "fire")
	return m
}

func hideModel(material string) *mesh {
	m := &mesh{}
	m.box(vec3{0, 0.07, 0}, vec3{1.4, 0.14, 1.0}, material)
	for _, point := range []vec3{{-0.62, 0.08, -0.55}, {0, 0.08, -0.62}, {0.62, 0.08, -0.55}, {-0.62, 0.08, 0.55}, {0.62, 0.08, 0.55}} {
		m.cone(point, 0.18, 0.42, 4, material)
	}
	return m
}

func featherModel() *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, 0.72, 0}, 0.035, 1.44, 6, "silver")
	for index := range 5 {
		y := 0.5 + float64(index)*0.22
		width := 0.46 - float64(index)*0.055
		m.polygon("raven", vec3{0, y - 0.12, 0}, vec3{-width, y, 0}, vec3{0, y + 0.18, 0})
		m.polygon("raven", vec3{0, y - 0.12, 0.01}, vec3{width, y, 0.01}, vec3{0, y + 0.18, 0.01})
	}
	return m
}

func coinModel() *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, 0.08, 0}, 0.48, 0.16, 12, "gold")
	m.octahedron(vec3{0, 0.18, 0}, 0.2, "sand")
	return m
}

func backpackModel() *mesh {
	m := &mesh{}
	m.box(vec3{0, 0.65, 0}, vec3{0.95, 1.1, 0.5}, "leather")
	m.cylinder(vec3{0, 1.22, 0}, 0.48, 0.18, 8, "leather_light")
	m.box(vec3{0, 0.64, 0.29}, vec3{0.52, 0.42, 0.16}, "leather_light")
	m.box(vec3{0, 0.64, 0.39}, vec3{0.34, 0.1, 0.08}, "gold")
	return m
}

func hatModel() *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, 0.08, 0}, 0.76, 0.14, 12, "leather")
	m.cone(vec3{0, 0.66, 0}, 0.5, 1.15, 8, "leather_light")
	m.box(vec3{0, 0.38, 0.46}, vec3{0.95, 0.12, 0.08}, "red")
	return m
}

func tunicModel() *mesh {
	m := &mesh{}
	m.box(vec3{0, 0.72, 0}, vec3{1.1, 1.25, 0.55}, "blue")
	m.box(vec3{-0.69, 0.9, 0}, vec3{0.28, 0.82, 0.48}, "blue")
	m.box(vec3{0.69, 0.9, 0}, vec3{0.28, 0.82, 0.48}, "blue")
	m.box(vec3{0, 0.62, 0.31}, vec3{1.12, 0.13, 0.1}, "leather")
	return m
}

func bootsModel() *mesh {
	m := &mesh{}
	for _, x := range []float64{-0.3, 0.3} {
		m.box(vec3{x, 0.44, 0}, vec3{0.42, 0.84, 0.48}, "leather")
		m.box(vec3{x, 0.14, 0.22}, vec3{0.46, 0.26, 0.84}, "leather_light")
	}
	return m
}

func merchantStallModel() *mesh {
	m := &mesh{}
	for _, x := range []float64{-1.7, 1.7} {
		m.box(vec3{x, 1.5, 0}, vec3{0.18, 3, 0.18}, "wood")
	}
	m.box(vec3{0, 1.05, 0}, vec3{3.7, 0.2, 1.2}, "wood")
	m.roof(vec3{0, 2.9, 0}, 4.2, 1.8, 0.65, "red")
	for _, x := range []float64{-1.05, 0, 1.05} {
		m.box(vec3{x, 1.25, -0.1}, vec3{0.8, 0.35, 0.7}, "sand")
	}
	return m
}

func forgeModel() *mesh {
	m := &mesh{}
	m.box(vec3{-0.8, 0.75, 0}, vec3{1.5, 1.5, 1.4}, "stone")
	m.box(vec3{-0.8, 1.9, 0}, vec3{0.6, 1.0, 0.6}, "stone")
	m.box(vec3{-0.8, 0.65, 0.72}, vec3{0.85, 0.7, 0.16}, "fire")
	m.box(vec3{1.05, 0.65, 0}, vec3{1.4, 0.28, 0.52}, "iron")
	m.box(vec3{1.05, 0.32, 0}, vec3{0.38, 0.65, 0.38}, "iron")
	m.cone(vec3{1.7, 0.74, 0}, 0.34, 0.65, 4, "iron")
	return m
}

func arenaModel() *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, -0.12, 0}, 7.2, 0.24, 12, "sand")
	for index := range 12 {
		angle := float64(index) * 2 * 3.141592653589793 / 12
		x, z := 6.8*cos(angle), 6.8*sin(angle)
		m.cylinder(vec3{x, 0.7, z}, 0.28, 1.4, 6, "stone")
	}
	return m
}

func villageModel() *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, -0.12, 0}, 8.5, 0.24, 12, "grass")
	for _, house := range []vec3{{-4.6, 0, -2.8}, {4.6, 0, -2.4}, {-3.7, 0, 4.0}} {
		m.box(vec3{house.x, 1.2, house.z}, vec3{2.6, 2.4, 2.3}, "stone")
		m.roof(vec3{house.x, 2.4, house.z}, 3.1, 2.8, 1.0, "red")
		m.box(vec3{house.x, 0.75, house.z + 1.18}, vec3{0.72, 1.45, 0.12}, "wood")
	}
	for _, tree := range []vec3{{4.7, 0, 4.4}, {1.8, 0, -5.4}, {-6.0, 0, 1.3}} {
		m.cylinder(vec3{tree.x, 1.0, tree.z}, 0.25, 2.0, 7, "wood")
		m.cone(vec3{tree.x, 2.65, tree.z}, 1.15, 2.4, 7, "green")
	}
	return m
}

func burstModel(material string) *mesh {
	m := &mesh{}
	m.octahedron(vec3{0, 0, 0}, 0.32, material)
	for _, point := range []vec3{{0.75, 0, 0}, {-0.75, 0, 0}, {0, 0.75, 0}, {0, -0.75, 0}, {0, 0, 0.75}, {0, 0, -0.75}} {
		m.cone(point, 0.12, 0.7, 5, material)
	}
	return m
}

func fireballModel() *mesh {
	m := &mesh{}
	m.octahedron(vec3{0, 0, 0}, 0.58, "fire")
	for _, point := range []vec3{{0.4, 0.4, 0}, {-0.4, 0.3, 0.1}, {0.1, -0.3, 0.4}, {0, 0.1, -0.45}} {
		m.cone(point, 0.16, 0.75, 5, "gold")
	}
	return m
}

func auraModel(material string) *mesh {
	m := &mesh{}
	m.cylinder(vec3{0, 0.04, 0}, 0.95, 0.08, 12, material)
	for index := range 8 {
		angle := float64(index) * 2 * 3.141592653589793 / 8
		x, z := 0.72*cos(angle), 0.72*sin(angle)
		m.cone(vec3{x, 0.65, z}, 0.1, 1.3, 5, material)
	}
	return m
}

func poisonCloudModel() *mesh {
	m := &mesh{}
	for _, point := range []vec3{{0, 0.25, 0}, {0.45, 0.45, 0.1}, {-0.42, 0.38, 0.2}, {0.1, 0.7, -0.35}, {-0.1, 0.82, 0.32}} {
		m.octahedron(point, 0.42, "poison")
	}
	return m
}

func sin(value float64) float64 {
	return math.Sin(value)
}

func cos(value float64) float64 {
	return math.Cos(value)
}
