package settings

type Tricks map[string]bool

func DefaultTricks() Tricks {
	return Tricks{
		"visible_collisions":          true,
		"grottos_without_agony":       true,
		"fewer_tunic_requirements":    true,
		"rusted_switches":             true,
		"man_on_roof":                 true,
		"windmill_poh":                true,
		"crater_bean_poh_with_hovers": true,
		"dc_jump":                     true,
		"lens_botw":                   true,
		"child_deadhand":              true,
		"forest_vines":                true,
		"lens_shadow":                 true,
		"lens_shadow_platform":        true,
		"lens_bongo":                  true,
		"lens_spirit":                 true,
		"lens_gtg":                    true,
		"lens_castle":                 true,
	}
}
