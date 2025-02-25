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
        `${backendUrl}/api/release-notes/${orgId}/status?for=widget`,
      );
      if (!response.ok) {
        throw new Error("Failed to fetch release note status");
      }
      const json = await response.json();
      return (json.data || []) as ReleaseNoteStatus[];
    },
    staleTime: Infinity,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
  });
}
