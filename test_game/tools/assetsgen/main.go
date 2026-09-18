// Command assetsgen creates the editable low-poly placeholder asset pack.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type assetEntry struct {
	Path        string `json:"path"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

type catalog struct {
	root    string
	entries []assetEntry
}

func main() {
	root, err := projectRoot()
	must(err)
	assets := &catalog{root: filepath.Join(filepath.Dir(root), "assets")}

	must(generatePalette(assets))
	must(generateModels(assets))
	must(generateIcons(assets))
	must(generateUI(assets))
	must(generateAnimations(assets))
	must(assets.writeManifest())

	fmt.Printf("%d assets generes dans %s\n", len(assets.entries), assets.root)
}

func projectRoot() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("go.mod introuvable")
		}
		directory = parent
	}
}

func (c *catalog) write(relativePath, kind, description string, data []byte) error {
	path := filepath.Join(c.root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	c.entries = append(c.entries, assetEntry{Path: filepath.ToSlash(relativePath), Kind: kind, Description: description})
	return nil
}

func (c *catalog) writeManifest() error {
	sort.Slice(c.entries, func(i, j int) bool { return c.entries[i].Path < c.entries[j].Path })
	data, err := json.MarshalIndent(struct {
		Style  string       `json:"style"`
		Assets []assetEntry `json:"assets"`
	}{
		Style:  "RPG fantasy low-poly, couleurs plates, placeholders editables",
		Assets: c.entries,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.root, "manifest.json"), append(data, '\n'), 0o644)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}
}
