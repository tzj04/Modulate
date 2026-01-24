import { Link } from "react-router-dom";
import { useAuth } from "@/auth/useAuth";
import { useState, useRef, useEffect } from "react";

export const Navbar = () => {
  const { user, logout } = useAuth();
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setShowMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link to="/" className="navbar-logo">
          Modulate
        </Link>

        <div className="navbar-right">
          {user ? (
            <div className="user-profile">
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
