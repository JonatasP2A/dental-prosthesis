import { api } from '@/lib/api-client'
import type { Client, CreateClientRequest, UpdateClientRequest } from '@/types'

export const clientService = {
  list: (laboratoryId: string) =>
    api.get<Client[]>(`/laboratories/${laboratoryId}/clients`),

  get: (laboratoryId: string, id: string) =>
    api.get<Client>(`/laboratories/${laboratoryId}/clients/${id}`),

  create: (laboratoryId: string, data: CreateClientRequest) =>
    api.post<Client>(`/laboratories/${laboratoryId}/clients`, data),

  update: (laboratoryId: string, id: string, data: UpdateClientRequest) =>
    api.put<Client>(`/laboratories/${laboratoryId}/clients/${id}`, data),

  delete: (laboratoryId: string, id: string) =>
    api.delete<void>(`/laboratories/${laboratoryId}/clients/${id}`),
}
