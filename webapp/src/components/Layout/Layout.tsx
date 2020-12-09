import React from 'react';

type Props = {
  sidebar: React.ReactNode
};

const Layout: React.FC<Props> = ({ children, sidebar }) => (
  <div className="">
    <div className="">{sidebar}</div>
    <div className="">
      {children}

    </div>
  </div>
);

export default Layout;
