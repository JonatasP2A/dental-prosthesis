import { createFileRoute } from '@tanstack/react-router'
import { Laboratories } from '@/features/laboratories'

export const Route = createFileRoute('/_authenticated/laboratories/')({
  component: Laboratories,
})
