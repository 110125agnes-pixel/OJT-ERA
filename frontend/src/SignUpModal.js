import React, { useState } from 'react';
import './Login.css';

function SignUpModal({ isOpen, onClose, onSignUp }) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const validatePassword = (pwd) => {
    const minLength = /.{8,}/;
    const hasUpper = /[A-Z]/;
    const hasNumber = /[0-9]/;
    return minLength.test(pwd) && hasUpper.test(pwd) && hasNumber.test(pwd);
  };

  const validateEmail = (email) => {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const handleSignUp = async (e) => {
    e.preventDefault();
    if (!username || !email || !password || !confirmPassword) {
      setError('All fields are required.');
      return;
    }
    if (!validateEmail(email)) {
      setError('Please enter a valid email address.');
      return;
    }
    if (!validatePassword(password)) {
      setError('Password must be at least 8 characters, include an uppercase letter and a number.');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match.');
      return;
    }
    const signUpResult = await onSignUp(username.trim(), email.trim(), password);
    if (!signUpResult?.success) {
      setError(signUpResult?.message || 'Unable to create account.');
      return;
    }
    setError('');
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h3>Sign Up</h3>
          <button className="close-btn" onClick={onClose}>&times;</button>
        </div>
        <form onSubmit={handleSignUp} className="signup-form">
          <div className="form-group">
            <label htmlFor="signup-username">Username</label>
            <input
              type="text"
              id="signup-username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              autoComplete="username"
            />
          </div>
          <div className="form-group">
            <label htmlFor="signup-email">Email</label>
            <input
              type="email"
              id="signup-email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              autoComplete="email"
            />
          </div>
          <div className="form-group">
            <label htmlFor="signup-password">Password</label>
            <div className="password-input-row">
              <input
                type={showPassword ? 'text' : 'password'}
                id="signup-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                autoComplete="new-password"
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
          <div className="form-group">
            <label htmlFor="signup-confirm-password">Confirm Password</label>
            <div className="password-input-row">
              <input
                type={showConfirmPassword ? 'text' : 'password'}
                id="signup-confirm-password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Confirm your password"
                autoComplete="new-password"
                className="password-input"
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword((prev) => !prev)}
                className="toggle-password-btn"
              >
                {showConfirmPassword ? 'Hide' : 'Show'}
              </button>
            </div>
          </div>
          {error && <div className="error-message">{error}</div>}
          <button type="submit" className="signup-button">Create Account</button>
        </form>
      </div>
    </div>
  );
}

export default SignUpModal;
