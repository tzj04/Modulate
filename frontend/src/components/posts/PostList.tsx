import { useNavigate } from "react-router-dom";

export const PostList = ({ posts }: { posts: any[] }) => {
  const navigate = useNavigate();

  if (posts.length === 0) {
    return (
      <div style={{ textAlign: "center", marginTop: "40px" }}>
        No comments yet!
      </div>
    );
  }

  return (
    <div className="post-feed">
      {posts.map((post) => (
        <div
          key={post.id}
          className="post-card"
          onClick={() => navigate(`/posts/${post.id}`)}
          style={{ cursor: "pointer" }} // Changes mouse to a hand icon
        >
          <div className="post-header">
            <div>
              <span>
                Posted by{" "}
                <span className="post-author">
                  {post.username || "Unknown"}
                </span>
              </span>
            </div>
            <div className="post-date-right">
              {new Date(post.created_at).toLocaleDateString(undefined, {
                month: "short",
                day: "numeric",
                year: "numeric",
              })}
            </div>
          </div>
          <h3 className="post-title">{post.title}</h3>
          {/* Use a snippet here so the card doesn't get too long */}
          <p className="post-content">
            {post.content.length > 150
              ? post.content.substring(0, 150) + "..."
              : post.content}
          </p>
        </div>
      ))}
    </div>
  );
};
