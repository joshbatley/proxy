import React, { useState } from 'react';
import { Collections } from 'types';

type Props = {
  data: Collections
};

const Chevron: React.FC<{ isOpen: boolean }> = ({ isOpen }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" height="18px" className={`mr-2 self-center ${isOpen ? 'transform rotate-90' : ''}`}>
    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
  </svg>
);

type BtnProps = {
  name: string;
  endpoints: number;
  isOpen: boolean;
  toggle: () => void;
};

const CollectionBtn: React.FC<BtnProps> = ({
  name, endpoints, isOpen, toggle,
}) => (
  <button
    type="button"
    className={`px-2 py-4 flex content-center hover:bg-gray-100 w-full ${isOpen && 'bg-gray-100 border-b '}`}
    onClick={toggle}
  >
    <Chevron isOpen={isOpen} />
    <div className="text-left">
      <span>{name}</span>
      <span className="text-xs text-gray-600 block">
        {endpoints} endpoints
      </span>
    </div>
  </button>
);

const Collection: React.FC<Props> = ({ data, children }) => {
  let [isOpen, setOpen] = useState<boolean>(false);

  function toggle() {
    setOpen(!isOpen);
  }

  return (
    <div className="text-sm border-b">
      <CollectionBtn
        name={data.name}
        endpoints={data.endpoints?.length || 0}
        toggle={toggle}
        isOpen={isOpen}
      />
      {isOpen && children}
    </div>
  );
};

export default Collection;
