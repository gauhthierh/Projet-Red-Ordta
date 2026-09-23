# Assets low-poly

Ce dossier contient un pack de placeholders simples et retouchables pour le RPG :

```text
assets/
├── animations/          keyframes indicatifs pour le joueur et le gobelin
├── backgrounds/         village, forge et arène en SVG
├── materials/           palette MTL commune aux modèles
├── models/
│   ├── characters/      Humain, Elfe, Nain et Gobelin
│   ├── environment/     village, marchand, forge et arène
│   ├── equipment/       chapeau, tunique et bottes
│   ├── items/           potions, ressources, livre, pièce et sac
│   └── vfx/             attaque, feu, soin, poison et montée de niveau
├── pixel_art/
│   ├── atlases/         3 planches originales de 16 éléments
│   └── icons/           48 icônes détaillées détourées en PNG transparent
├── ui/
│   ├── components/      panneaux, boutons, cases et barres
│   ├── icons/            versions PNG et SVG
│   └── screens/          maquettes des différents écrans
└── manifest.json         catalogue complet
```

Les modèles utilisent le format Wavefront OBJ et la palette `materials/low_poly.mtl`. Ils sont volontairement composés de peu de formes et n'utilisent pas de textures.

Les SVG sont les sources éditables. Les PNG de 256 × 256 pixels sont immédiatement utilisables pour les icônes du jeu.

Le nouveau pack `pixel_art` propose une direction artistique plus détaillée : gros pixels, contours nets, matières lisibles et éclairage homogène. Son propre `manifest.json` permet de charger chaque icône individuellement. Les planches originales sont conservées pour pouvoir les retoucher puis relancer, depuis la racine du dépôt, `go run ./test_game/tools/atlascrop/main.go`.
