import { WidgetInit, type WidgetConfig } from "@/lib/types";
import PopoverWidget from "./popover";
import ModalWidget from "./modal";
import SidebarWidget from "./sidebar";

export interface WidgetProps {
  config: WidgetConfig;
  init: WidgetInit;
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

export default function Widget(props: WidgetProps) {
  const widgetType = props.config.widget_type || "popover";
  switch (widgetType) {
    case "sidebar":
      return <SidebarWidget {...props}>{props.children}</SidebarWidget>;
    case "modal":
      return <ModalWidget {...props}>{props.children}</ModalWidget>;
    default:
      return <PopoverWidget {...props}>{props.children}</PopoverWidget>;
  }
}
