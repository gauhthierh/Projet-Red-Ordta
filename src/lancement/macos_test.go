package lancement

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheminsAudioMac(t *testing.T) {
	for _, racine := range []string{"/opt/homebrew/opt", "/usr/local/opt", "/installation avec espaces/opt"} {
		prefixes := map[string]string{}
		for _, p := range []string{"openal-soft", "libvorbis", "libogg"} {
			prefixes[p] = filepath.Join(racine, p)
		}
		includes, liens := drapeauxAudioMac(prefixes)
		for _, p := range []string{"openal-soft", "libvorbis", "libogg"} {
			if !strings.Contains(includes, fmt.Sprintf("%q", "-I"+filepath.Join(prefixes[p], "include"))) || !strings.Contains(liens, fmt.Sprintf("%q", "-L"+filepath.Join(prefixes[p], "lib"))) {
				t.Fatalf("chemin absent pour %s : %s / %s", p, includes, liens)
			}
		}
	}
}
