// Génère une petite bande-son RPG originale, sans bibliothèque supplémentaire.
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

const rate = 32000

type morceau struct {
	samples []float64
	boucle  bool
}

func nouveau(secondes float64, boucle bool) *morceau {
	return &morceau{make([]float64, int(secondes*rate)), boucle}
}

// Instruments synthétiques doux : corde pincée, flûte, nappe et cloche.
func (m *morceau) note(debut, duree float64, midi int, volume float64, instrument string) {
	freq := 440 * math.Pow(2, float64(midi-69)/12)
	for n := 0; n < int(duree*rate); n++ {
		t := float64(n) / rate
		phase := 2 * math.Pi * freq * t
		env := math.Min(1, t/.025) * math.Min(1, (duree-t)/.08)
		v := math.Sin(phase)
		switch instrument {
		case "corde":
			// Luth : les harmoniques aiguës s'éteignent avant la fondamentale.
			v = 0
			for h := 1; h <= 6; h++ {
				v += math.Sin(float64(h)*phase) * math.Exp(-t*(2+float64(h)*1.8)) / float64(h*h)
			}
			env = math.Min(1, t/.004) * math.Min(1, (duree-t)/.08)
		case "cordes":
			v = 0
			for h := 1; h <= 5; h++ {
				v += (math.Sin(float64(h)*phase) + .4*math.Sin(float64(h)*phase*1.0015)) / float64(h*h)
			}
			env *= math.Min(1, t/.12)
		case "cor":
			v = .8*v + .35*math.Sin(2*phase) + .18*math.Sin(3*phase) + .07*math.Sin(4*phase)
			env *= math.Min(1, t/.07)
		case "flute":
			v = math.Sin(phase+.025*math.Sin(2*math.Pi*5*t)) + .12*math.Sin(2*phase)
		case "nappe":
			env *= math.Min(1, t/.3)
			v = .6*v + .25*math.Sin(phase*1.002) + .15*math.Sin(phase*2)
		case "cloche":
			v = (v + .3*math.Sin(phase*2.01) + .12*math.Sin(phase*3.98)) * math.Exp(-5*t/duree)
		case "tambour":
			v = math.Sin(2*math.Pi*(65*t+2*(1-math.Exp(-20*t)))) * math.Exp(-18*t)
		}
		m.ajouter(int(debut*rate)+n, v*env*volume)
	}
}

func (m *morceau) ajouter(i int, v float64) {
	if m.boucle {
		i %= len(m.samples)
	}
	if i >= 0 && i < len(m.samples) {
		m.samples[i] += v
	}
}

