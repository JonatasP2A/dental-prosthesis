import { api } from '@/lib/api-client'
import type {
  Order,
  CreateOrderRequest,
  UpdateOrderRequest,
  UpdateOrderStatusRequest,
} from '@/types'

export const orderService = {
  list: (laboratoryId: string) =>
    api.get<Order[]>(`/laboratories/${laboratoryId}/orders`),

  get: (laboratoryId: string, id: string) =>
    api.get<Order>(`/laboratories/${laboratoryId}/orders/${id}`),

  create: (laboratoryId: string, data: CreateOrderRequest) =>
    api.post<Order>(`/laboratories/${laboratoryId}/orders`, data),

  update: (laboratoryId: string, id: string, data: UpdateOrderRequest) =>
    api.put<Order>(`/laboratories/${laboratoryId}/orders/${id}`, data),

  updateStatus: (
    laboratoryId: string,
    id: string,
    data: UpdateOrderStatusRequest
  ) =>
    api.patch<Order>(`/laboratories/${laboratoryId}/orders/${id}/status`, data),

  cancel: (laboratoryId: string, id: string) =>
    api.delete<void>(`/laboratories/${laboratoryId}/orders/${id}`),
}
