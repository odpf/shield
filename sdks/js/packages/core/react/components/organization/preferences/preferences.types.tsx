import { ReactNode } from 'react';

export type PreferencesSelectionTypes = {
  label: string;
  name: string;
  text: string;
  defaultValue?: string;
  values: { title: ReactNode; value: string }[];
  onSelection?: (value: string) => void;
};
