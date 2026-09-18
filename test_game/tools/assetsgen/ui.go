package main

import (
	"fmt"
	"html"
)

type screenSpec struct {
	name        string
	title       string
	description string
	accent      string
	mode        string
}

func generateUI(c *catalog) error {
	screens := []screenSpec{
		{"main_menu", "PROJET RED", "Ecran du menu principal", "#3b82f6", "menu"},
		{"character_creation", "CREATION DU PERSONNAGE", "Ecran de creation du personnage", "#22c55e", "character"},
		{"character_info", "INFORMATIONS", "Ecran d'informations du personnage", "#06b6d4", "character"},
		{"inventory", "INVENTAIRE", "Ecran d'inventaire", "#d6a25e", "grid"},
		{"merchant", "MARCHAND", "Ecran du marchand", "#f97316", "shop"},
		{"blacksmith", "FORGERON", "Ecran du forgeron", "#94a3b8", "shop"},
		{"crafting", "FABRICATION", "Ecran de fabrication", "#f59e0b", "craft"},
		{"combat", "COMBAT D'ENTRAINEMENT", "Interface de combat", "#ef4444", "combat"},
		{"victory", "VICTOIRE", "Ecran de victoire", "#facc15", "result"},
		{"defeat", "DEFAITE", "Ecran de defaite", "#64748b", "result"},
		{"who_are_they", "QUI SONT-ILS ?", "Ecran des artistes caches", "#c084fc", "artists"},
	}
	for _, screen := range screens {
		if err := c.write("ui/screens/"+screen.name+".svg", "ui-screen-svg", screen.description, []byte(screenSVG(screen))); err != nil {
			return err
		}
	}

	components := []struct {
		name, description, data string
	}{
		{"panel", "Panneau d'interface", panelSVG()},
		{"button_primary", "Bouton principal", buttonSVG("ACTION", "#3b82f6")},
		{"button_danger", "Bouton de danger", buttonSVG("QUITTER", "#b91c1c")},
		{"inventory_slot", "Case d'inventaire", slotSVG()},
		{"health_bar", "Barre de vie", barSVG("PV", "#ef4444", 0.72)},
		{"mana_bar", "Barre de mana", barSVG("MANA", "#38bdf8", 0.58)},
		{"experience_bar", "Barre d'experience", barSVG("XP", "#eab308", 0.43)},
	}
	for _, component := range components {
		if err := c.write("ui/components/"+component.name+".svg", "ui-component-svg", component.description, []byte(component.data)); err != nil {
			return err
		}
	}

	backgrounds := []struct {
		name, description, data string
	}{
		{"village", "Arriere-plan du village", villageBackgroundSVG()},
		{"forge", "Arriere-plan de la forge", forgeBackgroundSVG()},
		{"arena", "Arriere-plan de l'arene", arenaBackgroundSVG()},
	}
	for _, background := range backgrounds {
		if err := c.write("backgrounds/"+background.name+".svg", "background-svg", background.description, []byte(background.data)); err != nil {
			return err
		}
	}
	return nil
}

