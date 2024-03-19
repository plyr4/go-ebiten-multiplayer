package images

import (
	_ "embed"
)

var (
	//go:embed gopher.png
	Gopher []byte
	//go:embed runner.png
	Runner []byte
)
