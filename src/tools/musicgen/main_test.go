package main

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestWAVGeneres(t *testing.T) {
	for _, nom := range []string{"menu", "monde", "combat", "boss", "attaque", "impact", "sort", "potion", "victoire", "defaite"} {
		t.Run(nom, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("../../../assets/audio", nom+".wav"))
			if err != nil {
				t.Fatal(err)
			}
			if len(data) < 46 || string(data[:4]) != "RIFF" || string(data[8:16]) != "WAVEfmt " || string(data[36:40]) != "data" {
				t.Fatal("WAV invalide")
			}
			if int(binary.LittleEndian.Uint32(data[40:])) != len(data)-44 || binary.LittleEndian.Uint32(data[24:]) != rate {
				t.Fatal("taille ou fréquence incorrecte")
			}
			energy := 0.0
			for i := 44; i < len(data); i += 2 {
				x := float64(int16(binary.LittleEndian.Uint16(data[i:])))
				if math.Abs(x) > 26000 {
					t.Fatal("niveau trop élevé")
				}
				energy += x * x
			}
			if energy/float64((len(data)-44)/2) < 100 {
				t.Fatal("piste silencieuse")
			}
			if nom == "menu" || nom == "monde" || nom == "combat" || nom == "boss" {
				first := int(int16(binary.LittleEndian.Uint16(data[44:])))
				last := int(int16(binary.LittleEndian.Uint16(data[len(data)-2:])))
				if math.Abs(float64(first-last)) > 2500 {
					t.Fatal("raccord de boucle trop abrupt")
				}
			}
		})
	}
}
