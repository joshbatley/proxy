import React, { createContext, useContext, useState } from 'react';
import RulesModal from 'components/RulesModal';
import useGetRules from 'api/useGetRules';
import type { Wrapped, Rule } from 'types';

interface RulesModalManager {
  isOpen: boolean;
  closeModal: () => void;
  openModal: (id: string) => void;
  data: Wrapped<Rule>[];
  loading: boolean | null;
  error: Error | null;
}

export const RulesModalCtx = createContext<RulesModalManager | void>(undefined);

export function useRulesModal(): RulesModalManager {
  let context = useContext(RulesModalCtx);
  if (context === undefined) {
    throw new Error('useRulesModal must be used within a useRulesModalProvider');
  }

  return context;
}

export const RulesModalProvider: React.FC = ({ children }) => {
  let [isOpen, setOpen] = useState<string | undefined>(undefined);
  let { data, loading, error } = useGetRules({ limit: 20, id: isOpen });

  return (
    <RulesModalCtx.Provider value={{
      isOpen: Boolean(isOpen),
      closeModal: () => setOpen(undefined),
      openModal: (id) => setOpen(id),
      data,
      loading,
      error,
    }}
    >
      {children}
      {Boolean(isOpen) && <RulesModal data={data} closeModal={() => setOpen(undefined)} />}
    </RulesModalCtx.Provider>
  );
};
