import React from 'react';

type Props = {
  sidebar: React.ReactNode
};

const Layout: React.FC<Props> = ({ children, sidebar }) => (
  <div className="h-full w-screen grid grid-cols-4">
    <div className="h-auto border-r col-span-2 sm:col-span-1">{sidebar}</div>
    <div className="max-h-full col-span-2 sm:col-span-3 p-3">
      {children}
    </div>
  </div>
);

export default Layout;
