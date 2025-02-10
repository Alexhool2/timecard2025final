import { ConflictError, UnauthorizedError } from "../errors/http_errors"

const API_URL = import.meta.env.VITE_API_URL
console.log("API URL:", API_URL)

async function fetchData(input: RequestInfo, init?: RequestInit) {
  const response = await fetch(input, {
    ...init,
    credentials: "include",
  })
  if (response.ok) {
    return response
  } else {
    const errorBody = await response.json()
    const errorMessage = errorBody.error
    if (response.status === 401) {
      throw new UnauthorizedError(errorMessage)
    } else if (response.status === 409) {
      throw new ConflictError(errorMessage)
    } else {
      throw Error(
        "Request failed with status " +
          response.status +
          "message:" +
          errorMessage
      )
    }
  }
}

export interface LoginCredentials {
  userName: string
  password: string
}

export interface LoggedUser {
  id: number
  userName: string
  role: string
}

export async function login(
  credentials: LoginCredentials
): Promise<LoggedUser> {
  const response = await fetchData(`${API_URL}/api/v1/users/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },

    body: JSON.stringify(credentials),
  })
  const data = await response.json()
  const user: LoggedUser = {
    id: data.user.id,
    userName: data.user.userName,
    role: data.user.role,
  }

  return user
}
export interface LoggedInUser {
  id: number
  firstName: string
  lastName: string
  userName: string
  email: string
  isAdmin: boolean
  createdAt: string
  role: string
}

export async function getLoggedInUser(): Promise<{
  user: LoggedUser | null
  eventID: number | null
  role: string
}> {
  const response = await fetchData(`${API_URL}/api/v1/users/me`, {
    method: "GET",
    credentials: "include",
  })
  const data = await response.json()

  return {
    user: data.user || null,
    eventID: data.eventID === -1 ? null : data.eventID,
    role: data.user.role,
  }
}

export async function logout() {
  const response = await fetchData(`${API_URL}/api/v1/users/logout`, {
    method: "POST",
    credentials: "include",
  })

  return response.json()
}

export async function getUserEvents(userId: number, date: string) {
  const response = await fetch(
    `${API_URL}/api/v1/event/user/${userId}?date=${date}`,
    {
      method: "GET",
      credentials: "include",
    }
  )
  if (!response.ok) {
    throw new Error("error finding event")
  }
  return response.json()
}

export interface GetAllUsersInterface {
  id: number
  firstName: string
  lastName: string
  userName: string
  email: string
  IsAdmin: boolean
  role: string
}

export async function getAllUsers(): Promise<GetAllUsersInterface[]> {
  const response = await fetchData(`${API_URL}/api/v1/users/`, {
    method: "GET",
    credentials: "include",
  })
  if (!response.ok) {
    throw new Error("error finding users")
  }
  const users: GetAllUsersInterface[] = await response.json()
  return users
}

export interface CreateUserInterface {
  firstName: string
  lastName: string
  userName: string
  password: string
  email: string
  role: string
}
export async function createUser(
  user: CreateUserInterface
): Promise<CreateUserInterface> {
  const response = await fetchData(`${API_URL}/api/v1/users/create`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(user),
  })
  if (!response.ok) {
    throw new Error("error creating user")
  }

  if (response.ok) {
    return user
  }
  return await response.json()
}

export async function changePassword(userId: number, newPassword: string) {
  const response = await fetchData(
    `${API_URL}/api/v1/users/${userId}/reset-password`,
    {
      method: "PUT",
      credentials: "include",
      body: JSON.stringify({ newPassword }),
    }
  )
  if (!response.ok) {
    throw new Error("error updating user")
  }
  return await response.json()
}
