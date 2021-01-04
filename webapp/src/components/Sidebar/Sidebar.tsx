import React from 'react';
import useGetCollections from 'api/useGetCollections';
import SearchBar from 'components/SearchBar';
import Endpoint from './Endpoint';
import Collection from './Collection';

const EmptyCollection = () => (
  <div className="p-2 text-gray-700">
    <p className="break-words whitespace-normal text-gray-800">This collection is empty. Manually add a request or hook up your API to pass through here.</p>
  </div>
);

const PlusIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16px" height="100%" className="mr-1 mt-0.5">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
  </svg>
);

const CollectionIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16px" height="100%" className="text-gray-700 mr-1.5">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
  </svg>
);

const TrashIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16px" height="100%" className="text-gray-500 mr-1.5">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
  </svg>
);

const Sidebar: React.FC = () => {
  let { data, loading, next } = useGetCollections({ limit: 10 });

  return (
    <div className="h-screen w-full border-r col-start-1 col-end-3 row-start-2 sm:col-end-2">
      <SearchBar />
      <div className="border-b flex">
        <button type="button" className="w-1/2 border-b-4 border-blue-500 pb-1 flex justify-center">
          <CollectionIcon />Collections
        </button>
        <button type="button" className="text-gray-500 w-1/2 pb-1 flex justify-center">
          <TrashIcon /> Trash
        </button>
      </div>
      <div className="border-b">
        <button type="button" className="text-sm text-blue-500 p-3 flex content-center">
          <PlusIcon />
          <div>New collection</div>
        </button>
      </div>
      {data?.map(({ data: collections }) => collections.map((c) => (
        <Collection key={c.id} data={c}>
          {c.endpoints && (
          <div className="py-2 pl-2 shadow-inset">
            {c.endpoints.map((e) => (
              <Endpoint key={e.id} data={e} />))}
          </div>
          )}
          {!c.endpoints && <EmptyCollection />}
        </Collection>
      )))}
    </div>
  );
};

export default Sidebar;
