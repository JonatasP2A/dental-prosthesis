import { useEffect } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useAuth } from '@clerk/clerk-react'
import { useApiAuth } from '@/hooks/use-api-auth'
import { AuthenticatedLayout } from '@/components/layout/authenticated-layout'

function AuthGuard() {
  const { isLoaded, isSignedIn } = useAuth()
  const navigate = useNavigate()

  // Sync Clerk token with API client
  useApiAuth()

  useEffect(() => {
    if (isLoaded && !isSignedIn) {
      navigate({ to: '/sign-in' })
    }
  }, [isLoaded, isSignedIn, navigate])

  if (!isLoaded) {
    return (
      <div className='flex h-screen items-center justify-center'>
        <div className='animate-pulse'>Loading...</div>
      </div>
    )
  }

  if (!isSignedIn) {
    return null
  }

  return <AuthenticatedLayout />
}

export const Route = createFileRoute('/_authenticated')({
  component: AuthGuard,
})
