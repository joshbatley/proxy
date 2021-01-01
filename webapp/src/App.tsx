import React from 'react';
import { Route, Switch } from 'react-router-dom';

import Layout from 'components/Layout';
import Sidebar from 'components/Sidebar';

import Placeholder from 'pages/Placeholder';
import ResponseList from 'pages/ResponseList';
import ResponseDetails from 'pages/ResponseDetails';

const App: React.FC = () => (
  <Layout sidebar={<Sidebar />}>
    <Switch>
      <Route component={ResponseDetails} path="/:collection/:endpointId/:response" />
      <Route component={ResponseList} path="/:collection/:endpointId" />
      <Route path="/">
        <Placeholder />
      </Route>
    </Switch>
  </Layout>
);

export default App;
