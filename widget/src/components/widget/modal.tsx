import { Dialog } from "../ui/dialog";
import { ScrollArea } from "../ui/scroll-area";
import { WidgetProps } from ".";
import { ExternalLinkIcon, XIcon } from "lucide-react";
import { Button } from "../ui/button";

export default function ModalWidget({
  config,
  init,
  children,
  isOpen,
  onClose,
}: WidgetProps) {
  return (
    <Dialog
      title={config.title}
      description={config.description}
      style={{
        borderRadius: config.widget_border_radius,
        borderColor: config.widget_border_color,
        borderWidth: config.widget_border_width,
        backgroundColor: config.widget_bg_color,
        color: config.widget_font_color,
        fontFamily: init.font_family?.join(","),
      }}
      isOpen={isOpen}
      onClose={onClose}
      actions={[
        <Button size="icon" variant="ghost" asChild>
          <a href={config.release_page_baseurl} target="_blank">
            <ExternalLinkIcon className="w-4 h-4" />
          </a>
        </Button>,
        <Button size="icon" variant="ghost" onClick={onClose}>
          <XIcon className="w-4 h-4" />
        </Button>,
      ]}
    >
      <ScrollArea className="h-[32rem]">{children}</ScrollArea>
    </Dialog>
  );
}
