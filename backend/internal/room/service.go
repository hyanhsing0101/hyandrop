package room

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"
)

const minDeviceIDLength = 16

var (
	ErrInvalidDeviceID = errors.New("device ID must contain at least 16 characters")
	ErrInvalidRoomTTL  = errors.New("room TTL must be positive")
)

type CreateInput struct {
	DeviceID string
}

type CreateResult struct {
	RoomID      string
	ExpiresAt   time.Time
	DeleteToken string
}

type Service struct {
	repo    Repository
	roomTTL time.Duration
	clock   func() time.Time
}

func NewService(repo Repository, roomTTL time.Duration, clock func() time.Time) *Service {
	return &Service{
		repo:    repo,
		roomTTL: roomTTL,
		clock:   clock,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (CreateResult, error) {
	deviceID := strings.TrimSpace(input.DeviceID)
	if len(deviceID) < minDeviceIDLength {
		return CreateResult{}, ErrInvalidDeviceID
	}

	if s.roomTTL <= 0 {
		return CreateResult{}, ErrInvalidRoomTTL
	}

	roomID, err := randomHex(16)
	if err != nil {
		return CreateResult{}, fmt.Errorf("generate room ID: %w", err)
	}

	deleteToken, err := randomHex(32)
	if err != nil {
		return CreateResult{}, fmt.Errorf("generate delete token: %w", err)
	}

	now := s.clock().UTC()
	tokenHash := sha256.Sum256([]byte(deleteToken))
	room := Room{
		ID:                   roomID,
		ExpiresAt:            now.Add(s.roomTTL),
		OwnerDeleteTokenHash: tokenHash[:],
		UsedBytes:            0,
		ReservedBytes:        0,
		CreatedAt:            now,
	}
	creator := RoomDevice{
		RoomID:       roomID,
		DeviceID:     deviceID,
		JoinedAt:     now,
		LastActiveAt: now,
	}

	if err := s.repo.Create(ctx, room, creator); err != nil {
		return CreateResult{}, fmt.Errorf("store room: %w", err)
	}

	return CreateResult{
		RoomID:      room.ID,
		ExpiresAt:   room.ExpiresAt,
		DeleteToken: deleteToken,
	}, nil
}
