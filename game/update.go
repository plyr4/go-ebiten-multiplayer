package game

func (g *Game) Update() error {
	if g.error != nil {
		return g.error
	}
	g.Frame++

	g.Input.Update()

	for _, e := range g.entities {
		e.Update()
	}

	return g.error
}
