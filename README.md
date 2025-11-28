# Clarkson - Vehicle Expense Tracker

Clarkson is a comprehensive, self-hosted vehicle expense and maintenance tracker designed for multi-vehicle households. It combines the efficiency of Hammond with advanced maintenance reminder features inspired by MyCar, all packaged in a lightweight Docker container optimized for Unraid deployment.

## Features

### Multi-Vehicle Management
- Track multiple vehicles simultaneously
- Support for different unit systems (miles/km)
- Vehicle sharing between users
- Detailed vehicle profiles

### Expense & Fuel Tracking
- Log fuel entries with price, gallons, and odometer
- Track expenses by category (maintenance, insurance, parking, etc.)
- Attach photos and receipts to entries
- Location tracking for fuel purchases
- Monthly trends and statistics

### Maintenance Reminders
- Schedule reminders by mileage intervals (e.g., oil change every 5,000 miles)
- Schedule reminders by time intervals (e.g., tire rotation every 1 year)
- Automatic notifications when service is due or overdue
- Mark reminders as complete with new mileage/date
- Color-coded alerts (overdue: red, due soon: yellow)

### Reporting & Analytics
- Vehicle-specific reports with cost breakdowns
- Multi-vehicle comparison reports
- Fuel economy trends and analysis
- Expense category breakdowns
- Overall statistics and summaries
- CSV and JSON export

### Data Management
- Import from Hammond vehicle tracker
- Import from Fuelly CSV exports
- Backup and restore functionality
- Search and filter entries

### User Management
- Multi-user support with authentication
- Share vehicles with other users
- Role-based access (admin/user)
- User preferences (units, currency)

### Modern UI
- Responsive design (mobile-first)
- Dark mode support
- Real-time notifications
- PWA support for offline access
- Intuitive navigation

## System Requirements

- Docker installed
- 500MB free disk space
- 256MB RAM (typical)
- Modern web browser

## Quick Start (Unraid)

### 1. Prepare Directories

\`\`\`bash
mkdir -p /mnt/user/appdata/clarkson/config
mkdir -p /mnt/user/appdata/clarkson/assets
chmod 755 /mnt/user/appdata/clarkson/config
chmod 755 /mnt/user/appdata/clarkson/assets
\`\`\`

### 2. Generate JWT Secret

\`\`\`bash
JWT_SECRET=$(openssl rand -base64 32)
echo $JWT_SECRET
\`\`\`

### 3. Deploy Container

Via Unraid Docker UI:
- Go to Docker tab
- Add Container
- Repository: `clarkson:latest`
- Container Name: `clarkson`
- Port Mapping: 3000:3000 (see [Port Configuration](#port-configuration) to customize)
- Volumes:
  - `/config` → `/mnt/user/appdata/clarkson/config`
  - `/assets` → `/mnt/user/appdata/clarkson/assets`
- Environment Variables:
  - `JWT_SECRET`: (paste generated secret)
  - `PORT`: 3000 (internal port, usually keep as 3000)
  - `PORT_MAPPING`: 3000 (external port visible to browser, change as needed)

### 4. Access Application

Open browser to `http://unraid-ip:3000` (or your configured external port)

Create admin account on first visit.

## Port Configuration

Clarkson allows flexible port configuration to avoid conflicts or use custom ports like Unraid's 43535.

### How Ports Work

- **PORT**: Internal port used by the application (default: 3000)
- **PORT_MAPPING**: External port visible from your browser (docker-compose only)

When deploying on Unraid via Docker UI, the port mapping is configured in the UI itself.

### Examples

#### Default Setup (Port 3000)
\`\`\`env
PORT=3000
PORT_MAPPING=3000
\`\`\`
Access: `http://unraid-ip:3000`

#### Unraid Custom Port (43535)
\`\`\`env
PORT=3000
PORT_MAPPING=43535
\`\`\`
Access: `http://unraid-ip:43535`

