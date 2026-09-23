package main

import (
	"fmt"
	"ordta/lancement"
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
	if _, err := exec.LookPath("gcc"); err != nil && runtime.GOOS != "darwin" {
		return fmt.Errorf("GCC est introuvable. Installez GCC via https://www.msys2.org/, ajoutez C:\\msys64\\ucrt64\\bin au Path Windows, puis rouvrez le terminal (voir README)")
	}
	commande := exec.Command("go", "run", "./3D")
	commande.Stdin, commande.Stdout, commande.Stderr = os.Stdin, os.Stdout, os.Stderr
	commande.Env = append(os.Environ(), "CGO_ENABLED=1")
	if runtime.GOOS == "darwin" {
		if err := lancement.PreparerMac(commande); err != nil {
			return err
		}
	}
	if runtime.GOOS == "windows" {
		cibleGCC, err := exec.Command("gcc", "-dumpmachine").Output()
		if err != nil || !strings.HasPrefix(strings.TrimSpace(string(cibleGCC)), "x86_64-") {
			return fmt.Errorf("la 3D demande GCC 64 bits : placez C:\\msys64\\ucrt64\\bin en premier dans le Path, puis rouvrez le terminal")
		}
		// G3N v0.2.0 utilise des constantes qui débordent les int en 32 bits.
		// Ne changer que le processus 3D, même si GOARCH=386 est enregistré dans Go.
		commande.Env = append(commande.Env, "GOOS=windows", "GOARCH=amd64", "CC=gcc")
		// Au premier lancement, G3N peut ne pas encore être dans le cache Go.
		// Ce téléchargement fournit également ses DLL audio ; ensuite le cache est réutilisé.
		telechargement := exec.Command("go", "mod", "download", "github.com/g3n/engine")
		telechargement.Stdout, telechargement.Stderr = os.Stdout, os.Stderr
		if err := telechargement.Run(); err != nil {
			return fmt.Errorf("téléchargement de G3N impossible : vérifiez votre connexion Internet : %w", err)
		}
		dossier, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/g3n/engine").Output()
		if err != nil {
			return fmt.Errorf("dépendance G3N introuvable : %w", err)
		}
		dll := filepath.Join(strings.TrimSpace(string(dossier)), "audio", "windows", "bin")
		if _, err := os.Stat(filepath.Join(dll, "OpenAL32.dll")); err != nil {
			return fmt.Errorf("bibliothèques audio G3N introuvables dans %s : %w", dll, err)
		}
		commande.Env = append(commande.Env, "PATH="+dll+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	return commande.Run()
}
