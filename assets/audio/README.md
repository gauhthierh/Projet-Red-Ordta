# Audio RPG

Compositions synthétiques originales générées en Go, sans samples externes.
Sonorité médiévale-fantastique synthétique : luth, flûte, cordes frottées,
cors, cloches, tambours et percussions bruitées. Ce ne sont pas des
enregistrements orchestraux. Les thèmes comprennent quatre phrases de huit
mesures, avec variations de rythme et un passage plus calme.

- `menu.wav` : accueil et pause, environ 101 s, luth et flûte.
- `monde.wav` : exploration, 80 s, couleur dorienne et percussion légère.
- `combat.wav` : vagues ordinaires, environ 61 s, cors et cordes rythmées.
- `boss.wav` : dernière vague, environ 53 s, percussions et cors graves.
- `victoire.wav`, `defaite.wav` : fanfare et conclusion sombre, 4 s chacune.
- `attaque.wav` : souffle de lame puis choc métallique, 0,95 s.
- `impact.wav` : choc sourd et bref bruit d'armure, 0,8 s.
- `sort.wav`, `potion.wav` : montée magique et carillon de soin.

Les quatre ambiances bouclent et se croisent avec un fondu. En pause les
effets sont arrêtés. Les volumes se règlent dans `src/3D/monde/audio.go` :
0,22 pour la musique et 0,4 pour les effets. Aucun fichier backend ne gère l'audio.

Pour régénérer depuis `src` : `go run ./tools/musicgen`.
Format attendu : WAV PCM mono, 16 bits, 32 kHz, en-tête de 44 octets.
Un fichier manquant est signalé dans la console sans arrêter le jeu.
