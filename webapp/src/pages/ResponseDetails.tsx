import React from 'react';
import { useParams } from 'react-router-dom';
import useGetResponses from 'api/useGetResponses';

const Response: React.FC = () => {
  let { endpoint } = useParams<{ endpoint:string }>();
  let { data, loading } = useGetResponses({ limit: 5, id: endpoint });

  if (loading) {
    return <>loading</>;
  }

  return (
    <>
      <div>Response Detail</div>
      {data && data.map(d => (
        d.data.map((r) => Object.entries(r).map(([k, v]) => (
          <p key={k} className="break-words"><b>{k}:</b>{v}</p>
        )))
      ))}
    </>
  );
};

export default Response;
