import { jwtDecode } from "jwt-decode"

export interface DecodedToken {
  id: number
  userName: string
  email: string
  role: string
  exp: number
}

export const getRoleFromToken = (token: string): string => {
  try {
    const decoded: DecodedToken = jwtDecode(token)
    if (!decoded.role) {
      throw new Error("The token does not contain the 'role' field.")
    }
    return decoded.role
  } catch (error) {
    console.error("failed to decode token", error)
    throw new Error("Invalid token")
  }
}
