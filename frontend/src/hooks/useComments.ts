import { useState, useEffect } from "react";
import { Comment } from "../types/comment";
import { commentApi } from "../api/comments";

export const useComments = (postId: number | undefined) => {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchComments = async () => {
    if (!postId) return;
    try {
      const data = await commentApi.getThread(postId);
      setComments(buildCommentTree(data));
    } catch (err) {
      console.error("Failed to load comments", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComments();
  }, [postId]);

  // The logic to nest comments based on parent_comment_id
  const buildCommentTree = (flat: Comment[] | null): Comment[] => {
    if (!flat) return [];

    const map: Record<number, Comment> = {};
    const roots: Comment[] = [];
    const added = new Set<number>();

    // Initialize the map using Number IDs
    flat.forEach((item) => {
      const id = Number(item.id);
      map[id] = { ...item, children: [] };
    });

    // Link children to parents using Number conversion
    flat.forEach((item) => {
      const nodeId = Number(item.id);
      const node = map[nodeId];

      if (item.parent_comment_id !== null) {
        const parentId = Number(item.parent_comment_id);

        if (map[parentId]) {
          map[parentId].children?.push(node);
          added.add(nodeId); // Mark as added to a parent
        } else {
          // if parentId exists but isn't in our map, treat as root
          if (!added.has(nodeId)) {
            roots.push(node);
            added.add(nodeId);
          }
        }
      } else {
        // if truly no parent
        if (!added.has(nodeId)) {
          roots.push(node);
          added.add(nodeId);
        }
      }
    });

    return roots;
  };

  return {
    comments,
    loading,
    refresh: fetchComments,
  };
};
