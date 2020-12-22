import React from 'react';

type Props = {
  method: string;
};

type Map = {
  name?: string;
  color?: string;
  bg?: string;
};

const methodValues = {
  GET: {
    name: 'GET', color: 'text-green-900', bg: 'bg-green-100',
  },
  POST: {
    name: 'POST', color: 'text-orange-900', bg: 'bg-orange-100',
  },
  PUT: {
    name: 'PUT', color: 'text-blue-900', bg: 'bg-blue-100',
  },
  PATCH: 'PATCH',
  DELETE: {
    name: 'DEL', color: 'text-red-900', bg: 'bg-red-100',
  },
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
} as Record<string, Map | string>;

function mapper(m: string) {
  let color = 'text-gray-700';
  let bg = 'bg-gray-300';

  let mapped = methodValues[m];
  if (typeof mapped === 'string') {
    return {
      name: mapped || m,
      color,
      bg,
    };
  }

  return {
    name: mapped.name || m,
    color: mapped.color || color,
    bg: mapped.bg || bg,
  };
}

let classes = 'min-w-min w-11 font-bold text-xxs p-1 mr-2 rounded leading-3 float-left text-center flex-shrink-0';
const MethodTag: React.FC<Props> = ({ method }) => {
  let mapped = mapper(method);
  return (
    <div className={`${classes} ${mapped.color} ${mapped.bg}`}>{mapped.name}</div>
  );
};

export default MethodTag;
