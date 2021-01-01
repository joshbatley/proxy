import React from 'react';
import { Rule, Wrapped } from 'types';

type Props ={
  data: Wrapped<Rule>[]
};

const RuleBar: React.FC<Props> = ({ data }) => {
  let a = '1';
  return (
    <div className="h-screen w-full z-10 row-start-2 col-start-3 col-span-2 sm:col-start-2 sm:col-end-3 bg-white shadow-md">
      <div>
        {data.map(r => (<></>))}
      </div>
    </div>
  );
};

export default RuleBar;
