package main

import (
	"fmt"
	"ordta/library"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Depuis src : go run . propose le choix entre les deux interfaces.
func main() {
	fmt.Println("=== ORDTA ===\n1 — Jouer en 3D\n2 — Jouer dans le terminal\n0 — Quitter")
	switch library.ReadChoiceEntre("Votre choix : ", 2) {
	case 1:
		if err := lancer3D(); err != nil {
			fmt.Println("Impossible de lancer la 3D :", err)
		}
	case 2:
		personnage := library.CharacterCreation()
		personnage.MainMenu()
	}
}

// Le CLI reste indépendant de G3N et de ses dépendances natives.
func lancer3D() error {
	if _, err := exec.LookPath("gcc"); err != nil {
		return fmt.Errorf("GCC doit être accessible dans ce terminal (vérifiez gcc --version)")
	}
	commande := exec.Command("go", "run", "./3D")
	commande.Stdin, commande.Stdout, commande.Stderr = os.Stdin, os.Stdout, os.Stderr
	commande.Env = append(os.Environ(), "CGO_ENABLED=1")
	if runtime.GOOS == "windows" {
		dossier, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/g3n/engine").Output()
		if err != nil {
			return fmt.Errorf("dépendance G3N introuvable : %w", err)
		}
		dll := filepath.Join(strings.TrimSpace(string(dossier)), "audio", "windows", "bin")
		commande.Env = append(commande.Env, "PATH="+dll+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	return commande.Run()
}
