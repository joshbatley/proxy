import React, { useCallback } from 'react';

type Props = {
  close: () => void;
};

const Overlay: React.FC<Props> = ({ children, close }) => {
  const clickModal = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    e.stopPropagation();
    close();
  }, [close]);

  const escModal = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      e.stopPropagation();
      close();
    }
  }, [close]);

  return (
    <div
      className="flex flex-wrap bg-gray-800 bg-opacity-50 w-screen h-screen overflow-hidden absolute top-0 left-0 z-50 justify-center content-center"
      onClick={clickModal}
      onKeyUp={escModal}
      tabIndex={-1}
      role="dialog"
    >
      {children}
    </div>
  );
};

export default Overlay;
