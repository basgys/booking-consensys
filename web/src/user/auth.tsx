import React, { useContext, useEffect, useState } from "react";
// import detectEthereumProvider from '@metamask/detect-provider'

declare global {
  interface Window {
    ethereum: any
  }
}

export interface Auth {
  token?: string

  auth: (token: string) => void
  fetcher: (input: RequestInfo, init?: RequestInit) => Promise<any>
}

export class ResponseError extends Error {
  info: any
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

const apiFetcher = async (input: RequestInfo, init?: RequestInit) => {
  const baseUrl = process.env.REACT_APP_BASE_URL || "http://localhost:8484"
  const url = new URL(input.toString(), baseUrl)

  const res = await fetch(url.toString(), init)

  // If the status code is not in the range 200-299,
  // we still try to parse and throw it.
  if (!res.ok) {
    let error: ResponseError
    switch (res.status) {
      case 400:
        error = new ResponseError('Bad request', res.status)
        break;
      case 401:
        error = new ResponseError('Authentication required', res.status)
        break;
      case 403:
        error = new ResponseError('Permission denied', res.status)
        break;
      case 404:
        error = new ResponseError('Not found', res.status)
        break;
      case 409:
        error = new ResponseError('Conflict', res.status)
        break;
      default:
        error = new ResponseError('An error occurred while fetching the data.', res.status)
        break;
    }
    error.info = await res.json() // Attach extra info to the error object.
    throw error
  }
  return res
}

const defaultFetcher = async (input: RequestInfo, init?: RequestInit) => {
  return (await apiFetcher(input, init)).json()
}

const authFetcher = (token: string) => {
  return (input: RequestInfo, init?: RequestInit) => {
    init = {
      ...init,
      headers: {
        'Authorization': token,
      }
    }
    input = input + "#" + token // For cache invalidation
    return defaultFetcher(input, init)
  }
}

export const AuthContext = React.createContext<Auth>({
  auth: (token: string) => {},
  fetcher: defaultFetcher
});

const AuthProvider = ({ children }: React.HTMLAttributes<any>) => {
  const token = process.env.REACT_APP_AUTH_TOKEN || ""

  const [auth, setAuth] = useState<Auth>({
    auth: (token: string) => {
      setAuth((a) => ({
        ...a,
        token: token,
        fetcher: authFetcher(token),
      }))
    },
    token: token,
    fetcher: authFetcher(token),
  })

  return (
    <AuthContext.Provider value={auth}>
      {children}
    </AuthContext.Provider>
  );
}

export const AuthForm = () => {
  const ctx = useContext(AuthContext)
  const [error, setError] = useState<string>("")

  useEffect(() => {
    const auth = async () => {
      // FIXME: Detect ethereum returns an "unknown" provider
      // const provider = await detectEthereumProvider();
      const provider = window.ethereum
      if (provider === undefined) {
        console.log("Please install MetaMask!")
        return null
      }
      const accounts = await provider.request({ method: 'eth_requestAccounts' })
      const account = accounts[0] // Pick first account for now
      console.log("Account", account)

      // 1. Request auth challenge
      const { challenge } = await (
        await apiFetcher('/auth/challenge', {
          method: "POST",
          body: JSON.stringify({
            address: account,
          })
        })
      ).json()

      // 2. Request account signature
      const sig = await provider.request({ method: 'personal_sign', params: [challenge, account] })
      console.log("Signature", sig)

      // 3. Finalise authentication process
      const res = await apiFetcher('/auth/authorise', {
        method: "POST",
        body: JSON.stringify({
          address: account,
          signature: sig,
        })
      })
      const token: string = res.headers.get("Authorization") || ""
      if (token) {
        ctx.auth(token)
      } else {
        throw new Error("Auth failed")
      }
    }

    auth().then(() => {
      // Redirect after login
      window.location.pathname = "/"
    }).catch((err) => {
      // TODO: Handle auth error
      setError(err.message)
    })
  }, [ctx])

  return (
    <div className="section">
      <h1>Authentication...</h1>
      <p>{error}</p>
    </div>
  );
}

export default AuthProvider;
