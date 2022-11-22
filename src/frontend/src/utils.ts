import { ErrorResponse } from "@remix-run/router"

export const handleErrors = (response: Response) => {
  if (!response.ok) {
    throw new Error(response.statusText)
  }
  return response
}