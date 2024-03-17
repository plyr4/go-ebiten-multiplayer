package game

func (g *Game) Update() error {
	if g.error != nil {
		return g.error
	}

	g.Input.Update()

	for _, e := range g.entities {
		e.Update(g.Input)
	}

	g.Debug.Frame++

	return g.error
}
