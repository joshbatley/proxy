import React from 'react';
import Sidebar from 'components/Sidebar';
import { RulebarProvider } from 'contexts/RulebarManager';
import Logo from 'components/Logo';

const height = {
  height: 'calc(100vh - calc(4rem + 1px))',
};

const Layout: React.FC = ({ children }) => (
  <>
    <div className="border-b border-gray-200 row-start-1 col-span-4">
      <Logo />
    </div>
    <div className="h-full w-screen grid grid-cols-4 grid-rows-1 overflow-hidden" style={height}>
      <RulebarProvider>
        <div className="h-auto border-r col-start-1 col-end-3 row-start-2 sm:col-end-2"><Sidebar /></div>
        <div className="p-3 bg-gray-50 row-start-2 col-start-3 col-span-2 sm:col-start-2 sm:col-span-3 sm:flex">
          {children}
        </div>
      </RulebarProvider>
    </div>
  </>
);

export default Layout;
