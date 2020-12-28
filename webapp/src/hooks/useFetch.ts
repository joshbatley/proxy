import { useState, useEffect, useCallback } from 'react';

type UseFetch = {
  input: RequestInfo;
  init?: RequestInit;
  runOnMount: boolean;
};

function useFetch<T>({
  input, init, runOnMount,
}: UseFetch) {
  let [data, setData] = useState<T>();
  let [loading, setLoading] = useState<boolean | null>(null);
  let [error, setError] = useState<Error | null>(null);

  let req = useCallback(async () => {
    setLoading(true);
    try {
      let response = await fetch(input, init);
      if (response.ok) {
        let d: T = await response.json();
        setData(d);
      }
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [input, init]);

  useEffect(() => {
    if (!loading && runOnMount) {
      req();
    }
  }, [req]);

  return {
    loading,
    error,
    data,
    req,
  };
}

export default useFetch;
