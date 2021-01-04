import React from 'react';
import Logo from 'components/Logo';

const Header: React.FC = () => (
  <div className="border-b border-gray-200 row-start-1 col-span-4 text-white bg-gray-900">
    <Logo />
  </div>
);

export default Header;
