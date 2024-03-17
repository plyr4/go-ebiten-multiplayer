package player

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plyr4/go-ebiten-multiplayer/entity"
	"github.com/plyr4/go-ebiten-multiplayer/input"
	"github.com/plyr4/go-ebiten-multiplayer/resources/images"
)

type Player struct {
	entity.Entity
	speed float64
	img   *ebiten.Image
}

func New() (*Player, error) {
	p := new(Player)
	p.X = 0
	p.Y = 0
	p.speed = 2

	img, _, err := image.Decode(bytes.NewReader(images.Gopher))
	if err != nil {
		return nil, err
	}

	p.img = ebiten.NewImageFromImage(img)

	return p, nil
}

func (p *Player) Draw(screen *ebiten.Image) error {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(p.img, opts)

	return nil
}

func (p *Player) Update(i *input.Input) error {
	if i.Up {
		p.Y -= p.speed
	}

	if i.Down {
		p.Y += p.speed
	}

	if i.Left {
		p.X -= p.speed
	}

	if i.Right {
		p.X += p.speed
	}

	return nil
}
