package postgres

import (
	"context"
	"fmt"

	"github.com/hyanhsing/hyandrop/backend/internal/room"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, room room.Room, creator room.RoomDevice) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin room transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO rooms (
			id, expires_at, owner_delete_token_hash, used_bytes, reserved_bytes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, room.ID, room.ExpiresAt, room.OwnerDeleteTokenHash, room.UsedBytes, room.ReservedBytes, room.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert room: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO room_devices (
			room_id, device_id, nickname_ciphertext, joined_at, last_active_at
		) VALUES ($1, $2, $3, $4, $5)
	`, creator.RoomID, creator.DeviceID, creator.NicknameCiphertext, creator.JoinedAt, creator.LastActiveAt)
	if err != nil {
		return fmt.Errorf("insert room device: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit room transaction: %w", err)
	}

	return nil
}
