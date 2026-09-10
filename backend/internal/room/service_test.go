package room

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"
	"time"
)

type recordingRepository struct {
	room    Room
	creator RoomDevice
	err     error
	called  bool
}

func (r *recordingRepository) Create(_ context.Context, room Room, creator RoomDevice) error {
	r.called = true
	r.room = room
	r.creator = creator
	return r.err
}

func TestServiceCreateStoresRoomAndCreator(t *testing.T) {
	fixedNow := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time {
		return fixedNow
	}
	repo := &recordingRepository{}
	service := NewService(repo, 24*time.Hour, clock)

	result, err := service.Create(context.Background(), CreateInput{DeviceID: "device-1234567890"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if !repo.called {
		t.Fatal("Create() did not call the repository")
	}
	if result.RoomID == "" || len(result.RoomID) != 32 {
		t.Fatalf("RoomID = %q, want a 32-character random hex ID", result.RoomID)
	}
	if len(result.DeleteToken) != 64 {
		t.Fatalf("DeleteToken length = %d, want 64", len(result.DeleteToken))
	}
	if repo.creator.DeviceID != "device-1234567890" {
		t.Fatalf("creator DeviceID = %q", repo.creator.DeviceID)
	}
	if repo.creator.RoomID != result.RoomID {
		t.Fatalf("creator RoomID = %q, want %q", repo.creator.RoomID, result.RoomID)
	}

	wantHash := sha256.Sum256([]byte(result.DeleteToken))
	if string(repo.room.OwnerDeleteTokenHash) != string(wantHash[:]) {
		t.Fatal("stored token hash does not match the returned delete token")
	}
	wantExpiresAt := fixedNow.Add(24 * time.Hour)
	if !result.ExpiresAt.Equal(wantExpiresAt) {
		t.Fatalf("ExpiresAt = %s, want %s", result.ExpiresAt, wantExpiresAt)
	}

	if !repo.room.CreatedAt.Equal(fixedNow) {
		t.Fatalf("CreatedAt = %s, want %s", repo.room.CreatedAt, fixedNow)
	}
}

func TestServiceCreateRejectsShortDeviceID(t *testing.T) {

	fixedNow := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time {
		return fixedNow
	}

	repo := &recordingRepository{}
	service := NewService(repo, 24*time.Hour, clock)

	_, err := service.Create(context.Background(), CreateInput{DeviceID: "too-short"})
	if !errors.Is(err, ErrInvalidDeviceID) {
		t.Fatalf("Create() error = %v, want ErrInvalidDeviceID", err)
	}
	if repo.called {
		t.Fatal("Create() called the repository for an invalid device ID")
	}
}
