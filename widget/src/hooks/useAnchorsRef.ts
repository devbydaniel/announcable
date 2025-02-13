import { useEffect, useRef, useState } from "react";

interface Props {
  querySelector?: string;
}

export default function useAnchorsRef({ querySelector }: Props) {
  const [anchorElements, setAnchorElements] =
    useState<NodeListOf<HTMLElement> | null>(null);
  const anchorsRef = useRef<NodeListOf<HTMLElement> | null>(null);

  useEffect(() => {
    if (!querySelector) return;

    const elements = document.querySelectorAll(
      querySelector,
    ) as NodeListOf<HTMLElement>;
    if (elements && elements.length > 0) {
      setAnchorElements(elements);
    }
  }, [querySelector]);

  useEffect(() => {
    if (anchorElements && anchorElements.length > 0) {
      anchorsRef.current = anchorElements;
    }
  }, [anchorElements]);

  return anchorsRef.current;
}
