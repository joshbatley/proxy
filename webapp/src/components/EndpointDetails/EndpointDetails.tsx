import React from 'react';
import { Endpoint } from 'types';

type Props = {
  data?: Endpoint
};

const EndpointDetails: React.FC<Props> = ({ data }) => (
  <>{data?.url}</>
);

export default EndpointDetails;
