import React from 'react';
import { AuthProvider } from './src/contexts/AuthContext';

// Wrap the app with AuthProvider for SSR
export const wrapRootElement = ({ element }) => {
  return <AuthProvider>{element}</AuthProvider>;
};
