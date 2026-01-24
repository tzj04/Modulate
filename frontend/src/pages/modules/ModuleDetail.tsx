import { useParams, Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { Module } from "../../types/module";
import { usePosts } from "../../hooks/usePosts";
import { moduleApi } from "../../api/modules";
import { PostList } from "../../components/posts/PostList";

export const ModuleDetail = () => {
  const { id } = useParams<{ id: string }>();
  const moduleId = Number(id);

  const [moduleInfo, setModuleInfo] = useState<Module | null>(null);
  const { posts, loading, error } = usePosts(moduleId);

  useEffect(() => {
    if (id) {
      moduleApi
        .getById(id)
        .then(setModuleInfo)
        .catch((err) => console.error("Failed to fetch module details:", err));
    }
  }, [id]);

  if (loading) return <div style={containerStyle}>Loading...</div>;

  return (
    <div style={containerStyle}>
      {/* Module Header */}
      {moduleInfo && (
        <div style={{ marginBottom: "20px" }}>
          <h1 style={{ color: "#003D7C" }}>
            {moduleInfo.code}: {moduleInfo.title}
          </h1>
          <p style={{ color: "#666" }}>{moduleInfo.description}</p>
        </div>
      )}

      <div style={headerActionStyle}>
        <h2>Recent Posts</h2>
        <Link to={`/modules/${id}/create`} style={buttonStyle}>
          Create New Post
        </Link>
      </div>

      <hr style={{ margin: "20px 0", border: "0.5px solid #eee" }} />

      {/* 2. USE the component here instead of mapping manually */}
      {error ? (
        <div style={{ color: "red" }}>{error}</div>
      ) : (
        <PostList posts={posts} />
      )}
    </div>
  );
};

// Styles to keep the code clean
const containerStyle = { maxWidth: "800px", margin: "0 auto", padding: "20px" };
const headerActionStyle = {
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
};
const buttonStyle = {
  background: "#EF7C00",
  color: "white",
  padding: "8px 16px",
  textDecoration: "none",
  borderRadius: "4px",
};
