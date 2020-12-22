import React, { useState } from 'react';

import useSelector from 'api/selector';
import useResponse from 'api/response';

import Layout from 'components/Layout';
import Sidebar from 'components/Sidebar';

const App: React.FC = () => {
  const [id, setId] = useState<string | null>(null);

  const { data: res, loading: resLoading } = useResponse({ limit: 5, id });
  const { data, loading, next } = useSelector({ limit: 5 });

  console.log(res);
  return (
    <Layout sidebar={<Sidebar data={data} handleClick={setId} />}>
      <div>
        <header className="App-header">
          <button type="button" onClick={() => next()}>next page</button>
        </header>
      </div>
    </Layout>
  );
};

export default App;
