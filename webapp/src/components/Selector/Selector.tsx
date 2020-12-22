import React from 'react';
import type { Collections } from 'types';
import Endpoint from './Endpoint';
import Collection from './Collection';

type Props = {
  collections: Collections[];
  handleClick: (id: string) => void;
};

const Selector: React.FC<Props> = ({ collections, handleClick }) => (
  <div className="min-w-full max-w-full">
    {collections.map((c) => (
      <Collection key={c.id} data={c}>
        {c.endpoints && (
          <div className="py-2 pl-2 shadow-inset">
            {c.endpoints.map((e) => (
              <Endpoint key={e.id} data={e} handleClick={handleClick} />))}
          </div>
        )}
        {!c.endpoints && (
          <div className="p-2 text-gray-700">
            <p className="break-words whitespace-normal text-gray-800">This collection is empty. Manually add a request or hook up your API to pass through here.</p>
          </div>
        )}
      </Collection>
    ))}
  </div>
);

export default Selector;
