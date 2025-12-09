'use client';

import { useAuth } from '@clerk/nextjs';
import { useState } from 'react';

/**
 * Client Component example showing how to access JWT token and make API calls
 */
export default function ClientApiExample() {
  const { getToken, userId, isLoaded } = useAuth();
  const [token, setToken] = useState<string | null>(null);
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGetToken = async () => {
    const jwtToken = await getToken();
    setToken(jwtToken);
  };

  const handleFetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      const jwtToken = await getToken();
      if (!jwtToken) {
        setError('No token available. Please sign in.');
        return;
      }

      const API_URL =
        process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${API_URL}/api/v1/laboratories`, {
        headers: {
          Authorization: `Bearer ${jwtToken}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }

      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  if (!isLoaded) {
    return <div>Loading...</div>;
  }

  if (!userId) {
    return (
      <div className="text-gray-600">Please sign in to use this example.</div>
    );
  }

  return (
    <div className="space-y-4 p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-xl font-semibold">Client Component API Example</h2>

      <div className="space-y-2">
        <button
          onClick={handleGetToken}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Get JWT Token
        </button>

        {token && (
          <div className="bg-gray-100 p-3 rounded text-xs break-all">
            <strong>Token:</strong> {token.substring(0, 50)}...
          </div>
        )}
      </div>

      <div className="space-y-2">
        <button
          onClick={handleFetchData}
          disabled={loading}
          className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:bg-gray-400"
        >
          {loading ? 'Loading...' : 'Fetch Laboratories from API'}
        </button>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded">
            {error}
          </div>
        )}

        {data && (
          <div className="bg-green-50 p-4 rounded">
            <pre className="text-xs overflow-auto">
              {JSON.stringify(data, null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
