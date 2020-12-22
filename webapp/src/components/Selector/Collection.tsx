import React, { useState } from 'react';
import { Collections } from 'types';

type Props ={
  data: Collections
};

const Chev: React.FC<{ isOpen: boolean }> = ({ isOpen }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" height="16px" className={`m-0.5 self-center ${isOpen ? 'transform rotate-90' : ''}`}>
    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
  </svg>
);

const Collection: React.FC<Props> = ({ data, children }) => {
  let [isOpen, setOpen] = useState<boolean>(false);

  return (
    <div key={data.id} className="text-sm min-w-full max-w-full">
      <button
        type="button"
        className={`px-2 py-4 border-b flex content-center hover:bg-gray-100 w-full ${isOpen && 'bg-gray-100'}`}
        onClick={() => setOpen(!isOpen)}
      >
        <div className="flex content-center"><Chev isOpen={isOpen} /></div>
        <div>
          <div>{data.name}</div>
          <div className="text-xs text-gray-600">
            {(data.endpoints && data.endpoints.length) || 0}
            {' '}
            endpoints
          </div>
        </div>
      </button>
      {isOpen && children}
    </div>
  );
};

export default Collection;
