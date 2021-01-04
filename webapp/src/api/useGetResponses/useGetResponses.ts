import usePagination from 'hooks/usePagination';
import routes from 'services/app-settings';
import type { Response } from 'types';

type Params = {
  limit: number;
  id?: string | null;
};

const useGetResponses = ({ limit, id }: Params) => {
  let {
    data, error, loading, next,
  } = usePagination<Response>({ input: `${routes.response}/${id}`, limit, resetOn: [id] });
  return {
    data, error, loading, next,
  };
};

export default useGetResponses;
