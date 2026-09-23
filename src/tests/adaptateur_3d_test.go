package tests

import (
	"ordta/library"
	"testing"
	"time"
)

func TestCreationClassesEtNoms(t *testing.T) {
	for _, classe := range []string{"Humain", "Elfe", "Nain"} {
		p, err := library.CreerPersonnage3D("éLoDiE", classe)
		c, _ := library.ClasseDepuisNom3D(classe)
		if err != nil || p.Nom != "Élodie" || p.PvActuel != c.PvMax/2 || p.PvMaxTotal != c.PvMax || p.TotalInventaire() != 3 || p.CapaciteInventaire != 10 {
			t.Fatalf("création %s : %+v, %v", classe, p, err)
		}
	}
	for _, nom := range []string{"", "Jean Paul", "Test1", "A-b"} {
		if _, err := library.CreerPersonnage3D(nom, "Humain"); err == nil {
			t.Errorf("nom accepté : %q", nom)
		}
	}
	if _, err := library.CreerPersonnage3D("Test", "Inconnue"); err == nil {
		t.Fatal("classe inconnue acceptée")
	}
}

func TestCatalogueEtApprentissage(t *testing.T) {
	if len(library.BoutiqueMarchand3D()) != len(library.Boutique) {
		t.Fatal("catalogue incomplet")
	}
	for index, article := range library.Boutique {
		p := library.NouveauPersonnage3D("Test", "Humain")
		p.Argent = 1000
		if article.NiveauMin > 1 {
			if p.AcheterMarchand3D(index).Reussite || p.Argent != 1000 {
				t.Fatal("niveau minimum ignoré", article.Nom)
			}
		}
		p.Niveau = 10
		if !p.AcheterMarchand3D(index).Reussite {
			t.Fatal("achat refusé", article.Nom)
		}
		sort, livre := library.SortDuLivre(article.Nom)
		attaque, manuel := library.AttaqueDuManuel(article.Nom)
		if !livre && !manuel {
			continue
		}
		if !p.UtiliserObjet3D(article.Nom).Reussite || p.Inventaire[article.Nom] != 0 {
			t.Fatal("apprentissage", article.Nom)
		}
		p.AddInventory(article.Nom)
		if p.UtiliserObjet3D(article.Nom).Reussite || p.Inventaire[article.Nom] != 1 {
			t.Fatal("doublon consommé", article.Nom)
		}
		if livre && p.Skill[len(p.Skill)-1] != sort {
			t.Fatal("sort manquant")
		}
		if manuel && p.AttaquesPhysiques[len(p.AttaquesPhysiques)-1] != attaque {
			t.Fatal("attaque manquante")
		}
	}
}

func TestInitiativeEtNumerotation(t *testing.T) {
	vus := map[bool]bool{}
	for essai := 0; essai < 200; essai++ {
		p := library.NouveauPersonnage3D("Test", "Humain")
		p.PvActuel = 100
		c, err := library.NouveauCombat3D(&p, library.ModeEntrainement, library.TypeGobelin)
		if err != nil {
			t.Fatal(err)
		}
		premier := p.Initiative >= c.Ennemis[0].Monstre.Initiative
		vus[premier] = true
		if (c.Phase == library.PhaseTourJoueur) != premier {
			t.Fatal("initiative ignorée")
		}
		for tour := 1; tour <= 3; tour++ {
			if c.Tour != tour {
				t.Fatalf("tour %d au lieu de %d", c.Tour, tour)
			}
			if premier {
				c.Attaquer(0)
			}
			r := c.ProchaineActionMonstre()
			attendu := 5
			if tour == 3 {
				attendu = 10
			}
			if !r.Reussite || r.Degats != attendu {
				t.Fatalf("pattern tour %d : %+v", tour, r)
			}
			if !premier {
				c.Attaquer(0)
			}
			if c.Tour != tour+1 {
				t.Fatal("compteur de tours incorrect")
			}
		}
	}
	if len(vus) != 2 {
		t.Fatal("les deux ordres d'initiative n'ont pas été exercés")
	}
}

