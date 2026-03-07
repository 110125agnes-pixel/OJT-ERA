import React, { useState } from 'react';
import SignUpModal from './SignUpModal';
import './Login.css';


function Login({ onLogin, onSignUp }) {
const [username, setUsername] = useState('');
const [password, setPassword] = useState('');
const [error, setError] = useState('');
const [showSignUp, setShowSignUp] = useState(false);
const [showPassword, setShowPassword] = useState(false);

const validatePassword = (pwd) => {
// Example conditions: min 8 chars, 1 uppercase, 1 number
const minLength = /.{8,}/;
const hasUpper = /[A-Z]/;
const hasNumber = /[0-9]/;
return minLength.test(pwd) && hasUpper.test(pwd) && hasNumber.test(pwd);
};

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!username.trim() || !password) {
      setError('Username and password are required.');
      return;
    }
    if (!validatePassword(password)) {
      setError('Password must be at least 8 characters, include an uppercase letter and a number.');
      return;
    }
    const loginResult = onLogin(username.trim(), password);
    if (!loginResult?.success) {
      setError(loginResult?.message || 'Login failed.');
      return;
    }
    setError('');
  };

const handleSignUpClick = (e) => {
e.preventDefault();
setShowSignUp(true);
};

return (
<div className="login-container">
<div className="login-box">
<div className="login-header">
<h1>Employee Portal</h1>
<p>Sign in to continue</p>
</div>
<form onSubmit={handleSubmit} className="login-form">
<div className="form-group">
<label htmlFor="username">Username</label>
<input
type="text"
id="username"
value={username}
onChange={(e) => setUsername(e.target.value)}
placeholder="Enter your username"
autoComplete="username"
/>
</div>
<div className="form-group">
<label htmlFor="password">Password</label>
<div className="password-input-row">
<input
type={showPassword ? 'text' : 'password'}
id="password"
value={password}
onChange={(e) => setPassword(e.target.value)}
placeholder="Enter your password"
autoComplete="current-password"
className="password-input"
/>
<button
type="button"
onClick={() => setShowPassword((prev) => !prev)}
className="toggle-password-btn"
>
{showPassword ? 'Hide' : 'Show'}
</button>
</div>
</div>
{error && <div className="error-message">{error}</div>}
<div className="auth-actions">
<button type="submit" className="login-button">
Login
</button>
<button type="button" className="signup-button secondary" onClick={handleSignUpClick}>
Sign Up
</button>
</div>
</form>
<div className="login-footer">
<p>For authorized employees only</p>
</div>
</div>
<SignUpModal
isOpen={showSignUp}
onClose={() => setShowSignUp(false)}
onSignUp={onSignUp}
/>
</div>
);
}

export default Login;