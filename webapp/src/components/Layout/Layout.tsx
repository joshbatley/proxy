import React from 'react';

type Props = {
  sidebar: React.ReactNode
};

const Layout: React.FC<Props> = ({ children, sidebar }) => (
  <div className="flex h-full">
    <div className="w-1/2 sm:w-1/3 md:w-1/4 xl:w-1/5 h-auto overflow-hidden border-r">{sidebar}</div>
    <div className="max-h-full flex-auto">
      {children}
    </div>
  </div>
);

export default Layout;
