import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ExternalLinkIcon, XIcon } from "lucide-react";
import { WidgetProps } from ".";
import { cn } from "@/lib/utils";

export default function SidebarWidget({
  config,
  init,
  children,
  onClose,
  isOpen,
}: WidgetProps) {
  const { title, description } = config;
  return (
    <Card
      className={cn(
        "w-[32rem] fixed top-0 right-0 h-screen transition-all duration-300 ease-in-out",
        isOpen ? "translate-x-0" : "translate-x-[100%]",
      )}
      style={{
        borderRadius: 0,
        backgroundColor: config.widget_bg_color,
        color: config.widget_font_color,
        fontFamily: init.font_family?.join(","),
        zIndex: 9999,
      }}
    >
      <div className="fixed top-0 right-0 p-2 flex">
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
      <CardContent className="h-[calc(100vh-120px)] overflow-y-auto">
        {children}
      </CardContent>
    </Card>
  );
}
