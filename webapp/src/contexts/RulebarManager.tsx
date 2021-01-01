import React, { createContext, useContext, useState } from 'react';
import RuleBar from 'components/RuleBar';
import useRules from 'api/rules';
import type { Wrapped, Rule } from 'types';

interface RulebarManager {
  isOpen: boolean;
  closeModal: () => void;
  openModal: (id: string) => void;
  data: Wrapped<Rule>[];
  loading: boolean | null;
  error: Error | null;
}

export const RulebarCtx = createContext<RulebarManager | void>(undefined);

export function useRulebar(): RulebarManager {
  let context = useContext(RulebarCtx);
  if (context === undefined) {
    throw new Error('useRulebar must be used within a useRulebarProvider');
  }

  return context;
}

export const RulebarProvider: React.FC = ({ children }) => {
  let [isOpen, setOpen] = useState<string | undefined>(undefined);
  let { data, loading, error } = useRules({ limit: 20, id: isOpen });

  return (
    <RulebarCtx.Provider value={{
      isOpen: Boolean(isOpen),
      closeModal: () => setOpen(undefined),
      openModal: (id) => setOpen(id),
      data,
      loading,
      error,
    }}
    >
      {children}
      {Boolean(isOpen) && <RuleBar data={data} />}
    </RulebarCtx.Provider>
  );
};
