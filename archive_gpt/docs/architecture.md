# Structure 3D - premiere tranche

`test_game` est un exemple fige : ne pas le modifier. Le vrai module Go est dans `src`. Les assets partages restent a la racine dans `assets`.

- `src/cmd/red3d` : entree graphique.
- `src/internal/game3d` : fenetre G3N, camera, lumiere, geometrie, controles et menus 2D.
- `src/internal/game` : contrat Backend, actions et donnees de presentation ; moteur Demo temporaire.

Le rendu appelle `Backend.Execute(Action)` puis lit `Backend.Snapshot()`. Les fonctions finales seront branchees via un nouvel adaptateur implementant Backend. Le futur CLI utilisera le meme contrat, sans importer G3N.

## Monde semi-ouvert

Le village n'est plus enferme dans une arene. La camera suit le personnage et le monde procedural se prolonge par secteurs de 12 unites, y compris dans les coordonnees negatives. Le moteur garde 25 secteurs actifs autour du joueur, dont 9 avec des objets plus detailles. Les secteurs hors de cette zone sont retires du rendu et leurs donnees mises en cache sont liberees. Le contenu est deterministe : revenir dans un secteur regenere les memes obstacles.

Les donnees spatiales vivent dans `src/internal/worldmap`, independamment des objets G3N. `StarterMap` definit le village, trois routes et leurs destinations : camp des loups a l'ouest, cache des collines a l'est et ancien sanctuaire au sud. Les chemins sont reserves dans la vegetation, visibles en 3D et representes dans la carte 2D. Chaque destination peut etre fouillee une fois (`G`) pour obtenir des materiaux de fabrication. Le terrain au-dela reste procedural : prairie, foret et hautes terres. Les collisions utilisent le meme placement deterministe que le rendu. Il s'agit d'une petite carte de depart, pas encore de zones de quete completes.

Le moteur Go garde les regles, la position logique et la sauvegarde. G3N ne garde que les objets visuels ; le futur CLI pourra lire les memes donnees sans OpenGL.

## Qualite visuelle visee

- Les materiaux de sol, pierre, bois, tuiles, tissu et metal utilisent maintenant des textures procedurales generees en Go. L'eclairage directionnel, les points lumineux de forge et de lampes, ainsi que des ombres de contact legeres sont presents.
- Ajouter des UV corrects aux modeles, des variantes de materiaux, davantage de details de silhouettes et des animations completes.
- Il reste a realiser de vraies ombres projetees, du brouillard de distance, des animations completes et des materiaux peints a la main si la direction artistique le demande.
- Prevoir plusieurs niveaux de detail et limiter le nombre de maillages/draw calls dans les zones visibles. Tester les performances sur la machine cible.
- Garder les icones 2D de `assets` separees des textures et modeles 3D ; un sprite pixel art n'est pas automatiquement une texture de surface 3D.

Le monde semi-ouvert constitue une base jouable, pas un jeu fini. Les regles commerciales, la forge et le combat utilisent encore un backend temporaire ; quetes, sauvegardes, son, transitions de zones et equilibrage restent a faire.

Le moteur provisoire ne constitue pas la realisation complete des quetes. Il permet de creer un personnage avec nom et classe, d'acheter potions et materiaux, de fabriquer puis d'equiper les trois pieces d'armure, et de combattre le gobelin au tour par tour. L'inventaire est une grille de type coffre RPG sans barre rapide : deux clics deplacent ou permutent les piles, les emplacements d'armure servent a equiper et retirer les pieces. Les bonus de PV maximum respectent le sujet : chapeau +10, tunique +25, bottes +15. Les icones PNG de `assets/ui/icons/png` sont chargees par les menus G3N. Aucun degat ne depend du temps de rendu. Les personnages et le village 3D sont generes en Go ; les OBJ ne sont pas encore charges.

Controles : ZQSD/WASD ou fleches pour marcher ; M marchand, F forge, T arene (a proximite) ; G fouiller une destination ; I inventaire et armure, C personnage, B carte ; E potion, Espace attaque en combat, R fuite ; Echap ouvre/ferme le menu. Les panneaux sont cliquables a la souris. Le menu d'accueil permet de creer un personnage ou de quitter.

Compilation depuis `src` : `go build -o bin/red3d.exe ./cmd/red3d`. G3N requiert gcc amd64, CGO_ENABLED=1, OpenGL et les DLL audio. La verification de compilation ne valide pas le rendu sur une machine sans execution graphique.

La premiere compilation est disponible dans `src/bin/red3d.exe`, accompagnee des DLL audio. Lancer cet executable pour verifier visuellement la scene. Le README racine reste vide, comme demande precedemment.
