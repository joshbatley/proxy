import React from 'react';

type Props = {
  title: string
};

const Panel: React.FC<Props> = ({ children, title }) => (
  <div className="mx-auto justify-self-center mt-10 w-3/4">
    <h3 className="text-lg font-bold">{title}</h3>
    <div className="bg-white rounded shadow p-4">{children}</div>
  </div>
);

export default Panel;
