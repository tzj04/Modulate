import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import { RequireAuth } from "./auth/RequireAuth";

// Layout & Pages
import { Navbar } from "./components/layout/Navbar";
import { LoginPage } from "./pages/Login";
import { RegisterPage } from "./pages/Register";
import { ModuleList } from "./pages/modules/ModuleList";
import { ModuleDetail } from "./pages/modules/ModuleDetail";
import { CreatePostPage } from "./pages/posts/CreatePost";
import { PostDetail } from "./pages/posts/PostDetail";

function AppContent() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/" element={<ModuleList />} />
        <Route path="/modules/:id" element={<ModuleDetail />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/posts/:postId" element={<PostDetail />} />

        {/* Protected Routes */}
        <Route
          path="/modules/:id/create"
          element={
            <RequireAuth>
              <CreatePostPage />
            </RequireAuth>
          }
        />
      </Routes>
    </>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppContent />
      </BrowserRouter>
    </AuthProvider>
  );
}
