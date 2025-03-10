import React from "react";
import type { WidgetConfig } from "@/lib/types";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import type { ReleaseNote } from "@/lib/types";
import { Skeleton } from "./skeleton";
import useReleaseNoteMetrics from "@/hooks/useReleaseNoteMetrics";
import { getEmbedUrl } from "@/lib/media";

interface ReleaseNotesListProps {
  children: React.ReactNode;
}

export function ReleaseNotesList({ children }: ReleaseNotesListProps) {
  return (
    <div className="flex flex-col gap-6">
      {React.Children.toArray(children).map((child, i) => (
        <div key={i}>{child}</div>
      ))}
    </div>
  );
}

interface ReleaseNoteEntryProps {
  config: WidgetConfig;
  releaseNote: ReleaseNote;
}

export function ReleaseNoteEntry({
  config,
  releaseNote,
}: ReleaseNoteEntryProps) {
  const { elementRef, trackCtaClick } = useReleaseNoteMetrics({
    releaseNoteId: releaseNote.id,
    orgId: config.org_id,
  });

  const ctaLabel = releaseNote.cta_label_override
    ? releaseNote.cta_label_override
    : config.cta_text;
  const baseUrl = config.release_page_baseurl;
  const ctaHref =
    releaseNote.cta_href_override || `${baseUrl}#${releaseNote.id}`;

  // Get embedded media URL if available
  const embeddedMedia = releaseNote.media_link
    ? getEmbedUrl(releaseNote.media_link)
    : null;
  console.log("embeddedMedia", embeddedMedia);

  return (
    <Card
      ref={elementRef}
      style={{
        borderRadius: config.release_note_border_radius,
        borderColor: config.release_note_border_color,
        borderWidth: config.release_note_border_width,
        color: config.release_note_font_color,
        backgroundColor: config.release_note_bg_color,
      }}
    >
      <CardHeader className="pb-4">
        <CardTitle>{releaseNote.title}</CardTitle>
        <CardDescription style={{ color: config.release_note_font_color }}>
          {releaseNote.date || ""}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="w-full flex flex-col gap-4">
          {embeddedMedia?.embedUrl ? (
            <div className="relative w-full aspect-video">
              <iframe
                src={embeddedMedia.embedUrl}
                className="absolute top-0 left-0 w-full h-full"
                allow="fullscreen"
                allowFullScreen
                loading="lazy"
                referrerPolicy="no-referrer"
                sandbox="allow-scripts allow-presentation"
                title={releaseNote.title}
              />
            </div>
          ) : releaseNote.imageSrc ? (
            <div>
              <img
                src={releaseNote.imageSrc}
                alt={releaseNote.title}
                onError={(e) => {
                  console.error(
                    `Image failed to load for ${releaseNote.title}`,
                    releaseNote.imageSrc,
                    e,
                  );
                  e.currentTarget.style.display = "none";
                }}
              />
            </div>
          ) : null}
          {releaseNote.text && (
            <div className="whitespace-pre-wrap">{releaseNote.text}</div>
          )}
          {!releaseNote.hide_cta && (
            <div className="w-full flex justify-center">
              <a href={ctaHref} target="_blank" onClick={trackCtaClick}>
                {ctaLabel}
              </a>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export function ReleaseNoteSkeleton(props: { config: WidgetConfig }) {
  const { config } = props;
  const skeletonBgColor = config.widget_bg_color;
  const skeletonBorderRadius = config.release_note_border_radius;
  return (
    <Card
      style={{
        borderRadius: config.release_note_border_radius,
        borderColor: config.release_note_border_color,
        borderWidth: config.release_note_border_width,
        color: config.release_note_font_color,
        backgroundColor: config.release_note_bg_color,
      }}
    >
      <CardHeader className="pb-4">
        <div className="space-y-2">
          <Skeleton
            className={`h-7 w-3/4`}
            style={{
              backgroundColor: skeletonBgColor,
              borderRadius: skeletonBorderRadius,
            }}
          />
          <Skeleton
            className={`h-4 w-1/4`}
            style={{
              backgroundColor: skeletonBgColor,
              borderRadius: skeletonBorderRadius,
            }}
          />
        </div>
      </CardHeader>
      <CardContent>
        <div className="w-full flex flex-col gap-4">
          <Skeleton
            className={`h-48 w-full `}
            style={{
              backgroundColor: skeletonBgColor,
              borderRadius: skeletonBorderRadius,
            }}
          />
          <div className="space-y-2">
            <Skeleton
              className={`h-4 w-full`}
              style={{
                backgroundColor: skeletonBgColor,
                borderRadius: skeletonBorderRadius,
              }}
            />
            <Skeleton
              className={`h-4 w-full`}
              style={{
                backgroundColor: skeletonBgColor,
                borderRadius: skeletonBorderRadius,
              }}
            />
            <Skeleton
              className={`h-4 w-full`}
              style={{
                backgroundColor: skeletonBgColor,
                borderRadius: skeletonBorderRadius,
              }}
            />
            <Skeleton
              className={`h-4 w-3/4`}
              style={{
                backgroundColor: skeletonBgColor,
                borderRadius: skeletonBorderRadius,
              }}
            />
          </div>
          <div className="w-full flex justify-center">
            <Skeleton
              className={`h-4 w-24`}
              style={{
                backgroundColor: skeletonBgColor,
                borderRadius: skeletonBorderRadius,
              }}
            />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
