export interface User {
  id: number
  firstName: string
  lastName: string
  userName: string
  password: string
  email: string
  isAdmin: boolean
  createdAt: string
  role: string
}

export interface LoginCredentials {
  userName: string
  password: string
}
