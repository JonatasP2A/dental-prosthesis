import { api } from '@/lib/api-client'
import type {
  Prosthesis,
  CreateProsthesisRequest,
  UpdateProsthesisRequest,
} from '@/types'

export const prosthesisService = {
  list: (laboratoryId: string) =>
    api.get<Prosthesis[]>(`/laboratories/${laboratoryId}/prostheses`),

  get: (laboratoryId: string, id: string) =>
    api.get<Prosthesis>(`/laboratories/${laboratoryId}/prostheses/${id}`),

  create: (laboratoryId: string, data: CreateProsthesisRequest) =>
    api.post<Prosthesis>(`/laboratories/${laboratoryId}/prostheses`, data),

  update: (laboratoryId: string, id: string, data: UpdateProsthesisRequest) =>
    api.put<Prosthesis>(`/laboratories/${laboratoryId}/prostheses/${id}`, data),

  delete: (laboratoryId: string, id: string) =>
    api.delete<void>(`/laboratories/${laboratoryId}/prostheses/${id}`),
}
