package monde

import (
	"path/filepath"
	"testing"

	"github.com/g3n/engine/graphic"
)

func TestChargerMap3DAvecTextureRepetee(t *testing.T) {
	dossierCarte := filepath.Join("..", "..", "..", "assets", "maps", "red_world")
	carte, err := ChargerMap3D(
		filepath.Join(dossierCarte, "red_world_map_3d.obj"),
		filepath.Join(dossierCarte, "red_world_map_3d.mtl"),
	)
	if err != nil {
		t.Fatalf("chargement de la carte : %v", err)
	}

	terrain, ok := carte.Children()[0].(*graphic.Mesh)
	if !ok {
		t.Fatal("le premier objet de la carte n'est pas le terrain")
	}
	if terrain.GetMaterial(0).GetMaterial().TextureCount() != 1 {
		t.Fatal("la texture répétitive n'est pas attachée au terrain")
	}
}
