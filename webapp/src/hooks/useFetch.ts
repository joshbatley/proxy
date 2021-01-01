import { useState, useEffect, useCallback } from 'react';

type UseFetch = {
  input: RequestInfo;
  init?: RequestInit;
  runOnMount?: boolean;
  runOnNullParams?: boolean;
};

function checkForNulls(url: string): boolean {
  let u = new URL(url);
  let hasNullValue = false;
  u.searchParams.forEach((value) => {
    if (value === 'null' || value === 'undefined' || value === null || value === undefined) {
      hasNullValue = true;
    }
  });

  return hasNullValue;
}

function checkNullValues(input: RequestInfo): boolean {
  return typeof input === 'string' ? checkForNulls(input) : checkForNulls(input.url);
}

function useFetch<T>({
  input, init, runOnMount, runOnNullParams,
}: UseFetch) {
  let [data, setData] = useState<T>();
  let [loading, setLoading] = useState<boolean | null>(null);
  let [error, setError] = useState<Error | null>(null);
  let allowedToRun = checkNullValues(input) && !runOnNullParams;

  let req = useCallback(async () => {
    if (allowedToRun) {
      return;
    }

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
