BEGIN;

ALTER TABLE messages DROP CONSTRAINT messages_room_id_fkey;
ALTER TABLE files DROP CONSTRAINT files_room_id_fkey;
ALTER TABLE upload_sessions DROP CONSTRAINT upload_sessions_room_id_fkey;

ALTER TABLE rooms ALTER COLUMN id DROP DEFAULT;
ALTER TABLE rooms ALTER COLUMN id TYPE TEXT USING replace(id::text, '-', '');
ALTER TABLE messages ALTER COLUMN room_id TYPE TEXT USING replace(room_id::text, '-', '');
ALTER TABLE files ALTER COLUMN room_id TYPE TEXT USING replace(room_id::text, '-', '');
ALTER TABLE upload_sessions ALTER COLUMN room_id TYPE TEXT USING replace(room_id::text, '-', '');

ALTER TABLE rooms DROP CONSTRAINT rooms_code_key;
ALTER TABLE rooms DROP COLUMN code;
ALTER TABLE rooms DROP COLUMN pin_hash;
ALTER TABLE rooms DROP COLUMN last_active_at;
ALTER TABLE rooms ADD COLUMN owner_delete_token_hash BYTEA;
ALTER TABLE rooms ADD COLUMN used_bytes BIGINT NOT NULL DEFAULT 0;
ALTER TABLE rooms ADD COLUMN reserved_bytes BIGINT NOT NULL DEFAULT 0;

UPDATE rooms
SET owner_delete_token_hash = gen_random_bytes(32)
WHERE owner_delete_token_hash IS NULL;

ALTER TABLE rooms ALTER COLUMN owner_delete_token_hash SET NOT NULL;
ALTER TABLE rooms ADD CONSTRAINT rooms_id_length CHECK (char_length(id) = 32);
ALTER TABLE rooms ADD CONSTRAINT rooms_delete_token_hash_length
    CHECK (octet_length(owner_delete_token_hash) = 32);
ALTER TABLE rooms ADD CONSTRAINT rooms_used_bytes_nonnegative CHECK (used_bytes >= 0);
ALTER TABLE rooms ADD CONSTRAINT rooms_reserved_bytes_nonnegative CHECK (reserved_bytes >= 0);

ALTER TABLE messages ADD CONSTRAINT messages_room_id_fkey
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;
ALTER TABLE files ADD CONSTRAINT files_room_id_fkey
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;
ALTER TABLE upload_sessions ADD CONSTRAINT upload_sessions_room_id_fkey
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE;

CREATE TABLE room_devices (
    room_id TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    device_id TEXT NOT NULL CHECK (char_length(device_id) >= 16),
    nickname_ciphertext BYTEA,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (room_id, device_id)
);

CREATE INDEX idx_room_devices_device_id ON room_devices (device_id);

COMMIT;
