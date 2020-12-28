import React from 'react';
import { Route, Switch } from 'react-router-dom';

import Layout from 'components/Layout';
import Sidebar from 'components/Sidebar';

import Placeholder from 'pages/placeholder';
import Response from 'pages/response';

const App: React.FC = () => (
  <Layout sidebar={<Sidebar />}>
    <Switch>
      <Route component={Response} path="/:collection/:endpoint" />
      <Route path="/">
        <Placeholder />
      </Route>
    </Switch>
  </Layout>
);

export default App;
