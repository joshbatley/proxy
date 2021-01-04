import React, { useCallback } from 'react';

type Props = {
  title: string;
  close: () => void;
};

const CloseIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" width="16px" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
  </svg>
);

const Modal: React.FC<Props> = ({ title, close, children }) => {
  const stopPropagation = useCallback((e) => e.stopPropagation(), []);
  return (
    <div className="bg-white w-2/3 h-1/3 rounded-md p-5" onClick={stopPropagation}>
      <div>
        {title}
        <button type="button" onClick={close} onFocus={close}>
          <CloseIcon />
        </button>
      </div>
      {children}
    </div>
  );
};

export default Modal;
