# Ordta — Projet RED

Ordta est un RPG médiéval-fantastique développé en Go pour le Projet RED. Il propose un jeu dans le terminal (CLI), demandé par le sujet, et une extension 3D avec G3N : exploration, marchands, fabrication, équipements visibles et combats au tour par tour.

Les deux interfaces utilisent le package `library`. L’adaptateur `library_3d.go` expose les actions à la 3D sans attendre de saisie terminal ni bloquer la boucle graphique. Les règles et catalogues du backend font référence.

## Lore et histoire

Le [lore d’Ordta est disponible ici](docs/Lore_Ordta.txt). Sa lecture est importante pour comprendre l’histoire, l’univers et les enjeux de l’aventure. Ce document constitue la référence narrative du projet ; il ne signifie pas que tous les événements décrits sont déjà jouables.

## Installation et lancement

### Windows : le minimum pour jouer en 3D

Installer **Go 64 bits** et **GCC/MinGW-w64 64 bits**, puis récupérer le projet. Cette installation se fait **une seule fois**. Ensuite, depuis le dossier `src`, la seule commande à taper sera :

```powershell
go run .
```

Choisir ensuite **1** dans le terminal pour ouvrir le jeu 3D.

### 1. Installer Go

Télécharger l’installateur **Windows x86-64** sur la [page officielle de téléchargement de Go](https://go.dev/dl/), puis l’installer avec les options par défaut. Le projet demande Go 1.21 minimum ; une version stable récente convient.

Si Go est déjà installé, passer à l’étape suivante. On vérifiera son fonctionnement à l’étape 4.

### 2. Installer GCC/MinGW-w64 64 bits

Utiliser **MSYS2**, qui permet d’installer GCC pour Windows :

1. Ouvrir la [page officielle de téléchargement de MSYS2](https://www.msys2.org/).
2. Télécharger l’installateur **x86_64**, l’exécuter et conserver le dossier proposé : `C:\msys64`.
3. Ouvrir **MSYS2 UCRT64** depuis le menu Démarrer de Windows.
4. Dans cette fenêtre MSYS2, copier la commande suivante puis appuyer sur Entrée :

```sh
pacman -S --needed mingw-w64-ucrt-x86_64-gcc
```

Accepter l’installation lorsqu’une confirmation est demandée et attendre sa fin.

**Installer MSYS2 seul ne suffit pas.** La commande ci-dessus installe réellement GCC et ses dépendances. Elle se lance dans **MSYS2 UCRT64**, pas dans PowerShell. Il n’est pas nécessaire d’installer VS Code, Visual Studio ou CMake pour ce projet.

Cette méthode est documentée par [MinGW-w64 : installation avec MSYS2](https://www.mingw-w64.org/getting-started/msys2/).

### 3. Ajouter GCC aux variables d’environnement Windows

Cette étape permet à PowerShell et à Go de trouver GCC.

1. Dans la recherche du menu Démarrer, taper **variables d’environnement**.
2. Ouvrir **Modifier les variables d’environnement système**, puis cliquer sur **Variables d’environnement…**.
3. Dans la partie **Variables utilisateur**, sélectionner **Path**, puis **Modifier…**.
4. Cliquer sur **Nouveau** et ajouter cette ligne :

```text
C:\msys64\ucrt64\bin
```

5. Valider toutes les fenêtres avec **OK**.
6. Fermer puis rouvrir PowerShell. Si le terminal est ouvert dans un éditeur, fermer puis rouvrir aussi cet éditeur.

**Ne pas supprimer les lignes déjà présentes dans Path.** Ajouter le dossier `bin`, pas le fichier `gcc.exe`. Si MSYS2 a été installé ailleurs, adapter le début du chemin. Ne pas utiliser `C:\msys64\usr\bin` à la place de `ucrt64\bin` : ce n’est pas le compilateur Windows choisi ici.

### 4. Vérifier l’installation dans PowerShell

Dans une **nouvelle fenêtre PowerShell**, exécuter :

```powershell
go version
where.exe gcc
gcc --version
```

Résultat attendu :

- `go version` affiche une version de Go et `windows/amd64`.
- `where.exe gcc` affiche `C:\msys64\ucrt64\bin\gcc.exe` si vous avez conservé le dossier par défaut.
- `gcc --version` affiche la version de GCC, sans erreur.

Si GCC fonctionne dans MSYS2 mais pas dans PowerShell, **le Path Windows n’est pas correctement configuré ou le terminal n’a pas été rouvert**. Reprendre l’étape 3. Si plusieurs GCC apparaissent, vérifier que celui de `ucrt64\bin` est utilisé en premier.

### 5. Récupérer le projet et jouer

Sur le [dépôt du projet](https://github.com/gauhthierh/Projet-Red-Ordta), cliquer sur **Code → Download ZIP**, puis **extraire toute l’archive**. Git n’est pas obligatoire si vous utilisez cette méthode.

Ouvrir le dossier extrait, puis son sous-dossier **src**. Dans l’Explorateur Windows, cliquer sur la barre d’adresse, taper `powershell` et appuyer sur Entrée : le terminal s’ouvre directement dans ce dossier.

Exécuter :

```powershell
go run .
```

Le menu propose :

```text
1 — Jouer en 3D
2 — Jouer dans le terminal
0 — Quitter
```

Taper **1**, puis Entrée. Le premier lancement peut prendre plus de temps : Go télécharge les dépendances et compile le jeu. **Une connexion Internet est nécessaire pour ces premiers téléchargements.**

Le lanceur active automatiquement CGO pour la 3D, impose la compilation Windows **64 bits**, vérifie que GCC est en 64 bits, récupère G3N et ajoute ses bibliothèques audio au processus du jeu. Cela évite notamment l’erreur G3N `uintUndef ... overflows` causée par une cible Go 32 bits. **Pas de commande CGO à taper, pas de DLL audio à télécharger ou copier à la main.** Seul le Path de GCC doit avoir été configuré à l’étape 3. Ces réglages concernent le processus 3D, sans changer la configuration globale de Go.

Pour les lancements suivants : ouvrir un terminal dans `src`, taper **`go run .`**, puis choisir **1**. Ne pas déplacer `src` ou supprimer `assets` : le jeu a besoin de l’arborescence complète.

Pour le CLI seulement, choisir **2** : GCC et les dépendances graphiques ne sont pas nécessaires.

### À savoir

- La 3D nécessite une carte graphique et un pilote compatibles OpenGL. Si la création de la fenêtre échoue, vérifier le pilote graphique.
- L’interface est principalement conçue pour 1920 × 1080.
- Le code est en Go, mais G3N utilise des composants natifs : Go seul ne suffit pas pour compiler la 3D.
- Cette procédure concerne Windows 64 bits sur processeur Intel/AMD. Pour macOS, suivre la section suivante. La configuration Linux n’est pas automatisée ici.

### Installation sur Mac (Intel ou Apple Silicon)

**Prise en charge préparée dans le code, mais pas encore vérifiée sur un Mac réel.** Les tests effectués sous Windows ne garantissent pas le fonctionnement du rendu et de l’audio sur toutes les versions de macOS.

Sur Mac, **ne pas installer MinGW/MSYS2** : le compilateur utilisé est **Clang**, fourni par Apple.

1. Installer Go depuis [go.dev/dl](https://go.dev/dl/) : choisir **macOS ARM64** pour une puce Apple Silicon (M1, M2, etc.), ou **macOS x86-64** pour un Mac Intel. Le type de puce est indiqué dans le menu Apple → À propos de ce Mac.
2. Ouvrir Terminal et installer les outils de compilation Apple :

```sh
xcode-select --install
```

Attendre la fin de l’installation. Si les outils sont déjà installés, passer à la suite.

3. Installer [Homebrew depuis son site officiel](https://brew.sh/). Suivre les indications **Next steps** affichées par l’installateur pour rendre la commande `brew` accessible, puis rouvrir Terminal.
4. Installer les dépendances audio :

```sh
brew install openal-soft libvorbis
```

Homebrew installe également libogg, nécessaire à Vorbis. Ces dépendances sont celles indiquées par la [documentation G3N pour macOS](https://github.com/g3n/engine/wiki#macos).

5. Vérifier :

```sh
go version
xcrun --find clang
brew --prefix openal-soft
brew --prefix libvorbis
```

Go doit indiquer `darwin/arm64` sur Apple Silicon ou `darwin/amd64` sur Intel. **Ne pas mélanger Go Intel sous Rosetta avec les bibliothèques Homebrew ARM64.** Utiliser un Terminal natif et les installations correspondant à la puce.

6. Télécharger et extraire le projet complet, comme sous Windows. Dans Terminal, taper `cd ` (avec un espace), glisser le dossier **src** depuis Finder, puis appuyer sur Entrée.
7. Lancer :

```sh
go run .
```

Choisir **1** pour la 3D ou **2** pour le CLI. Le lanceur active CGO, choisit Clang et récupère les chemins des bibliothèques avec Homebrew : pas de DLL Windows à copier ni de variables audio à configurer à chaque lancement.

Le CLI seul nécessite uniquement Go. Pour la 3D, le Mac doit aussi fournir un contexte OpenGL compatible avec G3N. Si une erreur de compilation ou de création de fenêtre persiste, transmettre le message complet avec la version de macOS et le type de puce.

<details>
<summary>Développeurs uniquement : Git, lancement direct et environnement de test</summary>

Si Git est déjà installé, il peut remplacer le téléchargement ZIP :

```powershell
git clone https://github.com/gauhthierh/Projet-Red-Ordta.git
cd Projet-Red-Ordta/src
go run .
```

**Ne pas suivre les commandes suivantes pour jouer normalement.** Elles servent au lancement direct de la 3D et aux tests graphiques, sans passer par le lanceur. Depuis `src` :

```powershell
$env:CGO_ENABLED = "1"
$env:GOARCH = "amd64"
go mod download github.com/g3n/engine
$dossierG3N = go list -m -f '{{.Dir}}' github.com/g3n/engine
$env:PATH = "$dossierG3N/audio/windows/bin;$env:PATH"
go run ./3D
```

Un exécutable du lanceur peut être compilé avec `go build -o ordta.exe .`, mais le choix 3D appelle encore Go : ce n’est pas une distribution autonome.

</details>

## Démarrer une partie

### CLI

Après le choix 2, saisir un nom et une classe. Les menus donnent accès au personnage, à l’inventaire, au marchand, au forgeron et à l’entraînement. Le choix 0 permet de revenir ou quitter. Le CLI ne conserve pas de sauvegarde entre deux exécutions.

### 3D

Au premier démarrage, cliquer sur **Démarrer**, saisir le nom puis choisir Humain, Elfe ou Nain. Le nom ne doit contenir que des lettres, accents acceptés. Il est normalisé avec une majuscule initiale et le reste en minuscules.

La validation crée et sauvegarde le personnage. Revenir à l’accueil sans valider ne crée pas de partie. Si une sauvegarde existe, **Continuer** reprend ce personnage : aucun nouveau choix de classe n’est demandé.

Une sauvegarde illisible est conservée et le démarrage est bloqué avec un message, plutôt que de l’écraser silencieusement.

### Commandes 3D

| Commande | Action |
|---|---|
| Souris | Regarder autour de soi en exploration |
| ZQSD sur AZERTY / WASD sur QWERTY | Déplacement relatif à la caméra |
| Maj gauche | Sprint |
| TAB | Ouvrir ou fermer l’inventaire |
| M | Ouvrir ou fermer la carte 2D avec votre position |
| E | Interagir avec un PNJ proche |
| Échap | Pause / reprise après le démarrage |
| F5 | Sauvegarder hors combat |
| Clic gauche | Boutons, objets, choix d’attaque et cible |

La souris est libérée dans les menus, le commerce, l’inventaire et l’arène. La pause suspend les déplacements, les animations et les effets. Le menu permet aussi de quitter et d’afficher « Qui sont-ils ? ».

### Carte du monde

Appuyer sur **M** après le démarrage pour afficher la carte pixel art. Le carré rouge bordé de blanc indique votre position, avec vos coordonnées X/Y en bas. Le nord est en haut ; le marqueur utilise le repère et la taille du monde définis dans le JSON de la carte.

Fermer avec **M**, **Échap** ou le bouton **Fermer**. Pendant la consultation, les déplacements, combats et effets sont suspendus et la souris est libre. L’inventaire ou le commerce précédemment ouvert est retrouvé à la fermeture. L’illustration reste un plan artistique : les détails dessinés ne représentent pas les collisions au pixel près.

## Personnage, objets et équipements

| Classe | PV maximum initiaux | Mana maximum initial | Bonus d’attaque |
|---|---:|---:|---:|
| Humain | 100 | 100 | 5 |
| Elfe | 80 | 120 | 0 |
| Nain | 120 | 80 | 10 |

Le personnage commence au niveau 1 avec la moitié de ses PV, tout son mana, 100 or, trois potions de vie, l’attaque basique et le sort Coup de poing.

- Potion de vie : jusqu’à +50 PV, sans dépasser le maximum.
- Potion de mana : jusqu’à +40 mana.
- Potion de poison : inflige au buveur 10 dégâts par seconde pendant trois secondes, ou jusqu’à sa mort. En 3D, le chronomètre ne bloque pas l’affichage.
- Mort : résurrection à la moitié des PV maximum ; après une défaite en arène, sortir ou choisir un autre combat.
- Livres et manuels : apprennent un sort ou une attaque une seule fois. Un doublon reste dans l’inventaire.
- L’expérience excédentaire est conservée après un niveau. Le seuil suivant est multiplié par 13/10, avec calcul entier, et les gains de statistiques dépendent de la classe.

### Inventaire

Un exemplaire occupe une place : quatre potions occupent quatre places. La capacité commence à 10, puis augmente de 10 par extension, au maximum trois fois, soit 40.

L’inventaire 3D affiche **dix cases par page**, avec navigation pour les extensions. Les statistiques, l’aperçu du personnage et les trois emplacements d’équipement restent visibles. Sélectionner un objet puis cliquer sur « Utiliser / équiper ». Après une action, sélectionner à nouveau l’objet souhaité.

Les emplacements du backend sont tête, torse et pieds. La protection des jambes du modèle est un complément visuel, pas un quatrième emplacement ou un bonus de statistiques indépendant.

### Marchand

La boutique 3D reprend tous les articles de `library.Boutique`, avec leur prix et leur niveau minimum, sur plusieurs pages. Elle contient potions, matériaux, extensions, livres de sorts et manuels d’attaques physiques. La première potion de vie est offerte ; les suivantes coûtent 3 or.

Les matériaux peuvent donc être **achetés ou obtenus en combat 3D**. Le niveau, l’or et la place disponible sont vérifiés avant l’achat. Le niveau exigé est affiché sur chaque article. Comme dans le CLI, un livre ou manuel déjà possédé/appris et une extension devenue inutile sont refusés, sans dépense d’or.

### Forgeron

Chaque fabrication coûte 5 or en plus des matériaux :

| Équipement | Matériaux | Bonus |
|---|---|---:|
| Chapeau de l’aventurier | 1 plume de corbeau + 1 cuir de sanglier | +10 PV maximum |
| Tunique de l’aventurier | 2 fourrures de loup + 1 peau de troll | +25 PV maximum |
| Bottes de l’aventurier | 1 fourrure de loup + 1 cuir de sanglier | +15 PV maximum |

Équiper retire l’objet du sac ; remplacer une pièce y remet l’ancienne. Le bonus augmente le maximum de vie, pas les PV actuels. Hors combat, cliquer sur une pièce équipée permet de la retirer si le sac a une place libre.

## Combats

L’entrée dans l’arène ouvre le choix du mode. La caméra devient latérale et surélevée ; sélectionner les actions avec la souris.

- **Entraînement** : gobelin au niveau du joueur, expérience et or selon le backend. Au niveau 1 : 40 PV, attaque 5, 50 XP et 5 à 10 or. Les dégâts, le mana dépensé et les objets consommés sont conservés, comme dans le CLI. Aucun matériau.
- **Arène** : quatre vagues — gobelin ; loup et gobelin ; gobelin cuirassé et chaman ; troll.
- **Duel** : un corbeau, sanglier, loup ou troll, pour gagner de l’expérience, de l’or et éventuellement des matériaux.

L’initiative est tirée entre 1 et 10. Le camp à la valeur la plus élevée commence, égalité en faveur du joueur ; en groupe, la plus haute initiative ennemie représente le camp adverse. Le compteur avance après les deux camps, quel que soit le premier attaquant.

Le gobelin et le troll doublent leurs dégâts tous les trois tours. Le chaman peut soigner ses alliés. La défense 3D réduit de moitié, arrondi au supérieur, la prochaine attaque reçue.

Attaquer, lancer un sort ou utiliser un objet valide consomme le tour. Une action refusée ne dépense ni tour ni mana. Les livres et manuels peuvent aussi être utilisés depuis le sac pendant le tour du joueur. Équiper ou retirer une armure reste réservé à l’exploration. Le poison finit ses dégâts avant l’action ennemie.

Le journal présente les actions, dégâts, PV, dépenses de mana et récompenses. Les animations sont visuelles : les calculs restent dans l’adaptateur.

### Sorts disponibles

| Sort | Dégâts | Mana |
|---|---:|---:|
| Coup de poing | 8 | 5 |
| Grosse boule de feu | 26 | 30 |
| Lame du destin | 12 | 10 |
| Éclats du gardien | 16 | 15 |
| Flèche de lumière | 20 | 20 |
| Jugement des géants | 60 | 60 |

Coup de poing est connu au départ ; les cinq autres sorts s’apprennent grâce aux livres vendus selon le niveau requis. Les anciens sorts supprimés du backend ne sont plus exécutables.

Les attaques physiques disponibles sont l’attaque basique puis Pichenette, Claquounette, Coups de pied, Morsure et Uppercut, apprises avec les manuels. Les dégâts sont ceux de `InfosAttaquePhysique` plus le bonus d’attaque du personnage.

### Niveau des ennemis et récompenses 3D

Le niveau des ennemis est fixé au début de chaque combat ou vague d’après le niveau actuel du joueur. Une montée de niveau ne change pas les ennemis déjà présents. Le menu des duels affiche les PV et l’or adaptés ; les fiches de combat affichent le niveau.

Le gobelin utilise directement `InitGoblin(niveau)`. Les autres créatures conservent leurs bases propres et reprennent les gains du gobelin pour chaque niveau supplémentaire : **+10 PV, +2 attaque, +20 XP, +3 or minimum et +5 or maximum**. Ces règles complémentaires restent dans `library_3d.go`.

Fourchettes au **niveau 1**, tirage uniforme avec les deux bornes incluses :

| Créature | Or | Matériau possible | Probabilité |
|---|---:|---|---:|
| Corbeau | 2–4 | Plume de corbeau | 75 % |
| Sanglier | 3–6 | Cuir de sanglier | 70 % |
| Gobelin | 5–10 | Aucun | — |
| Loup | 5–10 | Fourrure de loup | 60 % |
| Chaman | 7–14 | Aucun | — |
| Gobelin cuirassé | 8–16 | Aucun | — |
| Troll | 20–40 | Peau de troll | 45 % |

L’or est attribué une seule fois par créature vaincue, **entraînement compris**, et n’occupe aucune case. Un matériau est perdu si le sac est plein, avec un message dans le journal. Quitter l’arène est possible à tout moment ; seules les créatures déjà vaincues rapportent quelque chose.

## Sauvegarde 3D

Fichier : `src/sauvegardes/partie_3d.json`, exclu de Git.

Sont conservés : identité, classe, statistiques, expérience, or, inventaire, équipements, attaques, sorts, position et orientation de la caméra.

- Première sauvegarde après validation de la création.
- Toutes les 15 secondes hors arène, avec F5, à la fermeture du commerce ou de l’inventaire et à la sortie de l’arène.
- À la fermeture normale du jeu après démarrage, hors combat.
- Pas de sauvegarde au milieu d’un effet de poison en exploration ; à la fermeture normale, ses dégâts restants sont résolus avant l’écriture.
- Écriture dans un fichier temporaire puis remplacement du fichier de partie.
- Migration des anciens noms courants d’objets, du sort Boule de feu et de l’attaque basique à la lecture.

**Le combat en cours n’est pas sauvegardé. Quitter l’arène avant de fermer le jeu pour conserver sa progression.** Une interruption brutale reprend la dernière sauvegarde réussie.

Pour recommencer, fermer le jeu puis déplacer ou renommer le fichier de sauvegarde en conservant une copie. Il n’existe pas encore de gestion de plusieurs parties dans le menu.

## Organisation

```text
Projet-Red/
├── README.md
├── Projet RED - Sujet.pdf
├── docs/                     Sujet et documents du projet
├── assets/                   Carte, modèles, textures, icônes et audio
├── src/
│   ├── main.go               Choix terminal CLI / 3D
│   ├── lancement/            Préparation macOS (macos.go) et ses tests
│   ├── go.mod, go.sum        Dépendances
│   ├── library/              Règles du jeu et interface CLI
│   │   └── library_3d.go     Adaptateur non bloquant pour la 3D
│   ├── 3D/
│   │   ├── main.go           Entrée graphique
│   │   ├── monde/            Scène, interfaces, caméra, audio, sauvegarde
│   │   ├── personnages/      Personnage et équipements visibles
│   │   └── monstres/         Modèles et articulations des ennemis
│   ├── tests/                Vérifications de compatibilité de l’adaptateur
│   ├── tools/                Générateurs de carte et de musique
│   └── sauvegardes/          Créé à l’exécution, non versionné
└── test_game/                Ancien exemple, non utilisé et laissé intact
```

Les fonctions terminal interactives du backend ne doivent pas être appelées depuis la boucle graphique. L’adaptateur utilise les structures, catalogues et fonctions de calcul disponibles, puis renvoie des résultats affichables. Les changements de noms ou de règles doivent également être répercutés dans cet adaptateur.

La préparation macOS est regroupée dans `src/lancement/`, pas dans des fichiers à la racine de `src`. Les deux points d’entrée sont conservés : `src/main.go` pour le choix CLI/3D et `src/3D/main.go` pour le jeu graphique. Le lancement habituel reste `go run .` depuis `src/`. Les tests de préparation macOS se lancent avec `go test ./lancement`.

## Vérifications

### À quoi servent les fichiers `_test.go` ?

Les fichiers dont le nom se termine par `_test.go` contiennent des **tests automatiques**. Ils vérifient qu’une fonction donne le résultat attendu et permettent de repérer une régression : une modification qui casse un comportement auparavant correct.

Un test prépare une situation, appelle le code du jeu puis compare le résultat obtenu au résultat attendu. Par exemple :

- Donner une potion à un personnage blessé, l’utiliser et vérifier que ses PV augmentent sans dépasser leur maximum.
- Tenter un achat avec un niveau insuffisant et vérifier que l’objet n’est pas acheté et que l’or est conservé.
- Lancer un sort et vérifier les dégâts ainsi que le mana consommé.
- Enregistrer une partie dans un dossier temporaire, la relire et vérifier que les informations sont conservées.

Ces fichiers ne constituent pas une autre version du jeu. **Ils ne sont pas exécutés avec `go run .` et ne sont pas inclus dans l’exécutable produit par `go build`.** Go les prend en compte lors de la commande `go test`, qui construit et exécute un programme de vérification séparé.

Il est utile de les conserver et de les relancer après une modification. Lorsqu’une règle change volontairement, il faut également mettre à jour les résultats attendus par les tests concernés.

### Où se trouvent-ils ?

| Emplacement | Rôle |
|---|---|
| `src/tests/` | Vérifier la compatibilité de l’adaptateur 3D avec les règles du backend, sans charger le moteur graphique. |
| Fichiers `_test.go` dans `src/3D/` | Vérifier notamment les modèles, animations, sauvegardes et données du monde. |
| Fichiers `_test.go` dans `src/library/` | Anciens tests du backend et de l’adaptateur ; certains nécessitent une mise à jour, expliquée ci-dessous. |
| `test_game/` | Ancien exemple de jeu indépendant : ce dossier n’est pas une suite de tests automatiques et n’est pas utilisé par le jeu actuel. |

### Comment les lancer ?

Depuis `src/` :

```powershell
go test ./tests
```

Cette suite vérifie notamment création, classes, niveaux d’achat, apprentissage, doublons, initiative, tours, mana, poison et capacité, sans importer G3N. `progression_3d_test.go` couvre aussi les sept créatures à plusieurs niveaux, l’or dans les trois modes, l’absence de double récompense, la progression entre vagues et l’interdiction d’équiper en combat.

Après préparation de l’environnement natif décrite plus haut :

```powershell
go test ./3D/...
go build -o "$env:TEMP/ordta-verification.exe" ./3D
```

Ces tests couvrent notamment modèles, animations, sauvegardes et données du monde, ainsi que la visibilité des commandes selon la phase et le rejet des fichiers audio tronqués. Ils ne remplacent pas un essai visuel et sonore.

Pour afficher le nom et le résultat de chaque test, ajouter `-v` :

```powershell
go test -v ./tests
```

Pour forcer une nouvelle exécution sans réutiliser un résultat de test en cache :

```powershell
go test -count=1 ./tests
```

### Comment lire le résultat ?

- `ok` : les tests du package ont réussi.
- `PASS` : réussite d’un test ou de la suite, notamment dans l’affichage détaillé.
- `FAIL` : un test a échoué ou le programme de test n’a pas pu être compilé. Lire les messages précédents pour identifier le fichier, la ligne et la cause.
- `[no test files]` : aucun fichier de test dans ce package ; cela ne signifie pas que son fonctionnement a été vérifié.
- `(cached)` : Go a réutilisé un résultat réussi encore valide.

Un échec ne signifie donc pas toujours que le jeu est inutilisable : le test peut lui-même employer un ancien nom de champ, comme dans le cas ci-dessous. Inversement, des tests réussis ne garantissent pas l’absence de tous les bugs, notamment visuels ou sonores.

**Limite actuelle des anciens tests :** les tests de `src/library/` compilent, mais plusieurs attendent encore les anciennes règles : joueur toujours premier, entraînement sans récompense avec restauration du personnage, marchand sans matériaux et or fixe. Ils ont été laissés intacts. `go test ./library` et donc `go test ./...` échouent tant que ces attentes ne sont pas adaptées à l’initiative aléatoire et aux règles actuelles du backend ; cela n’empêche pas la compilation du jeu ni les suites ciblées ci-dessus.

### Outils de génération : à ne pas confondre avec les tests

Les assets sont déjà livrés. Pour les développeurs seulement :

```powershell
go run ./tools/mapgen
go run ./tools/musicgen
```

Attention : ces commandes remplacent les fichiers générés de la carte ou de la musique. Sauvegarder ses retouches avant de les utiliser. Elles ne sont pas nécessaires pour jouer.

## Lire et comprendre le code 3D

Les fichiers Go de `src/3D/` et l'adaptateur `src/library/library_3d.go` sont commentés en français : rôle des fonctions, étapes importantes et raisons des calculs moins évidents. Les commentaires des tests expliquent aussi ce qu'ils vérifient.

Pour commencer, suivre cet ordre :

1. [`src/3D/main.go`](src/3D/main.go) appelle le démarrage de la 3D.
2. [`monde.go`](src/3D/monde/monde.go) construit la scène, branche les boutons et met à jour le jeu à chaque image.
3. Les autres fichiers de [`monde`](src/3D/monde/) séparent les caméras, collisions, menus, sons, sauvegardes et affichages.
4. [`personnages`](src/3D/personnages/) et [`monstres`](src/3D/monstres/) construisent ou chargent les modèles et leurs membres articulés.
5. [`library_3d.go`](src/library/library_3d.go) fournit les actions sans saisie terminal : il calcule les résultats que la scène affiche et anime.

Quelques repères pour lire les calculs : X et Y forment le sol, Z est la hauteur. Un membre tourne autour de son pivot local, tandis que le nœud racine déplace le personnage entier. Les angles G3N sont en radians ; les rotations du JSON de collisions sont en degrés. Les animations reçoivent un temps en secondes, et les effets backend utilisent `time.Duration`.

L’adaptation aux nouvelles règles reste dans `library_3d.go` : niveau des créatures, or et utilisation des constantes backend. La visibilité des commandes est regroupée dans `InterfaceCombat.ActualiserVisibilite`, et les branches de sorts de soutien absents du catalogue ont été retirées de l’adaptateur. Les autres fichiers backend et les assets ne sont pas modifiés par cette adaptation.

Le chargement audio vérifie séparément le format des WAV livrés. Si OpenAL échoue lors du chargement, les buffers déjà créés sont libérés et la partie continue sans audio, sans répéter l’erreur pour chaque piste.

## Correspondance avec le sujet

Le [sujet](docs/Projet%20RED%20-%20Sujet.pdf) demande un jeu CLI, un dossier `src`, un dossier `docs` avec le document de gestion de projet et un README présentant le jeu, son installation et son lancement. La 3D complète le CLI, elle ne le remplace pas.

| Exigences | Implémentation principale |
|---|---|
| Personnage, création, classes, statistiques | `library/character.go`, `creation.go` |
| Inventaire, potions, limite et extensions | `inventory.go`, `items.go` |
| Argent et marchand | `merchant.go`, `constantes.go` |
| Fabrication et équipement | `forgeron.go`, `equipement.go` |
| Gobelin, tours, initiative | `monster.go`, `combat.go` |
| Expérience, attaques, sorts et mana | `experience.go`, `attaquesphysiques.go`, `sort.go` |
| « Qui sont-ils ? » : ABBA et Steven Spielberg | Menu CLI et menu 3D |
| Interface graphique des règles | `library_3d.go` et `3D/monde/` |

La page 6 autorise des valeurs différentes tant que le principe des tâches est respecté. Ce README décrit les valeurs du code, pas une garantie de note ou de conformité totale.

À compléter par l’équipe : le document de gestion de projet n’a pas été trouvé dans `docs/` (qui contient le sujet et le lore). Le lore ne remplace pas ce livrable. Vérifier également le nom de dépôt demandé, `projet-red_NOM-DU-PROJET`, et le dépôt du lien sur Moodle avant l’échéance fixée par les encadrants.

## Limites et dépannage

- La carte est finie ; pas de génération infinie, saut, nage ou physique verticale.
- Le plan pixel art est accessible avec M et un marqueur ; il n’affiche pas les ennemis ni les collisions.
- De nombreux panneaux sont dimensionnés pour 1920 × 1080. Le plein écran utilise les dimensions de l’écran, sans imposer cette résolution.
- Les nouveaux livres et manuels réutilisent une icône de livre ; ils n’ont pas chacun une illustration dédiée.
- Les sorts Soin du cœur et Bouclier ont encore des constantes dans le backend mais ne sont pas proposés par `InfosSort` ni vendus sous forme de livres ; ils ne sont pas des sorts jouables de cette version.

| Problème | À vérifier |
|---|---|
| GCC introuvable | Vérifier `where.exe gcc` dans le terminal de lancement. |
| `build constraints exclude all Go files` | En lancement direct 3D, activer CGO ; le lanceur le fait automatiquement. |
| `exit status 0xc0000135` | DLL native absente : PATH audio G3N et dépendances GCC. |
| Asset introuvable | Lancer depuis `src/`, avec `assets/` à la racine. |
| Échec OpenGL | Pilote graphique et session graphique locale disponibles. |
| Mac : Clang introuvable | Installer les outils Apple avec `xcode-select --install`. |
| Mac : `al.h`, `codec.h` ou bibliothèque audio introuvable | Installer `brew install openal-soft libvorbis`, puis utiliser le lanceur `go run .`. |
| Mac : architecture incompatible | Go, Terminal et Homebrew doivent utiliser la même architecture : ARM64 sur Apple Silicon, x86-64 sur Intel. |
| G3N : `uintUndef ... overflows` | Compilation en 32 bits. Le lanceur corrigé impose `GOARCH=amd64` uniquement pour la 3D Windows et vérifie GCC 64 bits. Relancer avec `go run .`. En lancement direct, définir aussi `$env:GOARCH = "amd64"`. |
| Sauvegarde illisible | Ne pas la supprimer : conserver une copie avant de la renommer ou de la réparer. |
| Audio indisponible | Vérifier les WAV de `assets/audio/`. |
| Échec des anciens tests library | Adapter leurs anciennes références côté backend ; voir « Vérifications ». |

## Équipe et ressources

Contributeurs du projet :

- **Gautier** : frontend 3D.
- **Guillaume** : backend.
- **Quentin** : backend.
- **Lalie** : lore et backend.

G3N v0.2.0 utilise notamment GLFW pour la fenêtre, OpenGL pour le rendu et OpenAL pour l’audio. Les dépendances sont déclarées dans `src/go.mod` et `src/go.sum`.

**G3N** est le moteur 3D écrit en Go : il fournit les outils pour construire la scène, placer les modèles, gérer les caméras, matériaux, lumières et interfaces graphiques. Les règles de combat et d’inventaire restent dans notre propre code.

**GLFW** est une bibliothèque utilisée par le moteur pour créer la fenêtre et recevoir les événements du clavier et de la souris. Elle permet notamment de gérer le plein écran et la capture du curseur ; ce n’est pas elle qui dessine les modèles 3D.

Aucune licence globale n’est déclarée dans le dépôt. Les licences des dépendances ne constituent pas une autorisation de redistribution de tous les assets.
