package library

// Ce fichier adapte le backend textuel à la scène 3D.
// Il ne lit jamais le terminal et ne bloque jamais la boucle graphique.

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// PhaseCombat : Indique quel camp peut agir ou quelle transition reste à effectuer.
type PhaseCombat string

const (
	PhaseTourJoueur   PhaseCombat = "tour_joueur"
	PhaseTourMonstres PhaseCombat = "tour_monstres"
	PhaseEntreVagues  PhaseCombat = "entre_vagues"
	PhaseVictoire     PhaseCombat = "victoire"
	PhaseDefaite      PhaseCombat = "defaite"
	PhaseAbandon      PhaseCombat = "abandon"
)

// TypeMonstre : Identifie une créature pour ses règles, son butin et ses animations.
type TypeMonstre string

const (
	TypeGobelin         TypeMonstre = "gobelin"
	TypeGobelinCuirasse TypeMonstre = "gobelin_cuirasse"
	TypeChaman          TypeMonstre = "chaman"
	TypeLoup            TypeMonstre = "loup"
	TypeTroll           TypeMonstre = "troll"
	TypeSanglier        TypeMonstre = "sanglier"
	TypeCorbeau         TypeMonstre = "corbeau"
)

// ModeCombat : Distingue l'entraînement, les vagues d'arène et les duels.
type ModeCombat string

const (
	ModeEntrainement ModeCombat = "Entraînement"
	ModeArene        ModeCombat = "Arène"
	ModeDuel         ModeCombat = "Duel"
)

// ButinMonstre : Associe un matériau à sa probabilité de tomber, exprimée en pourcentage.
type ButinMonstre struct {
	Objet       string
	Pourcentage int
}

// L'or récompense chaque victoire sur une créature, indépendamment du butin.
func OrMonstre3D(monstre TypeMonstre) int {
	switch monstre {
	case TypeCorbeau:
		return 2
	case TypeSanglier:
		return 3
	case TypeGobelin:
		return OrMinGobelin
	case TypeLoup:
		return 5
	case TypeChaman:
		return 7
	case TypeGobelinCuirasse:
		return 8
	case TypeTroll:
		return 20
	default:
		return 0
	}
}

// Les mêmes matériaux sont utilisés par les recettes de l'Armurerie.
func ButinPossible3D(monstre TypeMonstre) ButinMonstre {
	switch monstre {
	case TypeLoup:
		return ButinMonstre{ItemFourrureDeLoup, 60}
	case TypeSanglier:
		return ButinMonstre{ItemCuirDeSanglier, 70}
	case TypeCorbeau:
		return ButinMonstre{ItemPlumeDeCorbeau, 75}
	case TypeTroll:
		return ButinMonstre{ItemPeauDeTroll, 45}
	default:
		return ButinMonstre{}
	}
}

// ResultatAction contient les informations nécessaires aux interfaces 3D.
type ResultatAction struct {
	// Réussir une action et consommer un tour sont deux décisions différentes.
	Reussite     bool
	TourConsomme bool
	Type         string
	Source       string
	Cible        string
	Message      string
	// Valeurs utiles au journal : distinguer le changement réel du coût annoncé.
	Degats           int
	Soin             int
	PVAvant          int
	PVApres          int
	ManaAvant        int
	ManaApres        int
	CibleVaincue     bool
	VagueTerminee    bool
	CombatTermine    bool
	Victoire         bool
	ExperienceGagnee int
	OrGagne          int
	NiveauxGagnes    int
}

// EnnemiCombat : Relie les données de combat à un fichier visuel, sans dépendre du moteur G3N.
type EnnemiCombat struct {
	Type     TypeMonstre
	Modele3D string
	Monstre  Monster
}

// CombatArene : Porte l'état d'une session ; la scène demande les actions une à une au rythme des animations.
type CombatArene struct {
	Mode           ModeCombat
	Butins         []string
	joueurCommence bool
	tirageButin    func(int) int
	// Pointeur vers le joueur du monde : les gains persistent après la sortie.
	Joueur  *Character
	Vagues  [][]EnnemiCombat
	Ennemis []EnnemiCombat
	// NumeroVague commence à 1 pour l'affichage ; les indices de slices commencent à 0.
	NumeroVague       int
	Tour              int
	Phase             PhaseCombat
	IndexMonstreActif int
	ExperienceTotale  int
	OrTotal           int
	DefenseActive     bool

	// Le temps s'accumule entre les images ; aucun sleep ne bloque le rendu.
	poisonSecondes      int
	poisonTempsAccumule time.Duration
}

// ObjetInventaire3D est une version directement affichable d'un objet.
type ObjetInventaire3D struct {
	Nom        string
	Quantite   int
	Icone      string
	Utilisable bool
	Equipable  bool
}

// ClasseDepuisNom3D : Retrouve la définition d'une classe et signale les noms inconnus.
func ClasseDepuisNom3D(nom string) (Classe, bool) {
	switch nom {
	case "Humain":
		return Classe{Nom: "Humain", PvMax: 100, ManaMax: 100, Attaque: 5, GainPvMax: 10, GainManaMax: 10, GainAttaque: 1}, true
	case "Elfe":
		return Classe{Nom: "Elfe", PvMax: 80, ManaMax: 120, Attaque: 0, GainPvMax: 6, GainManaMax: 15, GainAttaque: 0}, true
	case "Nain":
		return Classe{Nom: "Nain", PvMax: 120, ManaMax: 80, Attaque: 10, GainPvMax: 15, GainManaMax: 5, GainAttaque: 2}, true
	default:
		return Classe{}, false
	}
}

// NouveauPersonnage3D : Prépare un personnage par défaut ; la saisie du joueur passe par CreerPersonnage3D.
func NouveauPersonnage3D(nom string, nomClasse string) Character {
	classe, existe := ClasseDepuisNom3D(nomClasse)
	if !existe {
		classe, _ = ClasseDepuisNom3D("Humain")
	}
	return InitCharacter(nom, 1, classe)
}

