export const handleErrors = (response: Response) => {
  if (!response.ok) {
    throw new Error(response.statusText)
  }
  return response
}
