# Personnage 3D

- `modele.go` assemble le corps et les quatre pivots : deux jambes et deux bras.
- `animations.go` contient l'animation de marche écrite par le développeur du jeu.
- `chargement.go` lit les maillages fixes et les place dans la scène.
- `materiaux.go` charge les couleurs et les textures fixes.
- `assets/personnage.json` contient les sommets du modèle, directement en Z vertical.
- `assets/*.png` contient les quatre textures du modèle.

La scène, le modèle et les pivots utilisent tous Z comme axe vertical. Les formes et textures ne sont plus fabriquées pendant l'exécution. Le personnage garde la même géométrie, les mêmes couleurs et les mêmes textures que la version précédente. Pour animer un membre, modifier son pivot dans `animations.go` ; ne pas déplacer les sommets du fichier JSON.
