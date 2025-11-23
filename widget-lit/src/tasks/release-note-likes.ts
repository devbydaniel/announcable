import { Task, TaskStatus } from '@lit/task';
import { ReactiveController, ReactiveControllerHost } from 'lit';
import { backendUrl } from '@/lib/config';

type LikeState = {
  is_liked: boolean;
};

interface LikesTaskOptions {
  releaseNoteId: string;
  orgId: string;
  clientId: string;
}

/**
 * Controller for managing release note likes
 * Replaces useReleaseNoteLikes hook from React version
 */
export class ReleaseNoteLikesController implements ReactiveController {
  host: ReactiveControllerHost;
  
  private releaseNoteId: string;
  private orgId: string;
  private clientId: string;
  
  task: Task<[string, string, string], LikeState>;
  isPending = false;
  error: Error | null = null;

  constructor(host: ReactiveControllerHost, options: LikesTaskOptions) {
    this.host = host;
    this.releaseNoteId = options.releaseNoteId;
    this.orgId = options.orgId;
    this.clientId = options.clientId;
    host.addController(this);
    
    // Task to fetch current like state
    this.task = new Task(
      host,
      async ([releaseNoteId, orgId, clientId]) => {
        if (!clientId) {
          return { is_liked: false };
        }
        
        const response = await fetch(
          `${backendUrl}/api/release-notes/${orgId}/${releaseNoteId}/like?clientId=${clientId}`,
          { method: 'GET' }
        );
        
        if (!response.ok) {
          throw new Error('Failed to get like state');
        }
        
        return await response.json() as LikeState;
      },
      () => [this.releaseNoteId, this.orgId, this.clientId]
    );
  }

  hostConnected() {}
  hostDisconnected() {}

  /**
   * Toggle like state for this release note
   */
  async toggleLike(): Promise<void> {
    if (!this.clientId) {
      return;
    }

    this.isPending = true;
    this.error = null;
    this.host.requestUpdate();

    try {
      const response = await fetch(
        `${backendUrl}/api/release-notes/${this.orgId}/${this.releaseNoteId}/like`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            release_note_id: this.releaseNoteId,
            client_id: this.clientId,
          }),
        }
      );

      if (!response.ok) {
        throw new Error('Failed to toggle like');
      }

      // Manually trigger a re-run of the task to refresh the like state
      this.task.run();
    } catch (error) {
      this.error = error instanceof Error ? error : new Error('Unknown error');
      console.error('Failed to toggle like:', error);
    } finally {
      this.isPending = false;
      this.host.requestUpdate();
    }
  }

  /**
   * Get current like state
   */
  get isLiked(): boolean {
    if (this.task.status === TaskStatus.COMPLETE && this.task.value) {
      return this.task.value.is_liked;
    }
    return false;
  }
}
