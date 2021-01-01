import React from 'react';

type Props = {
  status: number;
};

type Map = {
  color?: string;
  bg?: string;
};

const StatusValues = {
  1: { color: 'text-blue-900', bg: 'bg-blue-100' },
  2: { color: 'text-green-900', bg: 'bg-green-100' },
  3: { color: 'text-orange-900', bg: 'bg-orange-100' },
  4: { color: 'text-red-900', bg: 'bg-red-100' },
} as Record<string, Map>;

function mapper(m: number) {
  let color = 'text-gray-700';
  let bg = 'bg-gray-300';

  let mapped = StatusValues[m.toString().charAt(0)];
  if (!mapped) {
    return {
      color,
      bg,
    };
  }

  return {
    color: mapped.color || color,
    bg: mapped.bg || bg,
  };
}

let classes = 'min-w-min w-11 font-bold text-xs p-1 mr-2 rounded leading-3 text-center flex-shrink-0 inline';
const Status: React.FC<Props> = ({ status }) => {
  let mapped = mapper(status);
  return (
    <div className={`${classes} ${mapped.color} ${mapped.bg}`}>{status}</div>
  );
};

export default Status;
