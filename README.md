# Full Stack Application: React + Go + MySQL

A modern full-stack web application with React frontend, Go (Golang) backend, and MySQL database.

## 🏗️ Project Structure

```
NEW ERA/
├── frontend/          # React application
│   ├── public/
│   ├── src/
│   ├── package.json
│   └── Dockerfile
├── backend/           # Go API server
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── database/          # Database scripts
│   └── init.sql
└── docker-compose.yml # Docker orchestration
```

## 🚀 Features

- **Frontend**: React 18 with hooks, Axios for API calls, responsive UI
- **Backend**: Go REST API with Gorilla Mux router, CORS support
- **Database**: MySQL 8.0 with automatic schema initialization
- **Docker**: Full containerization with Docker Compose
- **CRUD Operations**: Create, Read, and Delete items

## 📋 Prerequisites

- [Node.js](https://nodejs.org/) (v18 or higher)
- [Go](https://golang.org/) (v1.21 or higher)
- SQLite (included with Go driver, no separate installation needed)

## 🚀 Getting Started

### 1. Setup Backend (Go)

```powershell
cd backend

# Install dependencies
go mod download
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
go get github.com/rs/cors

# Create .env file (optional)
copy .env.example .env
# Edit .env with your database credentials

# Run the server
go run main.go
```

The backend will start on http://localhost:8080 and create `app.db` in the backend folder.

### 2. Setup Frontend (React)

```powershell
cd frontend

# Install dependencies
npm install
npm install axios

# Start the development server
npm start
```

The frontend will start on http://localhost:3000

## 🔧 Configuration

### Backend Environment Variables (Optional)

Create a `.env` file in the `backend/` directory to customize:

```env
DB_PATH=./app.db
PORT=8080
```

### Frontend Configuration

The frontend is configured to proxy API requests to `http://localhost:8080` (see `package.json`).

## 📡 API Endpoints

| Method | Endpoint          | Description          |
|--------|-------------------|----------------------|
| GET    | /api/health       | Health check         |
| GET    | /api/items        | Get all items        |
| POST   | /api/items        | Create a new item    |
| DELETE | /api/items/{id}   | Delete an item by ID |

### Example API Requests

**Get all items:**
```bash
curl http://localhost:8080/api/items
```

**Create an item:**
```bash
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"New Item"}'
```

**Delete an item:**
```bash
curl -X DELETE http://localhost:8080/api/items/1
```

## 🗄️ Database Schema

```sql
CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🏃‍♂️ Development

### Frontend Development

```powershell
cd frontend
npm start          # Start development server
npm test           # Run tests
npm run build      # Build for production
```

### Backend Development

```powershell
cd backend
go run main.go     # Run with hot reload (use air or fresh)
go build           # Build binary
go test ./...      # Run tests
```

### Database Management

```powershell
# Backup database
### Database Management

The SQLite database file is `backend/app.db`. You can:
- **Backup**: Simply copy the `app.db` file
- **View data**: Use [DB Browser for SQLite](https://sqlitebrowser.org/)
- **Reset**: Delete `app.db` and restart the backend Port Already in Use
If ports 3000, 8080, or 3306 are already in use:
- Stop other services using these ports
- Change the PORT in backend `.env` or frontend package.json proxy

### Database Connection Issues
- Ensure MySQL service is running
- Check database credentials in backend `.env` file
- Verify database `appdb` exists

### Frontend Can't Connect to Backend
- Verify backend is running on port 8080
### Database Connection Issues
- The SQLite database file `app.db` is created automatically
- Check file permissions in the backend directory
- Delete `app.db` and restart if corrupted

### Build frontend for production
```powershell
cd frontend
npm run build
```

### Build backend binary
```powershell
cd backend
go build -o backend.exe main.go
```

## 🔐 Security Notes

For production deployments:
## 🔐 Security Notes

For production deployments:
- Store `app.db` in a secure location with proper permissions
- Use environment variables for sensitive data
- Enable HTTPS/TLS
- Implement authentication and authorization
- Update CORS origins to match your domain
- Use prepared statements (already implemented)
- Consider using PostgreSQL or MySQL for production
- **Frontend**: React 18, Axios, CSS3
- **Backend**: Go 1.21, Gorilla Mux, MySQL Driver
- **Frontend**: React 18, Axios, CSS3
- **Backend**: Go 1.21, Gorilla Mux, SQLite Driver
- **Database**: SQLite 3

This project is open source and available under the [MIT License](LICENSE).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Support

For issues and questions, please open an issue in the repository.

---

**Happy Coding! 🚀**