// Inventaire3D : Déplie les quantités : quatre potions produisent quatre entrées, donc quatre cases.
func (c Character) Inventaire3D() []ObjetInventaire3D {
	objets := make([]ObjetInventaire3D, 0, len(c.Inventaire))
	for _, nom := range c.SortedItems() {
		_, equipable := TrouverEquipement(nom)
		// Un exemplaire occupe une case, même si plusieurs portent le même nom.
		for exemplaire := 0; exemplaire < c.Inventaire[nom]; exemplaire++ {
			objets = append(objets, ObjetInventaire3D{
				Nom:        nom,
				Quantite:   1,
				Icone:      IconeObjet3D(nom),
				Utilisable: objetUtilisable3D(nom),
				Equipable:  equipable,
			})
		}
	}
	return objets
}

// IconeObjet3D : Associe un nom d'objet à son fichier d'icône, avec une image de secours.
func IconeObjet3D(nom string) string {
	if _, ok := SortDuLivre(nom); ok {
		return "fireball_spellbook.png"
	}
	if _, ok := AttaqueDuManuel(nom); ok {
		return "fireball_spellbook.png"
	}
	switch nom {
	case ItemPotionDeVie:
		return "health_potion.png"
	case ItemPotionDePoison:
		return "poison_potion.png"
	case ItemPotionDeMana:
		return "mana_potion.png"
	case ItemLivreGrosseBouleDeFeu:
		return "fireball_spellbook.png"
	case ItemFourrureDeLoup:
		return "wolf_fur.png"
	case ItemPeauDeTroll:
		return "troll_hide.png"
	case ItemCuirDeSanglier:
		return "boar_leather.png"
	case ItemPlumeDeCorbeau:
		return "raven_feather.png"
	case ItemChapeauAventurier:
		return "adventurer_hat.png"
	case ItemTuniqueAventurier:
		return "adventurer_tunic.png"
	case ItemBottesAventurier:
		return "adventurer_boots.png"
	case ItemAugmentationInventaire:
		return "inventory_upgrade.png"
	default:
		return "inventory.png"
	}
}

// objetUtilisable3D : Classe les objets utilisables sans vérifier ici leur possession ni l'état du combat.
func objetUtilisable3D(nom string) bool {
	if _, ok := SortDuLivre(nom); ok {
		return true
	}
	if _, ok := AttaqueDuManuel(nom); ok {
		return true
	}
	if _, equipable := TrouverEquipement(nom); equipable {
		return true
	}
	switch nom {
	case ItemPotionDeVie, ItemPotionDePoison, ItemPotionDeMana,
		ItemLivreGrosseBouleDeFeu, ItemAugmentationInventaire:
		return true
	default:
		return false
	}
}

// BoutiqueMarchand3D copie le catalogue pour que l'interface ne modifie pas
// directement la liste partagée avec le marchand textuel.
func BoutiqueMarchand3D() []Item {
	return append([]Item(nil), Boutique...)
}

// AcheterMarchand3D : Vérifie les conditions d'achat avant d'ajouter l'objet et de retirer son prix.
func (c *Character) AcheterMarchand3D(index int) ResultatAction {
	articles := BoutiqueMarchand3D()
	if c == nil || index < 0 || index >= len(articles) {
		return actionRefusee3D("Article inconnu.")
	}
	article := articles[index]
	if c.Niveau < article.NiveauMin {
		return actionRefusee3D(fmt.Sprintf("Niveau %d requis.", article.NiveauMin))
	}
	if !c.ArticleUtile(article.Nom) {
		return actionRefusee3D("Cet article est déjà appris, possédé ou devenu inutile.")
	}
	prix := c.PrixPour(article)
	if c.Argent < prix {
		return actionRefusee3D(fmt.Sprintf("Il manque %d pièces d'or.", prix-c.Argent))
	}
	if !c.AddInventory(article.Nom) {
		return actionRefusee3D("L'inventaire est plein.")
	}
	c.Argent -= prix
	if article.Nom == ItemPotionDeVie && !c.PotionGratuitePrise {
		c.PotionGratuitePrise = true
	}
	if prix == 0 {
		return actionReussie3D(fmt.Sprintf("%s est offert par le marchand.", article.Nom))
	}
	return actionReussie3D(fmt.Sprintf("%s acheté pour %d pièces d'or.", article.Nom, prix))
}

// FabriquerForgeron3D applique les prix et matériaux de l'Armurerie.
func (c *Character) FabriquerForgeron3D(index int) ResultatAction {
	if c == nil || index < 0 || index >= len(Armurerie) {
		return actionRefusee3D("Équipement inconnu.")
	}
	equipement := Armurerie[index]
	if c.Argent < equipement.Prix {
		return actionRefusee3D(fmt.Sprintf("Il manque %d pièces d'or.", equipement.Prix-c.Argent))
	}
	manquants := c.MateriauxManquants(equipement)
	if len(manquants) > 0 {
		noms := make([]string, 0, len(manquants))
		for nom := range manquants {
			noms = append(noms, nom)
		}
		sort.Strings(noms)
		message := "Matériaux manquants :"
		for _, nom := range noms {
			message += fmt.Sprintf(" %s x%d;", nom, manquants[nom])
		}
		return actionRefusee3D(message)
	}

	// Vérifier la place après consommation avant de retirer le moindre matériau.
	placesLiberees := 0
	for _, quantite := range equipement.Materiaux {
		placesLiberees += quantite
	}
	if c.TotalInventaire()-placesLiberees+1 > c.CapaciteInventaire {
		return actionRefusee3D("L'inventaire est plein.")
	}
	for materiau, quantite := range equipement.Materiaux {
		for compteur := 0; compteur < quantite; compteur++ {
			c.RemoveInventory(materiau)
		}
	}
	if !c.AddInventory(equipement.Nom) {
		return actionRefusee3D("L'inventaire est plein.")
	}
	c.Argent -= equipement.Prix
	return actionReussie3D(fmt.Sprintf("%s a été fabriqué.", equipement.Nom))
}

