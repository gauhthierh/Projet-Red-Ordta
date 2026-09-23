package personnages

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

const (
	cuirArmure   uint32 = 0x68452e
	cuirSombre   uint32 = 0x38291f
	acierArmure  uint32 = 0x8297a1
	acierClair   uint32 = 0xb4c4c9
	laitonArmure uint32 = 0xb59350
	tissuArmure  uint32 = 0x30536b
)

// Pièces fixes en Z vertical. Les arrondis sont des ellipsoïdes à facettes.
// Seuls les métaux brillent légèrement ; cuir et tissu restent mats.
func pieceArmure(parent *core.Node, ronde bool, x, y, z, w, d, h float32, couleur uint32) *graphic.Mesh {
	mat := material.NewStandard(&math32.Color{R: float32(couleur>>16&255) / 255, G: float32(couleur>>8&255) / 255, B: float32(couleur&255) / 255})
	mat.SetSpecularColor(&math32.Color{R: .03, G: .03, B: .03})
	if couleur == acierArmure || couleur == acierClair || couleur == laitonArmure {
		mat.SetSpecularColor(&math32.Color{R: .28, G: .28, B: .28})
		mat.SetShininess(24)
	}
	var mesh *graphic.Mesh
	if ronde {
		mesh = graphic.NewMesh(geometry.NewSphere(.5, 12, 8), mat)
		mesh.SetScale(w, d, h)
	} else {
		mesh = graphic.NewMesh(geometry.NewBox(w, d, h), mat)
	}
	mesh.SetPosition(x, y, z)
	parent.Add(mesh)
	return mesh
}

// creerChapeauDetaille : Assemble les éléments visibles de la protection de tête.
func creerChapeauDetaille(parent *core.Node) {
	// Bord large, calotte, bande de cuir et boucle frontale ajourée.
	pieceArmure(parent, true, 0, 0, 2.035, .79, .68, .085, cuirSombre)
	pieceArmure(parent, true, 0, 0, 2.064, .75, .64, .075, cuirArmure)
	pieceArmure(parent, true, 0, .025, 2.17, .51, .46, .35, cuirArmure)
	pieceArmure(parent, false, 0, -.203, 2.125, .39, .045, .085, cuirSombre)
	for _, x := range []float32{-.07, .07} {
		pieceArmure(parent, false, x, -.236, 2.125, .018, .023, .085, laitonArmure)
	}
	for _, z := range []float32{2.087, 2.163} {
		pieceArmure(parent, false, 0, -.236, z, .15, .023, .016, laitonArmure)
	}
	pieceArmure(parent, false, 0, -.249, 2.125, .08, .016, .012, laitonArmure)
	// Plume bleutée inclinée sur le côté, avec nervure et plusieurs barbes.
	plume := core.NewNode()
	plume.SetPosition(.25, .08, 2.13)
	plume.SetRotationY(.35)
	parent.Add(plume)
	pieceArmure(plume, false, 0, 0, .16, .018, .024, .36, laitonArmure)
	pieceArmure(plume, true, 0, 0, .23, .12, .038, .33, 0x416b85)
	for _, z := range []float32{.12, .18, .24, .30} {
		pieceArmure(plume, false, .025, -.023, z, .065, .015, .013, 0x779cab).SetRotationY(-.35)
	}
}

