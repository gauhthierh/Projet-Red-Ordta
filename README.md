# Ordta — Projet RED

Ordta est un RPG médiéval-fantastique développé en **Go** dans le cadre du Projet RED d’Ynov. Le joueur gère un personnage, son inventaire et ses équipements, achète des objets, fabrique des armures et combat au tour par tour.

Le dépôt propose deux versions partageant les structures du package `library` :

- **CLI** : jeu dans le terminal, avec création du personnage et menus numérotés, correspondant au format demandé dans le sujet.
- **3D** : extension graphique avec **G3N v0.2.0**, exploration à la première personne, village, PNJ, inventaire graphique et arène vue de côté.

Le code du jeu et des outils est en Go. La 3D utilise néanmoins des composants natifs via CGO : compilateur C, OpenGL et bibliothèques audio. Elle n’est donc pas indépendante des bibliothèques système.

## Sommaire

- [Installation et lancement](#installation-et-lancement)
- [Comment jouer](#comment-jouer)
- [Personnage et inventaire](#personnage-et-inventaire)
- [Marchand et forgeron](#marchand-et-forgeron)
- [Combats 3D](#combats-3d)
- [Sauvegarde](#sauvegarde)
- [Organisation du dépôt](#organisation-du-dépôt)
- [Tests et outils](#tests-et-outils)
- [Correspondance avec le sujet](#correspondance-avec-le-sujet)
- [Limites actuelles](#limites-actuelles)
- [Dépannage](#dépannage)

## Installation et lancement

### Récupérer le projet

Prérequis communs : **Go 1.21 minimum**, conformément à `src/go.mod`, et Git pour cloner le dépôt. Une connexion est nécessaire au téléchargement initial des dépendances.

```powershell
git clone https://github.com/gauhthierh/Projet-Red-Ordta.git
cd Projet-Red-Ordta/src
go version
```

Les commandes suivantes sont à exécuter depuis **`src/`**, sauf indication contraire. La 3D charge notamment `../assets/` : conserver ce dossier à côté de `src/`.

### Lancer le CLI

```powershell
go run .
```

Le point d’entrée est `src/main.go`. Le CLI n’importe pas G3N et ne nécessite ni GCC, ni fenêtre OpenGL, ni bibliothèques audio natives.

Pour compiler et lancer un exécutable Windows :

```powershell
go build -o ordta-cli.exe .
.\ordta-cli.exe
```

### Préparer la 3D sous Windows

Prévoir Windows 64 bits, Go pour `windows/amd64`, un compilateur **GCC compatible MinGW-w64 en 64 bits** et un pilote graphique prenant en charge OpenGL. L’interface actuelle est conçue principalement pour **1920 × 1080**.

Installer GCC/MinGW-w64, par exemple via MSYS2, puis rendre son dossier `bin` accessible au terminal utilisé pour lancer Go. Vérifier dans PowerShell :

```powershell
go env GOOS GOARCH
where.exe gcc
gcc --version
```

Si GCC est installé mais introuvable, ajouter son véritable dossier `bin` au `PATH`. Exemple uniquement pour une installation MSYS2 UCRT64 située à cet emplacement :

```powershell
$env:Path = "C:\msys64\ucrt64\bin;$env:Path"
```

Activer CGO dans ce terminal, télécharger les dépendances déclarées par le projet et ajouter les DLL audio fournies par G3N :

```powershell
$env:CGO_ENABLED = "1"
$env:CC = "gcc"
go mod download
$dossierMoteur = go list -m -f '{{.Dir}}' github.com/g3n/engine
$dossierAudio = Join-Path $dossierMoteur "audio\windows\bin"
$env:Path = "$dossierAudio;$env:Path"
go env CGO_ENABLED
```

La dernière commande doit afficher `1`. Le dossier audio de G3N v0.2.0 contient `OpenAL32.dll`, `libogg.dll`, `libvorbis.dll` et `libvorbisfile.dll`.

Ces réglages `$env:...` concernent uniquement le terminal courant : les refaire après sa fermeture, ou configurer les variables Windows correspondantes. Ne pas mélanger les outils 32 et 64 bits.

### Lancer la 3D

Depuis `src/`, dans le terminal préparé ci-dessus :

```powershell
go run ./3D
```

Le point d’entrée est `src/3D/main.go`, qui appelle `monde.Lancer()`. Aucun ancien lanceur de `archive_gpt/` n’est nécessaire.

Pour compiler puis lancer :

```powershell
go build -o ordta-3d.exe ./3D
.\ordta-3d.exe
```

Garder `src/` comme répertoire de travail, même avec l’exécutable : tous les assets ne sont pas intégrés au binaire.

G3N fournit aussi des instructions pour Linux et macOS dans son README. Les exemples de ce dépôt ciblent Windows ; le fonctionnement graphique sur les autres systèmes reste à valider.

## Comment jouer

### Version terminal

Saisir le numéro du choix puis **Entrée**. Le menu principal propose les informations du personnage, l’inventaire, le marchand, le forgeron, l’entraînement et « Qui sont-ils ? ». `0` permet de revenir ou de quitter selon le menu.

La création demande un nom composé uniquement de lettres, accents acceptés, sans espaces, chiffres ou tirets. Il est normalisé avec une majuscule initiale. Choisir ensuite Humain, Elfe ou Nain.

### Version 3D

Cliquer sur **Démarrer** dans l’accueil. La sauvegarde existante est chargée automatiquement ; sans sauvegarde, le personnage est un Humain nommé « Joueur ».

| Commande | Action |
|---|---|
| W / A / S / D | Avancer, aller à gauche, reculer, aller à droite, relativement à la caméra |
| Maj gauche | Courir en maintenant une direction |
| Souris | Orienter la vue en exploration |
| TAB | Ouvrir ou fermer l’inventaire |
| E | Parler au marchand ou au forgeron proche ; fermer le commerce ouvert |
| Échap | Mettre en pause ou reprendre une partie démarrée |
| F5 | Sauvegarder hors arène |
| Clic sur les boutons | Acheter, fabriquer, utiliser un objet, choisir une cible ou une action |

Le code utilise `KeyW`, `KeyA`, `KeyS`, `KeyD` ; aucun réglage des touches n’est proposé dans les menus. La souris est capturée en exploration et libérée dans les interfaces et l’arène.

Pour découvrir le jeu :

1. Ouvrir l’inventaire avec TAB et utiliser une potion : le personnage démarre à la moitié de ses PV.
2. Rejoindre le marché à l’est de la place centrale. Marchand et forgeron occupent deux cabanons distincts ; s’approcher jusqu’à l’indication **E**.
3. Rejoindre l’arène au sud-est. Entrer dans sa zone centrale ouvre le choix du combat.
4. Essayer l’entraînement, puis les duels pour obtenir matériaux et or.
5. Revenir au forgeron, fabriquer une pièce et l’équiper depuis l’inventaire.

## Personnage et inventaire

### Statistiques de départ

| Classe | PV maximum | Mana maximum | Bonus d’attaque physique |
|---|---:|---:|---:|
| Humain | 100 | 100 | 5 |
| Elfe | 80 | 120 | 0 |
| Nain | 120 | 80 | 10 |

Le personnage commence au niveau 1, avec 50 % de ses PV, son mana plein, 100 pièces d’or, 3 potions de vie, l’attaque basique et Coup de poing. L’attaque basique inflige **5 + le bonus d’attaque**.

Le premier seuil d’expérience est de 100. Il augmente ensuite de 30 %, avec arrondi entier ; l’excédent est conservé. Les gains de PV maximum, de mana maximum et d’attaque dépendent de la classe.

### Objets et capacité

- **Un exemplaire = une place** : quatre potions occupent quatre cases.
- La capacité initiale est de 10 objets. Trois améliorations de 10 places sont prévues par les règles, jusqu’à 40. Voir la limite d’affichage 3D plus bas.
- Une potion de vie restaure jusqu’à 50 PV ; une potion de mana jusqu’à 40 points, sans dépasser le maximum.
- Le livre apprend Boule de feu une seule fois. Un livre supplémentaire n’est pas consommé si le sort est déjà connu.
- La potion de poison **blesse son utilisateur** de 10 PV par seconde pendant trois secondes. Ce n’est pas une attaque à lancer sur l’ennemi. En 3D, elle s’utilise uniquement en combat.

## Marchand et forgeron

### Achats

| Article commun aux deux marchands | Prix en or |
|---|---:|
| Potion de vie | Première offerte, puis 3 |
| Potion de poison | 6 |
| Potion de mana | 5 |
| Livre de Boule de feu | 25 |
| Augmentation d’inventaire | 30 |

Le marchand CLI vend également la fourrure de loup (4 or), la peau de troll (7), le cuir de sanglier (3) et la plume de corbeau (1). **En 3D, ces quatre matériaux viennent exclusivement du butin des combats.**

### Fabrication

| Équipement | Matériaux | Coût | Bonus de PV maximum |
|---|---|---:|---:|
| Chapeau de l’aventurier | 1 plume de corbeau + 1 cuir de sanglier | 5 or | +10 |
| Tunique de l’aventurier | 2 fourrures de loup + 1 peau de troll | 5 or | +25 |
| Bottes de l’aventurier | 1 fourrure de loup + 1 cuir de sanglier | 5 or | +15 |

La fabrication consomme les matériaux. Équiper une pièce la retire du sac ; remplacer une pièce remet l’ancienne dans l’inventaire.

En 3D, sélectionner l’objet puis **Utiliser / Équiper**. Cliquer sur un emplacement équipé permet de retirer la pièce si une place est libre. Les changements d’équipement sont interdits pendant le combat.

Le personnage dispose d’une tenue sans armure et de pièces 3D visibles : chapeau, plastron, épaulières, brassards, bottes, genouillères et jambières. Les protections de jambes sont liées aux bottes, sans quatrième emplacement indépendant.

## Combats 3D

### Modes

| Mode | Déroulement | Progression |
|---|---|---|
| Entraînement | Un gobelin de 40 PV | Aucun gain ; personnage, PV, mana et inventaire restaurés lorsqu’on quitte ce combat |
| Arène | Quatre vagues | XP, or et éventuels matériaux par ennemi vaincu |
| Duel | Un corbeau, sanglier, loup ou troll au choix | Récolte répétable d’XP, d’or et de matériaux |

Les vagues sont : **gobelin → loup et gobelin → gobelin cuirassé et chaman → troll final**.

La caméra prend une vue latérale surélevée. Le joueur commence chaque vague en 3D. Choisir une cible sur sa fiche, puis une attaque physique, un sort, une potion ou **Défendre**. Une action valide consomme le tour ; une action refusée ne le consomme pas.

La défense réduit de moitié les dégâts de la prochaine attaque reçue, avec arrondi supérieur. Gobelins et troll doublent leur attaque tous les trois tours ; le chaman tente alors de soigner un allié blessé.

Les animations différencient attaques, sorts et créatures. Le journal consultable conserve actions, dégâts, soins, variations de mana, protections, récompenses et changements de tour. Les fiches affichent les PV et statistiques.

**Quitter l’arène** reste possible pendant le combat. Hors entraînement, les récompenses déjà obtenues sont conservées ; un ennemi vivant ne rapporte rien. Après une défaite, quitter le combat ou en choisir un autre ressuscite le personnage à la moitié de ses PV maximum.

### Récompenses hors entraînement

| Créature | Or par élimination | Matériau possible | Probabilité |
|---|---:|---|---:|
| Corbeau | 2 | Plume de corbeau | 75 % |
| Sanglier | 3 | Cuir de sanglier | 70 % |
| Gobelin | 4 | Aucun | — |
| Loup | 5 | Fourrure de loup | 60 % |
| Chaman | 7 | Aucun | — |
| Gobelin cuirassé | 8 | Aucun | — |
| Troll | 20 | Peau de troll | 45 % |

Le matériau est perdu si le sac est plein, ce que le journal signale. L’or n’occupe aucune case.

### Sorts pris en charge par l’adaptateur 3D

| Sort | Mana | Effet |
|---|---:|---|
| Coup de poing | 5 | 8 dégâts |
| Boule de feu | 30 | 22 dégâts |
| Lame du destin | 10 | 10 dégâts |
| Éclats du gardien | 15 | 12 dégâts |
| Flèche de lumière | 20 | 12 dégâts |
| Foudre Célèste | 20 | 30 dégâts immédiats actuellement |
| Soin du cœur | 25 | Restaure jusqu’à 20 % des PV maximum |
| Bouclier | 30 | Protège contre la prochaine attaque |
| Dévotion | 35 | 20 à 30 dégâts ; 25 % de risque de subir 8 dégâts en retour |
| Dernier espoir | 50 | 40 dégâts, à 30 % des PV maximum ou moins |
| Jugement des géants | 60 | 60 dégâts |

Cette table décrit les effets codés, **pas des sorts tous obtenables actuellement**. Coup de poing est connu au départ et Boule de feu s’apprend grâce au livre. Les autres sorts nécessitent encore un système d’apprentissage raccordé au jeu.

## Sauvegarde

La sauvegarde concerne **uniquement la 3D**. Depuis le répertoire de lancement `src/`, elle se trouve dans `src/sauvegardes/partie_3d.json`.

Elle contient le personnage, ses statistiques, inventaire, équipements et sorts, ainsi que sa position et l’orientation de la caméra. Elle est chargée automatiquement au démarrage.

- Sauvegarde automatique toutes les 15 secondes de jeu hors arène.
- Sauvegarde manuelle avec F5 hors arène.
- Sauvegarde à la fermeture de l’inventaire ou du commerce hors arène, à la sortie de l’arène et à la fermeture normale du jeu hors arène, après démarrage de la partie.
- **Aucun combat en cours n’est sauvegardé. Quitter l’arène avant de fermer le jeu pour conserver la progression du combat.**

Le dossier est exclu de Git. Pour recommencer, fermer le jeu et déplacer ou renommer `partie_3d.json` en conservant une copie. Le CLI n’a pas de sauvegarde persistante.

## Organisation du dépôt

```text
Projet-Red-Ordta/
├── README.md
├── Projet RED - Sujet.pdf
├── docs/
│   └── Projet RED - Sujet.pdf
├── assets/
│   ├── maps/red_world/       OBJ/MTL, plan 2D, textures et collisions JSON
│   ├── models/monstres/     Sept modèles articulés JSON et textures
│   ├── ui/                 Icônes, composants et maquettes
│   ├── pixel_art/          Planches et icônes pixel art
│   ├── audio/              Musiques et effets WAV
│   └── ...                 Autres modèles et illustrations
├── src/
│   ├── go.mod / go.sum     Module ordta et dépendances
│   ├── main.go             Entrée CLI
│   ├── library/            Règles et menus textuels
│   │   └── library_3d.go   Adaptateur sans saisie terminal pour la 3D
│   ├── 3D/
│   │   ├── main.go         Entrée graphique
│   │   ├── monde/         Scène, caméra, collisions, interfaces, audio, sauvegarde
│   │   ├── personnages/   Joueur, PNJ, armures et marche
│   │   └── monstres/      Chargement des modèles articulés
│   ├── tools/
│   │   ├── mapgen/        Générateur de fichiers statiques de carte
│   │   └── musicgen/      Générateur de musiques et effets en Go
│   └── sauvegardes/       Créé à l’exécution, non versionné
└── archive_gpt/            Ancienne implémentation, hors du jeu actuel
```

`library_3d.go` expose des résultats exploitables par les interfaces sans bloquer sur une saisie terminal. `src/3D/monde/monde.go` relie les entrées, les actions, l’affichage et les animations.

La carte est un monde **fini et semi-ouvert de 600 × 600 unités**, chargé depuis des fichiers statiques. Le déplacement se fait dans le plan XY ; Z représente la hauteur. Le JSON décrit zones, point de départ, obstacles et passages des ponts.

Village, remparts, guilde, forge extérieure, champs, forêt, falaises, marais, sanctuaire et bosquet de mana composent le décor. Les descriptions du JSON contiennent aussi des intentions de conception : elles ne déclenchent pas automatiquement une quête ou une interaction.

Les musiques et effets WAV sont synthétisés par l’outil Go. Les ressources graphiques comprennent des créations assistées par génération d’images ; certains prompts sont conservés dans `assets/pixel_art/PROMPTS.md` et la documentation des textures.

Les documents de `archive_gpt/` décrivent une ancienne version : leurs lanceurs et commandes ne s’appliquent pas au jeu actuel. Certains README locaux d’assets sont également historiques ; le présent README décrit les fonctionnalités raccordées aux points d’entrée actuels.

## Tests et outils

Depuis `src/`, pour vérifier les règles sans les dépendances graphiques :

```powershell
go test ./library
```

Après préparation de l’environnement 3D :

```powershell
go test ./...
```

La suite comprend des vérifications du combat, du butin, de l’or, de l’inventaire, des sauvegardes, des modèles, armures, animations, textures, chemins, collisions et fichiers audio. Les tests de carte peuvent être plus longs : ils reconstruisent des assets dans des répertoires temporaires.

Cette commande couvre le module `src/`, pas celui de l’archive. Ajouter `-v` pour voir les noms des tests. Les tests ne remplacent pas une vérification visuelle et sonore en jeu.

Compilation des tests du monde sans ouvrir le jeu sous Windows :

```powershell
go test -c -o "$env:TEMP/ordta-monde.test.exe" ./3D/monde
```

Les assets sont déjà livrés ; leur régénération n’est **pas nécessaire pour jouer**. Pour les développeurs, depuis `src/` :

```powershell
go run ./tools/mapgen
go run ./tools/musicgen
```

**Attention :** ces commandes remplacent respectivement les OBJ/MTL/JSON de la carte et les WAV de `assets/audio/`. Conserver ses retouches avant de les exécuter. Le générateur de carte ne recrée ni l’illustration 2D ni les textures PNG.

## Correspondance avec le sujet

Le [sujet fourni](docs/Projet%20RED%20-%20Sujet.pdf), page 3, demande :

- Un dossier `src` contenant le code source.
- Un dossier `docs` contenant le document de gestion de projet complété.
- Un `README.md` avec une courte présentation et les instructions d’installation et de lancement.
- Un dépôt nommé `projet-red_NOM-DU-PROJET`, dont le lien doit être déposé sur Moodle avant l’échéance donnée par les encadrants.

Ce README documente les deux versions sans présenter la 3D comme un remplacement du CLI demandé. La page 6 précise que les valeurs du sujet sont indicatives ; les tableaux ci-dessus décrivent celles du code actuel.

| Partie du sujet | Fichiers principaux dans src/library/ |
|---|---|
| Tâches 1 à 3 : structure, initialisation et informations | `character.go` |
| Tâches 4 à 7 : inventaire, potion, menu et marchand | `inventory.go`, `items.go`, `menu.go`, `merchant.go` |
| Tâches 8 à 10 : mort, poison et apprentissage | `character.go`, `items.go` |
| Tâches 11 à 14 : création, capacité, argent et achats | `creation.go`, `inventory.go`, `character.go`, `merchant.go` |
| Tâches 15 à 18 : fabrication, équipement et extensions | `forgeron.go`, `equipement.go`, `inventory.go` |
| Tâches 19 à 22 : monstres, comportement et combat | `monster.go`, `character.go`, `attaquesphysiques.go`, `combat.go` |
| Bonus initiative, expérience, sorts et mana | `combat.go`, `experience.go`, `sort.go`, `items.go` |
| « Qui sont-ils ? » | `menu.go` : ABBA et Steven Spielberg |
| Extension graphique, vagues, duels et butin | `library_3d.go` et `src/3D/` |

Cette correspondance indique où examiner les fonctionnalités, **pas une certification de conformité ou une note**. Actuellement, `docs/` contient une copie du sujet, mais aucun document de gestion de projet complété n’a été trouvé. Ce livrable reste à ajouter par l’équipe ; le README ne le remplace pas.

## Limites actuelles

- **Création 3D** : pas encore de saisie du nom ni de choix de classe dans l’accueil, contrairement au CLI. La 3D utilise « Joueur », Humain, ou le personnage sauvegardé.
- **Inventaire 3D** : les règles permettent 40 objets après extensions, mais l’interface n’affiche que les dix premiers exemplaires, sans pagination. Les suivants ne sont pas sélectionnables tant qu’ils ne reviennent pas dans ces dix cases.
- **Sorts** : plusieurs effets sont codés sans moyen d’acquisition raccordé. Le CLI ne prend pas en charge tous les sorts de soutien de la 3D.
- **Entraînement et initiative** : le CLI détermine le premier attaquant par l’initiative et donne de l’XP en cas de victoire. L’entraînement 3D restaure l’état à la sortie et ne récompense pas ; la 3D fait commencer le joueur à chaque vague.
- **Carte 2D** : l’image pixel art existe, mais aucun écran de carte avec marqueur du joueur n’est raccordé.
- **Monde** : pas de génération infinie, saut, nage ou physique verticale. Les combats et la récolte passent par l’arène, pas par des créatures libres dans les biomes.
- **Affichage** : beaucoup de panneaux ont des dimensions fixes adaptées au 1920 × 1080. Le plein écran utilise les dimensions de l’écran, sans forcer cette résolution ni garantir une interface adaptée aux autres formats.
- **Sauvegarde** : un seul emplacement 3D, sans sauvegarde de combat ni menu de gestion des parties. Une sauvegarde invalide interrompt actuellement le démarrage.

## Dépannage

| Problème | Vérification |
|---|---|
| `build constraints exclude all Go files` dans `audio/al` ou `audio/vorbis` | Vérifier que `go env CGO_ENABLED` vaut `1` et que GCC est accessible dans le même terminal. |
| `gcc` introuvable | Ajouter le bon dossier `bin` de MinGW-w64 au `PATH`, puis vérifier avec `where.exe gcc`. |
| `exit status 0xc0000135` | Une DLL native manque. Ajouter `audio/windows/bin` de G3N au `PATH`, ainsi que les éventuelles dépendances de la chaîne GCC. |
| OBJ, JSON, texture ou icône introuvable | Lancer depuis `src/` et conserver tout `assets/` à la racine. |
| Échec du contexte graphique | Vérifier le pilote OpenGL et l’accès à une session graphique locale. |
| Interface coupée | Utiliser si possible un affichage 1920 × 1080 ; la mise en page n’est pas entièrement adaptative. |
| Sauvegarde illisible | Jeu fermé, conserver une copie puis renommer `src/sauvegardes/partie_3d.json` pour recommencer. |
| Combat non retrouvé après fermeture | Revenir à l’exploration avant de quitter : les combats ne sont pas sauvegardés. |
| `Audio indisponible` | Vérifier les WAV dans `assets/audio/`. Un fichier manquant est signalé sans arrêter le jeu. |

## Équipe et ressources

L’ancien README mentionne **gauhthierh, GuillaumeLarre, lalie-droid et quent1206**. Ces crédits sont conservés sans attribuer de rôles non documentés.

Le moteur est **G3N v0.2.0** ; GLFW assure la fenêtre et les entrées, OpenGL le rendu et OpenAL l’audio. Les versions des dépendances Go sont déclarées dans `src/go.mod` et vérifiées par `src/go.sum`.

Aucune licence globale du projet n’est déclarée dans le dépôt actuel. Les licences des dépendances ne constituent pas une autorisation de redistribution de tous les assets du jeu.
