# Lancer le jeu

Depuis la racine du depot :

```powershell
go run ./src/lancer.go
```

Ou depuis `src` : `go run lancer.go`.

Go doit etre installe. Au premier lancement, le lanceur telecharge le compilateur portable Windows dans `src/.tools`, verifie sa signature SHA-256, compile le jeu, prepare les DLL audio et ouvre la fenetre. Une connexion internet est necessaire pour les dependances initiales. Aucun fichier de `test_game` n'est utilise par le lanceur.

Pour compiler sans ouvrir la fenetre : `go run ./src/lancer.go -run=false`.

Le code source et la modelisation procedurale sont en Go. G3N utilise neanmoins cgo, OpenGL et des DLL audio natives : cela ne signifie pas un executable sans dependances natives.

# Organisation

- `src/lancer.go` : lancement, preparation et compilation.
- `src/cmd/red3d/main.go` : point d'entree de l'application graphique.
- `src/internal/game` : contrat des regles et simulation temporaire.
- `src/internal/game3d/app.go` : coordination de la fenetre, camera et controles.
- `src/internal/game3d/menus.go` : accueil, marchand, forge, personnage, combat et carte 2D.
- `src/internal/game3d/inventory_ui.go` : inventaire en grille et equipement, sans barre rapide.
- `src/internal/game3d/creation_ui.go` : saisie du nom et choix de la classe.
- `src/internal/game3d/hero_appearance.go` : differences de silhouette des classes et armures sur le heros 3D.
- `src/internal/game3d/characters.go` : personnages articules et animations.
- `src/internal/game3d/village.go` : maisons, marche et forge du village.
- `src/internal/game3d/sectors.go` : terrain et vegetation charges autour du joueur avec deux niveaux de detail.
- `src/internal/game3d/materials.go` : textures procedurales des materiaux, generees en Go.
- `src/internal/game3d/shapes.go` : formes, palette et ombres de contact partagees.
- `src/internal/worldmap` : regions, placement deterministe, collisions et coordonnees des secteurs.
- `src/internal/worldmap/landmarks.go` : carte de depart, positions des lieux et sorties du village.
- `assets` : fichiers artistiques partages ; les modeles de cette scene sont actuellement generes en Go.
- `docs` : documentation.
- `test_game` : exemple fige.

Le decor actuel est une premiere passe low-poly avec secteurs charges dynamiquement, textures procedurales, collisions de base et ombres de contact. Les menus 2D chargent les icones existantes du dossier `assets`. Ce n'est pas une production finale : ombres projetees, collisions complexes, quetes et sauvegardes restent a realiser.

Le menu d'accueil propose la creation d'un personnage : nom compose seulement de lettres, classe Humain (100 PV max), Elfe (80) ou Nain (120), avec la moitie des PV au depart, niveau 1, 100 pieces d'or et 3 potions. Le marche (`M`) et la forge (`F`) demandent de s'en approcher ; l'entrainement (`T`) se declenche pres de l'arene. Les trois chemins du village menent chacun a un site a fouiller avec `G` pour obtenir des materiaux. `I` ouvre un inventaire en grille sans barre rapide : cliquer sur une pile puis sur sa destination permet de la deplacer ou de l'echanger ; cliquer sur une armure dans le sac puis sur son emplacement permet de l'equiper. `E` utilise une potion, `B` ouvre la carte, `C` la fiche du personnage. Pendant le combat, les boutons permettent d'attaquer, boire une potion ou fuir.
