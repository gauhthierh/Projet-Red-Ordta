package game

import "fmt"

// Demo is a playable temporary backend, replaceable independently of the UI.
type Demo struct {
	state Snapshot
	slots [36]string
}

func NewDemo() *Demo {
	d := &Demo{state: Snapshot{
		Name: "Aventurier", Class: "Humain", Level: 1, HP: 100, MaxHP: 100, Gold: 100, EnemyHP: 40,
		Skills: []string{"Coup de poing"}, Equipment: []string{"", "", ""},
		Inventory: []ItemStack{{"health_potion", 3}, {"wolf_fur", 2}, {"boar_leather", 1}, {"raven_feather", 1}},
		Message:   "Explorez le village. M : marchand, F : forge, T : arene.",
	}}
	d.slots[0], d.slots[9], d.slots[10], d.slots[11] = "health_potion", "wolf_fur", "boar_leather", "raven_feather"
	return d
}
func (d *Demo) Snapshot() Snapshot {
	s := d.state
	s.Inventory = append([]ItemStack(nil), s.Inventory...)
	s.Equipment = append([]string(nil), s.Equipment...)
	s.Skills = append([]string(nil), s.Skills...)
	s.Discovered = append([]string(nil), s.Discovered...)
	s.Slots = make([]ItemStack, len(d.slots))
	for i, id := range d.slots {
		if id != "" {
			s.Slots[i] = ItemStack{id, d.count(id)}
		}
	}
	s.Potions = d.count("health_potion")
	return s
}
func (d *Demo) count(id string) int {
	for _, v := range d.state.Inventory {
		if v.ID == id {
			return v.Count
		}
	}
	return 0
}
func (d *Demo) add(id string, n int) {
	for i := range d.state.Inventory {
		if d.state.Inventory[i].ID == id {
			d.state.Inventory[i].Count += n
			if d.state.Inventory[i].Count == 0 {
				d.state.Inventory = append(d.state.Inventory[:i], d.state.Inventory[i+1:]...)
				for j, slotID := range d.slots {
					if slotID == id {
						d.slots[j] = ""
					}
				}
			}
			return
		}
	}
	if n > 0 {
		d.state.Inventory = append(d.state.Inventory, ItemStack{id, n})
		for i := 9; i < len(d.slots); i++ {
			if d.slots[i] == "" {
				d.slots[i] = id
				return
			}
		}
		for i := 0; i < 9; i++ {
			if d.slots[i] == "" {
				d.slots[i] = id
				return
			}
		}
	}
}

