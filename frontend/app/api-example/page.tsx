import { auth } from '@clerk/nextjs/server';
import { authenticatedFetch, getAuthToken } from '@/lib/api-client';

export default async function ApiExamplePage() {
  const { userId } = await auth();
  const token = await getAuthToken();

  // Example: Fetch data from your backend API
  let laboratories = null;
  let error = null;

  if (token) {
    try {
      const response = await authenticatedFetch('/api/v1/laboratories');
      if (response.ok) {
        laboratories = await response.json();
      } else {
        error = `API Error: ${response.status} ${response.statusText}`;
      }
    } catch (err) {
      error = `Failed to fetch: ${
        err instanceof Error ? err.message : 'Unknown error'
      }`;
    }
  }

  return (
    <div className="space-y-6 p-8">
      <h1 className="text-2xl font-bold">API Example Page</h1>

      <div className="bg-gray-100 p-4 rounded">
        <h2 className="font-semibold mb-2">User Info</h2>
        <p>User ID: {userId || 'Not authenticated'}</p>
        <p className="mt-2">
          Token:{' '}
          {token ? (
            <span className="font-mono text-xs break-all">
              {token.substring(0, 50)}...
            </span>
          ) : (
            'No token available'
          )}
        </p>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          <p className="font-semibold">Error:</p>
          <p>{error}</p>
        </div>
      )}

      {laboratories && (
        <div className="bg-green-100 p-4 rounded">
          <h2 className="font-semibold mb-2">Laboratories (from API):</h2>
          <pre className="text-xs overflow-auto">
            {JSON.stringify(laboratories, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