func (m *morceau) enregistrer(nom string) error {
	// Réflexions courtes : espace discret, sans répéter le coup d'épée.
	sec := append([]float64(nil), m.samples...)
	for i, v := range sec {
		if m.boucle {
			m.ajouter(i+int(.047*rate), v*.13)
			m.ajouter(i+int(.083*rate), v*.08)
			m.ajouter(i+int(.127*rate), v*.04)
		} else {
			m.ajouter(i+int(.019*rate), v*.06)
		}
	}
	peak := .001
	for _, v := range m.samples {
		peak = math.Max(peak, math.Abs(v))
	}
	data := make([]byte, 44+len(m.samples)*2)
	copy(data, "RIFF")
	binary.LittleEndian.PutUint32(data[4:], uint32(len(data)-8))
	copy(data[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(data[16:], 16)
	binary.LittleEndian.PutUint16(data[20:], 1)
	binary.LittleEndian.PutUint16(data[22:], 1)
	binary.LittleEndian.PutUint32(data[24:], rate)
	binary.LittleEndian.PutUint32(data[28:], rate*2)
	binary.LittleEndian.PutUint16(data[32:], 2)
	binary.LittleEndian.PutUint16(data[34:], 16)
	copy(data[36:], "data")
	binary.LittleEndian.PutUint32(data[40:], uint32(len(data)-44))
	for i, v := range m.samples {
		if !m.boucle {
			v *= math.Min(1, float64(len(m.samples)-1-i)/float64(rate/20))
		}
		binary.LittleEndian.PutUint16(data[44+i*2:], uint16(int16(v/peak*.78*32767)))
	}
	return os.WriteFile(nom, data, 0644)
}

func musique(nom string, tempo float64, accords [][3]int, intense bool) *morceau {
	beat := 60 / tempo
	m := nouveau(128*beat, true)
	// Quatre phrases de huit mesures ; la troisième laisse respirer la mélodie.
	for mesure := 0; mesure < 32; mesure++ {
		accord := accords[mesure%len(accords)]
		root := accord[0]
		start := float64(mesure*4) * beat
		for _, n := range accord {
			m.note(start, 4.15*beat, n, .06, "cordes")
		}
		arp := []int{0, 2, 1, 2, 0, 1, 2, 1}
		for pas := 0; pas < 8; pas++ {
			m.note(start+float64(pas)*beat/2, beat*1.3, accord[arp[pas]]+12, .12, "corde")
			if intense {
				m.note(start+float64(pas)*beat/2, beat*.33, accord[pas%3], .10, "cordes")
			}
		}
		m.note(start, beat*3.6, root-12, .16, "cordes")
		// Questions/réponses avec rythmes pointés et silences, au lieu de notes régulières.
		motifs := [][]int{{0, 1, 2, 1}, {2, 1, 0, 2}, {1, 2, 1, 0}, {2, 0, 1, 0}}
		motif := motifs[mesure%4]
		for pas, offset := range []float64{0, 1.5, 2, 3} {
			if mesure >= 16 && mesure < 24 && pas > 1 {
				continue
			}
			n := accord[motif[pas]] + 12
			instrument := "flute"
			if intense {
				instrument = "cor"
			}
			length := []float64{1.35, .4, .85, .7}[pas]
			if mesure%4 == 3 && pas == 3 {
				continue
			}
			m.note(start+offset*beat, length*beat, n, .17, instrument)
		}
		if nom == "menu" && mesure%4 == 0 {
			m.note(start, 2.5*beat, accord[2]+24, .07, "cloche")
		}
		if intense {
			for _, pas := range []float64{0, 1.5, 2.5} {
				m.percussion(start+pas*beat, .5, .22, "tambour")
			}
			for _, pas := range []float64{1, 3} {
				m.percussion(start+pas*beat, .18, .10, "caisse")
			}
			if mesure%4 == 3 {
				for _, pas := range []float64{3, 3.5, 3.75} {
					m.percussion(start+pas*beat, .2, .14, "tambour")
				}
			}
			if nom == "boss" && mesure%4 == 0 {
				m.percussion(start, 1.4, .09, "cymbale")
				m.note(start, 3*beat, root-12, .13, "cor")
			}
		} else if nom == "monde" {
			m.percussion(start, .35, .11, "tambour")
			m.percussion(start+2*beat, .15, .04, "caisse")
		}
	}
	return m
}

func main() {
	dossier := "../assets/audio"
	if len(os.Args) > 1 {
		dossier = os.Args[1]
	}
	if err := os.MkdirAll(dossier, 0755); err != nil {
		panic(err)
	}
	pistes := map[string]*morceau{
		"menu":   musique("menu", 76, [][3]int{{57, 60, 64}, {53, 57, 60}, {48, 52, 55}, {55, 59, 62}, {57, 60, 64}, {50, 53, 57}, {52, 56, 59}, {52, 56, 59}}, false),
		"monde":  musique("monde", 96, [][3]int{{50, 53, 57}, {55, 59, 62}, {48, 52, 55}, {50, 53, 57}, {53, 57, 60}, {55, 59, 62}, {48, 52, 55}, {50, 53, 57}}, false),
		"combat": musique("combat", 126, [][3]int{{52, 55, 59}, {52, 55, 59}, {48, 52, 55}, {50, 54, 57}, {45, 48, 52}, {48, 52, 55}, {47, 51, 54}, {47, 51, 54}}, true),
		"boss":   musique("boss", 144, [][3]int{{50, 53, 57}, {46, 50, 53}, {43, 46, 50}, {45, 49, 52}, {50, 53, 57}, {43, 46, 50}, {46, 50, 53}, {45, 49, 52}}, true),
	}
	for _, nom := range []string{"attaque", "impact", "sort", "potion", "victoire", "defaite"} {
		m := nouveau(2, false)
		switch nom {
		case "attaque":
			m = nouveau(.95, false)
			m.percussion(0, .22, .6, "lame")
			m.percussion(.14, .25, .4, "tambour")
			m.percussion(.15, .7, .3, "metal")
		case "impact":
			m = nouveau(.8, false)
			m.percussion(0, .4, .6, "tambour")
			m.percussion(0, .09, .5, "caisse")
			m.percussion(.015, .5, .15, "metal")
		case "sort":
			m.percussion(0, .45, .16, "lame")
			for i, n := range []int{72, 76, 79, 84, 88, 91} {
				m.note(float64(i)*.08, .65, n, .3, "cloche")
			}
		case "potion":
			for i, n := range []int{76, 79, 84, 88} {
				m.note(float64(i)*.12, .7, n, .3, "cloche")
			}
		case "victoire":
			m = nouveau(4, false)
			for i, n := range []int{72, 76, 79, 84} {
				m.note(float64(i)*.27, 1.2, n, .22, "cor")
			}
			for _, n := range []int{48, 52, 55, 60} {
				m.note(.85, 2.8, n, .12, "cordes")
			}
			m.percussion(.85, 1.4, .09, "cymbale")
		case "defaite":
			m = nouveau(4, false)
			for i, n := range []int{64, 62, 59, 52} {
				m.note(float64(i)*.4, 1.5, n, .2, "cordes")
			}
			for _, n := range []int{40, 47, 55} {
				m.note(1.2, 2.5, n, .1, "nappe")
			}
		}
		pistes[nom] = m
	}
	for nom, m := range pistes {
		if err := m.enregistrer(filepath.Join(dossier, nom+".wav")); err != nil {
			panic(err)
		}
		fmt.Printf("%s : %.1f secondes\n", nom, float64(len(m.samples))/rate)
	}
}
