import { useQuery } from "@tanstack/react-query";
import { backendUrl } from "@/lib/config";

export interface ReleaseNoteStatus {
  last_update_on: string;
  attention_mechanism?: string;
}

export default function useReleaseNoteStatus({ orgId }: { orgId: string }) {
  return useQuery({
    queryKey: ["releaseNoteStatus", orgId],
    queryFn: async () => {
      const response = await fetch(
        `${backendUrl}/api/release-notes/${orgId}/status`,
      );
      if (!response.ok) {
        throw new Error("Failed to fetch release note status");
      }
      const json = await response.json();
      return (json.data || []) as ReleaseNoteStatus[];
    },
    // Refresh every 5 minutes
    refetchInterval: 5 * 60 * 1000,
  });
}
