import React, { useState } from 'react';
import { Route, Switch } from 'react-router-dom';

import useSelector from 'api/selector';
import useResponse from 'api/response';

import Layout from 'components/Layout';
import Sidebar from 'components/Sidebar';

import Placeholder from 'pages/placeholder';
import Response from 'pages/response';

const App: React.FC = () => {
  const [id, setId] = useState<string | null>(null);

  const { data: res, loading: resLoading } = useResponse({ limit: 1, id });
  const { data, loading, next } = useSelector({ limit: 3 });

  return (
    <Layout sidebar={<Sidebar data={data} handleClick={setId} />}>
      <Switch>
        <Route component={Response} path="/123" />
        <Route path="/">
          <Placeholder next={next} />
        </Route>
      </Switch>
    </Layout>
  );
};

export default App;
