import { createFileRoute } from '@tanstack/react-router'
import { Prostheses } from '@/features/prostheses'

export const Route = createFileRoute('/_authenticated/prostheses/')({
  component: Prostheses,
})
