import { Link } from "react-router-dom";
import { useAuth } from "@/auth/useAuth";

export const Navbar = () => {
  const { user, logout } = useAuth();

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
              <button onClick={logout} className="logout-button">
                Logout
              </button>
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
