package library

// ================================================== //
// === Ce fichier regroupe les réglages du jeu    === //
// ================================================== //

// ------------------------------------------------ //
// --- Personnage                               --- //
// ------------------------------------------------ //

const (
	CapaciteInventaireDepart    = 10
	ArgentDepart                = 100
	MaxAugmentationsInventaire  = 3
	BonusAugmentationInventaire = 10
)

// ------------------------------------------------ //
// --- Monstres                                 --- //
// ------------------------------------------------ //

// Gobelin : valeurs au niveau 1 //
const (
	PvGobelin         = 40
	AttaqueGobelin    = 5
	ExperienceGobelin = 50
	OrMinGobelin      = 5
	OrMaxGobelin      = 10
	// Gain par niveau supplémentaire //
	GainPvGobelin         = 10
	GainAttaqueGobelin    = 2
	GainExperienceGobelin = 20
	GainOrMinGobelin      = 3
	GainOrMaxGobelin      = 5
)

// ------------------------------------------------ //
// --- Expérience                               --- //
// ------------------------------------------------ //

// L'expérience requise est multipliée par 13/10 (x1,3) à chaque niveau //
const (
	Experienceinitiale                 = 0
	Experiencemaximale                 = 100
	AugmentationExperienceNumerateur   = 13
	AugmentationExperienceDenominateur = 10
)

// ------------------------------------------------ //
// --- Objets : consommables                    --- //
// ------------------------------------------------ //

const (
	ItemPotionDeVie    = "Potion de vie"
	ItemPotionDeMana   = "Potion de mana"
	ItemPotionDePoison = "Potion de poison"
)

// Effets des consommables //
const (
	SoinPotionDeVie        = 50
	ManaPotionDeMana       = 40
	DegatsPoisonParSeconde = 10
	DureePoison            = 3
)

// ------------------------------------------------ //
// --- Objets : matériaux de fabrication        --- //
// ------------------------------------------------ //

const (
	ItemFourrureDeLoup = "Fourrure de loup"
	ItemPeauDeTroll    = "Peau de troll"
	ItemCuirDeSanglier = "Cuir de sanglier"
	ItemPlumeDeCorbeau = "Plume de corbeau"
)

// ------------------------------------------------ //
// --- Objets : équipements                     --- //
// ------------------------------------------------ //

const (
	ItemChapeauAventurier = "Chapeau de l'aventurier"
	ItemTuniqueAventurier = "Tunique de l'aventurier"
	ItemBottesAventurier  = "Bottes de l'aventurier"
)

// Emplacements où un équipement peut être porté //
const (
	EmplacementTete  = "tête"
	EmplacementTorse = "torse"
	EmplacementPied  = "pied"
)

// ------------------------------------------------ //
// --- Objets : divers                          --- //
// ------------------------------------------------ //

const (
	ItemAugmentationInventaire = "Augmentation d'inventaire"
)

// ------------------------------------------------ //
// --- Livres de sorts                          --- //
// ------------------------------------------------ //

const (
	ItemLivreGrosseBouleDeFeu  = "Livre de sort : Grosse boule de feu"
	ItemLivreLameDuDestin      = "Livre de sort : Lame du destin"
	ItemLivreEclatsDuGardien   = "Livre de sort : Éclats du gardien"
	ItemLivreFlecheDeLumiere   = "Livre de sort : Flèche de lumière"
	ItemLivreJugementDesGeants = "Livre de sort : Jugement des géants"
)

// ------------------------------------------------ //
// --- Manuels de combat                        --- //
// ------------------------------------------------ //

const (
	ItemLivreAttaquePichenette   = "Manuel de combat : Pichenette"
	ItemLivreAttaqueClaquounette = "Manuel de combat : Claquounette"
	ItemLivreAttaqueCoupsDePied  = "Manuel de combat : Coups de pied"
	ItemLivreAttaqueMorsure      = "Manuel de combat : Morsure"
	ItemLivreAttaqueUppercut     = "Manuel de combat : Uppercut"
)

// ------------------------------------------------ //
// --- Sorts : nom, dégâts, coût en mana        --- //
// ------------------------------------------------ //

const (
	SortCoupDePoing     = "Coup de poing"
	DegatsCoupDePoing   = 8
	CoutManaCoupDePoing = 5

	SortLameDuDestin     = "Lame du destin"
	DegatsLameDuDestin   = 12
	CoutManaLameDuDestin = 10

	SortEclatsDuGardien     = "Éclats du gardien"
	DegatsEclatsDuGardien   = 16
	CoutManaEclatsDuGardien = 15

	SortFlecheDeLumiere     = "Flèche de lumière"
	DegatsFlecheDeLumiere   = 20
	CoutManaFlecheDeLumiere = 20

	SortGrosseBouleDeFeu     = "Grosse boule de feu"
	DegatsGrosseBouleDeFeu   = 26
	CoutManaGrosseBouleDeFeu = 30

	SortJugementDesGeants     = "Jugement des géants"
	DegatsJugementDesGeants   = 60
	CoutManaJugementDesGeants = 60
)

// Sorts prévus, pas encore implémentés //
const (
	SortSoinDuCoeur     = "Soin du coeur"
	CoutManaSoinDuCoeur = 25

	SortBouclier     = "Bouclier"
	CoutManaBouclier = 30
)

// ------------------------------------------------ //
// --- Attaques physiques : nom, dégâts de base --- //
// ------------------------------------------------ //

// Le bonus d'attaque du joueur s'ajoute à ces dégâts //
const (
	AttaqueBasique       = "Attaque basique"
	DegatsAttaqueBasique = 5

	AttaquePichenette = "Pichenette"
	DegatsPichenette  = 7

	AttaqueClaquounette = "Claquounette"
	DegatsClaquounette  = 9

	AttaqueCoupsDePied = "Coups de pied"
	DegatsCoupsDePied  = 11

	AttaqueMorsure = "Morsure"
	DegatsMorsure  = 13

	AttaqueUppercut = "Uppercut"
	DegatsUppercut  = 16
)
