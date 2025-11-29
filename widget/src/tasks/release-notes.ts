import { Task } from '@lit/task';
import { ReactiveController, ReactiveControllerHost } from 'lit';
import type { ReleaseNote } from '@/lib/types';
import { backendUrl } from '@/lib/config';

/**
 * Task for fetching release notes
 * Replaces useReleaseNotes hook from React version
 */
export class ReleaseNotesTask implements ReactiveController {
  host: ReactiveControllerHost;
  private orgId: string;
  
  task: Task<[string], ReleaseNote[]>;

  constructor(host: ReactiveControllerHost, orgId: string) {
    this.host = host;
    this.orgId = orgId;
    host.addController(this);
    
    this.task = new Task(
      host,
      async ([orgId]) => {
        const url = `${backendUrl}/api/release-notes/${orgId}?for=widget`;
        const res = await fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        if (!res.ok) {
          throw new Error('Failed to fetch release notes');
        }
        
        const { data } = await res.json();
        return (data || []) as ReleaseNote[];
      },
      () => [this.orgId]
    );
  }

  hostConnected() {}
  hostDisconnected() {}
}
