package library

import "fmt"

// =================================================== //
// === Ce fichier contient les textes de l'univers === //
// =================================================== //

// Raconte le réveil du gardien au sanctuaire, avant la création du personnage //
func AfficherIntroduction() {
	// Les accents graves (`) délimitent un texte sur plusieurs lignes //
	fmt.Print(`
==================================================
                ORDTA : LE RÉVEIL
==================================================

Vous ouvrez les yeux dans l'obscurité du sanctuaire des gardiens.
Les cristaux d'équilibre luisent faiblement, d'une lueur bleutée.
L'air est lourd, comme si le monde entier retenait son souffle.

Une voix grave résonne dans votre esprit :

  « Jeune gardien... Vous l'avez ressenti, n'est-ce pas ?
    Le Cœur d'Équilibre s'affaiblit.
    Une ombre s'étend, et l'un des nôtres a failli à ses devoirs. »

Devant vous se dresse le Géant de la Sagesse,
vêtu d'une robe tissée de fils d'argent.

  « Ordta n'est plus ce qu'il était. Le danger grandit.
    Mais avant tout... dites-moi qui vous êtes. »
`)
	ReadLine("Appuyez sur Entrée pour continuer...")
}

// Présente les trois peuples d'Ordta avant le choix de la classe //
func AfficherDescriptionClasses() {
	fmt.Print(`
Trois peuples veillent sur Ordta :

  1. Humain : polyvalent, il s'adapte à toutes les batailles.
  2. Elfe   : fragile, mais en communion avec les cristaux d'équilibre.
              Sa magie est la plus puissante des gardiens.
  3. Nain   : solide comme la pierre des anciens gardiens.
              Il frappe fort, mais maîtrise peu la magie.
`)
}

// Présente l'entraînement de Grei Thau avant le combat //
func AfficherIntroEntrainement() {
	fmt.Println(`
Grei Thau vous attend dans la salle d'armes du sanctuaire.

  « Le géant dit que des ennemis approchent. Avant qu'ils n'arrivent,
    montre-moi ce que tu vaux, gardien. Ce gobelin fera l'affaire. »`)
}

// Réplique de Grei Thau après une victoire à l'entraînement //
func AfficherVictoireEntrainement() {
	fmt.Println(`
Grei Thau hoche la tête.
  « Pas mal. Le Cœur d'Équilibre répond à ta force. »`)
}

// Réplique de Grei Thau après une défaite : elle explique la résurrection //
func AfficherDefaiteEntrainement() {
	fmt.Println(`
Grei Thau vous tend la main pour vous relever.
  « Le Cœur d'Équilibre ne laisse pas tomber ses gardiens.
    Debout. On recommence quand tu veux. »`)
}

// Accueil au marché, à l'ouverture du marchand //
func AfficherAccueilMarchand() {
	fmt.Println(`Une vieille marchande lève les yeux de ses étals.
  « Des potions, des grimoires, des matériaux...
    Tout ce qu'il faut à un gardien, si sa bourse le permet. »`)
}

// Accueil à la forge, à l'ouverture du forgeron //
func AfficherAccueilForgeron() {
	fmt.Println(`Le forgeron essuie son front noirci de suie.
  « Apporte-moi des matériaux, gardien,
    et je te forgerai de quoi survivre à ce qui arrive. »`)
}