func (d *Demo) MoveSlot(from, to int) error {
	if from < 0 || from >= len(d.slots) || to < 0 || to >= len(d.slots) {
		return fmt.Errorf("Emplacement invalide")
	}
	if d.slots[from] == "" {
		return fmt.Errorf("Cet emplacement est vide")
	}
	if from == to {
		return nil
	}
	d.slots[from], d.slots[to] = d.slots[to], d.slots[from]
	d.state.Message = "Objet deplace dans l'inventaire."
	return nil
}
func (d *Demo) buy(id string, cost int) error {
	if d.state.Fighting {
		return fmt.Errorf("Marchand indisponible pendant le combat")
	}
	if d.state.Gold < cost {
		return fmt.Errorf("Or insuffisant : %d necessaires", cost)
	}
	d.state.Gold -= cost
	d.add(id, 1)
	d.state.Message = fmt.Sprintf("Achat : %s (-%d or).", ItemName(id), cost)
	return nil
}
func (d *Demo) craft(id string, ingredients []ItemStack) error {
	if d.state.Fighting {
		return fmt.Errorf("Forge indisponible pendant le combat")
	}
	for _, item := range ingredients {
		if d.count(item.ID) < item.Count {
			return fmt.Errorf("Il manque %s x%d", ItemName(item.ID), item.Count-d.count(item.ID))
		}
	}
	if d.state.Gold < 5 {
		return fmt.Errorf("La forge demande 5 pieces d'or")
	}
	for _, item := range ingredients {
		d.add(item.ID, -item.Count)
	}
	d.state.Gold -= 5
	d.add(id, 1)
	d.state.Message = fmt.Sprintf("%s fabrique : equipez-le depuis l'inventaire.", ItemName(id))
	return nil
}
func (d *Demo) Execute(a Action) error {
	s := &d.state
	switch a {
	case Train:
		if s.Fighting {
			return fmt.Errorf("Combat deja en cours")
		}
		s.Fighting, s.EnemyHP, s.Turn = true, 40+s.Victories*5, 0
		s.Message = "Un gobelin entre dans l'arene. Attaquez ou utilisez une potion."
	case Return:
		if !s.Fighting {
			return fmt.Errorf("Aucun combat en cours")
		}
		s.Fighting = false
		s.Message = "Vous quittez l'arene."
	case Attack:
		if !s.Fighting {
			return fmt.Errorf("Entrez dans l'arene avec T")
		}
		s.EnemyHP -= 5 + equippedCount(s.Equipment)*2
		s.Message = "Vous frappez le gobelin."
	case Potion:
		if d.count("health_potion") == 0 {
			return fmt.Errorf("Aucune potion de soin")
		}
		if s.HP == s.MaxHP {
			return fmt.Errorf("Vie deja au maximum")
		}
		d.add("health_potion", -1)
		s.HP += 50
		if s.HP > s.MaxHP {
			s.HP = s.MaxHP
		}
		s.Message = "Vous buvez une potion de soin (+50 PV)."
	case Buy, BuyPotion:
		return d.buy("health_potion", 3)
	case BuyWolfFur:
		return d.buy("wolf_fur", 4)
	case BuyBoarLeather:
		return d.buy("boar_leather", 3)
	case BuyRavenFeather:
		return d.buy("raven_feather", 1)
	case BuyTrollHide:
		return d.buy("troll_hide", 7)
	case Craft, CraftHat:
		return d.craft("adventurer_hat", []ItemStack{{"raven_feather", 1}, {"boar_leather", 1}})
	case CraftTunic:
		return d.craft("adventurer_tunic", []ItemStack{{"wolf_fur", 2}, {"troll_hide", 1}})
	case CraftBoots:
		return d.craft("adventurer_boots", []ItemStack{{"wolf_fur", 1}, {"boar_leather", 1}})
	case ExploreWest, ExploreEast, ExploreSouth:
		if s.Fighting {
			return fmt.Errorf("Impossible d'explorer pendant un combat")
		}
		id := map[Action]string{ExploreWest: "wolf_camp", ExploreEast: "hill_cache", ExploreSouth: "old_shrine"}[a]
		for _, seen := range s.Discovered {
			if seen == id {
				return fmt.Errorf("Ce lieu a deja ete fouille")
			}
		}
		s.Discovered = append(s.Discovered, id)
		switch a {
		case ExploreWest:
			d.add("wolf_fur", 2)
			s.Message = "Camp des loups decouvert : 2 fourrures de loup."
		case ExploreEast:
			d.add("boar_leather", 2)
			s.Message = "Cache des collines decouverte : 2 cuirs de sanglier."
		case ExploreSouth:
			d.add("raven_feather", 2)
			s.Message = "Ancien sanctuaire decouvert : 2 plumes de corbeau."
		}
	default:
		return fmt.Errorf("Action inconnue : %s", a)
	}
	if s.Fighting && (a == Attack || a == Potion) {
		s.Turn++
		if s.EnemyHP <= 0 {
			s.EnemyHP, s.Fighting = 0, false
			s.Victories++
			s.Gold += 8
			s.Message = "Victoire ! Vous gagnez 8 pieces d'or."
			return nil
		}
		damage := 5
		if s.Turn%3 == 0 {
			damage = 10
		}
		s.HP -= damage
		if s.HP <= 0 {
			s.HP, s.Fighting = s.MaxHP/2, false
			s.Message = "Defaite : vous reprenez connaissance au village."
		} else {
			s.Message += fmt.Sprintf(" Gobelin : -%d PV. Tour %d.", damage, s.Turn)
		}
	}
	return nil
}

func ItemName(id string) string {
	switch id {
	case "health_potion":
		return "Potion de soin"
	case "wolf_fur":
		return "Fourrure de loup"
	case "boar_leather":
		return "Cuir de sanglier"
	case "troll_hide":
		return "Peau de troll"
	case "raven_feather":
		return "Plume de corbeau"
	case "adventurer_hat":
		return "Chapeau d'aventurier"
	case "adventurer_tunic":
		return "Tunique d'aventurier"
	case "adventurer_boots":
		return "Bottes d'aventurier"
	default:
		return id
	}
}