// creerTuniqueDetaillee : Assemble le vêtement et ses détails autour du torse.
func creerTuniqueDetaillee(parent *core.Node) {
	// Vêtement rembourré et coutures latérales.
	pieceArmure(parent, false, 0, 0, 1.17, .69, .40, .64, tissuArmure)
	for _, x := range []float32{-.30, .30} {
		pieceArmure(parent, false, x, -.211, 1.19, .024, .025, .52, cuirArmure)
	}
	// Plastron bombé ; bandes latérales et dorsale protègent aussi le dos.
	pieceArmure(parent, true, 0, -.215, 1.25, .59, .20, .46, acierArmure)
	pieceArmure(parent, false, 0, -.319, 1.25, .025, .018, .30, acierClair)
	pieceArmure(parent, true, 0, .207, 1.24, .58, .12, .46, acierArmure)
	for _, x := range []float32{-.34, .34} {
		pieceArmure(parent, false, x, 0, 1.23, .055, .34, .38, cuirArmure)
	}
	// Gorgerin en trois pièces, sans recouvrir le cou ni le visage.
	pieceArmure(parent, false, 0, -.208, 1.48, .38, .055, .055, acierClair)
	for _, x := range []float32{-.22, .22} {
		pieceArmure(parent, false, x, 0, 1.48, .07, .35, .055, acierArmure)
	}
	// Lames abdominales superposées et rivets visibles.
	for _, z := range []float32{1.07, 1.00, .93} {
		pieceArmure(parent, false, 0, -.225, z, .54, .07, .055, acierArmure)
		for _, x := range []float32{-.235, .235} {
			pieceArmure(parent, true, x, -.27, z, .025, .025, .025, laitonArmure)
		}
	}
	for _, x := range []float32{-.18, .18} {
		pieceArmure(parent, false, x, -.234, 1.41, .06, .04, .17, cuirSombre)
		pieceArmure(parent, true, x, -.264, 1.46, .03, .022, .03, laitonArmure)
	}
	// Ceinture, boucle et deux petites poches latérales.
	pieceArmure(parent, false, 0, 0, .875, .72, .45, .085, cuirSombre)
	for _, x := range []float32{-.068, .068} {
		pieceArmure(parent, false, x, -.25, .875, .018, .035, .078, laitonArmure)
	}
	for _, z := range []float32{.84, .91} {
		pieceArmure(parent, false, 0, -.25, z, .15, .035, .018, laitonArmure)
	}
	for _, x := range []float32{-.295, .295} {
		pieceArmure(parent, false, x, -.25, .88, .115, .12, .15, cuirArmure)
		pieceArmure(parent, false, x, -.317, .93, .12, .02, .038, cuirSombre)
		pieceArmure(parent, true, x, -.332, .91, .022, .014, .022, laitonArmure)
	}
	// Emblème discret en relief.
	pieceArmure(parent, false, 0, -.328, 1.30, .095, .018, .095, laitonArmure).SetRotationY(math32.Pi / 4)
}

// creerProtectionBras : Attache la protection au pivot du bras pour qu'elle suive son animation.
func creerProtectionBras(bras *core.Node) *core.Node {
	protection := core.NewNode()
	bras.Add(protection)
	// Épaulière arrondie, puis deux lames articulées.
	pieceArmure(protection, true, 0, 0, -.025, .34, .39, .25, acierArmure)
	for _, z := range []float32{-.10, -.16} {
		pieceArmure(protection, false, 0, -.005, z, .29, .34, .045, acierClair)
	}
	for _, x := range []float32{-.11, .11} {
		pieceArmure(protection, true, x, -.182, -.10, .027, .025, .027, laitonArmure)
	}
	// Brassard de cuir, plaque et sangles ; la main reste visible.
	pieceArmure(protection, false, 0, 0, -.30, .245, .30, .19, cuirArmure)
	pieceArmure(protection, false, 0, -.164, -.30, .18, .05, .17, acierArmure)
	for _, z := range []float32{-.24, -.37} {
		pieceArmure(protection, false, 0, 0, z, .26, .32, .028, cuirSombre)
		pieceArmure(protection, false, .08, -.177, z, .035, .02, .035, laitonArmure)
	}
	return protection
}

// creerBotteDetaillee : Attache la botte et la protection de jambe au même pivot que la jambe.
func creerBotteDetaillee(jambe *core.Node) *core.Node {
	botte := core.NewNode()
	jambe.Add(botte)
	// Cuisse renforcée, genou à facettes et jambière : même pivot que la jambe.
	pieceArmure(botte, false, 0, 0, -.16, .29, .32, .27, cuirArmure)
	pieceArmure(botte, true, 0, -.175, -.16, .25, .10, .25, acierArmure)
	pieceArmure(botte, true, 0, -.19, -.345, .30, .15, .17, acierClair)
	pieceArmure(botte, false, 0, 0, -.57, .30, .34, .40, cuirArmure)
	pieceArmure(botte, true, 0, -.195, -.55, .24, .085, .26, acierArmure)
	pieceArmure(botte, false, 0, -.238, -.54, .018, .018, .19, acierClair)
	// Semelle, talon et bout du pied capoté.
	pieceArmure(botte, false, 0, -.06, -.755, .32, .47, .07, cuirSombre)
	pieceArmure(botte, false, 0, .07, -.73, .30, .17, .10, cuirSombre)
	pieceArmure(botte, false, 0, -.08, -.69, .31, .45, .12, cuirArmure)
	pieceArmure(botte, true, 0, -.22, -.665, .32, .23, .13, acierArmure)
	// Deux sangles fermées par des boucles et un revers cousu.
	for _, z := range []float32{-.42, -.62} {
		pieceArmure(botte, false, 0, 0, z, .315, .36, .037, cuirSombre)
		pieceArmure(botte, false, .10, -.19, z, .055, .025, .05, laitonArmure)
		pieceArmure(botte, false, .10, -.208, z, .025, .012, .022, cuirSombre)
	}
	for _, x := range []float32{-.10, 0, .10} {
		pieceArmure(botte, true, x, -.175, -.385, .014, .017, .022, laitonArmure)
	}
	return botte
}
