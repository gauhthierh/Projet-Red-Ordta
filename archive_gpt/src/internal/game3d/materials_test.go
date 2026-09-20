package game3d

import (
	"projet-red/internal/worldmap"
	"testing"
)

func TestProceduralMaterialHasVariation(t *testing.T) {
	for _, kind := range []string{"wood", "stone", "roof", "plaster", "cloth", "metal", "skin"} {
		img := surfaceImage(kind)
		if img.Bounds().Dx() != 64 || img.Bounds().Dy() != 64 {
			t.Fatalf("wrong texture size for %s", kind)
		}
		first := img.RGBAAt(0, 0)
		varied := false
		for y := 0; y < 64 && !varied; y++ {
			for x := 0; x < 64; x++ {
				if img.RGBAAt(x, y) != first {
					varied = true
					break
				}
			}
		}
		if !varied {
			t.Fatalf("flat texture: %s", kind)
		}
	}
}

func TestTerrainPaletteChangesWithBiome(t *testing.T) {
	id := worldmap.SectorID{X: 5, Z: 5}
	a := terrainTexture(id, worldmap.Forest)
	b := terrainTexture(id, worldmap.Highland)
	if a.RGBAAt(20, 20) == b.RGBAAt(20, 20) {
		t.Fatal("biome terrain is identical")
	}
}
