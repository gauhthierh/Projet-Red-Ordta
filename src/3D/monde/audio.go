package monde

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/g3n/engine/audio/al"
)

// AudioJeu garde la lecture sur le fil principal du jeu, comme OpenAL.
// Les volumes sont volontairement bas pour laisser les effets audibles.
type AudioJeu struct {
	buffers       map[string]uint32
	musiques      map[string]uint32
	volumes       map[string]float32
	effets        []uint32
	prochainEffet int
	ambiance      string
}

func NouvelAudioJeu(dossier string) *AudioJeu {
	a := &AudioJeu{buffers: map[string]uint32{}, musiques: map[string]uint32{}, volumes: map[string]float32{}}
	for _, nom := range []string{"menu", "monde", "combat", "boss", "victoire", "defaite", "attaque", "impact", "sort", "potion"} {
		data, err := os.ReadFile(filepath.Join(dossier, nom+".wav"))
		// Format produit par tools/musicgen : PCM mono 16 bits, 32 kHz.
		if err != nil || len(data) < 46 {
			fmt.Println("Audio indisponible :", nom, err)
			continue
		}
		if string(data[:4]) != "RIFF" || string(data[8:16]) != "WAVEfmt " || string(data[36:40]) != "data" || binary.LittleEndian.Uint16(data[20:]) != 1 || binary.LittleEndian.Uint16(data[22:]) != 1 || binary.LittleEndian.Uint16(data[34:]) != 16 || binary.LittleEndian.Uint32(data[24:]) != 32000 || len(data[44:])%2 != 0 {
			fmt.Println("Format audio non pris en charge :", nom)
			continue
		}
		buffer := al.GenBuffers(1)[0]
		al.BufferData(buffer, al.FormatMono16, unsafe.Pointer(&data[44]), uint32(len(data)-44), 32000)
		if err := al.GetError(); err != nil {
			al.DeleteBuffers([]uint32{buffer})
			fmt.Println("Audio désactivé :", err)
			continue
		}
		a.buffers[nom] = buffer
	}
	for _, nom := range []string{"menu", "monde", "combat", "boss", "victoire", "defaite"} {
		if a.buffers[nom] == 0 {
			continue
		}
		source := al.GenSource()
		al.Sourcei(source, al.Buffer, int32(a.buffers[nom]))
		al.Sourcei(source, al.SourceRelative, 1)
		al.Sourcef(source, al.RolloffFactor, 0)
		if nom != "victoire" && nom != "defaite" {
			al.Sourcei(source, al.Looping, 1)
		}
		al.Sourcef(source, al.Gain, 0)
		a.musiques[nom] = source
	}
	if len(a.buffers) > 0 {
		for i := 0; i < 4; i++ {
			source := al.GenSource()
			al.Sourcei(source, al.SourceRelative, 1)
			al.Sourcef(source, al.RolloffFactor, 0)
			al.Sourcef(source, al.Gain, .4)
			a.effets = append(a.effets, source)
		}
	}
	return a
}

func (a *AudioJeu) Ambiance(nom string, delta float32) {
	if a.ambiance != nom {
		a.ambiance = nom
		if source := a.musiques[nom]; source != 0 {
			al.SourcePlay(source)
		}
	}
	for piste, source := range a.musiques {
		cible := float32(0)
		if piste == nom {
			cible = .22
		}
		volume := a.volumes[piste]
		pas := delta * .3
		if volume < cible {
			volume += pas
			if volume > cible {
				volume = cible
			}
		}
		if volume > cible {
			volume -= pas
			if volume < cible {
				volume = cible
			}
		}
		al.Sourcef(source, al.Gain, volume)
		if volume == 0 && piste != nom {
			al.SourceStop(source)
		}
		a.volumes[piste] = volume
	}
}

func (a *AudioJeu) Effet(nom string) {
	if a.buffers[nom] == 0 || len(a.effets) == 0 {
		return
	}
	source := a.effets[a.prochainEffet%len(a.effets)]
	a.prochainEffet++
	al.SourceStop(source)
	al.Sourcei(source, al.Buffer, int32(a.buffers[nom]))
	al.SourcePlay(source)
}

func (a *AudioJeu) ArreterEffets() {
	for _, s := range a.effets {
		al.SourceStop(s)
	}
}

func (a *AudioJeu) Fermer() {
	for _, s := range a.musiques {
		al.SourceStop(s)
		al.DeleteSource(s)
	}
	for _, s := range a.effets {
		al.SourceStop(s)
		al.DeleteSource(s)
	}
	for _, b := range a.buffers {
		al.DeleteBuffers([]uint32{b})
	}
}
