import React from 'react';
import { useParams } from 'react-router-dom';
import useResponse from 'api/response';
import StatusList from 'components/StatusList';
import EndpointDetails from 'components/EndpointDetails';
import useEndpoint from 'api/endpoint';

const Response: React.FC = () => {
  let { endpointId } = useParams<{ endpointId:string }>();
  let { data: responses, loading } = useResponse({ limit: 5, id: endpointId });
  let { data: endpoint, loading: eLoading } = useEndpoint({ id: endpointId });

  // if (loading) {
  //   return <>loading</>;
  // }
  console.log(endpoint);

  return (
    <>
      <EndpointDetails data={endpoint} />
      <StatusList data={responses} />
    </>
  );
};

export default Response;
