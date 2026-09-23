package library

import "testing"

func TestMarchand3DSansMateriaux(t *testing.T) {
	articles := BoutiqueMarchand3D()
	attendus := map[string]bool{ItemPotionDeVie: true, ItemPotionDePoison: true, ItemPotionDeMana: true, ItemLivreGrosseBouleDeFeu: true, ItemAugmentationInventaire: true}
	if len(articles) != len(attendus) {
		t.Fatal("catalogue 3D incorrect")
	}
	for index, article := range articles {
		if !attendus[article.Nom] {
			t.Fatalf("article inattendu : %s", article.Nom)
		}
		p := joueurCombatTest()
		p.Argent = 1000
		avant := p.Inventaire[article.Nom]
		prix := p.PrixPour(article)
		r := p.AcheterMarchand3D(index)
		if !r.Reussite || p.Inventaire[article.Nom] != avant+1 || p.Argent != 1000-prix {
			t.Fatalf("index d'achat décalé : %s", article.Nom)
		}
	}
	p := joueurCombatTest()
	if p.AcheterMarchand3D(len(articles)).Reussite {
		t.Fatal("ancien index accepté hors catalogue")
	}
}

func TestOrCombat3D(t *testing.T) {
	for _, genre := range []TypeMonstre{TypeCorbeau, TypeSanglier, TypeLoup, TypeTroll} {
		p := joueurCombatTest()
		p.Attaque = 1000
		// Une bourse n'occupe pas de case : récompense même avec le sac plein.
		p.CapaciteInventaire = p.TotalInventaire()
		avant := p.Argent
		c, err := NouveauCombat3D(&p, ModeDuel, genre)
		if err != nil {
			t.Fatal(err)
		}
		r := c.Attaquer(0)
		attendu := OrMonstre3D(genre)
		if !r.Victoire || r.OrGagne != attendu || c.OrTotal != attendu || p.Argent != avant+attendu {
			t.Fatalf("récompense incorrecte pour %s : %+v", genre, r)
		}
		c.Attaquer(0)
		c.Quitter3D()
		c.Quitter3D()
		if p.Argent != avant+attendu {
			t.Fatal("or perdu ou attribué deux fois")
		}
	}
}

func TestPasOrEntrainementOuAbandon(t *testing.T) {
	for _, mode := range []ModeCombat{ModeEntrainement, ModeDuel} {
		p := joueurCombatTest()
		p.Attaque = 1000
		avant := p.Argent
		c, _ := NouveauCombat3D(&p, mode, TypeLoup)
		if mode == ModeEntrainement {
			c.Attaquer(0)
		}
		if p.Argent != avant || c.OrTotal != 0 {
			t.Fatal("or gagné sans victoire réelle")
		}
		c.Quitter3D()
		if p.Argent != avant {
			t.Fatal("l'abandon donne de l'or")
		}
	}
}
