import React from 'react';

const Response = () => {
  function next() {
    console.log('next');
  }
  return (
    <div>
      <header className="App-header">
        <div>Response</div>
        <button type="button" onClick={() => next()}>next page</button>
      </header>
    </div>
  );
};

export default Response;
