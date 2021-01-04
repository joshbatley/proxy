import { useRulesModal } from 'contexts/RulesModalManager';
import React, { useState } from 'react';
import { useRouteMatch } from 'react-router-dom';
import { Collections } from 'types';
import SettingButtons from './SettingButtons';

type Props = {
  data: Collections
};

const Chevron: React.FC<{ isOpen: boolean }> = ({ isOpen }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" height="18px" className={`mr-2 self-center ${isOpen ? 'transform rotate-90' : ''}`}>
    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
  </svg>
);

type BtnProps = {
  name: string;
  endpoints: number;
  isOpen: boolean;
  toggle: () => void;
  openSettings: (e: React.MouseEvent) => void;
};

const CollectionBtn: React.FC<BtnProps> = ({
  name, endpoints, isOpen, toggle, openSettings,
}) => {
  let [ShowBtns, SetShowBtns] = useState(false);

  return (
    <button
      type="button"
      className={`pl-2 flex content-center hover:bg-gray-100 w-full h-16 ${isOpen && 'bg-gray-100 border-b'}`}
      onClick={toggle}
      onMouseLeave={() => SetShowBtns(false)}
      onMouseOver={() => SetShowBtns(true)}
      onFocus={() => SetShowBtns(true)}
    >
      <Chevron isOpen={isOpen} />
      <div className="text-left py-3">
        <span>{name}</span>
        <span className="text-xs text-gray-600 block">
          {endpoints} endpoints
        </span>
      </div>
      {ShowBtns && <SettingButtons openSettings={openSettings} />}
    </button>
  );
};

const Collection: React.FC<Props> = ({ data, children }) => {
  let matched = useRouteMatch<{ collection: string }>('/:collection/:endpoint');
  let { openModal, closeModal } = useRulesModal();
  let collection = matched?.params.collection;
  let [isOpen, setOpen] = useState<boolean>(data.id === collection);

  function toggle() {
    setOpen(!isOpen);
    if (isOpen) {
      closeModal();
    }
  }

  function openSettings(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    openModal(data.id);
  }

  return (
    <div className="text-sm border-b">
      <CollectionBtn
        name={data.name}
        endpoints={data.endpoints?.length || 0}
        toggle={toggle}
        isOpen={isOpen}
        openSettings={openSettings}
      />
      {isOpen && children}
    </div>
  );
};

export default Collection;
