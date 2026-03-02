import React, { useEffect, useState } from "react";
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

const REGISTERED_USER_STORAGE_KEY = "registeredUser";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [registeredUser, setRegisteredUser] = useState(null);

  useEffect(() => {
    const storedUser = localStorage.getItem(REGISTERED_USER_STORAGE_KEY);
    if (!storedUser) {
      return;
    }

    try {
      const parsedUser = JSON.parse(storedUser);
      if (parsedUser?.username && parsedUser?.password) {
        setRegisteredUser(parsedUser);
      }
    } catch (error) {
      localStorage.removeItem(REGISTERED_USER_STORAGE_KEY);
    }
  }, []);

  useEffect(() => {
    if (!registeredUser) {
      return;
    }

    localStorage.setItem(
      REGISTERED_USER_STORAGE_KEY,
      JSON.stringify(registeredUser)
    );
  }, [registeredUser]);

  const handleLogin = (username, password) => {
    if (!registeredUser) {
      return {
        success: false,
        message: "No account found. Please sign up first.",
      };
    }

    const isValidUser =
      username === registeredUser.username && password === registeredUser.password;

    if (!isValidUser) {
      return {
        success: false,
        message: "Invalid username or password.",
      };
    }

    setIsLoggedIn(true);
    return { success: true };
  };

  const handleSignUp = async (username, email, password) => {
    try {
      await authService.signup({ username, email, password });
      setRegisteredUser({ username, email, password });
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
