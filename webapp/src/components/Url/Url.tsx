import React, { useState } from 'react';
import Protocol from './Protocol';

type Props = {
  url: string;
};

const EditIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16px" height="100%" className="text-gray-500">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
  </svg>
);

const UrlComponent: React.FC<Props> = ({ url }) => {
  let urlObj = new URL(url);
  let [edit, setEdit] = useState(false);
  let [protocol, setProtocol] = useState<string>(urlObj.protocol);
  function onChange(e: React.ChangeEvent<HTMLSelectElement>) {
    setProtocol(e.target.value);
  }

  return (
    <div className="w-full leading-normal">
      { !edit
        && `${urlObj.protocol}//${urlObj.host}${urlObj.pathname}${urlObj.search}`}
      {edit && (
        <>
          <Protocol text={protocol} onChange={onChange} />//
          {urlObj.host}{urlObj.pathname}{urlObj.search}
        </>
      )}
      <button type="button" className="float-right opacity-0 h-full group-hover:opacity-100" onChange={() => setEdit(!edit)}>
        <EditIcon />
      </button>
    </div>
  );
};

export default UrlComponent;
