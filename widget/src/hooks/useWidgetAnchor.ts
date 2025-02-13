import { useState, useEffect, useRef } from "react";

interface Props {
  querySelector?: string;
  onClick: () => void;
}

export default function useWidgetAnchor({ querySelector, onClick }: Props) {
  const [anchorElement, setAnchorElement] = useState<HTMLElement | null>(null);
  const anchorRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    if (!querySelector) return;

    const element = document.querySelector(querySelector) as HTMLElement;
    if (element) {
      element.addEventListener("click", onClick);
      setAnchorElement(element);
    }

    return () => {
      if (element) {
        element.removeEventListener("click", onClick);
      }
    };
  }, [querySelector, onClick]);

  useEffect(() => {
    if (anchorElement) {
      anchorRef.current = anchorElement;
    }
  }, [anchorElement]);

  return anchorElement;
}
