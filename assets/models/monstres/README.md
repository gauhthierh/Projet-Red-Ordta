# Monstres de l'arène

Cinq modèles fixes en JSON : gobelin, gobelin_cuirasse, chaman, loup et troll.
Ils utilisent les textures du personnage (copies locales dans textures).
Les JSON décrivent des pièces placées à la main : boîtes et ellipsoïdes facettés.
Aucune génération aléatoire, statistique ou animation automatique.

Les pieds sont à Z=0 et le regard dirigé vers -Y. Les dimensions sont proches
du personnage : gobelins environ 2 unités, loup 1,7, troll environ 4.
La taille, la position et la couleur de chaque pièce sont modifiables dans le JSON.
Ce sont des modèles G3N assemblés, pas des fichiers OBJ ni des squelettes skinnés.

Depuis src, charger avec monstres.Charger("../assets/models/monstres/gobelin.json").
La fonction retourne le modèle et une erreur à vérifier. Ajouter ensuite son
Noeud() à la scène, comme pour le personnage. Plusieurs chargements donnent
des instances indépendantes. Les modèles ne sont pas ajoutés automatiquement
au monde : le futur système de vagues contrôlera leur apparition.

Membres des humanoïdes : corps, tete, bras_gauche, bras_droit, jambe_gauche,
jambe_droite et arme ; bouclier en plus pour le cuirassé.
Membres du loup : corps, tete, queue, patte_avant_gauche, patte_avant_droite,
patte_arriere_gauche, patte_arriere_droite.
Accéder à Membres["bras_droit"] pour tourner son pivot. Arme et bouclier
ont leurs propres pivots : leur donner la même rotation que le bras porteur.

Le package monde fournit PositionJoueurArene, PositionsMonstresArene et
PlacerCameraArene. Le joueur regarde +X avec une rotation Z de Pi/2 ; les
monstres regardent -X avec -Pi/2. Appeler la caméra d'arène dans la boucle
de combat, libérer la souris à l'entrée et suspendre le regard souris.
Ces changements de mode restent à raccorder au combat existant de library.

Les assets d'interface existants sont dans assets/ui/screens/combat.svg,
les icônes dans assets/ui/icons, et les effets dans assets/models/vfx.
Leur affichage et leur animation restent à intégrer avec le système de combat.
