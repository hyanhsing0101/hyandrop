import { openDB, type DBSchema } from 'idb';

export type SavedRoom = {
  roomId: string;
  roomKey: string;
  deleteToken: string;
  expiresAt: string;
  createdAt: string;
};

type MetaRecord = {
  key: string;
  value: string;
};

interface HyanDropDB extends DBSchema {
  meta: {
    key: string;
    value: MetaRecord;
  };
  rooms: {
    key: string;
    value: SavedRoom;
    indexes: {
      'by-expires-at': string;
    };
  };
}

const databasePromise = openDB<HyanDropDB>('hyandrop', 1, {
  upgrade(database) {
    database.createObjectStore('meta', {
      keyPath: 'key'
    });

    const rooms = database.createObjectStore('rooms', {
      keyPath: 'roomId'
    });

    rooms.createIndex('by-expires-at', 'expiresAt');
  }
});

export async function getOrCreateDeviceID(): Promise<string> {
  const database = await databasePromise;
  const existing = await database.get('meta', 'device-id');

  if (existing) {
    return existing.value;
  }

  const deviceID = crypto.randomUUID();

  await database.put('meta', {
    key: 'device-id',
    value: deviceID
  });

  return deviceID;
}

export function createRoomKey(): string {
  const bytes = crypto.getRandomValues(new Uint8Array(32));
  const binary = String.fromCharCode(...bytes);

  return window
    .btoa(binary)
    .replaceAll('+', '-')
    .replaceAll('/', '_')
    .replace(/=+$/, '');
}

export async function saveRoom(room: SavedRoom): Promise<void> {
  const database = await databasePromise;
  await database.put('rooms', room);
}

export async function loadRooms(): Promise<SavedRoom[]> {
  const database = await databasePromise;
  const rooms = await database.getAllFromIndex('rooms', 'by-expires-at');

  return rooms.reverse();
}