package lancement

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Ce fichier n'importe pas G3N : le choix CLI reste indépendant de l'audio.
func PreparerMac(commande *exec.Cmd) error {
	clang, err := exec.Command("xcrun", "--find", "clang").Output()
	if err != nil {
		return fmt.Errorf("installez les outils Apple avec xcode-select --install, puis relancez le terminal")
	}
	if runtime.GOARCH != "arm64" && runtime.GOARCH != "amd64" {
		return fmt.Errorf("architecture Mac non prise en charge : %s", runtime.GOARCH)
	}
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew est introuvable : installez-le depuis https://brew.sh et suivez ses Next steps pour le PATH")
	}
	prefixes := map[string]string{}
	for _, paquet := range []string{"openal-soft", "libvorbis", "libogg"} {
		sortie, err := exec.Command("brew", "--prefix", paquet).Output()
		if err != nil {
			return fmt.Errorf("dépendances audio manquantes : exécutez brew install openal-soft libvorbis")
		}
		prefixe := strings.TrimSpace(string(sortie))
		if _, err := os.Stat(filepath.Join(prefixe, "include")); err != nil {
			return fmt.Errorf("%s n'est pas installé : exécutez brew install openal-soft libvorbis", paquet)
		}
		prefixes[paquet] = prefixe
	}
	cflags, ldflags := drapeauxAudioMac(prefixes)
	commande.Env = append(commande.Env,
		"GOOS=darwin", "GOARCH="+runtime.GOARCH,
		"CC="+strings.TrimSpace(string(clang)),
		"CGO_CFLAGS="+strings.TrimSpace(os.Getenv("CGO_CFLAGS")+" "+cflags),
		"CGO_LDFLAGS="+strings.TrimSpace(os.Getenv("CGO_LDFLAGS")+" "+ldflags))
	return nil
}

// Les préfixes sont demandés à Homebrew, pas supposés identiques sur tous les Mac.
func drapeauxAudioMac(prefixes map[string]string) (string, string) {
	var includes, liens []string
	for _, paquet := range []string{"openal-soft", "libvorbis", "libogg"} {
		prefixe := prefixes[paquet]
		includes = append(includes, fmt.Sprintf("%q", "-I"+filepath.Join(prefixe, "include")))
		if paquet == "openal-soft" {
			includes = append(includes, fmt.Sprintf("%q", "-I"+filepath.Join(prefixe, "include", "AL")))
		}
		if paquet == "libvorbis" {
			includes = append(includes, fmt.Sprintf("%q", "-I"+filepath.Join(prefixe, "include", "vorbis")))
		}
		liens = append(liens, fmt.Sprintf("%q", "-L"+filepath.Join(prefixe, "lib")))
	}
	return strings.Join(includes, " "), strings.Join(liens, " ")
}
