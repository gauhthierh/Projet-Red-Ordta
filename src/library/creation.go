package library

import (
	"fmt"
	"unicode"
)

// ================================================= //
// === Ce fichier gère la création du personnage === //
// ================================================= //

// Demande le nom et la classe du joueur //
// Puis crée son personnage //
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

// Vérifie que le nom ne contient que des lettres et le met en forme //
// Renvoie false si le nom est invalide //
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

// Demande au joueur de choisir sa classe et renvoie sa fiche //
func ChooseClass() Classe {
	for {
		choix, _ := ReadChoice("Entrez la classe du personnage (1 pour Humain, 2 pour Elfe, 3 pour Nain) : ")
		switch choix {
		case 1:
			return Classe{Nom: "Humain", PvMax: 100, ManaMax: 100, Attaque: 5, GainPvMax: 10, GainManaMax: 10, GainAttaque: 1}
		case 2:
			return Classe{Nom: "Elfe", PvMax: 80, ManaMax: 120, Attaque: 0, GainPvMax: 6, GainManaMax: 15, GainAttaque: 0}
		case 3:
			return Classe{Nom: "Nain", PvMax: 120, ManaMax: 80, Attaque: 10, GainPvMax: 15, GainManaMax: 5, GainAttaque: 2}
		default:
			fmt.Println("Classe invalide. Veuillez entrer 1, 2 ou 3.")
		}
	}
}

// Représente les caractéristiques d'une classe //
type Classe struct {
	Nom         string
	PvMax       int
	ManaMax     int
	Attaque     int
	GainPvMax   int
	GainManaMax int
	GainAttaque int
}
