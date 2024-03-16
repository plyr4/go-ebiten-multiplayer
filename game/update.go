package game

func (g *Game) Update() error {
	if g.error != nil {
		return g.error
	}

	g.Debug.Frame++

	return g.error
}
