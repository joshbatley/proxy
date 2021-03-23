import React from 'react';
import { Route, Switch } from 'react-router-dom';

import Layout from 'components/Layout';
import Placeholder from 'pages/placeholder';
import ResponseList from 'pages/ResponseList';
import ResponseDetails from 'pages/ResponseDetails';

const App: React.FC = () => (
  <Layout>
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