// UtiliserObjet3D : Traite les effets immédiats. Les effets dans le temps passent par les contrôleurs de combat ou d'exploration.
func (c *Character) UtiliserObjet3D(nom string) ResultatAction {
	if c == nil {
		return actionRefusee3D("Le personnage est absent.")
	}

	if sort, ok := SortDuLivre(nom); ok {
		if contientTexte3D(c.Skill, sort) {
			return actionRefusee3D("Ce sort est déjà connu.")
		}
		if !c.RemoveInventory(nom) {
			return actionRefusee3D("Ce livre est absent de l'inventaire.")
		}
		c.SpellBook(sort)
		return actionReussie3D("Sort appris : " + sort)
	}
	if attaque, ok := AttaqueDuManuel(nom); ok {
		if contientTexte3D(c.AttaquesPhysiques, attaque) {
			return actionRefusee3D("Cette attaque est déjà connue.")
		}
		if !c.RemoveInventory(nom) {
			return actionRefusee3D("Ce manuel est absent de l'inventaire.")
		}
		c.ApprentissageAttaque(attaque)
		return actionReussie3D("Attaque apprise : " + attaque)
	}
	switch nom {
	case ItemPotionDeVie:
		return c.utiliserPotionVie3D()
	case ItemPotionDeMana:
		return c.utiliserPotionMana3D()
	case ItemPotionDePoison:
		return actionRefusee3D("La potion de poison s'utilise uniquement pendant un combat.")
	case ItemAugmentationInventaire:
		if c.Inventaire[nom] <= 0 {
			return actionRefusee3D("Cette amélioration n'est pas dans l'inventaire.")
		}
		if !c.UpgradeInventorySlot() {
			return actionRefusee3D("La capacité maximale est déjà atteinte.")
		}
		c.RemoveInventory(nom)
		return actionReussie3D(fmt.Sprintf("Capacité de l'inventaire : %d.", c.CapaciteInventaire))
	default:
		if equipement, existe := TrouverEquipement(nom); existe {
			return c.equiper3D(equipement)
		}
		return actionRefusee3D("Cet objet ne peut pas être utilisé.")
	}
}

// Desequiper3D : Replace une pièce portée dans l'inventaire, si la capacité le permet, puis recalcule les PV maximum.
func (c *Character) Desequiper3D(emplacement string) ResultatAction {
	if c == nil {
		return actionRefusee3D("Le personnage est absent.")
	}
	if !c.VerifPlaceInventaire() {
		return actionRefusee3D("L'inventaire est plein.")
	}

	var equipe *Stuff
	switch emplacement {
	case EmplacementTete:
		equipe = &c.Equipement.Tete
	case EmplacementTorse:
		equipe = &c.Equipement.Torse
	case EmplacementPied:
		equipe = &c.Equipement.Pied
	default:
		return actionRefusee3D("Emplacement d'équipement inconnu.")
	}

	if equipe.Nom == "" {
		return actionRefusee3D("Cet emplacement est déjà vide.")
	}

	nom := equipe.Nom
	c.Inventaire[nom]++
	*equipe = Stuff{}
	c.MettreAJourPvMax()
	return actionReussie3D(fmt.Sprintf("%s est retiré.", nom))
}

// equiper3D : Échange la nouvelle pièce avec l'ancienne et actualise les bonus sans soigner le joueur.
func (c *Character) equiper3D(nouveau Stuff) ResultatAction {
	if c.Inventaire[nouveau.Nom] <= 0 {
		return actionRefusee3D("Cet équipement n'est pas dans l'inventaire.")
	}

	var equipe *Stuff
	switch nouveau.Emplacement {
	case EmplacementTete:
		equipe = &c.Equipement.Tete
	case EmplacementTorse:
		equipe = &c.Equipement.Torse
	case EmplacementPied:
		equipe = &c.Equipement.Pied
	default:
		return actionRefusee3D("Emplacement d'équipement inconnu.")
	}

	ancien := equipe.Nom
	c.RemoveInventory(nouveau.Nom)
	if ancien != "" {
		c.Inventaire[ancien]++
	}
	*equipe = nouveau
	c.MettreAJourPvMax()

	if ancien == "" {
		return actionReussie3D(fmt.Sprintf("%s est équipé.", nouveau.Nom))
	}
	return actionReussie3D(fmt.Sprintf("%s remplace %s.", nouveau.Nom, ancien))
}

// utiliserPotionVie3D : Consomme une potion et renvoie le soin effectif, limité par les PV maximum.
func (c *Character) utiliserPotionVie3D() ResultatAction {
	if c.PvActuel <= 0 {
		return actionRefusee3D("Le personnage ne peut pas utiliser de potion.")
	}
	if c.PvActuel >= c.PvMaxTotal {
		return actionRefusee3D("Les points de vie sont déjà au maximum.")
	}
	if c.Inventaire[ItemPotionDeVie] <= 0 {
		return actionRefusee3D("Aucune potion de vie dans l'inventaire.")
	}

	avant := c.PvActuel
	c.RemoveInventory(ItemPotionDeVie)
	c.PvActuel += SoinPotionDeVie
	if c.PvActuel > c.PvMaxTotal {
		c.PvActuel = c.PvMaxTotal
	}
	return ResultatAction{
		Reussite: true, Type: "objet", Source: c.Nom, Cible: c.Nom,
		Soin: c.PvActuel - avant, PVAvant: avant, PVApres: c.PvActuel,
		Message: fmt.Sprintf("%s récupère %d PV.", c.Nom, c.PvActuel-avant),
	}
}

// utiliserPotionMana3D : Consomme une potion et renvoie le mana réellement récupéré.
func (c *Character) utiliserPotionMana3D() ResultatAction {
	if c.ManaActuel >= c.ManaMax {
		return actionRefusee3D("Le mana est déjà au maximum.")
	}
	if c.Inventaire[ItemPotionDeMana] <= 0 {
		return actionRefusee3D("Aucune potion de mana dans l'inventaire.")
	}

	avant := c.ManaActuel
	c.RemoveInventory(ItemPotionDeMana)
	c.ManaActuel += ManaPotionDeMana
	if c.ManaActuel > c.ManaMax {
		c.ManaActuel = c.ManaMax
	}
	return ResultatAction{
		Reussite: true, Type: "objet", Source: c.Nom, Cible: c.Nom,
		ManaAvant: avant, ManaApres: c.ManaActuel,
		Message: fmt.Sprintf("%s récupère %d mana.", c.Nom, c.ManaActuel-avant),
	}
}

// NouveauGobelinCombat : Prépare les données du gobelin et le chemin de son modèle, sans charger de scène.
func NouveauGobelinCombat() EnnemiCombat {
	return nouvelEnnemi3D(TypeGobelin, "gobelin.json", InitGoblin(1))
}

