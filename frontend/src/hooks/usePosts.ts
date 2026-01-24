import { useState, useEffect, useCallback } from "react";
import { postApi } from "../api/posts";
import { Post } from "../types/post";

export const usePosts = (moduleId: number | undefined) => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Define fetchPosts outside so it can be reused/returned
  const fetchPosts = useCallback(async () => {
    if (moduleId === undefined || isNaN(moduleId)) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const data = await postApi.getByModule(moduleId);
      setPosts(data || []); // Ensure there's always an array
      setError(null);
    } catch (err: any) {
      setError(err.message || "Failed to load posts");
    } finally {
      setLoading(false);
    }
  }, [moduleId]);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  return { posts, loading, error, refresh: fetchPosts };
};
