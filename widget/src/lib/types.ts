export interface ReleaseNote {
  id: string;
  title: string;
  date?: string;
  imageSrc?: string;
  text?: string;
  last_update_on: Date;
  cta_label_override?: string;
  cta_href_override?: string;
  hide_cta?: boolean;
  attention_mechanism?: null | "show_indicator" | "instant_open";
}

export interface WidgetConfig {
  org_id: string;
  title: string;
  description: string;
  cta_text?: string;
  widget_type: "popover" | "modal" | "sidebar";
  widget_border_radius: number;
  widget_border_color: string;
  widget_border_width: number;
  widget_bg_color: string;
  widget_font_color: string;
  release_note_border_radius: number;
  release_note_border_color: string;
  release_note_border_width: number;
  release_note_bg_color: string;
  release_note_font_color: string;
  release_page_baseurl: string | undefined;
}

export type WidgetInit = {
  org_id: string;
  anchor_query_selector?: string;
  hide_indicator?: boolean;
  font_family?: string[];
};
