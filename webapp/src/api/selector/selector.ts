import { useState } from 'react';
import useFetch from 'hooks/useFetch';
import routes from 'services/routes';

type Params = {
  limit: number;
};

type Data = {
  'count': any,
  'skip': any,
  'limit': any,
  'data': any,
};

function useSelector({ limit }: Params) {
  const [page, setPage] = useState(0);

  const {
    isLoading, error, data = [], fetch,
  } = useFetch<Data>(`${routes.selector}?limit=${limit}&skip=${page * limit}`, { onMount: false });

  return {
    loading: isLoading, error, data, next: fetch,
  };
}

export default useSelector;
