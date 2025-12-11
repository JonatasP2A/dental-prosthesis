// Common types
export interface Address {
  street: string
  city: string
  state: string
  postal_code: string
  country: string
}

// Laboratory types
export interface Laboratory {
  id: string
  name: string
  email: string
  phone: string
  address: Address
  created_at: string
  updated_at: string
}

export interface CreateLaboratoryRequest {
  name: string
  email: string
  phone: string
  address: Address
}

export interface UpdateLaboratoryRequest {
  name: string
  email: string
  phone: string
  address: Address
}

// Client types
export interface Client {
  id: string
  laboratory_id: string
  name: string
  email: string
  phone: string
  address: Address
  created_at: string
  updated_at: string
}

export interface CreateClientRequest {
  name: string
  email: string
  phone: string
  address: Address
}

export interface UpdateClientRequest {
  name: string
  email: string
  phone: string
  address: Address
}

// Order types
export type OrderStatus =
  | 'received'
  | 'in_production'
  | 'quality_check'
  | 'ready'
  | 'delivered'
  | 'revision'
  | 'cancelled'

export interface ProsthesisItem {
  type: string
  material: string
  shade?: string
  quantity: number
  notes?: string
}

export interface Order {
  id: string
  client_id: string
  laboratory_id: string
  status: OrderStatus
  prosthesis: ProsthesisItem[]
  created_at: string
  updated_at: string
}

export interface CreateOrderRequest {
  client_id: string
  prosthesis: ProsthesisItem[]
}

export interface UpdateOrderRequest {
  prosthesis: ProsthesisItem[]
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus
}

// Prosthesis (catalog) types
export interface Prosthesis {
  id: string
  laboratory_id: string
  type: string
  material: string
  shade: string
  specifications: string
  notes: string
  created_at: string
  updated_at: string
}

export interface CreateProsthesisRequest {
  type: string
  material: string
  shade?: string
  specifications?: string
  notes?: string
}

export interface UpdateProsthesisRequest {
  type: string
  material: string
  shade?: string
  specifications?: string
  notes?: string
}

// Technician types
export type TechnicianRole = 'junior' | 'senior' | 'lead' | 'manager'

export interface Technician {
  id: string
  laboratory_id: string
  name: string
  email: string
  phone: string
  role: TechnicianRole
  specializations: string[]
  created_at: string
  updated_at: string
}

export interface CreateTechnicianRequest {
  name: string
  email: string
  phone: string
  role: TechnicianRole
  specializations?: string[]
}

export interface UpdateTechnicianRequest {
  name: string
  email: string
  phone: string
  role: TechnicianRole
  specializations?: string[]
}

// API Error types
export interface ApiError {
  error: string
  details?: Record<string, string>
}
