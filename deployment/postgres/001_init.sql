CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    pin_hash TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_rooms_expires_at ON rooms (expires_at);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('text', 'link', 'system')),
    content TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_room_created ON messages (room_id, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_expires_at ON messages (expires_at);

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    original_name TEXT NOT NULL,
    stored_name TEXT NOT NULL,
    stored_path TEXT NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
    sha256 TEXT,
    status TEXT NOT NULL CHECK (status IN ('uploading', 'ready', 'expired', 'deleted')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_files_room_created ON files (room_id, created_at);
CREATE INDEX IF NOT EXISTS idx_files_expires_at ON files (expires_at);
CREATE INDEX IF NOT EXISTS idx_files_status ON files (status);

CREATE TABLE IF NOT EXISTS upload_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    file_id UUID REFERENCES files(id) ON DELETE CASCADE,
    total_size_bytes BIGINT NOT NULL CHECK (total_size_bytes >= 0),
    chunk_size_bytes INTEGER NOT NULL CHECK (chunk_size_bytes > 0),
    total_chunks INTEGER NOT NULL CHECK (total_chunks > 0),
    uploaded_chunks INTEGER NOT NULL DEFAULT 0 CHECK (uploaded_chunks >= 0),
    status TEXT NOT NULL CHECK (status IN ('pending', 'uploading', 'complete', 'aborted', 'expired')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_upload_sessions_room ON upload_sessions (room_id);
CREATE INDEX IF NOT EXISTS idx_upload_sessions_expires_at ON upload_sessions (expires_at);

CREATE TABLE IF NOT EXISTS file_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    upload_id UUID NOT NULL REFERENCES upload_sessions(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL CHECK (chunk_index >= 0),
    size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
    sha256 TEXT,
    stored_path TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (upload_id, chunk_index)
);
