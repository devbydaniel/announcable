/**
 * Parses YouTube and Loom URLs to get their embed URLs
 * @param url The original media URL from YouTube or Loom
 * @returns Object containing the media type and embed URL
 */
export const getEmbedUrl = (
  url: string,
): { type: "youtube" | "loom" | null; embedUrl: string | null } => {
  try {
    const urlObj = new URL(url);

    // YouTube URL parsing
    if (
      urlObj.hostname.includes("youtube.com") ||
      urlObj.hostname.includes("youtu.be")
    ) {
      let videoId = "";
      if (urlObj.hostname.includes("youtu.be")) {
        videoId = urlObj.pathname.slice(1);
      } else {
        videoId = urlObj.searchParams.get("v") || "";
      }
      if (videoId) {
        return {
          type: "youtube",
          embedUrl: `https://www.youtube.com/embed/${videoId}`,
        };
      }
    }

    // Loom URL parsing
    if (urlObj.hostname.includes("loom.com")) {
      const shareId = urlObj.pathname.split("/").pop();
      if (shareId) {
        return {
          type: "loom",
          embedUrl: `https://www.loom.com/embed/${shareId}`,
        };
      }
    }
  } catch (e) {
    console.error("Invalid URL:", url);
  }

  return { type: null, embedUrl: null };
};
