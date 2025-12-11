import { api } from '@/lib/api-client'
import type {
  Laboratory,
  CreateLaboratoryRequest,
  UpdateLaboratoryRequest,
} from '@/types'

const BASE_URL = '/laboratories'

export const laboratoryService = {
  list: () => api.get<Laboratory[]>(BASE_URL),

  get: (id: string) => api.get<Laboratory>(`${BASE_URL}/${id}`),

  create: (data: CreateLaboratoryRequest) =>
    api.post<Laboratory>(BASE_URL, data),

  update: (id: string, data: UpdateLaboratoryRequest) =>
    api.put<Laboratory>(`${BASE_URL}/${id}`, data),

  delete: (id: string) => api.delete<void>(`${BASE_URL}/${id}`),
}
