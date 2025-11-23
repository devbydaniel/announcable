import { createContext } from '@lit/context';
import type { WidgetConfig, WidgetInit } from './types';

/**
 * Context for widget configuration (fetched from backend)
 */
export const widgetConfigContext = createContext<WidgetConfig | undefined>(
  'widget-config'
);

/**
 * Context for widget initialization parameters (from user)
 */
export const widgetInitContext = createContext<WidgetInit>('widget-init');
