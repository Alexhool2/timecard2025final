class HttpError extends Error {
  constructor(message?: string) {
    super(message)
    this.name = this.constructor.name
  }
}

/**
 * Error thrown when the user is not authenticated
 */
export class UnauthorizedError extends HttpError {
  constructor(message: string = "User not authenticated") {
    super(message)
  }
}

/**
 * Error thrown when there is a conflict with existing data
 */
export class ConflictError extends HttpError {
  constructor(message: string = "Resource already exists") {
    super(message)
  }
}

/**
 * Error thrown when a resource is not found
 */
export class NotFoundError extends HttpError {
  constructor(message: string = "Resource not found") {
    super(message)
  }
}
