import { useCallback, useEffect, useRef } from "react";
import { backendUrl } from "@/lib/config";
import { getOrCreateClientId } from "@/lib/clientId";

type MetricType = "view" | "cta_click";

interface Props {
  releaseNoteId: string;
  orgId: string;
}

export default function useReleaseNoteMetrics({ releaseNoteId, orgId }: Props) {
  const hasTrackedView = useRef(false);
  const elementRef = useRef<HTMLDivElement | null>(null);

  const sendMetric = useCallback(
    async (type: MetricType) => {
      try {
        const clientId = getOrCreateClientId();
        await fetch(`${backendUrl}/api/release-notes/${orgId}/metrics`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            release_note_id: releaseNoteId,
            metric_type: type,
            client_id: clientId,
          }),
        });
      } catch (error) {
        console.error("Failed to send release note metric:", error);
      }
    },
    [releaseNoteId, orgId],
  );

  const trackCtaClick = useCallback(() => {
    sendMetric("cta_click");
  }, [sendMetric]);

  useEffect(() => {
    // Wait a bit for the ScrollArea to be fully mounted
    const timeoutId = setTimeout(() => {
      const viewport = document.querySelector(
        "[data-radix-scroll-area-viewport]",
      );

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (
              entry.isIntersecting &&
              entry.intersectionRatio >= 0.5 &&
              !hasTrackedView.current
            ) {
              hasTrackedView.current = true;
              sendMetric("view");
              observer.disconnect();
            }
          });
        },
        {
          threshold: 0.5,
          root: viewport,
          rootMargin: "0px",
        },
      );

      // Start observing if we have an element
      if (elementRef.current) {
        observer.observe(elementRef.current);
      }

      return () => {
        observer.disconnect();
      };
    }, 100);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [sendMetric]);

  return {
    elementRef,
    trackCtaClick,
  };
}
