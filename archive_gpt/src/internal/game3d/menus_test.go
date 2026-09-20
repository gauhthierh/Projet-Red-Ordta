package game3d

import "testing"

func TestMenuAssetPaths(t *testing.T) {
	for _, id := range []string{"main_menu", "health_potion", "gold_coin", "adventurer_hat", "training_goblin"} {
		if iconPath(id) == "" {
			t.Fatalf("asset introuvable : %s", id)
		}
	}
}
