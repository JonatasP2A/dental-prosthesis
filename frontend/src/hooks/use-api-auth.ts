import { useEffect } from 'react'
import { useAuth } from '@clerk/clerk-react'
import { setAuthToken } from '@/lib/api-client'

/**
 * Hook to synchronize Clerk authentication token with the API client.
 * Should be used in a component that wraps authenticated routes.
 */
export function useApiAuth() {
  const { getToken, isSignedIn } = useAuth()

  useEffect(() => {
    async function syncToken() {
      if (isSignedIn) {
        const token = await getToken()
        setAuthToken(token)
      } else {
        setAuthToken(null)
      }
    }

    syncToken()
  }, [getToken, isSignedIn])
}
