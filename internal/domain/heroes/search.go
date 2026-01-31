package heroes

// BestForResource returns a list of available heroes with a buff for the specified resource.
func (h Heroes) BestForResource(resource string) []Hero {
	var result []Hero
	key := resource + "_gathering_speed"

	for _, hero := range h.List {
		if !hero.State.IsAvailable {
			continue
		}
		if _, ok := hero.Buffs[key]; ok {
			result = append(result, hero)
		}
	}
	return result
}

// BestForDefense returns available heroes with a defense role.
func (h Heroes) BestForDefense() []Hero {
	var result []Hero
	for _, hero := range h.List {
		if !hero.State.IsAvailable {
			continue
		}
		for _, role := range hero.Roles {
			if role == "garrison_defense" || role == "defense" {
				result = append(result, hero)
				break
			}
		}
	}
	return result
}

// BestForAttack returns available heroes with a combat role.
func (h Heroes) BestForAttack() []Hero {
	var result []Hero

	for _, hero := range h.List {
		if !hero.State.IsAvailable {
			continue
		}
		for _, role := range hero.Roles {
			if role == "rally_leader" || role == "combat" {
				result = append(result, hero)
				break
			}
		}
	}

	return result
}

// Available returns a Heroes structure with only available heroes.
func (h Heroes) Available() Heroes {
	out := Heroes{
		IsNotify: h.IsNotify,
		List:     make(map[string]Hero),
	}

	for name, hero := range h.List {
		if hero.State.IsAvailable {
			out.List[name] = hero
		}
	}

	return out
}
