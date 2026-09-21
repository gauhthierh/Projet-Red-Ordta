# RED World

Cette carte reprend les fonctionnalités du sujet dans un monde semi-ouvert unique.

## Fichiers

- `red_world_map_2d_pixel.png` : carte pixel art détaillée à afficher dans l'interface du jeu. C'est la version principale prévue pour recevoir le marqueur du joueur.
- `red_world_map_2d.png` : version illustrée précédente, conservée comme référence.
- `red_world_map_3d.obj` : décor 3D statique complet, directement en Z vertical.
- `red_world_map_3d.mtl` : matériaux du décor.
- `textures/grass_pixel.png` : première texture pixel art, appliquée au matériau d'herbe grâce aux coordonnées UV du modèle OBJ.
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

Le tableau `collisions` du JSON reprend les empreintes solides du décor sur le plan XY. Les formes `circle` utilisent `radius`. Les formes `rectangle` utilisent `width`, `height` et, lorsque nécessaire, `rotation_degrees`. Les lampadaires sont solides ; les cultures, fleurs, bannières et autres petits ornements restent volontairement traversables afin de ne pas gêner les déplacements.

## Zones liées au sujet

Le village central contient le marché. La guilde sert à la création et à l'équipement du personnage. La forge correspond à la fabrication. L'arène accueille le gobelin d'entraînement. Le sanctuaire représente la résurrection. La forêt, les prés, les falaises et le marais correspondent respectivement aux ressources du loup, du sanglier, du corbeau et du troll. Le bosquet de mana prépare la mission optionnelle sur la magie.
