import usePagination from 'hooks/usePagination';
import routes from 'services/app-settings';
import type { Collections } from 'types';

type Params = {
  limit: number;
};

const useGetCollections = ({ limit }: Params) => {
  let {
    data, error, loading, next,
  } = usePagination<Collections>({ input: routes.collections, limit });
  return {
    data, error, loading, next,
  };
};

export default useGetCollections;
