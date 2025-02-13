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
  const ctaLabel = releaseNote.cta_label_override
    ? releaseNote.cta_label_override
    : config.cta_text;
  const baseUrl = config.release_page_baseurl
    ? config.release_page_baseurl
    : import.meta.env.VITE_RELEASE_PAGE_BASE_URL;
  const ctaHref =
    releaseNote.cta_href_override || `${baseUrl}#${releaseNote.id}`;
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
        <CardTitle>{releaseNote.title}</CardTitle>
        <CardDescription style={{ color: config.release_note_font_color }}>
          {releaseNote.date || ""}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="w-full flex flex-col gap-4">
          {releaseNote.imageSrc && (
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
          )}
          {releaseNote.text && (
            <div className="whitespace-pre-wrap">{releaseNote.text}</div>
          )}
          {!releaseNote.hide_cta && (
            <div className="w-full flex justify-center">
              <a href={ctaHref} target="_blank">
                {ctaLabel}
              </a>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
