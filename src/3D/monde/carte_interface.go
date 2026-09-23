package monde

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

// CarteInterface : Regroupe le plan, son marqueur et l'état d'ouverture, séparés de la carte 3D.
type CarteInterface struct {
	Panneau                               *gui.Panel
	Marqueur                              *gui.Panel
	Position                              *gui.Label
	Fermer                                *gui.Button
	Ouverte                               bool
	tailleMonde, cote, origineX, origineY float32
}

// La carte garde ses proportions et utilise le même repère XY que le JSON.
func NouvelleCarteInterface(scene *core.Node, largeur, hauteur, tailleMonde float32) (*CarteInterface, error) {
	if tailleMonde <= 0 {
		return nil, fmt.Errorf("taille du monde invalide pour la carte")
	}
	illustration, err := gui.NewImage("../assets/maps/red_world/red_world_map_2d_pixel.png")
	if err != nil {
		return nil, fmt.Errorf("carte 2D : %w", err)
	}
	cote := hauteur - 120
	if largeur-40 < cote {
		cote = largeur - 40
	}
	c := &CarteInterface{Panneau: gui.NewPanel(largeur, hauteur), tailleMonde: tailleMonde, cote: cote, origineX: (largeur - cote) / 2, origineY: 60}
	c.Panneau.SetColor4(&math32.Color4{R: .02, G: .03, B: .04, A: 1})
	illustration.SetSize(cote, cote)
	illustration.SetPosition(c.origineX, c.origineY)
	c.Panneau.Add(illustration)
	titre := gui.NewLabel("CARTE DU MONDE — Nord en haut")
	titre.SetPosition(c.origineX, 20)
	c.Panneau.Add(titre)
	c.Marqueur = gui.NewPanel(16, 16)
	c.Marqueur.SetColor(&math32.Color{R: 1, G: 1, B: 1})
	centre := gui.NewPanel(10, 10)
	centre.SetPosition(3, 3)
	centre.SetColor(&math32.Color{R: 1, G: .08, B: .08})
	c.Marqueur.Add(centre)
	c.Panneau.Add(c.Marqueur)
	c.Position = gui.NewLabel("")
	c.Position.SetPosition(c.origineX, hauteur-45)
	c.Panneau.Add(c.Position)
	c.Fermer = gui.NewButton("FERMER [M / ÉCHAP]")
	c.Fermer.SetSize(220, 35)
	c.Fermer.SetPosition(c.origineX+cote-220, hauteur-50)
	c.Panneau.Add(c.Fermer)
	c.Panneau.SetVisible(false)
	scene.Add(c.Panneau)
	return c, nil
}

// Actualiser : Déplace le marqueur en conservant une marge pour qu'il reste visible sur la carte.
func (c *CarteInterface) Actualiser(position math32.Vector3) {
	x, y := positionSurCarte(position.X, position.Y, c.tailleMonde, c.cote)
	// Le cadre reste entièrement visible même aux limites du monde.
	x = math32.Clamp(x, 8, c.cote-8)
	y = math32.Clamp(y, 8, c.cote-8)
	c.Marqueur.SetPosition(c.origineX+x-8, c.origineY+y-8)
	c.Position.SetText(fmt.Sprintf("Rouge : vous — X %.1f / Y %.1f", position.X, position.Y))
}

// Inverse de image_to_world : +Y (nord) correspond au haut de l'image.
func positionSurCarte(x, y, tailleMonde, cote float32) (float32, float32) {
	return (x/tailleMonde + .5) * cote, (.5 - y/tailleMonde) * cote
}
