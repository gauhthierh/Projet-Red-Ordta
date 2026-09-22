package library

import (
	"fmt"
	"unicode"
)

/*
Ce fichier gère la création du personnage, la validation de son nom et le choix de sa classe.
*/

func CharacterCreation() Character {
	var nom string
	for {
		saisie := ReadLine("Entrez le nom de votre personnage : ")
		if saisie == "" {
			fmt.Println("Le nom ne peut pas être vide.")
			continue
		}
		var ok bool
		nom, ok = FormatName(saisie)
		if ok {
			break
		}
		fmt.Println("Nom invalide : seules les lettres sont acceptées (sans espace, tiret, chiffre ni symbole).")
	}

	structureClasse := ChooseClass()
	return InitCharacter(nom, 1, structureClasse)
}

/* La fonction FormatName vérifie que le nom contient uniquement des lettres et le reformate avec une majuscule suivie de minuscules. */
func FormatName(saisie string) (string, bool) {
	runes := []rune(saisie)
	if len(runes) == 0 {
		return "", false
	}
	for _, r := range runes {
		if !unicode.IsLetter(r) {
			return "", false
		}
	}
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes), true
}

/* La fonction ChooseClass demande au joueur de choisir une classe et renvoie son nom ainsi que ses points de vie maximum. */
func ChooseClass() Classe {
	for {
		choix, _ := ReadChoice("Entrez la classe du personnage (1 pour Humain, 2 pour Elfe, 3 pour Nain) : ")

		switch choix {
		case 1:
			return Classe{Nom: "Humain", PvMax: 100, ManaMax: 100, Attaque: 5}
		case 2:
			return Classe{Nom: "Elfe", PvMax: 80, ManaMax: 120, Attaque: 0}
		case 3:
			return Classe{Nom: "Nain", PvMax: 120, ManaMax: 80, Attaque: 10}
		default:
			fmt.Println("Classe invalide. Veuillez entrer 1, 2 ou 3.")
		}
	}
}

type Classe struct {
	Nom     string
	PvMax   int
	ManaMax int
	Attaque int
}
