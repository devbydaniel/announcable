export const backendUrl =
  import.meta.env.MODE === "production"
    ? "https://release-notes.danielbenner.de"
    : "http://localhost:3000";
