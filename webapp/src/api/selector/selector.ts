import { useState } from 'react';
import useFetch from 'use-http';
import routes from 'services/app-settings';

type Params = {
  limit: number;
};

type Data = {
  id: string;
  name: string
  endpoints: Endpoints[]
};

type Endpoints = {
  id: string;
  status: number;
  url: string;
  method: string;
};

const useSelector = ({ limit }: Params) => {
  const [page, setPage] = useState(0);

  const { loading, error, data = [] } = useFetch<Data[]>(`${routes.selector}?limit=${limit}&skip=${page * limit}`, {
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
