import React from 'react';
import { useParams } from 'react-router-dom';
import useGetResponses from 'api/useGetResponses';
import StatusList from 'components/StatusList';
import EndpointDetails from 'components/EndpointDetails';
import useGetEndpoint from 'api/useGetEndpoint';

const Response: React.FC = () => {
  let { endpointId } = useParams<{ endpointId:string }>();
  let { data: responses, loading } = useGetResponses({ limit: 5, id: endpointId });
  let { data: endpoint, loading: eLoading } = useGetEndpoint({ id: endpointId });

  // if (loading) {
  //   return <>loading</>;
  // }
  console.log(endpoint);

  return (
    <>
      {endpoint && <EndpointDetails data={endpoint} />}
      <StatusList data={responses} />
    </>
  );
};

export default Response;
