export const backendUrl =
  import.meta.env.MODE === "production"
    ? "https://release-notes.danielbenner.de"
    : "http://localhost:3000";

export const clientIdKey = "announcable_client_id";
