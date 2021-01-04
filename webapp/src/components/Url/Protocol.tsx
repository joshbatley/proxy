import React from 'react';

type Props = {
  text: string;
  onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
};

const options = ['http:', 'https:'];

const Protocol: React.FC<Props> = ({ text, onChange }) => (
  <select value={text} onChange={onChange}>
    {options.map(o => (<option>{o}</option>))}
  </select>
);

export default Protocol;
