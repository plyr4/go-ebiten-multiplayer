package player

import (
	"bytes"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/plyr4/go-ebiten-multiplayer/entity"
	"github.com/plyr4/go-ebiten-multiplayer/input"
	"github.com/plyr4/go-ebiten-multiplayer/resources/images"
)

type Player struct {
	entity.Entity
	speed float64
	img   *ebiten.Image
	// debug random hue for each player
	Hue float64
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

	p.Hue = 2.0 * math.Pi * (rand.Float64() * 50)

	return p, nil
}

func (p *Player) Draw(screen *ebiten.Image) error {
	opts := &colorm.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)

	cm := colorm.ColorM{}
	cm.Reset()
	cm.RotateHue(p.Hue)
	colorm.DrawImage(screen, p.img, cm, opts)

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
