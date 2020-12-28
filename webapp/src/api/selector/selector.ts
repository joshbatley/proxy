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

// const useSelector = ({ limit }: Params) => {
//   let [page, setPage] = useState(0);
//   let [data, setData] = useState<Collections[]>([]);
//   let [loading, setLoading] = useState<boolean | null>(null);
//   let [error, setError] = useState<Error | null>(null);

//   useEffect(() => {
//     async function req() {
//       setLoading(true);
//       try {
//         let response = await fetch(`${routes.selector}?limit=${limit}&skip=${page * limit}`);
//         if (response.ok) {
//           let d: Wrapped = await response.json();
//           if (data === null) {
//             setData([d]);
//           }
//           let exists = data.filter(i => i.skip === page * limit).length > 0;
//           if (!exists) {
//             setData([...data, d]);
//           }
//         }
//       } catch (err) {
//         setError(error);
//       } finally {
//         setLoading(false);
//       }
//     }

//     if (loading === null || loading === false) {
//       req();
//     }
//   }, [limit, page]);

//   return {
//     loading,
//     error,
//     data,
//     next() {
//       setPage(page + 1);
//     },
//   };
// };

export default useSelector;
