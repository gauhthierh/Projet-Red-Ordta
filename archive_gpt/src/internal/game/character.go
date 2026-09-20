package game

import (
	"fmt"
	"strings"
	"unicode"
)

func NormalizeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("Saisissez un nom pour votre personnage")
	}
	runes := []rune(name)
	if len(runes) > 20 {
		return "", fmt.Errorf("Le nom ne peut pas depasser 20 lettres")
	}
	for _, r := range runes {
		if !unicode.IsLetter(r) {
			return "", fmt.Errorf("Le nom doit contenir uniquement des lettres")
		}
	}
	for i, r := range runes {
		if i == 0 {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
	}
	return string(runes), nil
}

func ClassHealth(class string) (int, error) {
	switch class {
	case "Humain":
		return 100, nil
	case "Elfe":
		return 80, nil
	case "Nain":
		return 120, nil
	default:
		return 0, fmt.Errorf("Choisissez Humain, Elfe ou Nain")
	}
}

func (d *Demo) CreateCharacter(name, class string) error {
	proper, err := NormalizeName(name)
	if err != nil {
		return err
	}
	maxHP, err := ClassHealth(class)
	if err != nil {
		return err
	}
	if d.state.Fighting {
		return fmt.Errorf("Impossible de creer un personnage pendant un combat")
	}
	d.state = Snapshot{
		Name: proper, Class: class, Level: 1, HP: maxHP / 2, MaxHP: maxHP,
		Gold: 100, EnemyHP: 40, Skills: []string{"Coup de poing"},
		Inventory: []ItemStack{{"health_potion", 3}}, Equipment: []string{"", "", ""},
		Message: fmt.Sprintf("Bienvenue %s, %s de niveau 1 !", proper, class),
	}
	d.slots = [36]string{}
	d.slots[0] = "health_potion"
	return nil
}

func equipmentSlot(id string) (int, int, bool) {
	switch id {
	case "adventurer_hat":
		return 0, 10, true
	case "adventurer_tunic":
		return 1, 25, true
	case "adventurer_boots":
		return 2, 15, true
	default:
		return 0, 0, false
	}
}

func equippedCount(equipment []string) int {
	n := 0
	for _, id := range equipment {
		if id != "" {
			n++
		}
	}
	return n
}

func (d *Demo) EquipSlot(slot int) error {
	if d.state.Fighting {
		return fmt.Errorf("Equipement indisponible pendant le combat")
	}
	if slot < 0 || slot >= len(d.slots) {
		return fmt.Errorf("Emplacement invalide")
	}
	id := d.slots[slot]
	place, bonus, ok := equipmentSlot(id)
	if !ok {
		return fmt.Errorf("Cet objet n'est pas une piece d'armure")
	}
	if d.count(id) == 0 {
		return fmt.Errorf("Objet absent du sac")
	}
	if len(d.state.Equipment) < 3 {
		d.state.Equipment = []string{"", "", ""}
	}
	previous := d.state.Equipment[place]
	if previous == id {
		return fmt.Errorf("Cette armure est deja equipee")
	}
	d.add(id, -1)
	if previous != "" {
		_, oldBonus, _ := equipmentSlot(previous)
		d.state.MaxHP -= oldBonus
		d.add(previous, 1)
	}
	d.state.Equipment[place] = id
	d.state.MaxHP += bonus
	if d.state.HP > d.state.MaxHP {
		d.state.HP = d.state.MaxHP
	}
	d.state.Message = fmt.Sprintf("%s equipe (+%d PV max).", ItemName(id), bonus)
	return nil
}

func (d *Demo) Unequip(place int) error {
	if d.state.Fighting {
		return fmt.Errorf("Equipement indisponible pendant le combat")
	}
	if place < 0 || place >= 3 {
		return fmt.Errorf("Emplacement d'armure invalide")
	}
	if len(d.state.Equipment) < 3 || d.state.Equipment[place] == "" {
		return fmt.Errorf("Aucune armure equipee ici")
	}
	id := d.state.Equipment[place]
	_, bonus, _ := equipmentSlot(id)
	d.state.Equipment[place] = ""
	d.state.MaxHP -= bonus
	if d.state.HP > d.state.MaxHP {
		d.state.HP = d.state.MaxHP
	}
	d.add(id, 1)
	d.state.Message = fmt.Sprintf("%s retire et replace dans le sac.", ItemName(id))
	return nil
}
