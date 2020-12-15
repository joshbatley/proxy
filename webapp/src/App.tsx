import React from 'react';
import useSelector from 'api/selector';
import Layout from 'components/Layout';
import Sidebar from 'components/Sidebar';

const App: React.FC = () => {
  const { data, loading, next } = useSelector({ limit: 5 });

  return (
    <Layout sidebar={<Sidebar />}>
      <div>
        <header className="App-header">
          <button type="button" onClick={() => next()}>next page</button>
        </header>
      </div>
    </Layout>
  );
};

export default App;
