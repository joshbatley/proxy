import React from 'react';
import { createPortal } from 'react-dom';

type Props = {
  forwardRef?: unknown;
  as?: string;
};

const Portal: React.FC<Props> = ({
  children, forwardRef, as, ...props
}) => (typeof window !== 'undefined'
  ? createPortal(
    as
      ? React.createElement(
        as,
        {
          ref: forwardRef,
          ...props,
        },
        children,
      )
      : children,
    document.getElementById('portal')!,
  )
  : null);

export default Portal;
