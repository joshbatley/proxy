import React from 'react';

type Props = {
  method: string;
};

const map = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PATCH',
  DELETE: 'DEL',
  COPY: 'COPY',
  HEAD: 'HEAD',
  OPTIONS: 'OPT',
  LINK: 'LINK',
  UNLINK: 'UNLNK',
  PURGE: 'PURGE',
  LOCK: 'LOCK',
  UNLOCK: 'UNLCK',
  PROPFIND: 'PROP',
  VIEW: 'VIEW',
} as Record<string, string>;

// {
//   name: ''
//   color: '',
//   background: '',
// }

const MethodTag: React.FC<Props> = ({ method }) => {
  let color = '';
  const isGet = method === 'GET';

  if (isGet) {
    color += 'bg-green-100 text-green-900';
  } else {
    color += 'bg-gray-300 text-gray-700';
  }

  return (
    <div className={`min-w-min w-11 font-bold text-xxs p-1 mr-2 rounded leading-3 float-left text-center ${color} flex-shrink-0`}>{map[method]}</div>
  );
};

export default MethodTag;
