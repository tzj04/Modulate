import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { userApi } from "../api/users";

export const RegisterPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [label, setLabel] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      // Call Go Backend /auth/register
      await userApi.register({
        username,
        password,
        label: label || undefined, // Send undefined if empty
      });

      // Redirect to login on success
      navigate("/login", {
        state: { message: "Account created! Please log in." },
      });
    } catch (err: any) {
      setError(err.message || "Registration failed. Username might be taken.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2 className="auth-title">Join Modulate</h2>
        <p className="auth-subtitle">
          Create an account to join conversations!
        </p>

        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label className="form-label">Username</label>
            <input
              type="text"
              className="form-input"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="e.g. nus_student_2026"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Label (Optional)</label>
            <input
              type="text"
              className="form-input"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="e.g. Year 2 CS Student"
            />
          </div>

          <div className="form-group">
            <label className="form-label">Password</label>
            <input
              type="password"
              className="form-input"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>

          <button type="submit" className="auth-button" disabled={loading}>
            {loading ? "Creating Account..." : "Create Account"}
          </button>
        </form>

        <div className="auth-footer">
          Already have an account? <Link to="/login">Login here</Link>
        </div>
      </div>
    </div>
  );
};
