package main

import "fmt"

// characterCreation demande le nom et la classe, puis construit le personnage.
func characterCreation() Character {
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
		var classe string
		var pvmax, pvactuel int

		for {
			fmt.Print("Entrez la classe du personnage (1 pour Humain, 2 pour Elfe, 3 pour Nain) : ")
			fmt.Scanln(&choixClasse)

			switch choixClasse {
			case 1:
				classe = "Humain"
				pvmax = 100
				pvactuel = 100
			case 2:
				classe = "Elfe"
				pvmax = 80
				pvactuel = 80
			case 3:
				classe = "Nain"
				pvmax = 120
				pvactuel = 120
			default:
				fmt.Println("Classe invalide. Veuillez entrer 1, 2 ou 3.")
				continue
			}
			break
		}

		// Retourner le personnage créé
		return initCharacter(nom, classe, 1, pvmax, pvactuel)
	}
}
