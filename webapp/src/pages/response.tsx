import React from 'react';
import { useParams } from 'react-router-dom';
import useResponse from 'api/response';

const Response: React.FC = () => {
  let { endpoint } = useParams<{ endpoint:string }>();
  let { data, loading } = useResponse({ limit: 5, id: endpoint });

  if (loading) {
    return <>loading</>;
  }

  return (
    <div>
      <header className="App-header">
        <div>Response</div>
        {data && data.map(d => (
          d.data.map((r) => Object.entries(r).map(([k, v]) => (
            <div key={k}><b>{k}:</b>{v}</div>
          )))
        ))}
      </header>
    </div>
  );
};

export default Response;