// NouveauGobelinCuirasseCombat : Prépare la variante résistante du gobelin.
func NouveauGobelinCuirasseCombat() EnnemiCombat {
	return nouvelEnnemi3D(TypeGobelinCuirasse, "gobelin_cuirasse.json", Monster{
		Nom: "Gobelin cuirassé", PvMax: 65, PvActuel: 65, Attaque: 8, ExperienceDonnee: 65,
	})
}

// NouveauChamanCombat : Prépare les statistiques et le modèle du chaman.
func NouveauChamanCombat() EnnemiCombat {
	return nouvelEnnemi3D(TypeChaman, "chaman.json", Monster{
		Nom: "Chaman gobelin", PvMax: 45, PvActuel: 45, Attaque: 7, ExperienceDonnee: 60,
	})
}

// NouveauLoupCombat : Prépare les statistiques et le modèle du loup.
func NouveauLoupCombat() EnnemiCombat {
	return nouvelEnnemi3D(TypeLoup, "loup.json", Monster{
		Nom: "Loup des bois", PvMax: 35, PvActuel: 35, Attaque: 7, ExperienceDonnee: 45,
	})
}

// NouveauTrollCombat : Prépare les statistiques et le modèle du boss troll.
func NouveauTrollCombat() EnnemiCombat {
	return nouvelEnnemi3D(TypeTroll, "troll.json", Monster{
		Nom: "Troll de l'arène", PvMax: 160, PvActuel: 160, Attaque: 14, ExperienceDonnee: 180,
	})
}

// nouvelEnnemi3D : Associe les règles du monstre à son type et à son fichier visuel.
func nouvelEnnemi3D(typeMonstre TypeMonstre, modele string, monstre Monster) EnnemiCombat {
	monstre.Niveau = 1
	// Compléter les récompenses des créatures absentes du backend.
	if monstre.OrMin == 0 && monstre.OrMax == 0 {
		monstre.OrMin = OrMonstre3D(typeMonstre)
		monstre.OrMax = monstre.OrMin * 2
	}
	return EnnemiCombat{Type: typeMonstre, Modele3D: modele, Monstre: monstre}
}

// AdapterEnnemiNiveau3D prépare une copie au niveau demandé, avant le combat.
// Le gobelin suit exactement le backend. Les autres gardent leurs bases et
// utilisent les mêmes gains par niveau. Ne pas appeler sur un ennemi en combat.
func AdapterEnnemiNiveau3D(ennemi EnnemiCombat, niveau int) EnnemiCombat {
	if niveau < 1 {
		niveau = 1
	}
	if ennemi.Type == TypeGobelin {
		ennemi.Monstre = InitGoblin(niveau)
		return ennemi
	}
	m := &ennemi.Monstre
	niveauActuel := m.Niveau
	if niveauActuel < 1 {
		niveauActuel = 1
	}
	bonus := niveau - niveauActuel
	m.Niveau = niveau
	m.PvMax += GainPvGobelin * bonus
	m.PvActuel = m.PvMax
	m.Attaque += GainAttaqueGobelin * bonus
	m.ExperienceDonnee += GainExperienceGobelin * bonus
	m.OrMin += GainOrMinGobelin * bonus
	m.OrMax += GainOrMaxGobelin * bonus
	return ennemi
}

// VaguesAreneParDefaut : Décrit la succession des groupes ennemis de l'arène.
func VaguesAreneParDefaut() [][]EnnemiCombat {
	return [][]EnnemiCombat{
		{NouveauGobelinCombat()},
		{NouveauLoupCombat(), NouveauGobelinCombat()},
		{NouveauGobelinCuirasseCombat(), NouveauChamanCombat()},
		{NouveauTrollCombat()},
	}
}

// NouveauCombatArene : Démarre la préparation du mode arène avec les vagues par défaut.
func NouveauCombatArene(joueur *Character) (*CombatArene, error) {
	return NouveauCombat3D(joueur, ModeArene, TypeGobelin)
}

// EnnemiDuel3D : Prépare la créature demandée pour un duel ou signale un type non disponible.
func EnnemiDuel3D(genre TypeMonstre) (EnnemiCombat, error) {
	switch genre {
	case TypeLoup:
		return NouveauLoupCombat(), nil
	case TypeTroll:
		return NouveauTrollCombat(), nil
	case TypeSanglier:
		return nouvelEnnemi3D(genre, "sanglier.json", Monster{Nom: "Sanglier", PvMax: 30, PvActuel: 30, Attaque: 5, ExperienceDonnee: 30}), nil
	case TypeCorbeau:
		return nouvelEnnemi3D(genre, "corbeau.json", Monster{Nom: "Corbeau", PvMax: 20, PvActuel: 20, Attaque: 4, ExperienceDonnee: 20}), nil
	default:
		return EnnemiCombat{}, fmt.Errorf("créature de duel inconnue")
	}
}

// NouveauCombat3D : Prépare le mode choisi et son premier tour, sans exécuter d'attaque automatiquement.
func NouveauCombat3D(joueur *Character, mode ModeCombat, genre TypeMonstre) (*CombatArene, error) {
	if joueur == nil {
		return nil, fmt.Errorf("le personnage du combat est absent")
	}
	if joueur.PvActuel <= 0 {
		return nil, fmt.Errorf("le personnage doit être vivant pour entrer dans l'arène")
	}

	combat := &CombatArene{
		Joueur: joueur, Mode: mode, NumeroVague: 1, Tour: 1, tirageButin: rand.Intn,
	}
	switch mode {
	case ModeArene:
		combat.Vagues = VaguesAreneParDefaut()
	case ModeDuel:
		ennemi, err := EnnemiDuel3D(genre)
		if err != nil {
			return nil, err
		}
		combat.Vagues = [][]EnnemiCombat{{ennemi}}
	case ModeEntrainement:
		combat.Vagues = [][]EnnemiCombat{{NouveauGobelinCombat()}}
	default:
		return nil, fmt.Errorf("mode de combat inconnu")
	}
	combat.preparerVague(0)
	return combat, nil
}

// Quitter n'accorde rien pour un ennemi vivant et bloque toute nouvelle action.
func (c *CombatArene) Quitter3D() {
	if c == nil || c.Joueur == nil || c.Phase == PhaseAbandon {
		return
	}
	if c.Phase == PhaseDefaite {
		c.RessusciterApresDefaite()
	}
	c.poisonSecondes = 0
	c.Phase = PhaseAbandon
}

