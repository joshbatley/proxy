import React from 'react';
import Method from 'components/Method';
import Url from 'components/Url';
import { Endpoint } from 'types';

type Props = {
  data: Endpoint
};

const EndpointDetails: React.FC<Props> = ({ data }) => (
  <div className="w-100 bg-white p-4 border-b flex group">
    <Method method={data.method} />
    <Url url={data.url} />
  </div>
);

export default EndpointDetails;
