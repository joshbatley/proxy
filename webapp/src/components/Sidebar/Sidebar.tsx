import React from 'react';
import Logo from 'components/Logo';
import { Collections } from 'types';
import Selector from 'components/Selector';

type Props = {
  data: Collections[]
};

const Sidebar: React.FC<Props> = ({ data }) => (
  <aside className="border-r h-screen w-full">
    <div className="border-b border-gray-200">
      <Logo />
    </div>
    <Selector collections={data} />
  </aside>
);

export default Sidebar;
