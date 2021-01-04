import React from 'react';
import Sidebar from 'components/Sidebar';
import { RulesModalProvider } from 'contexts/RulesModalManager';
import Header from './Header';

const height = {
  height: 'calc(100vh - calc(4rem + 1px))',
};

const Layout: React.FC = ({ children }) => (
  <>
    <Header />
    <div className="h-full w-screen grid grid-cols-4 grid-rows-1 overflow-hidden" style={height}>
      <RulesModalProvider>
        <Sidebar />
        <div className="bg-gray-50 row-start-2 col-start-3 col-span-2 sm:col-start-2 sm:col-span-3">
          {children}
        </div>
      </RulesModalProvider>
    </div>
  </>
);

export default Layout;
