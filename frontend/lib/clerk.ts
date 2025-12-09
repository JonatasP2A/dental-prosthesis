import { auth } from '@clerk/nextjs/server';

/**
 * Get the Clerk JWT session token for API calls
 * This function can be used in Server Components or API routes
 */
export async function getClerkToken(): Promise<string | null> {
  const { getToken } = await auth();
  return await getToken();
}
