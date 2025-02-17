import { Button } from "./components/ui/button";
import { Gift } from "lucide-react";
import Widget from "./components/widget";
import { Skeleton } from "./components/ui/skeleton";
import { ErrorPanel } from "./components/ui/errorPanel";
import {
  ReleaseNotesList,
  ReleaseNoteEntry,
} from "./components/ui/releaseNotesList";
import { Indicator, AnchorIndicator } from "./components/ui/indicator";
import withSkeleton from "./components/hoc/withSkeleton";
import type { WidgetInit } from "./lib/types";
import type { ReleaseNote } from "./lib/types";
import useReleaseNotes from "./hooks/useReleaseNotes";
import useWidgetConfig from "./hooks/useConfig";
import { useEffect } from "react";
import useWidgetToggle from "./hooks/useWidgetToggle";
import useAnchorsRef from "./hooks/useAnchorsRef";

interface Props {
  backendUrl: string;
  init: WidgetInit;
}

export default function App({ init, backendUrl }: Props) {
  const { isOpen, setIsOpen, lastOpened } = useWidgetToggle({
    querySelector: init.anchor_query_selector,
  });

  const anchorsRef = useAnchorsRef({
    querySelector: init.anchor_query_selector,
  });

  const {
    data: releaseNotes,
    isLoading: releaseNotesAreLoading,
    error: releaseNotesError,
  } = useReleaseNotes({ orgId: init.org_id, backendUrl });

  const {
    data: widgetConfig,
    isLoading: widgetConfigIsLoading,
    error: widgetConfigError,
  } = useWidgetConfig({ orgId: init.org_id, backendUrl });

  const isReadyToMount =
    !releaseNotesAreLoading &&
    !widgetConfigIsLoading &&
    !releaseNotesError &&
    !widgetConfigError;

  const hasUnseenReleaseNotesValue = hasUnseenReleaseNotes({
    lastOpened,
    releaseNotes,
  });

  const shouldDisplayIndicator =
    anchorsRef &&
    !init.hide_indicator &&
    hasUnseenReleaseNotesValue &&
    isReadyToMount;

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

  useEffect(() => {
    anchorsRef?.forEach((ref) => {
      updateIndicatorDataset(ref, hasUnseenReleaseNotesValue);
    });
  }, [hasUnseenReleaseNotesValue, anchorsRef]);

  useEffect(() => {
    if (shouldInstantOpen({ lastOpened, releaseNotes }) && !isOpen) {
      setIsOpen(true);
      anchorsRef?.forEach((ref) => {
        updateInstantOpenDataset(ref, true);
      });
    }
  }, [lastOpened, releaseNotes, setIsOpen]);

  return (
    <>
      {shouldDisplayIndicator &&
        Array.from(anchorsRef).map((ref, i) => (
          <AnchorIndicator key={i} anchorElement={ref} />
        ))}
      {!init.anchor_query_selector && isReadyToMount && (
        <Button
          className="fixed z-50 bottom-4 right-4"
          variant={isOpen ? "default" : "outline"}
          size="icon"
          onClick={() => setIsOpen(!isOpen)}
        >
          {hasUnseenReleaseNotes({ lastOpened, releaseNotes }) && (
            <Indicator className="absolute top-0 right-0" />
          )}
          <Gift className="w-4 h-4" />
        </Button>
      )}
      {isOpen && isReadyToMount && (
        <div className="relative">
          <Widget
            config={widgetConfig!}
            init={init}
            isOpen={isOpen}
            onClose={() => setIsOpen(false)}
          >
            {releaseNotesError || widgetConfigError ? (
              <ErrorPanel />
            ) : (
              withSkeleton(() => (
                <ReleaseNotesList>
                  {releaseNotes!.map((item, i) => (
                    <ReleaseNoteEntry
                      config={widgetConfig!}
                      key={i}
                      releaseNote={item}
                    />
                  ))}
                </ReleaseNotesList>
              ))({
                skeleton: (
                  <ReleaseNotesList>
                    {Array.from({ length: 3 }).map((_, i) => (
                      <Skeleton key={i} className="h-24" />
                    ))}
                  </ReleaseNotesList>
                ),
                isLoading: releaseNotesAreLoading,
              })
            )}
          </Widget>
        </div>
      )}
    </>
  );
}

function hasUnseenReleaseNotes({
  lastOpened,
  releaseNotes,
}: {
  lastOpened?: string | null;
  releaseNotes?: ReleaseNote[];
}) {
  if (!releaseNotes) return false;
  if (!lastOpened) return true;
  const lastOpenedDate = new Date(parseInt(lastOpened));
  return releaseNotes.some(
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
  releaseNotes,
}: {
  lastOpened?: string | null;
  releaseNotes?: ReleaseNote[];
}) {
  if (!releaseNotes) return false;
  if (!lastOpened) return true;
  const lastOpenedDate = new Date(parseInt(lastOpened));
  return releaseNotes.some(
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
