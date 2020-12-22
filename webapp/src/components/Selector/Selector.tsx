import React from 'react';
import type { Collections } from 'types';
import Endpoint from './Endpoint';
import Collection from './Collection';

type Props = {
  collections: Collections[]
};

const Selector: React.FC<Props> = ({ collections }) => (
  <div className="min-w-full max-w-full">
    {collections.map((c) => (
      <Collection key={c.id} data={c}>
        <div className="py-2 pl-2 border-b">
          {c.endpoints && c.endpoints.map((e) => (
            <Endpoint key={e.id} data={e} />
          ))}
          {!c.endpoints && (
          <p className="break-words whitespace-normal">This collection is empty.  to this collection and create folders to organize them</p>
          )}
        </div>
      </Collection>
    ))}
  </div>
);

export default Selector;
