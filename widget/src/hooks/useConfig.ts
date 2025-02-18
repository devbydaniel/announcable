import { useQuery } from "@tanstack/react-query";
import type { WidgetConfig } from "@/lib/types";
import { backendUrl } from "@/lib/config";

interface Props {
  orgId: string;
}
export default function useWidgetConfig({ orgId }: Props) {
  const url = `${backendUrl}/api/widget-config/${orgId}`;
  const { data, isLoading, error } = useQuery({
    queryKey: ["widget-config", orgId],
    queryFn: async () => {
      const res = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const { data } = await res.json();
      return {
        ...data,
        border_radius: parseInt(data.border_radius),
        border_width: parseInt(data.border_width),
      } as WidgetConfig;
    },
  });

  return { data, isLoading, error };
}
