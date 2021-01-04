import useFetch from 'hooks/useFetch';
import routes from 'services/app-settings';
import type { Endpoint } from 'types';

type Params = {
  id?: string | null;
};

const useGetEndpoint = ({ id }: Params) => {
  let {
    data, error, loading,
  } = useFetch<Endpoint>({ input: `${routes.endpoint}/${id}`, runOnMount: true });
  return {
    data, error, loading,
  };
};

export default useGetEndpoint;
