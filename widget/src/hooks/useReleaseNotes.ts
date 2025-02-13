import { useQuery } from "@tanstack/react-query";
import { type ReleaseNote } from "@/lib/types";
import config from "@/lib/config";

const backendUrl = config.backendUrl;

interface Props {
  orgId: string;
}
export default function useReleaseNotes({ orgId }: Props) {
  const url = `${backendUrl}/release-notes/${orgId}`;
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
