import { useEffect, useState } from "react";

interface Props {
  querySelector?: string;
}

export default function useWidgetToggle({ querySelector }: Props) {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [lastOpened, setLastOpened] = useState<string | null>(
    localStorage.getItem("release_beacon_widget_last_opened"),
  );

  function setIsOpenAndUpdateLocalStorage(isOpen: boolean) {
    setIsOpen((isOpen) => !isOpen);
    if (!isOpen) {
      const now = Date.now().toString();
      setLastOpened(now);
      localStorage.setItem("release_beacon_widget_last_opened", now);
    }
  }

  const anchors =
    querySelector &&
    (document.querySelectorAll(querySelector) as NodeListOf<HTMLElement>);

  function toggleWidget() {
    setIsOpen((isOpen) => !isOpen);
    if (!isOpen) {
      const now = Date.now().toString();
      setLastOpened(now);
      localStorage.setItem("release_beacon_widget_last_opened", now);
    }
  }

  useEffect(() => {
    if (anchors && anchors.length > 0) {
      anchors.forEach((anchor) => {
        anchor.addEventListener("click", toggleWidget);
      });
      return () => {
        anchors.forEach((anchor) => {
          anchor.removeEventListener("click", toggleWidget);
        });
      };
    }
  });

  return { isOpen, setIsOpen: setIsOpenAndUpdateLocalStorage, lastOpened };
}
