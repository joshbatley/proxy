import usePagination from 'hooks/usePagination';
import routes from 'services/app-settings';
import type { Response } from 'types';

type Params = {
  limit: number;
  id?: string | null;
};

const useResponse = ({ limit, id }: Params) => {
  let {
    data, error, loading, next,
  } = usePagination<Response>({ input: `${routes.response}?endpoint=${id}`, limit });
  return {
    data, error, loading, next,
  };
};

export default useResponse;