import React, { useState } from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import "./App.css";
import Login from "./Login";
import Dashboard from "./Dashboard";
import { authService } from "./services/api";

// --- FIXED IMPORT ---
// The folder is 'PatientView' and the file inside is 'PatientView.js'
import PatientView from "./PatientView/PatientView";

const getApiErrorMessage = (error, fallbackMessage) => {
  if (typeof error?.response?.data === "string" && error.response.data.trim()) {
    return error.response.data;
  }

  if (error?.response?.data?.message) {
    return error.response.data.message;
  }

  return fallbackMessage;
};

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleLogin = async (username, password) => {
    try {
      const response = await authService.login({ username, password });

      if (!response?.success) {
        return {
          success: false,
          message: response?.message || "Invalid username or password.",
        };
      }

      setIsLoggedIn(true);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        message: getApiErrorMessage(error, "Invalid username or password."),
      };
    }
  };

  const handleSignUp = async (username, email, password) => {
    try {
      await authService.signup({ username, email, password });
      return { success: true };
    } catch (error) {
      return {
        success: false,
        message: getApiErrorMessage(error, "Unable to create account."),
      };
    }
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
  };

  if (!isLoggedIn) {
    return <Login onLogin={handleLogin} onSignUp={handleSignUp} />;
  }

  return (
    <Router>
      <Routes>
        <Route
          path="/dashboard"
          element={<Dashboard onLogout={handleLogout} />}
        />
        <Route path="/patient/:id" element={<PatientView />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
