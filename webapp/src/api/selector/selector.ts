import { useState } from 'react';
import useFetch from 'use-http';
import routes from 'services/app-settings';
import type { Collections } from 'types';

type Params = {
  limit: number;
};

const useSelector = ({ limit }: Params) => {
  const [page, setPage] = useState(0);

  const { loading, error, data = [] } = useFetch<Collections[]>(`${routes.selector}?limit=${limit}&skip=${page * limit}`, {
    onNewData: (acc = [], curr) => [...acc, ...curr.data],
    perPage: limit,
  }, [page]);

  return {
    loading,
    error,
    data,
    next() {
      setPage(page + 1);
    },
  };
};

export default useSelector;
