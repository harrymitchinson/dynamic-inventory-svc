package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
	"go.uber.org/zap"
)

const (
	// KeyFormatHostsFull is the key format for a full inventory.Host entry.
	KeyFormatHostsFull string = "hosts:%s:%s"
	// KeyFormatHostsWildcard is the key format for all inventory.Host entries for an environment.
	KeyFormatHostsWildcard string = "hosts:%s:*"
)

// HostService represents an implementation of inventory.HostService
type HostService struct {
	r *redis.Client
	l *zap.Logger
}

// NewHostService initialises a HostService.
func NewHostService(r *redis.Client, l *zap.Logger) HostService {
	return HostService{
		r: r,
		l: l,
	}
}

// SetHost sets a key in the redis database for the host, overwriting any previously existing key.
func (s HostService) SetHost(e inventory.Environment, h *inventory.Host) error {
	log := s.l.With(
		zap.String("environment", string(e)),
		zap.String("name", h.Name)).Named("SetHost")
	ctx := s.r.Context()

	key := fmt.Sprintf(KeyFormatHostsFull, e, h.Name)

	log.Info("setting key", zap.String("key", key))
	if err := s.r.Set(ctx, key, h, 0).Err(); err != nil {
		log.Error("failed to set key", zap.Error(err))
		return err
	}
	log.Info("key set", zap.String("key", key))
	return nil
}

// GetHosts gets all hosts for an environment.
func (s HostService) GetHosts(e inventory.Environment) ([]inventory.Host, error) {
	log := s.l.With(zap.String("environment", string(e))).Named("GetHosts")
	ctx := s.r.Context()
	hosts := make([]inventory.Host, 0)

	pattern := fmt.Sprintf(KeyFormatHostsWildcard, e)
	log = log.With(zap.String("pattern", pattern))
	log.Info("scanning for pattern")

	iter := s.r.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		iterLog := log.With(zap.String("key", key))
		iterLog.Info("getting key")

		cacheData, err := s.r.Get(ctx, key).Result()
		if err != nil {
			iterLog.Error("failed to read key", zap.Error(err))
			return hosts, err
		}

		iterLog.Info("unmarshalling value")
		var h inventory.Host
		if err := h.UnmarshalBinary([]byte(cacheData)); err != nil {
			iterLog.Error("Failed to unmarshal", zap.Error(err))
			return hosts, err
		}

		iterLog.Info("appending to hosts slice")
		hosts = append(hosts, h)
	}

	log.Info("loaded hosts successfully", zap.Int("results", len(hosts)))
	return hosts, nil
}
