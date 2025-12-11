import { createFileRoute } from '@tanstack/react-router'
import { Technicians } from '@/features/technicians'

export const Route = createFileRoute('/_authenticated/technicians/')({
  component: Technicians,
})
