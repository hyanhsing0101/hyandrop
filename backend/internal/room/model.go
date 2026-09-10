package room

import "time"

type Room struct {
	ID                   string
	ExpiresAt            time.Time
	OwnerDeleteTokenHash []byte
	UsedBytes            int64
	ReservedBytes        int64
	CreatedAt            time.Time
}

type RoomDevice struct {
	RoomID             string
	DeviceID           string
	NicknameCiphertext []byte
	JoinedAt           time.Time
	LastActiveAt       time.Time
}
