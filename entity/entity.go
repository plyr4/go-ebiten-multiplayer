package entity

type Entity struct {
	X, Y float64
}

func (e *Entity) Position() (float64, float64) {
	return e.X, e.Y
}
