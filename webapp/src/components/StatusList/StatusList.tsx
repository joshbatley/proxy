import React from 'react';
import type { Wrapped, Response } from 'types';
import Panel from 'components/Panel';
import Status from 'components/Status';

type Props = {
  data: Wrapped<Response>[]
};

const StatusList: React.FC<Props> = ({ data }) => (
  <Panel title="Statues">
    <div>
      {data.map(d => (
        d.data.map((r) => (
          <div key={r.id} className="border-b p-3 w-full">
            <Status status={r.status} />{r.url}
          </div>
        ))
      ))}
      <button type="button" className="py-2 px-4 rounded block mt-3 w-full hover:bg-gray-100 text-gray-700">Add new response</button>
    </div>
  </Panel>
);

export default StatusList;
