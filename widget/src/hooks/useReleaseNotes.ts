import { useQuery } from "@tanstack/react-query";
import { type ReleaseNote } from "@/lib/types";

interface Props {
  orgId: string;
  backendUrl: string;
}
export default function useReleaseNotes({ orgId, backendUrl }: Props) {
  const url = `${backendUrl}/api/release-notes/${orgId}`;
  const { data, isLoading, error } = useQuery({
    queryKey: ["release-notes", orgId],
    queryFn: async () => {
      const res = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const { data } = await res.json();
      return data as ReleaseNote[];
    },
  });

  return { data, isLoading, error };
}
