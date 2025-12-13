# Quick Start Guide - Surgery Profiling

## Prerequisites Check

✅ Node.js v18.15.0 - Installed
✅ npm 8.5.5 - Installed
✅ React dependencies - Installed
⚠️ Go - **NOT INSTALLED** (Required for backend)

## Setup Instructions

### 1. Install Go (Required)

Since Go is not yet installed on your system:

1. **Download Go:**
   - Visit: https://go.dev/dl/
   - Download: Windows installer (`.msi`) - Latest stable version
   - Run the installer (typically installs to `C:\Program Files\Go`)

2. **Verify Installation:**
   After installation, **restart VS Code** or open a new PowerShell terminal, then run:
   ```powershell
   go version
   ```
   You should see something like: `go version go1.21.x windows/amd64`

### 2. Install Go Dependencies

Once Go is installed:

```powershell
cd backend
go mod download
go mod tidy
```

This will install:
- Gorilla Mux (HTTP router)
- CORS middleware
- SQLite driver

### 3. Start the Backend Server

```powershell
cd backend
go run main.go
```

Expected output:
```
Successfully connected to SQLite database!
Table 'items' ready
Table 'surgeries' ready
Table 'inventory' ready
Server starting on port 8080...
```

The backend will be running at: **http://localhost:8080**

### 4. Start the Frontend (New Terminal)

Open a new terminal (keep the backend running):

```powershell
cd frontend
npm start
```

The React app will open automatically at: **http://localhost:3000**

## Using the Surgery Feature

1. **Login** to the application
2. Click **🏥 Surgery** in the navigation bar
3. Schedule a new surgery by filling out the form
4. View, edit, or delete surgery records

## Troubleshooting

### Backend won't start?
- Check if Go is installed: `go version`
- Check if port 8080 is available
- Check database file permissions

### Frontend won't connect to backend?
- Ensure backend is running on port 8080
- Check the proxy setting in `frontend/package.json`
- Clear browser cache

### Database errors?
- The app will auto-create `app.db` in the backend folder
- Delete `app.db` to reset the database

## File Structure

```
OJT-ERA/
├── backend/
│   ├── main.go                 # Server entry point
│   ├── go.mod                  # Go dependencies
│   ├── models/
│   │   ├── item.go            # Employee model
│   │   └── surgery.go         # Surgery model ✨ NEW
│   ├── controllers/
│   │   ├── item_controller.go # Employee controller
│   │   └── surgery_controller.go ✨ NEW
│   └── routes/
│       └── routes.go          # API routes (updated)
├── frontend/
│   └── src/
│       ├── SurgeryProfiling.js ✨ NEW
│       ├── SurgeryProfiling.css ✨ NEW
│       ├── Dashboard.js       # (updated)
│       └── services/
│           └── api.js         # (updated with surgery service)
└── database/
    └── init.sql               # Database schema (updated)
```

## Next Steps After Setup

1. Test the surgery scheduling feature
2. Add some sample surgery records
3. Test the search and filter functionality
4. Explore edit and delete operations

Need help? Check [SURGERY_FEATURE_README.md](SURGERY_FEATURE_README.md) for detailed documentation.
