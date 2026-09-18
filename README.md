# Projet RED

Mini jeu de rôle en ligne de commande, écrit en Go, réalisé dans le cadre du
Projet RED (Ynov). Le joueur crée son personnage, consulte ses informations,
utilise les objets de son inventaire et achète de nouveaux objets chez un
marchand.

Cette version couvre les tâches 1 à 13 du sujet (voir `docs/`).

## Prérequis

- [Go](https://go.dev/dl/) **1.21 ou plus récent** (`go version` pour vérifier)
- Git

Le jeu n'utilise que la bibliothèque standard de Go : aucune dépendance à
télécharger.

## Installation

```bash
git clone https://github.com/gauhthierh/Projet-Red.git
cd Projet-Red
```

## Lancement

Le code source et le fichier `go.mod` se trouvent dans `src/` : le jeu se
lance depuis ce dossier.

```bash
cd src
go run .
```

Pour produire un exécutable à la place :

```bash
cd src
go build
```

Dans les menus, on tape le numéro de son choix puis Entrée. Le choix
« Quitter » du menu principal ferme le jeu, tout comme Ctrl+D (Linux, macOS)
ou Ctrl+Z puis Entrée (Windows).

## Tests

Depuis le dossier `src/` :

```bash
go test ./...
```

Les tests de la potion de poison attendent réellement une seconde par tour de
dégâts. Pour les sauter et obtenir un résultat immédiat :

```bash
go test -short ./...
```

Ajouter `-v` affiche le détail de chaque test.

## Fonctionnalités

| Tâche | Fonctionnalité | Où la trouver |
|---|---|---|
| 1 | Structure `Character` : nom, classe, niveau, PV max, PV actuels, inventaire | `character.go` |
| 2 | `initCharacter` construit un personnage | `character.go` |
| 3 | `displayInfo` affiche les informations du personnage | `character.go` |
| 4 | `accessInventory` affiche les objets possédés, triés par nom | `inventory.go` |
| 5 | `takePot` : une potion de vie rend 50 PV, sans dépasser les PV max ; refusée à pleins PV | `items.go` |
| 6 | Menu principal (`switch`) avec choix « Retour » dans les sous-menus | `menu.go` |
| 7 | Marchand ; `addInventory` et `removeInventory` | `merchant.go`, `inventory.go` |
| 8 | `isDead` : à 0 PV ou moins, mort puis résurrection à 50 % des PV max | `character.go` |
| 9 | `poisonPot` : 10 dégâts par seconde pendant 3 s, arrêt à la mort ; potion vendue par le marchand | `items.go` |
| 10 | Sorts (`Skill`) : « Coup de poing » au départ ; `spellBook` apprend « Boule de Feu » une seule fois ; livre vendu par le marchand | `items.go` |
| 11 | `characterCreation` : nom en lettres uniquement (accents acceptés), remis en forme ; classe Humain (100 PV), Elfe (80) ou Nain (120) ; départ à 50 % des PV | `creation.go` |
| 12 | Inventaire limité à 10 objets (somme des quantités), capacité stockée dans le personnage | `inventory.go` |
| 13 | 100 pièces d'or au départ | `character.go` |

Les achats chez le marchand sont gratuits : les prix arrivent avec la tâche 14.

## Organisation du code

Tout le code appartient au package `main`, réparti en plusieurs fichiers :

```
src/
├── go.mod          module projet-red, Go 1.21
├── main.go         point d'entrée : création du personnage puis menu principal
├── character.go    structure Character, initCharacter, displayInfo, isDead
├── creation.go     characterCreation, validation du nom, choix de la classe
├── inventory.go    affichage, ajout, retrait, limite et menu de l'inventaire
├── items.go        noms des objets et des sorts, effets des objets
├── merchant.go     menu du marchand et achats
├── menu.go         menu principal et choix « Retour »
├── input.go        lecture des saisies clavier (ligne entière)
└── *_test.go       tests unitaires, un fichier par fichier source
```

## Équipe

Membres, tels qu'ils apparaissent comme auteurs dans l'historique Git :

- gauhthierh
- GuillaumeLarre
- lalie-droid
- quent1206
