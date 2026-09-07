import { useEffect, useState, type ReactNode } from 'react';
import { Cloud, Database, Radio, Server } from 'lucide-react';

type ApiVersion = {
  name: string;
  environment: string;
  room_ttl_hours: number;
  max_file_size_mb: number;
};

export function App() {
  const [version, setVersion] = useState<ApiVersion | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch('/api/version')
      .then((response) => {
        if (!response.ok) {
          throw new Error(`API returned ${response.status}`);
        }
        return response.json() as Promise<ApiVersion>;
      })
      .then(setVersion)
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Unable to reach API');
      });
  }, []);

  return (
    <main className="app-shell">
      <section className="intro">
        <div className="brand-mark">
          <Cloud size={30} strokeWidth={2.2} />
        </div>
        <div>
          <h1>HyanDrop</h1>
          <p>Temporary browser rooms for moving text and files across devices.</p>
        </div>
      </section>

      <section className="status-panel" aria-label="Development stack status">
        <StatusRow icon={<Server size={20} />} label="Backend" value={version ? 'Connected' : 'Waiting'} />
        <StatusRow icon={<Database size={20} />} label="PostgreSQL" value={version ? 'Configured' : 'Pending'} />
        <StatusRow icon={<Radio size={20} />} label="Redis" value={version ? 'Configured' : 'Pending'} />
      </section>

      <section className="runtime">
        {version ? (
          <>
            <p>
              API is online in <strong>{version.environment}</strong> mode.
            </p>
            <dl>
              <div>
                <dt>Room TTL</dt>
                <dd>{version.room_ttl_hours}h after the latest transfer</dd>
              </div>
              <div>
                <dt>MVP file limit</dt>
                <dd>{version.max_file_size_mb} MB</dd>
              </div>
            </dl>
          </>
        ) : (
          <p>{error ? `API check failed: ${error}` : 'Checking API connection...'}</p>
        )}
      </section>
    </main>
  );
}

function StatusRow({ icon, label, value }: { icon: ReactNode; label: string; value: string }) {
  return (
    <div className="status-row">
      <span className="status-icon">{icon}</span>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}
