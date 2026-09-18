# Detailed pixel-art asset pack

This folder contains the more detailed fantasy pixel-art pass requested for the RPG.

- `atlases/`: the three original 4x4 generated sheets, kept intact for retouching.
- `icons/`: 48 individually cropped PNG assets with transparent backgrounds.
- `manifest.json`: stable names and relative paths for loading the icons from Go.
- `PROMPTS.md`: the generation prompt set used to keep the pack stylistically consistent.

The earlier geometric SVG/PNG placeholders remain under `assets/ui`; this pack is separate so no previous work is lost.

To regenerate the crops after editing an atlas, run from `test_game`:

```text
go run ./tools/atlascrop
```
