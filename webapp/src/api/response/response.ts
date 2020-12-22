import { useEffect, useState } from 'react';
import useFetch, { CachePolicies } from 'use-http';
import routes from 'services/app-settings';
import type { Response } from 'types';

type Params = {
  limit: number;
  id?: string | null;
};

const useResponse = ({ limit, id }: Params) => {
  // const [data, setData] = useState([]);
  const [page, setPage] = useState(0);
  // const [id, setId] = useState<string | null>(null);

  const {
    loading, data, error,
  } = useFetch<Response[]>(`${routes.response}?limit=${limit}&skip=${page * limit}&endpoint=${id}`, {
    onNewData: (acc = [], curr) => [...acc, ...curr.data],
    cachePolicy: CachePolicies.NO_CACHE,
    perPage: limit,
  }, [page, id]);

  // useEffect(() => {
  //   if (id) {
  //     const t = get();
  //     t.then((i) => console.log(i));
  //   }
  // }, [get, id]);

  return {
    loading,
    error,
    data,
    next() {
      setPage(page + 1);
    },
  };
};

export default useResponse;
