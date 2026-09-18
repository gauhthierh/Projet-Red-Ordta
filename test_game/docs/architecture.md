# Architecture

```text
cmd/red3d             Point d'entrée du jeu
internal/game         Règles de jeu indépendantes du rendu
internal/game3d       Scène, modèles procéduraux, caméra, HUD et contrôles G3N
tools/dev             Préparation, compilation et lancement sous Windows
assets                Emplacement des futurs fichiers graphiques et audio
```

La séparation entre `internal/game` et `internal/game3d` permet de conserver une logique testable et de brancher plus tard une interface CLI sur les mêmes règles.
