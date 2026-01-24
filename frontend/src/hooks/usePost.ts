// src/hooks/usePost.ts
import { useState, useEffect } from "react";
import { postApi } from "../api/posts";
import { Post } from "../types/post";

export const usePost = (postId: number | undefined) => {
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!postId || isNaN(postId)) {
      setLoading(false);
      return;
    }

    const fetchPost = async () => {
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
    };

    fetchPost();
  }, [postId]);

  return { post, loading, error };
};
