import { auth } from '@clerk/nextjs/server';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * Create an authenticated fetch function that automatically includes the Clerk JWT token
 * Use this in Server Components or Server Actions
 */
export async function authenticatedFetch(
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> {
  const { getToken } = await auth();
  const token = await getToken();

  const headers = new Headers(options.headers);
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }
  headers.set('Content-Type', 'application/json');

  return fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });
}

/**
 * Get the Clerk JWT token as a string
 * Useful for debugging or custom API calls
 */
export async function getAuthToken(): Promise<string | null> {
  const { getToken } = await auth();
  return await getToken();
}