// donnerButin : Effectue le tirage du matériau hors entraînement et tente de le ranger dans l'inventaire.
func (c *CombatArene) donnerButin(genre TypeMonstre) string {
	butin := ButinPossible3D(genre)
	if butin.Objet == "" || c.Mode == ModeEntrainement {
		return ""
	}
	tirage := c.tirageButin
	if tirage == nil {
		tirage = rand.Intn
	}
	if tirage(100) >= butin.Pourcentage {
		return " Aucun matériau cette fois."
	}
	if !c.Joueur.AddInventory(butin.Objet) {
		message := "Perdu (inventaire plein) : " + butin.Objet
		c.Butins = append(c.Butins, message)
		return " " + message + "."
	}
	c.Butins = append(c.Butins, butin.Objet)
	return " Butin : " + butin.Objet + "."
}

// AttaquesDisponibles : Copie les attaques apprises pour que l'interface ne modifie pas la liste du personnage.
func (c *CombatArene) AttaquesDisponibles() []string {
	if c == nil || c.Joueur == nil {
		return nil
	}
	return append([]string(nil), c.Joueur.AttaquesPhysiques...)
}

// SortsDisponibles : Copie les sorts connus ; leur coût est vérifié au moment de les lancer.
func (c *CombatArene) SortsDisponibles() []string {
	if c == nil || c.Joueur == nil {
		return nil
	}
	return append([]string(nil), c.Joueur.Skill...)
}

// ObjetsUtilisables : Liste les potions disponibles ; ce raccourci ne remplace pas l'inventaire complet.
func (c *CombatArene) ObjetsUtilisables() []string {
	if c == nil || c.Joueur == nil {
		return nil
	}
	objets := []string{}
	for _, nom := range []string{ItemPotionDeVie, ItemPotionDeMana, ItemPotionDePoison} {
		if c.Joueur.Inventaire[nom] > 0 {
			objets = append(objets, nom)
		}
	}
	sort.Strings(objets)
	return objets
}

// Attaquer : Utilise la première attaque physique connue comme action par défaut.
func (c *CombatArene) Attaquer(indexCible int) ResultatAction {
	attaques := c.AttaquesDisponibles()
	if len(attaques) == 0 {
		return actionRefusee3D("Aucune attaque physique n'est disponible.")
	}
	return c.AttaquerPhysiquement(attaques[0], indexCible)
}

// AttaquerPhysiquement : Valide le tour et la cible, applique les dégâts, puis prépare la suite du combat.
func (c *CombatArene) AttaquerPhysiquement(nomAttaque string, indexCible int) ResultatAction {
	if resultat := c.verifierTourJoueur(indexCible); !resultat.Reussite {
		return resultat
	}
	if !contientTexte3D(c.Joueur.AttaquesPhysiques, nomAttaque) {
		return actionRefusee3D("Cette attaque physique n'est pas connue.")
	}
	degats, existe := c.Joueur.CalculerDegatsPhysique(nomAttaque)
	if !existe {
		return actionRefusee3D("Cette attaque physique n'est pas disponible.")
	}

	cible := &c.Ennemis[indexCible].Monstre
	avant := cible.PvActuel
	cible.PvActuel -= degats
	bornerPVMonstre3D(cible)

	resultat := ResultatAction{
		Reussite: true, TourConsomme: true, Type: "attaque",
		Source: c.Joueur.Nom, Cible: cible.Nom, Degats: avant - cible.PvActuel,
		PVAvant: avant, PVApres: cible.PvActuel, CibleVaincue: cible.PvActuel == 0,
		Message: fmt.Sprintf("%s utilise %s sur %s et inflige %d dégâts.", c.Joueur.Nom, nomAttaque, cible.Nom, avant-cible.PvActuel),
	}
	return c.terminerActionJoueur(indexCible, resultat)
}

// LancerSort : Vérifie le sort appris et le mana disponible avant d'appliquer son effet.
func (c *CombatArene) LancerSort(nomSort string, indexCible int) ResultatAction {
	if c == nil || c.Joueur == nil || c.Phase != PhaseTourJoueur {
		return actionRefusee3D("Ce n'est pas le tour du joueur.")
	}
	if !contientTexte3D(c.Joueur.Skill, nomSort) {
		return actionRefusee3D("Ce sort n'est pas connu.")
	}

	coutMana := coutManaSort3D(nomSort)
	if coutMana < 0 {
		return actionRefusee3D("Ce sort n'est pas disponible.")
	}
	if c.Joueur.ManaActuel < coutMana {
		return actionRefusee3D(fmt.Sprintf("Il manque %d mana.", coutMana-c.Joueur.ManaActuel))
	}

	if resultat := c.verifierTourJoueur(indexCible); !resultat.Reussite {
		return resultat
	}

	cible := &c.Ennemis[indexCible].Monstre
	degats, _, autorise := InfosSort(nomSort)
	if !autorise {
		return actionRefusee3D("Les conditions de ce sort ne sont pas remplies.")
	}

	manaAvant := c.Joueur.ManaActuel
	c.Joueur.ManaActuel -= coutMana
	avant := cible.PvActuel
	cible.PvActuel -= degats
	bornerPVMonstre3D(cible)

	resultat := ResultatAction{
		Reussite: true, TourConsomme: true, Type: "sort",
		Source: c.Joueur.Nom, Cible: cible.Nom, Degats: avant - cible.PvActuel,
		PVAvant: avant, PVApres: cible.PvActuel,
		ManaAvant: manaAvant, ManaApres: c.Joueur.ManaActuel,
		CibleVaincue: cible.PvActuel == 0,
		Message:      fmt.Sprintf("%s lance %s sur %s et inflige %d dégâts.", c.Joueur.Nom, nomSort, cible.Nom, avant-cible.PvActuel),
	}

	return c.terminerActionJoueur(indexCible, resultat)
}

// coutManaSort3D : Retourne le coût du sort, ou -1 si le nom est inconnu ; zéro reste un coût valide.
func coutManaSort3D(nomSort string) int {
	_, cout, connu := InfosSort(nomSort)
	if !connu {
		return -1
	}
	return cout
}

