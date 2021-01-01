import { useState, useEffect, useRef } from 'react';
import type { Wrapped } from 'types';
import useFetch from 'hooks/useFetch';

type UsePagination = {
  input: RequestInfo;
  init?: RequestInit;
  limit: number;
  resetOn?: unknown[];
};

function appendParams(url: string, limit: number, page: number): string {
  let u = new URL(url);
  u.searchParams.set('limit', limit.toString());
  u.searchParams.set('skip', (page * limit).toString());
  return u.toString();
}

function formatRequest(input: RequestInfo, limit: number, page: number) {
  let managedInput = input;
  if (typeof managedInput === 'string') {
    managedInput = appendParams(managedInput, limit, page);
  } else {
    managedInput = {
      ...managedInput,
      url: appendParams(managedInput.url, limit, page),
    };
  }

  return managedInput;
}

function usePagination<T>({
  input, init, limit, resetOn,
}: UsePagination) {
  let [page, setPage] = useState(0);
  let refresh = useRef<unknown[] | undefined>(undefined);
  let managedInput = formatRequest(input, limit, page);
  let [paginatedData, setPaginatedData] = useState<Wrapped<T>[]>([]);
  let [canFetchMore, setFetchMore] = useState(true);
  let {
    data, loading, error,
  } = useFetch<Wrapped<T>>({
    input: managedInput, init, runOnMount: true, runOnNullParams: false,
  });

  useEffect(() => {
    if (refresh.current === undefined) {
      refresh.current = resetOn;
    }

    if (refresh.current !== undefined) {
      let dataChange = refresh.current.map((v, i) => {
        if (resetOn?.[i] !== v) {
          return true;
        }
      }).every(Boolean);

      if (dataChange) {
        refresh.current = resetOn;
        setPaginatedData([]);
      }
    }

    if (error === null && loading === false && data) {
      let alreadySaved = paginatedData.some(i => i.skip === data?.skip);
      if (data.skip !== page * limit || alreadySaved) {
        return;
      }
      if (data.count < limit) {
        setFetchMore(false);
      }
      setPaginatedData([...paginatedData, data]);
    }
  }, [data, page, loading, resetOn]);

  return {
    loading,
    error,
    data: paginatedData,
    canFetchMore,
    next() {
      if (canFetchMore) {
        setPage(page + 1);
      }
    },
  };
}

export default usePagination;