#### Via docker-compose
Create `.env` file:
\`\`\`env
PORT=3000
PORT_MAPPING=43535
JWT_SECRET=your-generated-secret
\`\`\`

Then:
\`\`\`bash
docker-compose up -d
\`\`\`

#### Via Unraid Docker UI
Set Port Mapping to: `43535:3000`
- Left side (43535): External port - what you access in browser
- Right side (3000): Internal port - keep this as 3000

### Advanced Port Configuration

For different internal and external ports:
\`\`\`env
PORT=8080
PORT_MAPPING=43535
\`\`\`
This routes `http://unraid-ip:43535` to internal port 8080 (usually not needed, stick with PORT=3000).

## Architecture

### Backend (Go/Gin)
- Lightweight REST API
- SQLite database (embedded)
- JWT authentication
- GORM ORM
- Responsive to 3000 requests/sec typical

### Frontend (Vue 3)
- Single Page Application
- Pinia state management
- Tailwind CSS styling
- Responsive and accessible
- PWA support

### Database (SQLite)
- Single file: `/config/clarkson.db`
- Zero configuration
- ACID compliance
- Perfect for single-user/small team use

## File Structure

\`\`\`
clarkson/
├── backend/
│   ├── main.go           # Entry point
│   ├── models.go         # Database models
│   ├── handlers.go       # HTTP handlers
│   ├── routes.go         # Route setup
│   ├── notifications.go  # Reminder logic
│   ├── uploads.go        # File handling
│   ├── imports.go        # Import logic
│   ├── reports.go        # Report generation
│   └── go.mod            # Dependencies
│
├── frontend/
│   ├── src/
│   │   ├── App.vue       # Root component
│   │   ├── main.js       # Entry point
│   │   ├── router.js     # Routing
│   │   ├── style.css     # Global styles
│   │   ├── views/
│   │   │   ├── Dashboard.vue
│   │   │   ├── VehicleDetail.vue
│   │   │   ├── Reports.vue
│   │   │   ├── Settings.vue
│   │   │   ├── Login.vue
│   │   │   └── Register.vue
│   │   ├── components/
│   │   │   ├── VehicleCard.vue
│   │   │   ├── FileUpload.vue
│   │   │   ├── ImportData.vue
│   │   │   ├── NotificationBell.vue
│   │   │   ├── NotificationCenter.vue
│   │   │   └── MaintenanceAlerts.vue
│   │   └── stores/
│   │       ├── auth.js
│   │       ├── vehicles.js
│   │       ├── fuel.js
│   │       └── reminders.js
│   ├── package.json
│   ├── vite.config.js
│   └── index.html
│
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # Deployment config
├── .env.example         # Environment variables template
└── README.md            # This file
\`\`\`

## API Documentation

### Authentication

\`\`\`
POST /api/auth/register
POST /api/auth/login
\`\`\`

### Vehicles

\`\`\`
GET  /api/vehicles                    # List user's vehicles
POST /api/vehicles                    # Create vehicle
GET  /api/vehicles/:id                # Get vehicle details
PUT  /api/vehicles/:id                # Update vehicle
DELETE /api/vehicles/:id              # Delete vehicle
POST /api/vehicles/:id/share          # Share vehicle with user
\`\`\`

### Fuel Entries

\`\`\`
GET  /api/vehicles/:id/fuel           # List fuel entries
POST /api/vehicles/:id/fuel           # Add fuel entry
GET  /api/vehicles/:id/fuel-stats     # Get fuel statistics
PUT  /api/fuel/:id                    # Update fuel entry
DELETE /api/fuel/:id                  # Delete fuel entry
\`\`\`

### Expenses

\`\`\`
GET  /api/vehicles/:id/expenses       # List expenses
POST /api/vehicles/:id/expenses       # Add expense
GET  /api/vehicles/:id/expense-stats  # Get expense statistics
PUT  /api/expenses/:id                # Update expense
DELETE /api/expenses/:id              # Delete expense
\`\`\`

### Maintenance Reminders

\`\`\`
GET  /api/vehicles/:id/reminders      # List reminders
POST /api/vehicles/:id/reminders      # Create reminder
PUT  /api/reminders/:id               # Update reminder
DELETE /api/reminders/:id             # Delete reminder
POST /api/reminders/:id/complete      # Mark complete
GET  /api/reminders/check             # Check all reminders
GET  /api/reminders/overdue           # List overdue reminders
\`\`\`

### Reports & Export

\`\`\`
GET  /api/vehicles/:id/report         # Detailed vehicle report
GET  /api/report/overall              # Overall statistics
GET  /api/report/comparison           # Compare all vehicles
GET  /api/search?q=query              # Search entries
GET  /api/export/csv                  # Export CSV
GET  /api/export/json                 # Export JSON
\`\`\`

### File Management

\`\`\`
POST /api/upload                      # Upload file
GET  /api/download/:id                # Download file
DELETE /api/attachments/:id           # Delete attachment
GET  /api/attachments?type=&entry_id= # List attachments
\`\`\`

### Notifications

\`\`\`
GET  /api/notifications               # List notifications
GET  /api/notifications/summary       # Notification summary
POST /api/notifications/:id/read      # Mark as read
POST /api/notifications/:id/dismiss   # Dismiss notification
\`\`\`

## Configuration

### Environment Variables

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `JWT_SECRET` | (generated) | Yes | JWT signing secret (change in production!) |
| `PORT` | 3000 | No | API server port |
| `CONFIG_PATH` | /config | No | SQLite database directory |
| `ASSETS_PATH` | /assets | No | File uploads directory |

### Generate Strong JWT Secret

\`\`\`bash
openssl rand -base64 32
\`\`\`

## Backup & Recovery

### Backup

\`\`\`bash
# Backup database and assets
tar -czf clarkson-backup-$(date +%Y%m%d).tar.gz \
  /mnt/user/appdata/clarkson/config/clarkson.db \
  /mnt/user/appdata/clarkson/assets/
\`\`\`

### Restore

\`\`\`bash
# Extract backup
tar -xzf clarkson-backup-20240101.tar.gz -C /

# Restart container
docker-compose restart
\`\`\`

## Performance

- **Container Size**: ~180MB
- **Memory Usage**: 256MB typical
- **Startup Time**: <5 seconds
- **Database**: SQLite (embedded, zero-config)
- **Requests**: ~1000 requests/second typical

## Security Considerations

1. **Change JWT_SECRET** - Critical! Generate a strong secret.
2. **Use HTTPS** - Deploy behind reverse proxy (nginx) with SSL
3. **Set strong passwords** - On first login
4. **Regular backups** - Store backups securely
5. **Update regularly** - Monitor for security updates
6. **Restrict database access** - Keep `/config` permissions restrictive

## Troubleshooting

### Container won't start
\`\`\`bash
docker logs clarkson
\`\`\`

### Database locked
\`\`\`bash
docker-compose restart clarkson
\`\`\`

### Port already in use
Change `3000:3000` to `3001:3000` in docker-compose.yml

### High memory usage
Check for large attachments. Compress old photos.

### Slow queries
Add more retention policy for old entries. Consider archiving.

## License

Clarkson is licensed under the MIT License. See LICENSE file for details.

## Credits

Inspired by:
- Hammond (akhilrex/hammond) - Vehicle expense tracking architecture
- MyCar Android app - Maintenance reminder system
- Fuelly - Fuel economy tracking

## Support & Contributions

- GitHub Issues: Report bugs and feature requests
- Discussions: Community Q&A and ideas
- Pull Requests: Contributions welcome

## Roadmap

- PDF report generation
- Cloud sync (optional)
- Mobile app
- Advanced forecasting
- Integration with car APIs
- OAuth authentication
- Multi-language support

---

**Enjoy tracking your vehicles with Clarkson!**
