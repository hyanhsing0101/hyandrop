import { useState } from 'react';
import { Check, Cloud, Copy, Link as LinkIcon, Plus } from 'lucide-react';
import { createRoomKey, getOrCreateDeviceID, saveRoom } from './lib/room-storage';

type CreateRoomResponse = {
  room_id: string;
  expires_at: string;
  delete_token: string;
};

export function App() {
  const [roomURL, setRoomURL] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [isCopied, setIsCopied] = useState(false);

  async function createRoom() {
    setError(null);
    setIsCopied(false);
    setIsCreating(true);

    try {
      const deviceID = await getOrCreateDeviceID();
      const roomKey = createRoomKey();
      const response = await fetch('/api/rooms', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ device_id: deviceID })
      });
      const payload = (await response.json()) as CreateRoomResponse | { error: string };

      if (!response.ok || !('room_id' in payload)) {
        throw new Error('error' in payload ? payload.error : `API returned ${response.status}`);
      }

      const url = `${window.location.origin}/r/${payload.room_id}#${roomKey}`;
      await saveRoom({
        roomId: payload.room_id,
        roomKey,
        deleteToken: payload.delete_token,
        expiresAt: payload.expires_at,
        createdAt: new Date().toISOString()
      });
      setRoomURL(url);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Unable to create a room');
    } finally {
      setIsCreating(false);
    }
  }

  async function copyRoomURL() {
    if (!roomURL) {
      return;
    }

    await navigator.clipboard.writeText(roomURL);
    setIsCopied(true);
  }

  return (
    <main className="app-shell">
      <section className="intro">
        <div className="brand-mark">
          <Cloud size={30} strokeWidth={2.2} />
        </div>
        <div>
          <h1>HyanDrop</h1>
          <p>Private, temporary rooms for moving text and files across your devices.</p>
        </div>
      </section>

      <section className="create-card" aria-label="Create a temporary room">
        <span className="eyebrow">No account · End-to-end encryption ready</span>
        <h2>Start a private room</h2>
        <p>Share one link across devices. Rooms expire 24 hours after the latest successful transfer.</p>
        <button className="primary-button" type="button" onClick={createRoom} disabled={isCreating}>
          <Plus size={19} aria-hidden="true" />
          {isCreating ? 'Creating room…' : 'Create temporary room'}
        </button>
        {error && <p className="error-message">{error}</p>}
      </section>

      {roomURL && (
        <section className="share-card" aria-live="polite">
          <div className="share-heading">
            <span className="success-icon"><Check size={18} /></span>
            <div>
              <h2>Room created</h2>
              <p>Keep this link safe. Anyone with it can join this room.</p>
            </div>
          </div>
          <div className="share-link">
            <LinkIcon size={18} aria-hidden="true" />
            <code>{roomURL}</code>
            <button type="button" onClick={copyRoomURL}>
              {isCopied ? <Check size={18} aria-label="Copied" /> : <Copy size={18} aria-label="Copy link" />}
            </button>
          </div>
        </section>
      )}

      <section className="rules-card">
        <div><strong>100 MB</strong><span>per file</span></div>
        <div><strong>200 MB</strong><span>per room</span></div>
        <div><strong>24 hours</strong><span>after transfer</span></div>
      </section>
    </main>
  );
}