func screenSVG(spec screenSpec) string {
	content := genericCards(spec.accent)
	switch spec.mode {
	case "menu":
		content = `<g transform="translate(440 230)">
    <rect width="400" height="72" rx="14" fill="#263a56" stroke="#5d7ba4" stroke-width="3"/><text x="200" y="46" text-anchor="middle" class="button">NOUVELLE PARTIE</text>
    <rect y="94" width="400" height="72" rx="14" fill="#263a56" stroke="#5d7ba4" stroke-width="3"/><text x="200" y="140" text-anchor="middle" class="button">CONTINUER</text>
    <rect y="188" width="400" height="72" rx="14" fill="#3d2027" stroke="#8e4250" stroke-width="3"/><text x="200" y="234" text-anchor="middle" class="button">QUITTER</text>
  </g>`
	case "grid":
		content = inventoryGrid(spec.accent)
	case "combat":
		content = combatLayout()
	case "result":
		content = fmt.Sprintf(`<path d="M640 210L735 315L640 420L545 315Z" fill="%s" opacity=".85"/><rect x="440" y="485" width="400" height="76" rx="16" fill="#263a56"/><text x="640" y="533" text-anchor="middle" class="button">CONTINUER</text>`, spec.accent)
	case "artists":
		content = `<g fill="#283850" stroke="#7b5ca2" stroke-width="4"><circle cx="470" cy="300" r="78"/><path d="M350 520Q350 385 470 385Q590 385 590 520Z"/><circle cx="810" cy="300" r="78"/><path d="M690 520Q690 385 810 385Q930 385 930 520Z"/></g><text x="470" y="565" text-anchor="middle" class="small">ARTISTE 1</text><text x="810" y="565" text-anchor="middle" class="small">ARTISTE 2</text>`
	}
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="1280" height="720" viewBox="0 0 1280 720">
  <defs><style>.title{font:700 42px sans-serif;fill:#f8fafc;letter-spacing:3px}.button{font:700 22px sans-serif;fill:#f8fafc}.small{font:600 18px sans-serif;fill:#cbd5e1}.label{font:700 16px sans-serif;fill:#94a3b8}</style></defs>
  <rect width="1280" height="720" fill="#0b1220"/>
  <path d="M0 510L220 350L390 470L610 260L830 450L1050 300L1280 470V720H0Z" fill="#13253a"/>
  <path d="M0 610L280 470L520 610L780 420L1030 560L1280 430V720H0Z" fill="#173148"/>
  <rect x="34" y="28" width="1212" height="664" rx="22" fill="#101b2c" fill-opacity=".9" stroke="%s" stroke-width="4"/>
  <text x="640" y="102" text-anchor="middle" class="title">%s</text>
  <path d="M440 124H840" stroke="%s" stroke-width="5"/>
  %s
</svg>
`, spec.accent, html.EscapeString(spec.title), spec.accent, content)
}

func genericCards(accent string) string {
	return fmt.Sprintf(`<g transform="translate(110 180)">
    <rect width="470" height="410" rx="18" fill="#17263b" stroke="#334b68" stroke-width="3"/>
    <circle cx="235" cy="135" r="88" fill="%s" opacity=".7"/>
    <path d="M125 320Q135 230 235 230Q335 230 345 320Z" fill="#283d59"/>
    <rect x="85" y="345" width="300" height="24" rx="12" fill="#263a56"/>
  </g>
  <g transform="translate(635 180)" fill="#17263b" stroke="#334b68" stroke-width="3">
    <rect width="535" height="78" rx="14"/><rect y="96" width="535" height="78" rx="14"/><rect y="192" width="535" height="78" rx="14"/><rect y="288" width="535" height="122" rx="14"/>
  </g>`, accent)
}

func inventoryGrid(accent string) string {
	output := `<g transform="translate(120 175)"><rect width="730" height="430" rx="18" fill="#17263b" stroke="#334b68" stroke-width="3"/>`
	for row := range 3 {
		for column := range 5 {
			x, y := 35+column*134, 32+row*132
			output += fmt.Sprintf(`<rect x="%d" y="%d" width="106" height="106" rx="12" fill="#101a2a" stroke="%s" stroke-width="3"/>`, x, y, accent)
		}
	}
	return output + `</g><g transform="translate(900 175)"><rect width="260" height="430" rx="18" fill="#17263b" stroke="#334b68" stroke-width="3"/><circle cx="130" cy="115" r="68" fill="#2c4666"/><rect x="35" y="220" width="190" height="28" rx="10" fill="#263a56"/><rect x="35" y="270" width="190" height="28" rx="10" fill="#263a56"/><rect x="35" y="340" width="190" height="56" rx="12" fill="#31527b"/></g>`
}

func combatLayout() string {
	return `<g><path d="M130 520L315 225L500 520Z" fill="#285ca8"/><circle cx="315" cy="230" r="70" fill="#e5a978"/><path d="M780 520L965 225L1150 520Z" fill="#2f792e"/><circle cx="965" cy="230" r="76" fill="#69ad52"/></g>
  <g><rect x="115" y="575" width="390" height="32" rx="16" fill="#321c26"/><rect x="115" y="575" width="310" height="32" rx="16" fill="#ef4444"/><rect x="775" y="575" width="390" height="32" rx="16" fill="#321c26"/><rect x="775" y="575" width="205" height="32" rx="16" fill="#ef4444"/></g>
  <g transform="translate(480 490)"><rect width="320" height="62" rx="14" fill="#263a56"/><text x="160" y="40" text-anchor="middle" class="button">ATTAQUER</text><rect y="78" width="320" height="62" rx="14" fill="#263a56"/><text x="160" y="118" text-anchor="middle" class="button">INVENTAIRE</text></g>`
}

func panelSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="720" height="420"><rect x="8" y="8" width="704" height="404" rx="24" fill="#111d2f" fill-opacity=".94" stroke="#52749d" stroke-width="8"/><path d="M36 54H684" stroke="#2b4665" stroke-width="3"/></svg>`
}

func buttonSVG(label, accent string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="360" height="96"><rect x="5" y="5" width="350" height="86" rx="18" fill="%s" stroke="#dbeafe" stroke-width="5"/><text x="180" y="60" text-anchor="middle" fill="#fff" font-family="sans-serif" font-size="28" font-weight="700">%s</text></svg>`, accent, label)
}

func slotSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="128" height="128"><rect x="6" y="6" width="116" height="116" rx="15" fill="#111d2f" stroke="#52749d" stroke-width="6"/><path d="M22 100L64 30L106 100Z" fill="#243a55" opacity=".6"/></svg>`
}

func barSVG(label, accent string, ratio float64) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="520" height="72"><rect x="4" y="4" width="512" height="64" rx="18" fill="#241822" stroke="#dbeafe" stroke-width="4"/><rect x="12" y="12" width="%.0f" height="48" rx="12" fill="%s"/><text x="260" y="48" text-anchor="middle" fill="#fff" font-family="sans-serif" font-size="25" font-weight="700">%s</text></svg>`, 496*ratio, accent, label)
}

func villageBackgroundSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="1280" height="720"><rect width="1280" height="720" fill="#87bfe5"/><path d="M0 390L210 180L380 370L610 120L840 350L1040 190L1280 390V720H0Z" fill="#6e8a73"/><path d="M0 470L250 340L500 470L770 310L1030 460L1280 330V720H0Z" fill="#4d7548"/><g fill="#d0b17a" stroke="#755333" stroke-width="8"><rect x="130" y="370" width="260" height="190"/><rect x="820" y="350" width="300" height="210"/></g><g fill="#a34c39"><path d="M90 370L260 240L430 370Z"/><path d="M775 350L970 210L1165 350Z"/></g><path d="M0 650Q640 480 1280 650V720H0Z" fill="#b99964"/></svg>`
}

func forgeBackgroundSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="1280" height="720"><rect width="1280" height="720" fill="#171a20"/><path d="M0 560L220 410L440 540L670 360L900 520L1120 390L1280 510V720H0Z" fill="#262d35"/><rect x="120" y="150" width="410" height="440" fill="#4b5157"/><rect x="220" y="300" width="210" height="230" fill="#1f2227"/><path d="M245 500Q325 345 405 500Z" fill="#f4511e"/><rect x="760" y="430" width="300" height="90" fill="#59616c"/><path d="M710 430H1110L1030 360H790Z" fill="#77818c"/><rect x="880" y="520" width="70" height="120" fill="#454b53"/></svg>`
}

func arenaBackgroundSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="1280" height="720"><rect width="1280" height="720" fill="#7893ae"/><path d="M0 280L170 160L350 300L540 120L760 290L1000 150L1280 300V720H0Z" fill="#536878"/><ellipse cx="640" cy="590" rx="520" ry="170" fill="#b89a62" stroke="#6e5b3e" stroke-width="28"/><g fill="#6c747b"><rect x="120" y="290" width="48" height="300"/><rect x="1112" y="290" width="48" height="300"/><rect x="300" y="230" width="48" height="220"/><rect x="932" y="230" width="48" height="220"/></g></svg>`
}