func TestEntrainementExperienceSansRestauration(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	p.Attaque = 100
	c, _ := library.NouveauCombat3D(&p, library.ModeEntrainement, library.TypeGobelin)
	if c.Phase == library.PhaseTourMonstres {
		c.ProchaineActionMonstre()
	}
	r := c.Attaquer(0)
	pv := p.PvActuel
	if !r.Victoire || r.ExperienceGagnee != library.ExperienceGobelin || r.OrGagne < library.OrMinGobelin || r.OrGagne > library.OrMaxGobelin {
		t.Fatalf("%+v", r)
	}
	c.Quitter3D()
	if p.ExperienceActuelle != library.ExperienceGobelin || p.PvActuel != pv {
		t.Fatal("progression restaurée à tort")
	}
}

func TestPotionsPoisonEtExtensions(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	for tour := 0; tour < 3; tour++ {
		p.AddInventory(library.ItemAugmentationInventaire)
		if !p.UtiliserObjet3D(library.ItemAugmentationInventaire).Reussite {
			t.Fatal("extension refusée")
		}
	}
	p.AddInventory(library.ItemAugmentationInventaire)
	if p.UtiliserObjet3D(library.ItemAugmentationInventaire).Reussite || p.CapaciteInventaire != 40 {
		t.Fatal("limite des extensions")
	}
	p.Inventaire = map[string]int{library.ItemPotionDeVie: 4, library.ItemPotionDePoison: 1}
	if len(p.Inventaire3D()) != 5 {
		t.Fatal("les objets sont empilés")
	}
	effets := library.EffetsPersonnage3D{Joueur: &p}
	if !effets.Utiliser(library.ItemPotionDePoison).Reussite {
		t.Fatal("poison refusé")
	}
	if len(effets.MettreAJour(999*time.Millisecond)) != 0 {
		t.Fatal("tick anticipé")
	}
	if len(effets.MettreAJour(2001*time.Millisecond)) != 3 || p.PvActuel != 20 || effets.Actif() {
		t.Fatal("poison incomplet")
	}
	p.Inventaire[library.ItemPotionDePoison] = 1
	p.PvActuel = 5
	effets.Utiliser(library.ItemPotionDePoison)
	effets.MettreAJour(3 * time.Second)
	if p.PvActuel != 50 || effets.Actif() {
		t.Fatal("résurrection incorrecte")
	}
}

func TestCombatLivreEtMana(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	p.Inventaire[library.ItemLivreGrosseBouleDeFeu] = 1
	c, _ := library.NouveauCombat3D(&p, library.ModeEntrainement, library.TypeGobelin)
	if c.Phase == library.PhaseTourMonstres {
		c.ProchaineActionMonstre()
	}
	if r := c.UtiliserObjet(library.ItemLivreGrosseBouleDeFeu); !r.Reussite || !r.TourConsomme {
		t.Fatalf("%+v", r)
	}
	c.ProchaineActionMonstre()
	p.ManaActuel = 0
	pv := c.Ennemis[0].Monstre.PvActuel
	if c.LancerSort(library.SortGrosseBouleDeFeu, 0).Reussite || c.Ennemis[0].Monstre.PvActuel != pv || c.Phase != library.PhaseTourJoueur {
		t.Fatal("sort sans mana accepté")
	}
	p.ManaActuel = 100
	r := c.LancerSort(library.SortGrosseBouleDeFeu, 0)
	degats, cout, _ := library.InfosSort(library.SortGrosseBouleDeFeu)
	if !r.Reussite || r.Degats != degats || p.ManaActuel != 100-cout {
		t.Fatalf("%+v", r)
	}
}

func TestPoisonCombatAttendFinEffet(t *testing.T) {
	p := library.NouveauPersonnage3D("Test", "Humain")
	p.PvActuel = 100
	p.Inventaire[library.ItemPotionDePoison] = 1
	c, _ := library.NouveauCombat3D(&p, library.ModeEntrainement, library.TypeGobelin)
	if c.Phase == library.PhaseTourMonstres {
		c.ProchaineActionMonstre()
	}
	avant := p.PvActuel
	if !c.UtiliserObjet(library.ItemPotionDePoison).Reussite {
		t.Fatal("poison refusé")
	}
	if c.ProchaineActionMonstre().Reussite || p.PvActuel != avant {
		t.Fatal("attaque avant la fin du poison")
	}
	if len(c.MettreAJourEffets(3*time.Second)) != 3 || p.PvActuel != avant-30 {
		t.Fatal("effet incorrect")
	}
	if !c.ProchaineActionMonstre().Reussite {
		t.Fatal("tour adverse bloqué")
	}
}
