import React from 'react';

type Props = {
  sidebar: React.ReactNode
};

const Layout: React.FC<Props> = ({ children, sidebar }) => (
  <div className="flex h-full">
    <div className="w-60 h-auto">{sidebar}</div>
    <div className="w-full max-h-full flex-auto">
      {children}

    </div>
  </div>
);

export default Layout;
