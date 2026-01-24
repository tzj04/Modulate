import { useState, useEffect, useCallback } from "react";
import { postApi } from "../api/posts";
import { Post } from "../types/post";

export const usePost = (postId: number | undefined) => {
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Define fetchPost outside so it can be returned
  const fetchPost = useCallback(async () => {
    if (!postId || isNaN(postId)) {
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const data = await postApi.getPost(postId);
      setPost(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || "Failed to load post");
    } finally {
      setLoading(false);
    }
  }, [postId]);

  useEffect(() => {
    fetchPost();
  }, [fetchPost]);

  return { post, loading, error, refresh: fetchPost };
};
