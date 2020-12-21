import React from 'react';
import MethodTag from 'components/MethodTag';
import type { Collections } from 'types';

type Props = {
  collections: Collections[]
};

const Chev: React.FC<{ isOpen: boolean }> = ({ isOpen }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" height="16px" className={`m-0.5 self-center ${isOpen ? 'transform rotate-90' : ''}`}>
    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
  </svg>
);

const Selector: React.FC<Props> = ({ collections }) => {
  const [isOpen, setOpen] = React.useState<boolean>(false);
  return (
    <div className="min-w-full max-w-full truncate">
      {collections.map((c) => (
        <div key={c.id} className="text-sm min-w-full max-w-full truncate">
          <button type="button" className="px-2 py-4 border-b flex content-center hover:bg-gray-200 w-full" onClick={() => setOpen(!isOpen)}>
            <div className="flex content-center"><Chev isOpen={isOpen} /></div>
            <div>
              <div>{c.name}</div>
              <div className="text-xs text-gray-600">
                {(c.endpoints && c.endpoints.length) || 0}
                {' '}
                endpoints
              </div>
            </div>
          </button>
          {isOpen && c.endpoints && (
          <div className="py-2 pl-2 border-b">
            {c.endpoints.map((e) => (
              <div key={e.id} className="py-1 px-2 rounded-l hover:bg-gray-200 truncate cursor-pointer leading-normal flex content-center flex-auto">
                <MethodTag method={e.method} />
                <span className="truncate leading-4" title={e.url}>{e.url}</span>
              </div>
            ))}
          </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default Selector;
