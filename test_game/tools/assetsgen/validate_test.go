package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"image"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratedAssetPack(t *testing.T) {
	root, err := projectRoot()
	if err != nil {
		t.Fatal(err)
	}
	assetsRoot := filepath.Join(filepath.Dir(root), "assets")

	manifestData, err := os.ReadFile(filepath.Join(assetsRoot, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Assets []assetEntry `json:"assets"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Assets) < 100 {
		t.Fatalf("asset pack too small: %d entries", len(manifest.Assets))
	}

	for _, asset := range manifest.Assets {
		path := filepath.Join(assetsRoot, filepath.FromSlash(asset.Path))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("missing %s: %v", asset.Path, err)
			continue
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".obj":
			if !bytes.Contains(data, []byte("mtllib ../../materials/low_poly.mtl")) ||
				!bytes.Contains(data, []byte("\nv ")) || !bytes.Contains(data, []byte("\nf ")) {
				t.Errorf("invalid OBJ structure: %s", asset.Path)
			}
		case ".png":
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				t.Errorf("invalid PNG %s: %v", asset.Path, err)
				continue
			}
			if asset.Kind == "icon-png" && (decoded.Bounds().Dx() != 256 || decoded.Bounds().Dy() != 256) {
				t.Errorf("unexpected icon dimensions for %s: %v", asset.Path, decoded.Bounds())
			}
		case ".svg":
			decoder := xml.NewDecoder(bytes.NewReader(data))
			for {
				if _, err := decoder.Token(); err != nil {
					if err != io.EOF {
						t.Errorf("invalid SVG %s: %v", asset.Path, err)
					}
					break
				}
			}
		}
	}
}