// CoutSort3D expose le coût existant à l'interface, sans lancer le sort.
func CoutSort3D(nom string) int { return coutManaSort3D(nom) }

// UtiliserObjet : Traite l'objet pendant le tour du joueur et ne consomme le tour qu'en cas de réussite.
func (c *CombatArene) UtiliserObjet(nomObjet string) ResultatAction {
	if c == nil || c.Joueur == nil || c.Phase != PhaseTourJoueur {
		return actionRefusee3D("Ce n'est pas le tour du joueur.")
	}

	if _, equipement := TrouverEquipement(nomObjet); equipement {
		return actionRefusee3D("L'équipement ne peut pas être changé pendant le combat.")
	}
	var resultat ResultatAction
	switch nomObjet {
	case ItemPotionDeVie:
		resultat = c.Joueur.utiliserPotionVie3D()
	case ItemPotionDeMana:
		resultat = c.Joueur.utiliserPotionMana3D()
	case ItemPotionDePoison:
		if c.Joueur.Inventaire[ItemPotionDePoison] <= 0 {
			return actionRefusee3D("Aucune potion de poison dans l'inventaire.")
		}
		c.Joueur.RemoveInventory(ItemPotionDePoison)
		c.poisonSecondes = DureePoison
		c.poisonTempsAccumule = 0
		resultat = actionReussie3D(fmt.Sprintf("%s utilise une potion de poison.", c.Joueur.Nom))
	default:
		resultat = c.Joueur.UtiliserObjet3D(nomObjet)
	}

	if resultat.Reussite {
		resultat.TourConsomme = true
		c.commencerTourMonstres()
	}
	return resultat
}

// Defendre : Prépare une réduction sur la prochaine attaque reçue et termine l'action du joueur.
func (c *CombatArene) Defendre() ResultatAction {
	if c == nil || c.Joueur == nil || c.Phase != PhaseTourJoueur {
		return actionRefusee3D("Ce n'est pas le tour du joueur.")
	}
	c.DefenseActive = true
	resultat := actionReussie3D(fmt.Sprintf("%s se met en position de défense.", c.Joueur.Nom))
	resultat.Type = "defense"
	resultat.TourConsomme = true
	c.commencerTourMonstres()
	return resultat
}

// ProchaineActionMonstre : Joue au plus une action ennemie par appel pour laisser la 3D espacer les animations.
func (c *CombatArene) ProchaineActionMonstre() ResultatAction {
	if c == nil || c.Phase != PhaseTourMonstres || c.poisonSecondes > 0 {
		return actionRefusee3D("Les monstres ne peuvent pas jouer maintenant.")
	}
	for c.IndexMonstreActif < len(c.Ennemis) && c.Ennemis[c.IndexMonstreActif].Monstre.PvActuel <= 0 {
		c.IndexMonstreActif++
	}
	if c.IndexMonstreActif >= len(c.Ennemis) {
		c.finirTourMonstres()
		return actionReussie3D("Le tour des monstres est terminé.")
	}

	ennemi := &c.Ennemis[c.IndexMonstreActif]
	if ennemi.Type == TypeChaman && c.Tour%3 == 0 {
		resultat := c.soinChaman3D(ennemi)
		c.IndexMonstreActif++
		if c.IndexMonstreActif >= len(c.Ennemis) {
			c.finirTourMonstres()
		}
		return resultat
	}

	degats := ennemi.Monstre.Attaque
	if c.Tour%3 == 0 && (ennemi.Type == TypeGobelin || ennemi.Type == TypeTroll) {
		degats *= 2
	}
	if c.DefenseActive {
		degats = (degats + 1) / 2
		c.DefenseActive = false
	}

	avant := c.Joueur.PvActuel
	c.Joueur.PvActuel -= degats
	if c.Joueur.PvActuel < 0 {
		c.Joueur.PvActuel = 0
	}

	resultat := ResultatAction{
		Reussite: true, TourConsomme: true, Type: "attaque_monstre",
		Source: ennemi.Monstre.Nom, Cible: c.Joueur.Nom,
		Degats: avant - c.Joueur.PvActuel, PVAvant: avant, PVApres: c.Joueur.PvActuel,
		CibleVaincue: c.Joueur.PvActuel == 0,
		Message:      fmt.Sprintf("%s attaque %s et inflige %d dégâts.", ennemi.Monstre.Nom, c.Joueur.Nom, avant-c.Joueur.PvActuel),
	}

	c.IndexMonstreActif++
	if resultat.CibleVaincue {
		c.Phase = PhaseDefaite
		resultat.CombatTermine = true
		return resultat
	}
	if c.IndexMonstreActif >= len(c.Ennemis) {
		c.finirTourMonstres()
	}
	return resultat
}

// CommencerVagueSuivante : Passe au groupe suivant uniquement pendant la transition entre vagues.
func (c *CombatArene) CommencerVagueSuivante() error {
	if c == nil || c.Phase != PhaseEntreVagues {
		return fmt.Errorf("aucune vague n'est prête à commencer")
	}
	index := c.NumeroVague
	if index >= len(c.Vagues) {
		return fmt.Errorf("il ne reste aucune vague")
	}
	c.NumeroVague++
	c.Tour = 1
	c.preparerVague(index)
	return nil
}

// RessusciterApresDefaite : Rend la moitié des PV au joueur vaincu, sans lancer un nouveau combat.
func (c *CombatArene) RessusciterApresDefaite() ResultatAction {
	if c == nil || c.Joueur == nil || c.Phase != PhaseDefaite {
		return actionRefusee3D("Le personnage n'a pas besoin d'être ressuscité.")
	}
	avant := c.Joueur.PvActuel
	c.Joueur.PvActuel = c.Joueur.PvMaxTotal / 2
	return ResultatAction{
		Reussite: true, Type: "resurrection", Source: c.Joueur.Nom, Cible: c.Joueur.Nom,
		Soin: c.Joueur.PvActuel - avant, PVAvant: avant, PVApres: c.Joueur.PvActuel,
		Message: fmt.Sprintf("%s ressuscite avec %d PV.", c.Joueur.Nom, c.Joueur.PvActuel),
	}
}

