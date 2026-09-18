# Projet RED - Arène 3D

Petit prototype de combat 3D réalisé uniquement en Go avec le moteur G3N.

Le joueur contrôle un aventurier dans une arène, affronte un gobelin, peut utiliser trois potions et recommencer une partie après une victoire ou une défaite. Tous les modèles sont générés par le code : aucun asset externe n'est nécessaire.

## Lancer le jeu sous Windows

Prérequis : une installation de Go.

```powershell
go run ./tools/dev -run
```

La première exécution télécharge dans `.tools/` une chaîne de compilation C portable vérifiée par SHA-256, car G3N utilise OpenGL via `cgo`. Elle compile ensuite le jeu dans `bin/`, copie les DLL audio nécessaires puis ouvre la fenêtre.

Pour seulement compiler :

```powershell
go run ./tools/dev
```

## Contrôles

| Action | Touches |
| --- | --- |
| Déplacement | `ZQSD`, `WASD` ou flèches |
| Attaquer | `Espace` |
| Utiliser une potion | `E` |
| Recommencer | `R` |
| Quitter | `Échap` |

Il faut s'approcher du gobelin avant d'attaquer. Le gobelin poursuit automatiquement le joueur et frappe lorsqu'il est à portée.

## Organisation

```text
cmd/red3d/          point d'entrée
internal/game/      logique de jeu testable sans moteur graphique
internal/game3d/    rendu G3N, scène, HUD et contrôles
tools/dev/          outil Go de préparation et de compilation
assets/             futurs modèles, textures et sons
docs/               documentation du projet
```

## Tests

La logique indépendante du rendu possède des tests unitaires :

```powershell
go test ./internal/game
```
