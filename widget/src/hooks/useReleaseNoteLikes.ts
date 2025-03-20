import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { backendUrl } from "@/lib/config";

type Props = {
  releaseNoteId: string;
  orgId: string;
  clientId: string;
};

type LikeState = {
  is_liked: boolean;
};

export default function useReleaseNoteLikes({
  releaseNoteId,
  orgId,
  clientId,
}: Props) {
  const queryClient = useQueryClient();
  const { data: likeState, isLoading: isLoadingLikeState } =
    useQuery<LikeState>({
      queryKey: ["release-note-like", releaseNoteId, clientId],
      queryFn: async () => {
        if (!clientId) {
          return {
            is_liked: false,
          };
        }
        const response = await fetch(
          `${backendUrl}/api/release-notes/${orgId}/${releaseNoteId}/like?clientId=${clientId}`,
          {
            method: "GET",
          },
        );
        if (!response.ok) {
          throw new Error("Failed to get like state");
        }
        return response.json();
      },
    });

  const { error, isPending, mutate } = useMutation({
    mutationFn: async () => {
      if (!clientId) {
        return;
      }
      const response = await fetch(
        `${backendUrl}/api/release-notes/${orgId}/${releaseNoteId}/like`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            release_note_id: releaseNoteId,
            client_id: clientId,
          }),
        },
      );
      if (!response.ok) {
        throw new Error("Failed to toggle like");
      }
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["release-note-like", releaseNoteId, clientId],
      });
    },
  });

  return {
    toggleLike: mutate,
    isPending,
    error,
    isLiked: likeState?.is_liked ?? false,
    isLoadingLikeState,
  };
}
