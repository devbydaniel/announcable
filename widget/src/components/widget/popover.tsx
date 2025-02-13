import { Button } from "../ui/button";
import { ExternalLinkIcon, XIcon } from "lucide-react";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "../ui/card";
import { ScrollArea } from "../ui/scroll-area";
import { WidgetProps } from ".";

export default function PopoverWidget({
  config,
  init,
  children,
  onClose,
}: WidgetProps) {
  const { title, description } = config;
  return (
    <Card
      className="w-[32rem] fixed bottom-20 right-4"
      style={{
        borderRadius: config.widget_border_radius,
        borderWidth: config.widget_border_width,
        borderColor: config.widget_border_color,
        backgroundColor: config.widget_bg_color,
        color: config.widget_font_color,
        fontFamily: init.font_family?.join(","),
        zIndex: 9999,
      }}
    >
      <div className="absolute top-0 right-0 p-2 flex">
        <Button size="icon" variant="ghost">
          <ExternalLinkIcon className="w-4 h-4" />
        </Button>
        <Button size="icon" variant="ghost" onClick={onClose}>
          <XIcon className="w-4 h-4" />
        </Button>
      </div>
      <CardHeader>
        <CardTitle className="text-lg inline-flex items-center">
          {title}
        </CardTitle>
        {description && (
          <CardDescription style={{ color: config.widget_font_color }}>
            {description}
          </CardDescription>
        )}
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[32rem]">{children}</ScrollArea>
      </CardContent>
    </Card>
  );
}
