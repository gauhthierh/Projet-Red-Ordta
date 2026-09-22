package library

import (
	"fmt"
	"strings"
)

/*
Ce fichier définit le personnage et gère ses statistiques, son inventaire, ses équipements et ses actions pendant les combats.

	La structure Character contient toutes les informations et ressources appartenant au personnage contrôlé par le joueur.
*/
type Character struct {
	Nom                            string
	Classe                         string
	Niveau                         int
	PVMaxBase                      int
	PVActuel                       int
	PVMaxTotal                     int
	Inventaire                     map[string]int
	Skill                          []string
	Initiative                     int
	ExperienceActuelle             int
	ExperienceMax                  int
	CapaciteInventaire             int
	Argent                         int
	PotionGratuitePrise            bool
	Equipement                     Equipment
	AugmentationInventaireUtilisee int
	ManaActuel                     int
	ManaMax                        int
}

/* Ces constantes définissent les valeurs de départ du personnage et les caractéristiques de son attaque basique. */
const (
	capaciteInventaireDepart           = 10
	argentDepart                       = 100
	maxAugmentationsInventaire         = 3
	bonusAugmentationInventaire        = 10
	attaqueBasique                     = "Attaque Basique"
	degatsAttaqueBasique               = 5
	experienceinitiale                 = 0
	experiencemaximale                 = 100
	augmentationExperienceNumerateur   = 13
	augmentationExperienceDenominateur = 10
)

/* La fonction InitCharacter crée et initialise un personnage avec ses statistiques, ses objets, son sort et ses ressources de départ. */
func InitCharacter(nom string, classe string, niveau int, pvmax int, pvactuel int, manamax int) Character {
	inventaire := map[string]int{
		ItemPotionDeVie: 3,
	}

	personnage := Character{
		Nom:                nom,
		Classe:             classe,
		Niveau:             niveau,
		PVMaxBase:          pvmax,
		PVActuel:           pvactuel,
		Inventaire:         inventaire,
		Skill:              []string{SortCoupDePoing},
		CapaciteInventaire: capaciteInventaireDepart,
		Argent:             argentDepart,
		ExperienceActuelle: experienceinitiale,
		ExperienceMax:      experiencemaximale,
		ManaActuel:         manamax,
		ManaMax:            manamax,
	}
	personnage.MettreAJourPvMax()
	return personnage
}

/* La méthode displayInfo affiche les statistiques, les sorts, l'argent, l'inventaire et les équipements du personnage. */
func (c Character) displayInfo() {
	fmt.Println("\n=== PERSONNAGE ===")
	fmt.Printf("Nom : %s\n", c.Nom)
	fmt.Printf("Classe : %s\n", c.Classe)
	fmt.Printf("Niveau : %d\n", c.Niveau)
	fmt.Printf("Pv : %d / %d\n", c.PVActuel, c.PVMaxTotal)
	fmt.Printf("Sorts : %s\n", strings.Join(c.Skill, ", "))
	fmt.Printf("Argent : %d pièces d'or\n", c.Argent)
	fmt.Printf("Expérience : %d / %d\n", c.ExperienceActuelle, c.ExperienceMax)
	fmt.Println("\n=== INVENTAIRE ===")
	c.AccessInventory()
	fmt.Println("\n=== EQUIPEMENT ===")
	fmt.Printf("-- Tête --\n %s\n", c.Equipement.Tete.DescriptionEquipement())
	fmt.Printf("-- Torse --\n %s\n", c.Equipement.Torse.DescriptionEquipement())
	fmt.Printf("-- Pied --\n %s\n", c.Equipement.Pied.DescriptionEquipement())
}

/* La méthode isDead vérifie si le personnage est mort et le ressuscite avec la moitié de ses points de vie maximum. */
func (c *Character) isDead() bool {
	if c.PVActuel > 0 {
		return false
	}
	fmt.Printf("%s est mort !\n", c.Nom)
	c.PVActuel = c.PVMaxTotal / 2
	fmt.Printf("%s ressuscite avec %d / %d PV.\n", c.Nom, c.PVActuel, c.PVMaxTotal)
	return true
}

/* La méthode CalculPvAvecBonus calcule les points de vie maximum du personnage en ajoutant les bonus de ses équipements. */
func (c *Character) CalculPvAvecBonus() int {
	return c.Equipement.Tete.BonusPv + c.Equipement.Torse.BonusPv + c.Equipement.Pied.BonusPv + c.PVMaxBase
}

/* La méthode MettreAJourPvMax actualise les points de vie maximum du personnage après un changement d'équipement. */
func (c *Character) MettreAJourPvMax() {
	c.PVMaxTotal = c.CalculPvAvecBonus()
	if c.PVActuel > c.PVMaxTotal {
		c.PVActuel = c.PVMaxTotal
	}
}

/* La méthode CharacterTurn permet au personnage d'attaquer, d'utiliser un objet ou de lancer un sort pendant son tour. */
func (c *Character) CharacterTurn(m *Monster) {
	for {
		fmt.Println("=== COMBAT ===")
		fmt.Println("1. Attaquer")
		fmt.Println("2. Inventaire")
		fmt.Println("3. Sorts")
		choix, ok := ReadChoice("Entrez votre choix :")
		if !ok {
			fmt.Println("Choix invalide, veuillez entrer une saisie valide !")
			continue
		}
		switch choix {
		case 1:
			m.PVActuel -= degatsAttaqueBasique
			if m.PVActuel < 0 {
				m.PVActuel = 0
			}
			fmt.Printf("%s lance %s sur %s et lui inflige %d dégâts !\n", c.Nom, attaqueBasique, m.Nom, degatsAttaqueBasique)
			fmt.Printf("Pv de %s : %d / %d\n", m.Nom, m.PVActuel, m.PVMax)
			return
		case 2:
			if c.ChoixInventaire() {
				return
			}
		case 3:
			if c.ChoixSort(m) {
				return
			}
		default:
			fmt.Println("Choix invalide, veuillez entrer un choix valide")
		}
	}
}

/* La méthode ChoixInventaire affiche les objets utilisables en combat et indique si le joueur en a utilisé un. */
func (c *Character) ChoixInventaire() bool {
	objets := c.SortedItems()
	if len(objets) == 0 {
		fmt.Println("L'inventaire est vide !")
		return false
	}
	for {
		for indice, objet := range objets {
			fmt.Printf("%d. %s, quantité : %d\n", indice+1, objet, c.Inventaire[objet])
		}
		fmt.Println("0. Retour")
		choix, ok := ReadChoice("Votre choix : ")
		if !ok || choix < 0 || choix > len(objets) {
			fmt.Printf("Choix invalide. Entrez un nombre entre 0 et %d\n", len(objets))
			continue
		}
		if choix == 0 {
			return false
		}
		fmt.Printf("Vous utilisez %s\n", objets[choix-1])
		c.useItem(objets[choix-1])
		return true
	}
}
