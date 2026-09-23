package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/loader/obj"
)

func TestGeneratedMapAssets(t *testing.T) {
	dir := t.TempDir()
	mtlPath := filepath.Join(dir, "red_world_map_3d.mtl")
	objPath := filepath.Join(dir, "red_world_map_3d.obj")
	layoutPath := filepath.Join(dir, "red_world_layout.json")
	textureDir := filepath.Join(dir, "textures")
	if err := os.MkdirAll(textureDir, 0755); err != nil {
		t.Fatal(err)
	}
	textureFile, err := os.Create(filepath.Join(textureDir, "grass_pixel.png"))
	if err != nil {
		t.Fatal(err)
	}
	textureImage := image.NewRGBA(image.Rect(0, 0, 1, 1))
	textureImage.SetRGBA(0, 0, color.RGBA{R: 70, G: 130, B: 50, A: 255})
	if err := png.Encode(textureFile, textureImage); err != nil {
		textureFile.Close()
		t.Fatal(err)
	}
	if err := textureFile.Close(); err != nil {
		t.Fatal(err)
	}

	writeMaterials(mtlPath)
	writeWorld(objPath)
	writeLayout(layoutPath)

	decoder, err := obj.Decode(objPath, mtlPath)
	if err != nil {
		t.Fatalf("le modèle 3D ne peut pas être décodé : %v", err)
	}
	if len(decoder.Objects) < 15 {
		t.Fatalf("les groupes de matériaux semblent incomplets : %d", len(decoder.Objects))
	}
	if decoder.Vertices.Size()/3 < 80000 {
		t.Fatalf("la carte 3D manque de détails : %d sommets", decoder.Vertices.Size()/3)
	}
	if decoder.Materials["grass"].MapKd == "" {
		t.Fatal("la texture d'herbe n'est pas déclarée dans le MTL")
	}
	// Ce test valide la géométrie dans un dossier temporaire. Le chargeur du
	// jeu teste séparément les véritables textures du projet.
	for _, description := range decoder.Materials {
		description.MapKd = ""
	}
	group, err := decoder.NewGroup()
	if err != nil {
		t.Fatalf("G3N ne peut pas construire le groupe 3D : %v", err)
	}
	firstMesh, ok := group.Children()[0].(*graphic.Mesh)
	if !ok {
		t.Fatal("le terrain principal n'est pas un maillage G3N")
	}
	if firstMesh.GetMaterial(0) == nil {
		t.Fatal("le terrain principal n'a pas de matériau")
	}

	data, err := os.ReadFile(layoutPath)
	if err != nil {
		t.Fatal(err)
	}
	var plan layout
	if err := json.Unmarshal(data, &plan); err != nil {
		t.Fatalf("le plan de coordonnées est invalide : %v", err)
	}
	if plan.UpAxis != "+Z" || len(plan.Landmarks) < 10 {
		t.Fatal("le plan ne décrit pas toutes les zones en Z vertical")
	}
	if len(plan.Collisions) < 250 {
		t.Fatalf("le plan de collision semble incomplet : %d volumes", len(plan.Collisions))
	}
	for _, c := range plan.Collisions {
		if c.ID == "" || (c.Shape != "circle" && c.Shape != "rectangle") {
			t.Fatalf("collision invalide : %+v", c)
		}
		if c.Category == "fence" {
			angle := c.RotationDegrees * math.Pi / 180
			for j := 0; j <= 100; j++ {
				offset := (float64(j)/100 - .5) * c.Width
				x, y := c.X+math.Cos(angle)*offset, c.Y+math.Sin(angle)*offset
				if onRoad(x, y, 4.5) {
					t.Fatalf("la clôture %s bloque un chemin en %.2f, %.2f", c.ID, x, y)
				}
			}
		}
	}
}
