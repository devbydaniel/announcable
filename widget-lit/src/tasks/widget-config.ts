import { Task } from '@lit/task';
import { ReactiveController, ReactiveControllerHost } from 'lit';
import type { WidgetConfig } from '@/lib/types';
import { backendUrl } from '@/lib/config';

/**
 * Task for fetching widget configuration
 * Replaces useConfig hook from React version
 */
export class WidgetConfigTask implements ReactiveController {
  host: ReactiveControllerHost;
  private orgId: string;
  
  task: Task<[string], WidgetConfig>;

  constructor(host: ReactiveControllerHost, orgId: string) {
    this.host = host;
    this.orgId = orgId;
    host.addController(this);
    
    this.task = new Task(
      host,
      async ([orgId]) => {
        const url = `${backendUrl}/api/widget-config/${orgId}`;
        const res = await fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        if (!res.ok) {
          throw new Error('Failed to fetch widget config');
        }
        
        const { data } = await res.json();
        return {
          ...data,
          border_radius: parseInt(data.border_radius),
          border_width: parseInt(data.border_width),
        } as WidgetConfig;
      },
      () => [this.orgId]
    );
  }

  hostConnected() {}
  hostDisconnected() {}
}
