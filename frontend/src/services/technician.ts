import { api } from '@/lib/api-client'
import type {
  Technician,
  CreateTechnicianRequest,
  UpdateTechnicianRequest,
} from '@/types'

export const technicianService = {
  list: (laboratoryId: string) =>
    api.get<Technician[]>(`/laboratories/${laboratoryId}/technicians`),

  get: (laboratoryId: string, id: string) =>
    api.get<Technician>(`/laboratories/${laboratoryId}/technicians/${id}`),

  create: (laboratoryId: string, data: CreateTechnicianRequest) =>
    api.post<Technician>(`/laboratories/${laboratoryId}/technicians`, data),

  update: (laboratoryId: string, id: string, data: UpdateTechnicianRequest) =>
    api.put<Technician>(
      `/laboratories/${laboratoryId}/technicians/${id}`,
      data
    ),

  delete: (laboratoryId: string, id: string) =>
    api.delete<void>(`/laboratories/${laboratoryId}/technicians/${id}`),
}
