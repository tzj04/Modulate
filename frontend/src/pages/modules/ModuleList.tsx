import React, {
  useState,
  useMemo,
  memo,
  useDeferredValue,
  useEffect,
} from "react";
import { Link } from "react-router-dom";
import { useModules } from "../../hooks/useModules";

const SearchInput = memo(
  ({ onChange }: { onChange: (val: string) => void }) => {
    const [localValue, setLocalValue] = useState("");

    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
      const val = e.target.value;
      setLocalValue(val);
      onChange(val); // Notifies parent of the change
    };

    return (
      <input
        type="text"
        className="module-list-search"
        placeholder="Search by module code or title..."
        value={localValue}
        onChange={handleInput}
      />
    );
  },
);

/**
 * Sub-component: ModuleCard
 * Memoized to prevent unnecessary re-renders when other
 * parts of the list change.
 */
const ModuleCard = memo(({ mod }: { mod: any }) => (
  <Link to={`/modules/${mod.id}`} className="module-link-card">
    <span className="module-code">{mod.code}:</span> {mod.title}
  </Link>
));

export const ModuleList = () => {
  const { modules, loading, error } = useModules();
  const [searchQuery, setSearchQuery] = useState("");
  const [visibleCount, setVisibleCount] = useState(50);

  // High-priority typing vs Low-priority list filtering
  const deferredQuery = useDeferredValue(searchQuery);

  // Reset pagination when the user searches for something new
  useEffect(() => {
    setVisibleCount(50);
  }, [deferredQuery]);

  // Expensive filtering logic is memoized
  const filteredModules = useMemo(() => {
    const query = deferredQuery.toLowerCase().trim();
    if (!query) return modules;

    return modules.filter(
      (m) =>
        m.code.toLowerCase().includes(query) ||
        m.title.toLowerCase().includes(query),
    );
  }, [modules, deferredQuery]);

  const handleLoadMore = () => {
    setVisibleCount((prev) => prev + 50);
  };

  if (loading) return <div className="status-message">Fetching modules...</div>;
  if (error) return <div className="status-message error">Error: {error}</div>;

  return (
    <div className="module-list-page">
      <div className="module-search-container">
        <SearchInput onChange={setSearchQuery} />
      </div>

      <div className="module-grid">
        {filteredModules.slice(0, visibleCount).map((mod) => (
          <ModuleCard key={mod.id} mod={mod} />
        ))}

        {filteredModules.length === 0 && (
          <div className="no-results-message">
            No modules found matching "{deferredQuery}"
          </div>
        )}

        {filteredModules.length > visibleCount && (
          <div className="load-more-container">
            <button onClick={handleLoadMore} className="load-more-btn">
              Load More ({filteredModules.length - visibleCount} remaining)
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
