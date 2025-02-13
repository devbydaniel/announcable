import { cn } from "@/lib/utils";
import { useRef, useEffect } from "react";

interface DialogProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: React.ReactNode;
  actions?: React.ReactNode[];
  className?: string;
  style?: React.CSSProperties;
}

export function Dialog({
  isOpen,
  onClose,
  title,
  description,
  children,
  actions,
  className,
  style,
}: DialogProps) {
  const dialogRef = useRef<HTMLDialogElement | null>(null);

  useEffect(() => {
    if (isOpen) {
      dialogRef.current?.show();
      document.body.style.overflow = "hidden";
      function handleEscape(e: KeyboardEvent) {
        if (e.key === "Escape") {
          onClose();
        }
      }
      document.addEventListener("keydown", handleEscape);

      return () => {
        document.removeEventListener("keydown", handleEscape);
        document.body.style.overflow = "unset";
      };
    }
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  // Handle click outside
  function handleBackdropClick(e: React.MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }

  return (
    <div
      role="presentation"
      className="fixed inset-0 bg-black bg-opacity-70 backdrop-blur-sm flex items-center justify-center"
      onClick={handleBackdropClick}
      style={{ zIndex: 9999 }}
    >
      <dialog
        ref={dialogRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="dialog-title"
        className={cn(className, "shadow-lg p-4 max-w-[40rem] w-full")}
        tabIndex={-1}
        style={style}
      >
        <div className="grid gap-4">
          <div className="relative">
            <div className="absolute top-0 right-0 translate-x-[12px] -translate-y-[12px] flex gap-1">
              {actions?.map((action, i) => <div key={i}>{action}</div>)}
            </div>
            <h2
              id="dialog-title"
              className="font-semibold tracking-tight text-lg"
            >
              {title}
            </h2>
            {description && (
              <p className="mt-1.5 text-sm text-muted-foreground">
                {description}
              </p>
            )}
          </div>
          {children}
        </div>
      </dialog>
    </div>
  );
}
Dialog.displayName = "Dialog";
