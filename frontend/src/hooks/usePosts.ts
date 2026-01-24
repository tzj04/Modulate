import { useState, useEffect } from "react";
import { postApi } from "../api/posts";
import { Post } from "../types/post";

export const usePosts = (moduleId: number | undefined) => {
  // 1. Initialize as an empty array [] so .map() doesn't break
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // If no moduleId is provided, stop loading and return
    if (moduleId === undefined || isNaN(moduleId)) {
      setLoading(false);
      return;
    }

    const fetchPosts = async () => {
      try {
        setLoading(true);
        // 2. Use the correct API method that returns Post[]
        const data = await postApi.getByModule(moduleId);
        setPosts(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || "Failed to load posts");
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, [moduleId]);

  return { posts, loading, error };
};
