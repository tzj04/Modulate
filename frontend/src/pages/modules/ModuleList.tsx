// src/pages/modules/ModuleList.tsx
import { Link } from "react-router-dom";
import { useModules } from "../../hooks/useModules";

export const ModuleList = () => {
  const { modules, loading, error } = useModules();

  if (loading)
    return <div style={{ padding: "20px" }}>Fetching modules...</div>;
  if (error)
    return <div style={{ padding: "20px", color: "red" }}>Error: {error}</div>;

  return (
    <div style={{ padding: "20px" }}>
      <h1>Available Modules</h1>
      <div style={{ display: "grid", gap: "15px" }}>
        {modules.map((mod) => (
          <Link
            key={mod.id}
            to={`/modules/${mod.id}`}
            style={{
              padding: "15px",
              border: "1px solid #003D7C",
              borderRadius: "5px",
              textDecoration: "none",
              color: "#003D7C",
              fontWeight: "bold",
              display: "block",
            }}
          >
            {mod.code}: {mod.title}
          </Link>
        ))}
      </div>
    </div>
  );
};
