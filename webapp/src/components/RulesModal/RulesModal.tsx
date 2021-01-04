import React from 'react';
import { Rule, Wrapped } from 'types';
import Portal, { Overlay } from 'components/Portal';
import Modal from 'components/Modal';

type Props ={
  data: Wrapped<Rule>[]
  closeModal: () => void;
};

const RulesModal: React.FC<Props> = ({ data, closeModal }) => (
  <Portal>
    <Overlay close={closeModal}>
      <Modal close={closeModal} title="Rules">
        {data.map(r => (<></>))}
      </Modal>
    </Overlay>
  </Portal>
);

export default RulesModal;
