# RED World

Cette carte reprend les fonctionnalités du sujet dans un monde semi-ouvert unique.

## Fichiers

- `red_world_map_2d_pixel.png` : carte pixel art détaillée à afficher dans l'interface du jeu. C'est la version principale prévue pour recevoir le marqueur du joueur.
- `red_world_map_2d.png` : version illustrée précédente, conservée comme référence.
- `red_world_map_3d.obj` : décor 3D statique complet, directement en Z vertical.
- `red_world_map_3d.mtl` : matériaux du décor.
- `textures/grass_detailed.png` : herbe, sous-bois et variantes teintées.
- `textures/path_detailed.png` : chemins, sable et sols secs.
- `textures/stone_detailed.png` : remparts, maisons, ruines, rochers et neige teintée.
- `textures/wood_detailed.png` : ponts, clôtures, étals et accessoires.
- `textures/water_detailed.png` : rivières, fontaines et magie.
- `textures/roof_detailed.png` : tuiles teintées en rouge ou en bleu.
- `textures/crop_detailed.png` : champs cultivés.
- `textures/marsh_detailed.png` : sol humide du marais.
- `textures/metal_detailed.png` : forge, lampadaires et équipement.
- `textures/cloth_detailed.png` : tentes, bannières et auvents.
- `red_world_collision.obj` : représentation visuelle des volumes de collision, régénérée avec la carte pour faciliter leur contrôle.
- `red_world_layout.json` : dimensions, point d'apparition, coordonnées des zones, volumes de collision et formule de conversion entre l'image et le monde.

Le modèle 3D n'est pas généré pendant l'exécution. Le générateur Go situé dans `src/tools/mapgen` est uniquement un outil d'édition qui permet de reconstruire les fichiers statiques après une modification du plan.

## Repère commun

- Taille du monde : 600 x 600 unités.
- Nord : axe +Y.
- Hauteur : axe +Z.
- Centre de la carte et du village : X=0, Y=0.
- Conversion d'un pixel de la carte 2D vers le monde : voir `image_to_world` dans le fichier JSON.

## Collisions

Le tableau `collisions` du JSON reprend les empreintes solides du décor sur le plan XY. Les formes `circle` utilisent `radius`. Les formes `rectangle` utilisent `width`, `height` et, lorsque nécessaire, `rotation_degrees`. Les bâtiments, murs, clôtures, arbres, rochers, lampadaires, ruines, accessoires importants et zones d'eau sont couverts. Les collisions des rivières sont interrompues aux ponts. Les cultures, fleurs, bannières, petits buissons et autres ornements restent volontairement traversables afin de ne pas gêner les déplacements.

Le déplacement contrôle séparément les axes X et Y, ce qui permet au personnage de glisser contre les obstacles. Les rectangles orientés utilisent le même test que les rectangles droits après conversion de la position dans leur repère local.

## Zones liées au sujet

Le village central contient le marché, la fontaine, le château, les maisons et quatre portes. La guilde possède son terrain d'entraînement. La forge comprend sa cour de fabrication. L'arène accueille le gobelin d'entraînement. Le sanctuaire et le cimetière représentent la résurrection. La forêt, les prés, les falaises et le marais correspondent respectivement aux ressources du loup, du sanglier, du corbeau et du troll. Le bosquet de mana contient sa source, ses cristaux et ses ruines. Des chemins courbes relient toutes les zones et les sols de biome permettent de les distinguer en jeu.
