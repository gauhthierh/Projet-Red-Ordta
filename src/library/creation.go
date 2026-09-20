package library

import (
	"fmt"
	"unicode"
)

// characterCreation demande le nom et la classe, puis construit le personnage
// avec initCharacter : niveau 1, PV actuels à 50 % des PV max.
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

	classe, pvMax := ChooseClass()
	return InitCharacter(nom, classe, 1, pvMax, pvMax/2)
}

// formatName vérifie que le nom ne contient que des lettres, accentuées
// comprises, et le remet en forme : majuscule initiale, reste en minuscules.
// Le nom est traité rune par rune (caractère par caractère) et non octet par
// octet : en UTF-8, « é » occupe deux octets. Renvoie false si le nom est
// vide ou contient autre chose qu'une lettre.
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

// chooseClass demande la classe jusqu'à obtenir un choix valide, puis renvoie
// son nom et ses PV max.
func ChooseClass() (string, int) {
	for {
		// Une saisie non numérique donne 0, traité par le cas default.
		choix, _ := ReadChoice("Entrez la classe du personnage (1 pour Humain, 2 pour Elfe, 3 pour Nain) : ")

		switch choix {
		case 1:
			return "Humain", 100
		case 2:
			return "Elfe", 80
		case 3:
			return "Nain", 120
		default:
			fmt.Println("Classe invalide. Veuillez entrer 1, 2 ou 3.")
		}
	}
}
