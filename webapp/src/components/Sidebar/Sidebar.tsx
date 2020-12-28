import React from 'react';
import Logo from 'components/Logo';
import { Wrapped, Collections } from 'types';
import Selector from 'components/Selector';

type Props = {
  data: Wrapped<Collections>[] | undefined;
  handleClick: (id: string) => void;
};

const Sidebar: React.FC<Props> = ({ data, handleClick }) => (
  <aside className="h-screen w-full">
    <div className="border-b border-gray-200">
      <Logo />
    </div>
    <Selector data={data} handleClick={handleClick} />
  </aside>
);

export default Sidebar;
