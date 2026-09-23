package monde

import (
	"math"
	"ordta/library"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// poseCombat : Conserve une position et une rotation locales pour les restaurer après l'action.
type poseCombat struct {
	noeud              *core.Node
	position, rotation math32.Vector3
}

// Une seule animation à la fois. Aucune attente bloquante dans la boucle du jeu.
type AnimationCombat struct {
	Active        bool
	temps         float32
	action        string
	acteur, cible *core.Node
	membres       map[string]*core.Node
	poses         []poseCombat
	fin           func()
	effet         *graphic.Mesh
	matiere       *material.Standard
}

// NouvelleAnimationCombat : Prépare le contrôleur visuel ; les dégâts restent calculés par library.
func NouvelleAnimationCombat(scene *core.Node) *AnimationCombat {
	m := material.NewStandard(&math32.Color{R: .4, G: .75, B: 1})
	orbe := graphic.NewMesh(geometry.NewSphere(.18, 12, 8), m)
	orbe.SetVisible(false)
	scene.Add(orbe)
	return &AnimationCombat{effet: orbe, matiere: m}
}

// Demarrer : Mémorise la pose initiale puis lance l'animation de l'action choisie.
func (a *AnimationCombat) Demarrer(action string, acteur, cible *core.Node, membres map[string]*core.Node, fin func()) {
	a.Arreter()
	a.Active, a.temps, a.action = true, 0, action
	a.acteur, a.cible, a.membres, a.fin = acteur, cible, membres, fin
	couleur := math32.Color{R: .4, G: .75, B: 1}
	switch action {
	case library.SortGrosseBouleDeFeu:
		couleur = math32.Color{R: 1, G: .25, B: .04}
	case library.SortSoinDuCoeur, "soin_monstre":
		couleur = math32.Color{R: .2, G: 1, B: .3}
	case library.SortFlecheDeLumiere, library.SortJugementDesGeants:
		couleur = math32.Color{R: 1, G: .85, B: .25}
	}
	a.matiere.SetColor(&couleur)
	// Conserver les poses permet une restauration exacte, même après une sortie.
	connus := map[*core.Node]bool{}
	memoriser := func(n *core.Node) {
		if n != nil && !connus[n] {
			connus[n] = true
			a.poses = append(a.poses, poseCombat{n, n.Position(), n.Rotation()})
		}
	}
	memoriser(acteur)
	memoriser(cible)
	for _, n := range membres {
		memoriser(n)
	}
}

// restaurer : Remet les nœuds à leur pose de départ pour éviter que les mouvements s'accumulent.
func (a *AnimationCombat) restaurer() {
	for _, p := range a.poses {
		p.noeud.SetPositionVec(&p.position)
		p.noeud.SetRotationVec(&p.rotation)
	}
}

// Arreter : Interrompt l'animation et restaure les poses avant de retirer ses effets.
func (a *AnimationCombat) Arreter() {
	a.restaurer()
	a.Active = false
	a.poses = nil
	a.fin = nil
	if a.effet != nil {
		a.effet.SetVisible(false)
	}
}

// MettreAJour : Fait avancer l'animation avec delta en secondes, puis restaure la pose et appelle la fin.
func (a *AnimationCombat) MettreAJour(delta float32) {
	if !a.Active {
		return
	}
	a.temps += delta
	if a.temps >= 1.25 {
		fin := a.fin
		a.Arreter()
		if fin != nil {
			fin()
		}
		return
	}
	a.restaurer()
	// t parcourt l'animation de 0 à 1. Le sinus fait l'aller-retour :
	// pose de départ, amplitude maximale à mi-parcours, retour au repos.
	t := a.temps / 1.25
	v := float32(math.Sin(float64(t) * math.Pi))
	rotation := func(nom string, x, y, z float32) {
		if n := a.membres[nom]; n != nil {
			r := n.Rotation()
			n.SetRotation(r.X+x*v, r.Y+y*v, r.Z+z*v)
		}
	}
	magie := false
	a.effet.SetScale(1, 1, 1)
	// Chaque sort a sa propre gestuelle, sans modifier les statistiques.
	switch a.action {
	case library.AttaqueBasique, "attaque":
		// Le bras pend vers -Z et le personnage regarde vers -Y :
		// une rotation X négative fait partir la main vers l'avant.
		rotation("bras_droit", -1.8, 0, -.5)
		rotation("arme", -1.8, 0, -.5)
	case library.AttaqueCoupsDePied:
		rotation("jambe_droite", -1.4, 0, 0)
	case library.AttaqueMorsure:
		rotation("tete", .6, 0, 0)
	case library.AttaqueClaquounette:
		rotation("bras_droit", -1.2, 0, -1.4)
	case library.AttaquePichenette:
		rotation("bras_droit", -.9, -.4, 0)
	case library.AttaqueUppercut:
		rotation("bras_droit", -2.5, 0, 0)
	case library.SortCoupDePoing:
		rotation("bras_droit", 1.7, 0, 0)
	case library.SortLameDuDestin:
		rotation("bras_droit", 1.2, 0, -1.5)
	case library.SortGrosseBouleDeFeu:
		rotation("bras_droit", 1.6, .4, 0)
		magie = true
	case library.SortEclatsDuGardien:
		a.effet.SetScale(2, 2, .5)
		rotation("bras_droit", 1, 0, -1)
		rotation("bras_gauche", 1, 0, 1)
		magie = true
	case library.SortFlecheDeLumiere:
		a.effet.SetScale(3, .3, .3)
		rotation("bras_gauche", 1.6, 0, 0)
		rotation("bras_droit", .8, -1.2, 0)
		magie = true
	case library.SortSoinDuCoeur:
		a.effet.SetScale(2, 2, .25)
		rotation("bras_droit", .8, -.5, 0)
		rotation("bras_gauche", .8, .5, 0)
		magie = true
	case library.SortBouclier, "defense":
		rotation("bras_gauche", 1.5, .7, .5)
	case library.SortJugementDesGeants:
		a.effet.SetScale(3, 3, 3)
		rotation("bras_droit", 2.8, -.4, 0)
		rotation("bras_gauche", 2.8, .4, 0)
		magie = true
	case "objet":
		rotation("bras_droit", 2.3, -.5, 0)
	case "loup":
		rotation("tete", -.6, 0, 0)
		rotation("patte_avant_gauche", 1, 0, 0)
		rotation("patte_avant_droite", 1, 0, 0)
	case "sanglier":
		rotation("tete", .8, 0, 0)
		rotation("pattes", -.4, 0, 0)
	case "corbeau":
		rotation("ailes", 0, 1.2, 0)
		rotation("tete", .7, 0, 0)
	case "troll":
		rotation("bras_droit", 2.6, 0, 0)
		rotation("bras_gauche", 2.6, 0, 0)
		rotation("arme", 2.6, 0, 0)
	case "chaman", "soin_monstre":
		rotation("bras_gauche", 2.6, 0, 0)
		rotation("arme", 1.2, 0, 0)
		magie = true
	case "gobelin_cuirasse":
		rotation("bras_gauche", 1.2, 0, 0)
		rotation("bras_droit", 1, 0, -1.1)
		rotation("arme", 1, 0, -1.1)
	default: // Coup de gobelin.
		rotation("bras_droit", 1.8, 0, -.5)
		rotation("arme", 1.8, 0, -.5)
	}
	if a.acteur == nil || a.cible == nil {
		return
	}
	p, q := a.acteur.Position(), a.cible.Position()
	if magie {
		a.effet.SetVisible(true)
		a.effet.SetPosition(p.X+(q.X-p.X)*t, p.Y+(q.Y-p.Y)*t, p.Z+1.2+v*.5)
		if a.action == library.SortJugementDesGeants {
			a.effet.SetPosition(q.X, q.Y, q.Z+1+(1-t)*4)
		}
	} else if a.acteur != a.cible {
		// Avancer puis revenir, sans déplacer réellement le personnage dans le monde.
		a.acteur.SetPosition(p.X+(q.X-p.X)*v*.65, p.Y+(q.Y-p.Y)*v*.65, p.Z)
	}
	if a.acteur != a.cible && a.action != "soin_monstre" && t > .45 && t < .75 {
		a.cible.SetPosition(q.X+.12*float32(math.Sin(float64(t)*60)), q.Y, q.Z)
	}
}