// MettreAJourEffets : Fait progresser les effets temporaires avec le temps écoulé, sans bloquer la boucle de rendu.
func (c *CombatArene) MettreAJourEffets(delta time.Duration) []ResultatAction {
	if c == nil || c.poisonSecondes <= 0 || c.Phase == PhaseDefaite || c.Phase == PhaseVictoire || c.Phase == PhaseAbandon {
		return nil
	}
	c.poisonTempsAccumule += delta
	resultats := []ResultatAction{}
	for c.poisonTempsAccumule >= time.Second && c.poisonSecondes > 0 {
		c.poisonTempsAccumule -= time.Second
		c.poisonSecondes--
		avant := c.Joueur.PvActuel
		c.Joueur.PvActuel -= DegatsPoisonParSeconde
		if c.Joueur.PvActuel < 0 {
			c.Joueur.PvActuel = 0
		}
		resultat := ResultatAction{
			Reussite: true, Type: "poison", Source: ItemPotionDePoison, Cible: c.Joueur.Nom,
			Degats: avant - c.Joueur.PvActuel, PVAvant: avant, PVApres: c.Joueur.PvActuel,
			CibleVaincue: c.Joueur.PvActuel == 0,
			Message:      fmt.Sprintf("Le poison inflige %d dégâts à %s.", avant-c.Joueur.PvActuel, c.Joueur.Nom),
		}
		if resultat.CibleVaincue {
			c.Phase = PhaseDefaite
			resultat.CombatTermine = true
			c.poisonSecondes = 0
		}
		resultats = append(resultats, resultat)
	}
	return resultats
}

// CiblesVivantes : Retourne les indices encore valides, conservés pour retrouver les modèles 3D correspondants.
func (c *CombatArene) CiblesVivantes() []int {
	indices := []int{}
	if c == nil {
		return indices
	}
	for index := range c.Ennemis {
		if c.Ennemis[index].Monstre.PvActuel > 0 {
			indices = append(indices, index)
		}
	}
	return indices
}

// preparerVague : Copie les ennemis et compare les initiatives pour choisir le premier camp.
func (c *CombatArene) preparerVague(index int) {
	c.Ennemis = copierEnnemis3D(c.Vagues[index])
	c.IndexMonstreActif = 0
	c.Joueur.Initiative = rand.Intn(10) + 1
	initiativeEnnemie := 0
	for index := range c.Ennemis {
		c.Ennemis[index] = AdapterEnnemiNiveau3D(c.Ennemis[index], c.Joueur.Niveau)
		c.Ennemis[index].Monstre.Initiative = rand.Intn(10) + 1
		if c.Ennemis[index].Monstre.Initiative > initiativeEnnemie {
			initiativeEnnemie = c.Ennemis[index].Monstre.Initiative
		}
	}
	c.joueurCommence = c.Joueur.Initiative >= initiativeEnnemie
	if c.joueurCommence {
		c.Phase = PhaseTourJoueur
	} else {
		c.Phase = PhaseTourMonstres
	}
}

// verifierTourJoueur : Contrôle le tour et la cible sans modifier le combat.
func (c *CombatArene) verifierTourJoueur(indexCible int) ResultatAction {
	if c == nil || c.Joueur == nil {
		return actionRefusee3D("Le combat n'est pas initialisé.")
	}
	if c.Phase != PhaseTourJoueur {
		return actionRefusee3D("Ce n'est pas le tour du joueur.")
	}
	if indexCible < 0 || indexCible >= len(c.Ennemis) {
		return actionRefusee3D("La cible sélectionnée n'existe pas.")
	}
	if c.Ennemis[indexCible].Monstre.PvActuel <= 0 {
		return actionRefusee3D("La cible est déjà vaincue.")
	}
	return ResultatAction{Reussite: true}
}

// terminerActionJoueur : Attribue les gains d'une élimination puis choisit entre ennemis, vague suivante et victoire.
func (c *CombatArene) terminerActionJoueur(indexCible int, resultat ResultatAction) ResultatAction {
	if !resultat.Reussite || !resultat.TourConsomme {
		return resultat
	}
	if resultat.CibleVaincue {
		gain := c.Ennemis[indexCible].Monstre.ExperienceDonnee
		resultat.ExperienceGagnee = gain
		resultat.NiveauxGagnes = ajouterExperience3D(c.Joueur, gain)
		c.ExperienceTotale += gain
		// Même règle que FinCombat : borne supérieure comprise, entraînement inclus.
		monstre := c.Ennemis[indexCible].Monstre
		resultat.OrGagne = monstre.OrMin + rand.Intn(monstre.OrMax-monstre.OrMin+1)
		c.Joueur.Argent += resultat.OrGagne
		c.OrTotal += resultat.OrGagne
		resultat.Message += fmt.Sprintf(" +%d or.", resultat.OrGagne)
		resultat.Message += c.donnerButin(c.Ennemis[indexCible].Type)
	}
	if len(c.CiblesVivantes()) == 0 {
		resultat.VagueTerminee = true
		if c.NumeroVague >= len(c.Vagues) {
			c.Phase = PhaseVictoire
			resultat.CombatTermine = true
			resultat.Victoire = true
		} else {
			c.Phase = PhaseEntreVagues
		}
		return resultat
	}
	c.commencerTourMonstres()
	return resultat
}

// commencerTourMonstres : Passe aux ennemis et ajuste le numéro du round selon le camp qui avait commencé.
func (c *CombatArene) commencerTourMonstres() {
	if !c.joueurCommence {
		c.Tour++
	}
	c.Phase = PhaseTourMonstres
	c.IndexMonstreActif = 0
}

// finirTourMonstres : Rend la main au joueur en gardant un seul incrément de round pour les deux camps.
func (c *CombatArene) finirTourMonstres() {
	if c.joueurCommence {
		c.Tour++
	}
	c.Phase = PhaseTourJoueur
	c.IndexMonstreActif = 0
}

