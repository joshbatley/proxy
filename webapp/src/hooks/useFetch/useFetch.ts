import { useState, useEffect, useCallback, useRef } from 'react'

export interface Options extends RequestInit {
  onMount: boolean;
}

export type Return<T> = {
  data: T | undefined,
  isLoading: boolean;
  error: Error | undefined;
  fetch: (url?: RequestInfo, options?: Options) => void;
}

function useFetch<T>(url: RequestInfo, options: Options = { onMount: true }): Return<T> {
  const u = useRef<RequestInfo>()
  const o = useRef<Options | undefined>()
  const [data, setData] = useState<T | undefined>();
  const [error, setError] = useState<Error | undefined>();
  const [loading, setLoading] = useState<boolean>(false);
  const [mounted, setMounted] = useState<boolean>(false);

  u.current = url;
  o.current = options;

  const call = useCallback(async () => {
    try {
      setMounted(true)
      setLoading(true)
      const res = await fetch(u.current!, o.current);
      if (!res.ok) {
        throw Error(await res.json())
      }
      const data = await res.json()
      setData(data);
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!mounted && options?.onMount) {
      call()
    }
  }, [mounted, call, options])

  return {
    data,
    isLoading: loading,
    error,
    fetch(url?: RequestInfo, options?: Options) {
      if (url) {
        u.current = url
      }
      if (options) {
        o.current = options
      }
      call();
    }
  }
}

export default useFetch
