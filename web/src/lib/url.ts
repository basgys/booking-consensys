
export const parseURL = (path: string) => {
  return new URL(path, process.env.REACT_APP_BASE_URL)
}