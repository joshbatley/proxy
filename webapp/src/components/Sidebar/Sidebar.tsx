import React from 'react';
import Logo from 'components/Logo';
import Selector from 'components/Selector';

const Sidebar: React.FC = () => (
  <aside className="h-screen w-full">
    <div className="border-b border-gray-200">
      <Logo />
    </div>
    <Selector />
  </aside>
);

export default Sidebar;
