import { Task } from '@lit/task';
import { ReactiveController, ReactiveControllerHost } from 'lit';
import { backendUrl } from '@/lib/config';

export interface ReleaseNoteStatus {
  last_update_on: string;
  attention_mechanism?: string;
}

/**
 * Task for fetching release note status
 * Replaces useReleaseNoteStatus hook from React version
 */
export class ReleaseNoteStatusTask implements ReactiveController {
  host: ReactiveControllerHost;
  private orgId: string;
  
  task: Task<[string], ReleaseNoteStatus[]>;

  constructor(host: ReactiveControllerHost, orgId: string) {
    this.host = host;
    this.orgId = orgId;
    host.addController(this);
    
    this.task = new Task(
      host,
      async ([orgId]) => {
        const url = `${backendUrl}/api/release-notes/${orgId}/status?for=widget`;
        const response = await fetch(url);
        
        if (!response.ok) {
          throw new Error('Failed to fetch release note status');
        }
        
        const json = await response.json();
        return (json.data || []) as ReleaseNoteStatus[];
      },
      () => [this.orgId]
    );
  }

  hostConnected() {}
  hostDisconnected() {}
}
