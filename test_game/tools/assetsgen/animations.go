package main

import "encoding/json"

type animationClip struct {
	Name      string         `json:"name"`
	Duration  float64        `json:"duration_seconds"`
	Loop      bool           `json:"loop"`
	Keyframes []animationKey `json:"keyframes"`
}

type animationKey struct {
	Time     float64    `json:"time"`
	Position [3]float64 `json:"position,omitempty"`
	Rotation [3]float64 `json:"rotation_degrees,omitempty"`
	Scale    [3]float64 `json:"scale,omitempty"`
}

func generateAnimations(c *catalog) error {
	clips := []animationClip{
		{"idle", 1.2, true, []animationKey{{0, [3]float64{0, 0, 0}, [3]float64{}, [3]float64{1, 1, 1}}, {0.6, [3]float64{0, 0.04, 0}, [3]float64{}, [3]float64{1, 1.02, 1}}, {1.2, [3]float64{0, 0, 0}, [3]float64{}, [3]float64{1, 1, 1}}}},
		{"walk", 0.7, true, []animationKey{{0, [3]float64{0, 0, 0}, [3]float64{0, -8, 0}, [3]float64{1, 1, 1}}, {0.35, [3]float64{0, 0.07, 0}, [3]float64{0, 8, 0}, [3]float64{1, 1, 1}}, {0.7, [3]float64{0, 0, 0}, [3]float64{0, -8, 0}, [3]float64{1, 1, 1}}}},
		{"attack", 0.42, false, []animationKey{{0, [3]float64{}, [3]float64{0, -18, 0}, [3]float64{1, 1, 1}}, {0.18, [3]float64{0, 0, -0.18}, [3]float64{0, 34, 0}, [3]float64{1.05, 0.95, 1.05}}, {0.42, [3]float64{}, [3]float64{}, [3]float64{1, 1, 1}}}},
		{"hurt", 0.32, false, []animationKey{{0, [3]float64{}, [3]float64{}, [3]float64{1, 1, 1}}, {0.12, [3]float64{0, 0, 0.25}, [3]float64{-12, 0, 0}, [3]float64{0.94, 1.04, 0.94}}, {0.32, [3]float64{}, [3]float64{}, [3]float64{1, 1, 1}}}},
		{"death", 0.9, false, []animationKey{{0, [3]float64{}, [3]float64{}, [3]float64{1, 1, 1}}, {0.5, [3]float64{0, -0.3, 0}, [3]float64{0, 0, 65}, [3]float64{1, 1, 1}}, {0.9, [3]float64{0, -0.72, 0}, [3]float64{0, 0, 90}, [3]float64{1, 0.92, 1}}}},
	}
	for _, actor := range []string{"player", "goblin"} {
		data, err := json.MarshalIndent(struct {
			Actor string          `json:"actor"`
			Clips []animationClip `json:"clips"`
		}{Actor: actor, Clips: clips}, "", "  ")
		if err != nil {
			return err
		}
		if err := c.write("animations/"+actor+".json", "animation-spec", "Keyframes simples pour "+actor, append(data, '\n')); err != nil {
			return err
		}
	}
	return nil
}
