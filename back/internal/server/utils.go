package server

import (
	"fmt"
	"github.com/aakosarev/kanban-board/back/internal/config"
	"strings"
	"time"
)

const (
	waitShotDownDuration = 3 * time.Second
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}

func (s *Server) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
