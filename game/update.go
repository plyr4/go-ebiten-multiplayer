package game

func (g *Game) Update() error {
	if g.error != nil {
		return g.error
	}
	g.Frame++

	g.error = g.Input.Update()
	if g.error != nil {
		return g.error
	}

	for _, e := range g.entities {
		g.error = e.Update()
		if g.error != nil {
			return g.error
		}
	}

	return g.error
}
