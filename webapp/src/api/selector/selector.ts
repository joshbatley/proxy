import usePagination from 'hooks/usePagination';
import routes from 'services/app-settings';
import type { Collections } from 'types';

type Params = {
  limit: number;
};

const useSelector = ({ limit }: Params) => {
  const {
    data, error, loading, next,
  } = usePagination<Collections>({ input: routes.selector, limit });
  return {
    data, error, loading, next,
  };
};

export default useSelector;
