import React from 'react';

const Placeholder: React.FC<{ next: () => void; }> = ({ next }) => (
  <div>
    <header className="App-header">
      <button type="button" onClick={() => next()}>next page</button>
    </header>
  </div>
);

export default Placeholder;
