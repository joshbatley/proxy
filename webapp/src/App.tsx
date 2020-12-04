import React from 'react';
import { Provider } from 'use-http';
import useSelector from 'api/selector';

const App: React.FC = () => {
  const { data, loading, next } = useSelector({ limit: 5 });

  return (
    <Provider>
      <div className="bg-black text-white hover:text-red-500">
        <header className="App-header">
          <button type="button" onClick={() => next()}>next page</button>
        </header>
      </div>
    </Provider>
  );
};

export default App;
