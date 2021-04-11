package raspistill

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"time"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type Raspistill struct {
	timeout time.Duration
	quality uint8
}

func New(cfg *config.Raspistill) *Raspistill {
	return &Raspistill{
		timeout: time.Duration(cfg.Timeout),
		quality: cfg.Quality,
	}
}

func (r *Raspistill) GetPhoto(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	strQuality := strconv.FormatInt(int64(r.quality), 10)
	cmd := exec.CommandContext(ctx, "raspistill", "-q", strQuality, "-o", "-")

	out, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errors.New("timed out")
		}
		return nil, err
	}
	return out, nil
}
