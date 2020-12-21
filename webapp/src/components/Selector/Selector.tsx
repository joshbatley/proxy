import React from 'react';
import ReactDOM from 'react-dom';
import MethodTag from 'components/MethodTag';
import type { Collections, Endpoint } from 'types';

type Props = {
  collections: Collections[]
};

const Chev: React.FC<{ isOpen: boolean }> = ({ isOpen }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" height="16px" className={`m-0.5 self-center ${isOpen ? 'transform rotate-90' : ''}`}>
    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
  </svg>
);

const BaseEnd: React.FC<{ className: string, style: any, e: any }> = ({
  className, style, e, ...other
}) => (
  <div key={e.id} className={`py-1 px-2 truncate cursor-pointer leading-normal flex content-center flex-auto text-sm ${className}`} {...other} style={style}>
    <MethodTag method={e.method} />
    <span className="truncate leading-5">{e.url}</span>
  </div>
);

const EndpointComp: React.FC<{e: Endpoint}> = ({
  e, ...other
}) => {
  const [isHere, setHere] = React.useState<{e: Endpoint, x: number, y: number} | null>(null);

  let style = {};
  let classNames = '';
  if (isHere) {
    style = {
      top: isHere.y,
      left: isHere.x,
    };
    classNames += 'bg-gray-200 rounded absolute pointer-events-none';
  } else {
    classNames += 'hover:bg-gray-200 rounded-l';
  }

  function hover(ev: any, es: Endpoint) {
    if (ev.target.offsetWidth < ev.target.scrollWidth) {
      const t = ev.target.parentElement.getBoundingClientRect();
      const x = t.left + window.scrollX;
      const y = t.top + window.scrollY;
      if (x && y && ev.target.parentElement) {
        setHere({ e: es, x, y });
      }
    }
  }

  return (
    <>
      {isHere && ReactDOM.createPortal(<BaseEnd e={e} className={classNames} style={style} />, document.getElementById('portal')!)}
      <BaseEnd
        e={e}
      // @ts-ignore
        onMouseOver={(ev: any) => hover(ev, e)}
        onMouseLeave={() => setHere(null)}
      // @ts-ignore
        onFocus={(ev) => hover(ev, e)}
      />
    </>
  );
};

const Selector: React.FC<Props> = ({ collections }) => {
  const [isOpen, setOpen] = React.useState<boolean>(false);

  return (
    <>
      <div className="min-w-full max-w-full truncate">
        {collections.map((c) => (
          <div key={c.id} className="text-sm min-w-full max-w-full truncate">
            <button type="button" className="px-2 py-4 border-b flex content-center hover:bg-gray-200 w-full" onClick={() => setOpen(!isOpen)}>
              <div className="flex content-center"><Chev isOpen={isOpen} /></div>
              <div>
                <div>{c.name}</div>
                <div className="text-xs text-gray-600">
                  {(c.endpoints && c.endpoints.length) || 0}
                  {' '}
                  endpoints
                </div>
              </div>
            </button>
            {isOpen && c.endpoints && (
            <div className="py-2 pl-2 border-b">
              {c.endpoints.map((e) => (
                <EndpointComp
                  key={e.id}
                  e={e}
                />
              ))}
            </div>
            )}
          </div>
        ))}
      </div>
    </>
  );
};

export default Selector;
