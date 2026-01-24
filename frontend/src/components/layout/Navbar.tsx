import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "@/auth/useAuth";
import { useState, useRef, useEffect } from "react";
import { moduleApi } from "@/api/modules";
import { Module } from "@/types/module";

export const Navbar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [showMenu, setShowMenu] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchError, setSearchError] = useState("");
  const [modules, setModules] = useState<Module[]>([]);
  const [filteredModules, setFilteredModules] = useState<Module[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLFormElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  // Fetch all modules on component mount
  useEffect(() => {
    const fetchModules = async () => {
      try {
        const data = await moduleApi.getAll();
        setModules(data);
      } catch (err) {
        console.error("Failed to fetch modules:", err);
      }
    };
    fetchModules();
  }, []);

  // Focus search input when expanded
  useEffect(() => {
    if (isSearchExpanded && searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [isSearchExpanded]);

  // Filter modules as user types
  useEffect(() => {
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      const filtered = modules.filter(
        (m) =>
          m.code.toLowerCase().includes(query) ||
          m.title.toLowerCase().includes(query),
      );
      setFilteredModules(filtered.slice(0, 8)); // Limit to 8 suggestions
      setShowSuggestions(true);
    } else {
      setFilteredModules([]);
      setShowSuggestions(false);
    }
  }, [searchQuery, modules]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      const isSearchIconClick = target.closest(".search-icon-btn");
      const isSearchFormClick = target.closest(".navbar-search-form-inline");

      // Close search bar if clicking outside search icon and search form
      if (!isSearchIconClick && !isSearchFormClick && isSearchExpanded) {
        setIsSearchExpanded(false);
        setShowSuggestions(false);
        setSearchQuery("");
      }

      if (menuRef.current && !menuRef.current.contains(target)) {
        setShowMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isSearchExpanded]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchError("");

    if (!searchQuery.trim()) return;

    const found = modules.find(
      (m) => m.code.toLowerCase() === searchQuery.trim().toLowerCase(),
    );

    if (!found) {
      setSearchError("Module not found");
      setTimeout(() => setSearchError(""), 3000);
      return;
    }

    navigate(`/modules/${found.id}`);
    setSearchQuery("");
    setShowSuggestions(false);
  };

  const handleModuleSelect = (module: Module) => {
    navigate(`/modules/${module.id}`);
    setSearchQuery("");
    setShowSuggestions(false);
  };

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          <img
            src="/images/home_icon.png"
            alt="Home"
            className="navbar-logo-icon"
          />
          <span>Modulate</span>
        </Link>

        <div className="navbar-right">
          {user ? (
            <div className="user-profile">
              {/* Search Icon */}
              <button
                className="search-icon-btn"
                onClick={() => setIsSearchExpanded(!isSearchExpanded)}
                title="Search modules"
              >
                <img
                  src="images/search_icon.webp"
                  alt="Search"
                  className="search-icon-img"
                />
              </button>

              {/* Expandable Search Bar */}
              {isSearchExpanded && (
                <form
                  className="navbar-search-form-inline"
                  onSubmit={handleSearchSubmit}
                  ref={searchRef}
                >
                  <div className="navbar-search-wrapper">
                    <input
                      type="text"
                      className="navbar-search-input-inline"
                      placeholder="Search modules..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      onFocus={() => searchQuery && setShowSuggestions(true)}
                      ref={searchInputRef}
                    />
                    {searchError && (
                      <span className="navbar-search-error">{searchError}</span>
                    )}
                    {showSuggestions && filteredModules.length > 0 && (
                      <div className="navbar-search-suggestions">
                        {filteredModules.map((module) => (
                          <button
                            key={module.id}
                            type="button"
                            className="navbar-search-suggestion-item"
                            onClick={() => handleModuleSelect(module)}
                          >
                            <strong>{module.code}</strong>
                            <span className="module-title">{module.title}</span>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                </form>
              )}

              <span className="navbar-separator">|</span>

              <span className="user-greeting">
                Hi, <strong>{user.username}</strong>
              </span>
              <div className="navbar-menu-wrapper" ref={menuRef}>
                <button
                  className="navbar-kebab-btn"
                  onClick={() => setShowMenu(!showMenu)}
                >
                  ⋮
                </button>
                {showMenu && (
                  <div className="navbar-dropdown">
                    <button onClick={logout} className="navbar-logout-btn">
                      <img
                        src="/images/logout.png"
                        alt="logout"
                        className="logout-icon-img"
                      />
                      Logout
                    </button>
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="auth-links">
              <Link to="/login" className="login-link">
                Login
              </Link>
              <Link to="/register" className="register-button">
                Sign Up
              </Link>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
};