// soinChaman3D : Soigne l'allié vivant blessé qui possède le moins de PV actuels.
func (c *CombatArene) soinChaman3D(chaman *EnnemiCombat) ResultatAction {
	indexCible := -1
	for index := range c.Ennemis {
		monstre := &c.Ennemis[index].Monstre
		if monstre.PvActuel <= 0 || monstre.PvActuel >= monstre.PvMax {
			continue
		}
		if indexCible == -1 || monstre.PvActuel < c.Ennemis[indexCible].Monstre.PvActuel {
			indexCible = index
		}
	}
	if indexCible == -1 {
		return actionReussie3D("Le chaman tente un soin, mais aucun monstre n'est blessé.")
	}

	cible := &c.Ennemis[indexCible].Monstre
	avant := cible.PvActuel
	cible.PvActuel += 12
	bornerPVMonstre3D(cible)
	return ResultatAction{
		Reussite: true, TourConsomme: true, Type: "soin_monstre",
		Source: chaman.Monstre.Nom, Cible: cible.Nom, Soin: cible.PvActuel - avant,
		PVAvant: avant, PVApres: cible.PvActuel,
		Message: fmt.Sprintf("%s rend %d PV à %s.", chaman.Monstre.Nom, cible.PvActuel-avant, cible.Nom),
	}
}

// ajouterExperience3D : Applique les montées de niveau sans dialogue terminal et conserve l'expérience excédentaire.
func ajouterExperience3D(personnage *Character, gain int) int {
	if personnage == nil || gain <= 0 {
		return 0
	}
	niveaux := 0
	personnage.ExperienceActuelle += gain
	for personnage.ExperienceActuelle >= personnage.ExperienceMax {
		personnage.ExperienceActuelle -= personnage.ExperienceMax
		personnage.Niveau++
		personnage.ExperienceMax = personnage.ExperienceMax * AugmentationExperienceNumerateur / AugmentationExperienceDenominateur
		personnage.PvMaxBase += personnage.GainPvMax
		personnage.MettreAJourPvMax()
		personnage.ManaMax += personnage.GainManaMax
		personnage.Attaque += personnage.GainAttaque
		niveaux++
	}
	return niveaux
}

// copierEnnemis3D : Sépare les ennemis actifs des définitions de vagues pour ne pas conserver leurs blessures.
func copierEnnemis3D(source []EnnemiCombat) []EnnemiCombat {
	resultat := make([]EnnemiCombat, len(source))
	copy(resultat, source)
	return resultat
}

// bornerPVMonstre3D : Maintient les PV du monstre entre zéro et son maximum.
func bornerPVMonstre3D(monstre *Monster) {
	if monstre.PvActuel < 0 {
		monstre.PvActuel = 0
	}
	if monstre.PvActuel > monstre.PvMax {
		monstre.PvActuel = monstre.PvMax
	}
}

// contientTexte3D : Vérifie une correspondance exacte dans une liste de noms.
func contientTexte3D(liste []string, valeur string) bool {
	for _, element := range liste {
		if element == valeur {
			return true
		}
	}
	return false
}

// actionReussie3D : Crée un résultat positif ; l'appelant précise séparément si le tour est consommé.
func actionReussie3D(message string) ResultatAction {
	return ResultatAction{Reussite: true, Message: message}
}

// actionRefusee3D : Crée un refus lisible sans consommer de tour.
func actionRefusee3D(message string) ResultatAction {
	return ResultatAction{Reussite: false, Message: message}
}

// CreerPersonnage3D valide les mêmes noms et classes que la création CLI.
func CreerPersonnage3D(nom, nomClasse string) (Character, error) {
	nom, valide := FormatName(nom)
	if !valide {
		return Character{}, fmt.Errorf("Le nom doit contenir uniquement des lettres.")
	}
	classe, valide := ClasseDepuisNom3D(nomClasse)
	if !valide {
		return Character{}, fmt.Errorf("Choisissez Humain, Elfe ou Nain.")
	}
	return InitCharacter(nom, 1, classe), nil
}

// EffetEnCours3D : Indique si le poison doit encore être résolu avant la prochaine action ennemie.
func (c *CombatArene) EffetEnCours3D() bool { return c != nil && c.poisonSecondes > 0 }

// EffetsPersonnage3D conserve le chronomètre du poison hors combat.
// L'interface appelle MettreAJour à chaque image, jamais time.Sleep.
type EffetsPersonnage3D struct {
	Joueur   *Character
	secondes int
	temps    time.Duration
}

// Actif : Indique si un effet temporaire reste en cours pendant l'exploration.
func (e *EffetsPersonnage3D) Actif() bool { return e.secondes > 0 }

// Utiliser : Démarre le poison dans le temps ou délègue les objets à effet immédiat.
func (e *EffetsPersonnage3D) Utiliser(nom string) ResultatAction {
	if e.Joueur == nil {
		return actionRefusee3D("Personnage absent.")
	}
	if e.Actif() {
		return actionRefusee3D("Attendez la fin du poison.")
	}
	if nom != ItemPotionDePoison {
		return e.Joueur.UtiliserObjet3D(nom)
	}
	if !e.Joueur.RemoveInventory(nom) {
		return actionRefusee3D("Aucune potion de poison.")
	}
	e.secondes, e.temps = DureePoison, 0
	return actionReussie3D(fmt.Sprintf("Poison : %d dégâts par seconde pendant %d secondes.", DegatsPoisonParSeconde, DureePoison))
}

// MettreAJour : Résout les secondes de poison écoulées hors combat et gère la résurrection si nécessaire.
func (e *EffetsPersonnage3D) MettreAJour(delta time.Duration) []ResultatAction {
	if !e.Actif() || e.Joueur == nil {
		return nil
	}
	e.temps += delta
	var resultats []ResultatAction
	for e.temps >= time.Second && e.secondes > 0 {
		e.temps -= time.Second
		e.secondes--
		avant := e.Joueur.PvActuel
		e.Joueur.SubirDegats(DegatsPoisonParSeconde)
		r := ResultatAction{Reussite: true, Type: "poison", PVAvant: avant, PVApres: e.Joueur.PvActuel, Degats: avant - e.Joueur.PvActuel}
		r.Message = fmt.Sprintf("Poison : -%d PV (%d / %d).", r.Degats, e.Joueur.PvActuel, e.Joueur.PvMaxTotal)
		if e.Joueur.PvActuel == 0 {
			e.secondes = 0
			e.Joueur.PvActuel = e.Joueur.PvMaxTotal / 2
			r.Message += fmt.Sprintf(" Résurrection : %d PV.", e.Joueur.PvActuel)
		}
		resultats = append(resultats, r)
	}
	return resultats
}
