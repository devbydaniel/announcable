import { useEffect, useRef } from "react";
import { cn } from "@/lib/utils";

interface IndicatorProps {
  className?: string;
}

export function Indicator({ className }: IndicatorProps) {
  return (
    <div
      className={cn(
        "bg-red-500 rounded-full w-1.5 h-1.5 translate-x-1 -translate-y-1",
        className,
      )}
    />
  );
}

interface AnchorIndicatorProps {
  anchorElement: HTMLElement;
  className?: string;
}

export function AnchorIndicator({
  anchorElement,
  className,
}: AnchorIndicatorProps) {
  const indicatorRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    if (anchorElement) {
      const indicator = document.createElement("div");
      // set style of the indicator
      // through style property because of shadow DOM
      indicator.style.position = "absolute";
      indicator.style.top = "0";
      indicator.style.right = "0";
      indicator.style.backgroundColor = "red";
      indicator.style.borderRadius = "50%";
      indicator.style.width = "8px";
      indicator.style.height = "8px";
      indicator.style.transform = "translate(-50%, -50%)";

      indicatorRef.current = indicator;
      anchorElement.appendChild(indicator);

      // Clean up function
      return () => {
        if (indicatorRef.current) {
          anchorElement.removeChild(indicatorRef.current);
        }
      };
    }
  }, [anchorElement, className]);

  // The component doesn't render anything directly
  return null;
}
