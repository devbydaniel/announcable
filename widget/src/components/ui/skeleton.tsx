import { cn } from "@/lib/utils";

function Skeleton({
  className,
  style,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("animate-pulse rounded-md bg-muted", className)}
      style={style}
      {...props}
    />
  );
}

export { Skeleton };
