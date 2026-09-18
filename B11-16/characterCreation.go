package main

import (
	"fmt"
)

// Personnage représente un personnage avec ses attributs
type Personnage struct {
	Nom        string
	Classe     string
	Niveau     int
	PVMax      int
	PVActuel   int
	Inventaire map[string]int
	Argent     int
}

// Creer est une méthode qui permet de créer un personnage
func (p *Personnage) Creer() {
	inventaire := map[string]int{
		"potion de vie": 3,
	}
	p.Inventaire = inventaire

	// Demander le nom et le formater
	var nom string
	for {
		fmt.Print("Entrez le nom de votre personnage : ")
		fmt.Scanln(&nom)

		tab := []rune(nom)
		debut := true

		for i, cara := range tab {
			if debut {
				if 'a' <= cara && cara <= 'z' {
					tab[i] = cara - 32
				}
				debut = false
			} else {
				if 'A' <= cara && cara <= 'Z' {
					tab[i] = cara + 32
				}
				nom = string(tab)
			}
		}

		// Demander la classe
		var choixClasse int
		var pvmax, pvactuel int

		for {
			fmt.Print("Entrez la classe du personnage (1 pour Humain, 2 pour Elfe, 3 pour Nain) : ")
			fmt.Scanln(&choixClasse)

			switch choixClasse {
			case 1:
				p.Classe = "Humain"
				pvmax = 100
				pvactuel = 100
			case 2:
				p.Classe = "Elfe"
				pvmax = 80
				pvactuel = 80
			case 3:
				p.Classe = "Nain"
				pvmax = 120
				pvactuel = 120
			default:
				fmt.Println("Classe invalide. Veuillez entrer 1, 2 ou 3.")
				continue
			}
			break
		}

		// Initialiser les champs du personnage
		p.Nom = nom
		p.Niveau = 1
		p.PVMax = pvmax
		p.PVActuel = pvactuel
		p.Inventaire = inventaire
		p.Argent = 100
		break
	}
}

func main() {
	fmt.Println("=== Création de personnage ===")

	// Créer un personnage en utilisant la méthode Creer()
	var joueur Personnage
	joueur.Creer()

	// Afficher les informations du personnage
	fmt.Printf("\nPersonnage créé :\n")
	fmt.Printf("Nom : %s\n", joueur.Nom)
	fmt.Printf("Classe : %s\n", joueur.Classe)
	fmt.Printf("Niveau : %d\n", joueur.Niveau)
	fmt.Printf("PV Max : %d\n", joueur.PVMax)
	fmt.Printf("PV Actuel : %d\n", joueur.PVActuel)
	fmt.Printf("Inventaire : %v\n", joueur.Inventaire)
}
