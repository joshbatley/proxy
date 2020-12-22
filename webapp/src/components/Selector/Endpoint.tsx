import React from 'react';
import MethodTag from 'components/MethodTag';
import Portal from 'components/Portal';
import type { Endpoint } from 'types';

interface BaseProps extends React.ComponentPropsWithoutRef<'div'> {
  data: Endpoint
  truncate?: boolean;
}

let baseClasses = 'hover:shadow py-1.5 px-2 cursor-pointer leading-normal flex content-center flex-auto text-sm hover:bg-gray-100 rounded-l text-gray-700';
const BaseEnd: React.FC<BaseProps> = ({
  className, truncate, data, ...other
}) => (
  <div
    className={`${baseClasses} ${className}`}
    {...other}
  >
    <MethodTag method={data.method} />
    <span className={`leading-5 ${truncate ? 'truncate' : 'whitespace-nowrap overflow-clip'}`}>{data.url}</span>
  </div>
);

type Props = {
  data: Endpoint;
  handleClick: (id: string) => void;
};

const EndpointLink: React.FC<Props> = ({ data, handleClick }) => {
  let [showTooltip, setTooltip] = React.useState<{x: number, y: number} | null>(null);

  function hover(e: React.MouseEvent<HTMLDivElement>) {
    let text = e.currentTarget.lastElementChild!;
    let target = e.currentTarget;

    if (text.clientWidth < text.scrollWidth) {
      let t = target!.getBoundingClientRect();
      let x = t.left + window.scrollX;
      let y = t.top + window.scrollY;
      if (x && y && target) {
        setTooltip({ x, y });
      }
    }
  }
  return (
    <>
      {showTooltip && (
        <Portal>
          <BaseEnd
            data={data}
            className="bg-gray-100 shadow rounded absolute pointer-events-none z-5 "
            style={{
              top: showTooltip.y,
              left: showTooltip.x,
            }}
          />
        </Portal>
      )}
      <BaseEnd
        data={data}
        className={`z-10 relative ${showTooltip && 'hover:shadow-none'}`}
        truncate={!showTooltip}
        onClick={() => handleClick(data.id)}
        onMouseOver={hover}
        onMouseLeave={() => setTooltip(null)}
      />
    </>
  );
};

export default EndpointLink;
