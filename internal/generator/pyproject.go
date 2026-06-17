package generator

import "github.com/shv-ng/fastx/internal/config"

type generator struct {
	cfg *config.Config
}

func New(cfg *config.Config) *generator {
	return &generator{cfg: cfg}
}
