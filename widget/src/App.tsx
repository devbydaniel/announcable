import { Button } from "./components/ui/button";
import { Gift } from "lucide-react";
import Widget from "./components/widget";
import { ErrorPanel } from "./components/ui/errorPanel";
import {
  ReleaseNotesList,
  ReleaseNoteEntry,
  ReleaseNoteSkeleton,
} from "./components/ui/releaseNotesList";
import { Indicator, AnchorIndicator } from "./components/ui/indicator";
import type { WidgetInit } from "./lib/types";
import useReleaseNotes from "./hooks/useReleaseNotes";
import useWidgetConfig from "./hooks/useConfig";
import { useEffect, useMemo } from "react";
import useWidgetToggle from "./hooks/useWidgetToggle";
import useAnchorsRef from "./hooks/useAnchorsRef";
import useReleaseNoteStatus, {
  ReleaseNoteStatus,
} from "./hooks/useReleaseNoteStatus";

interface Props {
  init: WidgetInit;
}
export default function App({ init }: Props) {
  const { isOpen, setIsOpen, lastOpened } = useWidgetToggle({
    querySelector: init.anchor_query_selector,
  });

  const anchorsRef = useAnchorsRef({
    querySelector: init.anchor_query_selector,
  });

  const queryParams = useMemo(
    () => ({
      orgId: init.org_id,
    }),
    [init.org_id],
  );

  const { data: releaseNoteStatus } = useReleaseNoteStatus(queryParams);

  const hasUnseenValue = hasUnseenReleaseNotes({
    lastOpened,
    releaseNoteStatus,
  });

  const shouldDisplayIndicator =
    anchorsRef && !init.hide_indicator && hasUnseenValue;

  useEffect(() => {
    anchorsRef?.forEach((ref) => {
      updateIndicatorDataset(ref, hasUnseenValue);
    });
  }, [hasUnseenValue, anchorsRef]);

  useEffect(() => {
    if (shouldInstantOpen({ lastOpened, releaseNoteStatus }) && !isOpen) {
      setIsOpen(true);
      anchorsRef?.forEach((ref) => {
        updateInstantOpenDataset(ref, true);
      });
    }
  }, [lastOpened, isOpen, setIsOpen]);

  return (
    <>
      {shouldDisplayIndicator &&
        Array.from(anchorsRef).map((ref, i) => (
          <AnchorIndicator key={i} anchorElement={ref} />
        ))}
      {!init.anchor_query_selector && (
        <Button
          className="fixed z-50 bottom-4 right-4"
          variant={isOpen ? "default" : "outline"}
          size="icon"
          onClick={() => setIsOpen(!isOpen)}
        >
          {hasUnseenValue && <Indicator className="absolute top-0 right-0" />}
          <Gift className="w-4 h-4" />
        </Button>
      )}
      {isOpen && (
        <WidgetContent
          init={init}
          isOpen={isOpen}
          onClose={() => setIsOpen(false)}
        />
      )}
    </>
  );
}

interface WidgetContentProps {
  init: WidgetInit;
  isOpen: boolean;
  onClose: () => void;
}

function WidgetContent({ init, isOpen, onClose }: WidgetContentProps) {
  const {
    data: releaseNotes,
    isLoading: releaseNotesAreLoading,
    error: releaseNotesError,
  } = useReleaseNotes({ orgId: init.org_id });

  const {
    data: widgetConfig,
    isLoading: widgetConfigIsLoading,
    error: widgetConfigError,
  } = useWidgetConfig({ orgId: init.org_id });

  const isReadyToMount =
    !widgetConfigIsLoading && !releaseNotesError && !widgetConfigError;

  useEffect(() => {
    if (releaseNotesError) {
      console.error(releaseNotesError);
    }
  }, [releaseNotesError]);

  useEffect(() => {
    if (widgetConfigError) {
      console.error(widgetConfigError);
    }
  }, [widgetConfigError]);

  if (!isReadyToMount) return null;

  return (
    <div className="relative">
      <Widget
        config={widgetConfig!}
        init={init}
        isOpen={isOpen}
        onClose={onClose}
      >
        {releaseNotesError || widgetConfigError ? (
          <ErrorPanel />
        ) : (
          <ReleaseNotesList>
            {releaseNotesAreLoading ? (
              <ReleaseNoteSkeleton config={widgetConfig!} />
            ) : (
              releaseNotes!.map((releaseNote) => (
                <ReleaseNoteEntry
                  key={releaseNote.id}
                  config={widgetConfig!}
                  releaseNote={releaseNote}
                />
              ))
            )}
          </ReleaseNotesList>
        )}
      </Widget>
    </div>
  );
}

function hasUnseenReleaseNotes({
  lastOpened,
  releaseNoteStatus,
}: {
  lastOpened?: string | null;
  releaseNoteStatus?: ReleaseNoteStatus[];
}) {
  if (!releaseNoteStatus) return false;
  if (!lastOpened) return true;
  const lastOpenedDate = new Date(parseInt(lastOpened));
  return releaseNoteStatus.some(
    (note) =>
      note.last_update_on && new Date(note.last_update_on) > lastOpenedDate,
  );
}

function updateIndicatorDataset(
  anchorElement: HTMLElement | null,
  shouldDisplay: boolean,
) {
  if (!anchorElement) return;
  const newValue = shouldDisplay ? "true" : "false";
  anchorElement.dataset.newReleaseNotes = newValue;
}

function shouldInstantOpen({
  lastOpened,
  releaseNoteStatus,
}: {
  lastOpened?: string | null;
  releaseNoteStatus?: ReleaseNoteStatus[];
}) {
  if (!releaseNoteStatus) return false;
  if (!lastOpened) return true;
  const lastOpenedDate = new Date(parseInt(lastOpened));
  return releaseNoteStatus.some(
    (note) =>
      note.attention_mechanism === "instant_open" &&
      new Date(note.last_update_on) > lastOpenedDate,
  );
}

function updateInstantOpenDataset(
  anchorElement: HTMLElement | null,
  shouldDisplay: boolean,
) {
  if (!anchorElement) return;
  const newValue = shouldDisplay ? "true" : "false";
  anchorElement.dataset.instantOpen = newValue;
}
