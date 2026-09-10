package room

import "context"

type Repository interface {
	Create(ctx context.Context, room Room, creator RoomDevice) error
}
