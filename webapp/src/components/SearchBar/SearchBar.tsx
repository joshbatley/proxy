import React from 'react';

const SearchIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16px" height="100%" className="text-gray-500">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
  </svg>
);

const SearchBar = () => (
  <div className="rounded-full border flex py-0.5 pl-3 bg-gray-100 group hover:bg-gray-200 focus-within:bg-white m-3">
    <div className="mr-3">
      <SearchIcon />
    </div>
    <input type="text" placeholder="Search" className="placeholder-grey-500 outline-none bg-gray-100 group-hover:bg-gray-200 focus:bg-white" />
  </div>
);

export default SearchBar;
