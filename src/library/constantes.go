package library

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

const (
	ItemPotionDeVie            = "Potion de vie"
	ItemPotionDePoison         = "Potion de poison"
	ItemLivreBouleDeFeu        = "Livre de sort : Boule de feu"
	ItemFourrureDeLoup         = "Fourrure de loup"
	ItemPeauDeTroll            = "Peau de Troll"
	ItemCuirDeSanglier         = "Cuir de sanglier"
	ItemPlumeDeCorbeau         = "Plume de corbeau"
	ItemChapeauAventurier      = "Chapeau de l'aventurier"
	ItemTuniqueAventurier      = "Tunique de l'aventurier"
	ItemBottesAventurier       = "Bottes de l'aventurier"
	EmplacementTete            = "tête"
	EmplacementTorse           = "torse"
	EmplacementPied            = "pied"
	ItemAugmentationInventaire = "Augmentation d'inventaire"
	ItemPotionDeMana           = "Potion de Mana"
	ItemLivreLameDuDestin      = "Livre de sort : Lame du destin"
	ItemLivreEclatsDuGardien   = "Livre de sort : Éclats du gardien"
	ItemLivreFlecheDeLumiere   = "Livre de sort : Flèche de lumière"
	ItemLivreJugementDesGeants = "Livre de sort : Jugement des géants"
)

const (
	SortCoupDePoing       = "Coup de poing"
	SortBouleDeFeu        = "Boule de feu"
	SortLameDuDestin      = "Lame du destin"
	SortEclatDuGardien    = "Éclats du gardien"
	SortFlecheDeLumiere   = "Flèche de lumière"
	SortFoudreCeleste     = "Foudre Céleste"
	SortSoinDuCoeur       = "Soin du coeur"
	SortBouclier          = "Bouclier"
	SortDevotion          = "Dévotion"
	SortDernierEspoir     = "Dernier espoir"
	SortJugementDesGeants = "Jugement des géants"
	degatsCoupDePoing     = 8
	degatsBouleDeFeu      = 22
	degatsLameDuDestin    = 10
	degatsEclatDuGardien  = 12
	degatsFlecheDeLumiere = 12
	// degatFoudreCeleste = 10 par tour pendant 3 tours
	// soinDuCoeur = 20% de pv restauré
	// bouclier contre la prochaine attaque
	// devotion inflige entre 20 et 30 dégats et peut infliger 8 de dégats au lanceur
	// dernier espoir ne peut être lancé qu'a moins de 30% PV et fait 40 de dégâts
	degatsJugementDesGeants   = 60
	coutManaCoupDePoing       = 5
	coutManaBouleDeFeu        = 30
	coutManaLameDuDestin      = 10
	coutManaEclatDuGardien    = 15
	coutManaFlecheDeLumiere   = 20
	coutManaFoudreCeleste     = 20
	coutManaSoinDuCoeur       = 25
	coutManaBouclier          = 30
	coutManaDevotion          = 35
	coutManaDernierEspoir     = 50
	coutManaJugementDesGeants = 60
)
