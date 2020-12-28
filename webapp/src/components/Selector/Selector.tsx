import React from 'react';
import useSelector from 'api/selector';
import Endpoint from './Endpoint';
import Collection from './Collection';

const Selector: React.FC = () => {
  let { data, loading, next } = useSelector({ limit: 10 });

  return (
    <div className="min-w-full max-w-full">
      {data?.map(({ data: collections }) => collections.map((c) => (
        <Collection key={c.id} data={c}>
          {c.endpoints && (
          <div className="py-2 pl-2 shadow-inset">
            {c.endpoints.map((e) => (
              <Endpoint key={e.id} data={e} />))}
          </div>
          )}
          {!c.endpoints && (
          <div className="p-2 text-gray-700">
            <p className="break-words whitespace-normal text-gray-800">This collection is empty. Manually add a request or hook up your API to pass through here.</p>
          </div>
          )}
        </Collection>
      )))}
    </div>
  );
};

export default Selector;
