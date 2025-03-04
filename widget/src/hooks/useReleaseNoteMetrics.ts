import { useCallback, useEffect, useRef } from "react";
import { backendUrl } from "@/lib/config";

type MetricType = "view" | "cta_click";

interface Props {
  releaseNoteId: string;
  orgId: string;
}

const CLIENT_ID_KEY = 'announcable_client_id';

export default function useReleaseNoteMetrics({ releaseNoteId, orgId }: Props) {
  const hasTrackedView = useRef(false);
  const elementRef = useRef<HTMLDivElement | null>(null);

  // Get or create client ID
  const getClientId = useCallback(() => {
    let clientId = localStorage.getItem(CLIENT_ID_KEY);
    if (!clientId) {
      clientId = crypto.randomUUID();
      localStorage.setItem(CLIENT_ID_KEY, clientId);
    }
    return clientId;
  }, []);

  const sendMetric = useCallback(
    async (type: MetricType) => {
      try {
        const clientId = getClientId();
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
    [releaseNoteId, orgId, getClientId],
  );

  const trackCtaClick = useCallback(() => {
    sendMetric("cta_click");
  }, [sendMetric]);

  useEffect(() => {
    // Wait a bit for the ScrollArea to be fully mounted
    const timeoutId = setTimeout(() => {
      const viewport = document.querySelector('[data-radix-scroll-area-viewport]');
      
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