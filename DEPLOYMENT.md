# Clarkson Deployment Guide for Unraid

## Overview
Clarkson is a lightweight, full-stack vehicle expense tracker optimized for Unraid deployment in a single Docker container.

## Prerequisites
- Unraid 6.9+ with Docker support
- ~500MB free space in `/mnt/user/appdata/`

## Quick Start

### 1. Build the Docker Image
\`\`\`bash
# Clone or download Clarkson repository
cd clarkson

# Build the image
docker build -t clarkson:latest .
\`\`\`

### 2. Create Directories on Unraid
\`\`\`bash
# Via Unraid terminal or console
mkdir -p /mnt/user/appdata/clarkson/config
mkdir -p /mnt/user/appdata/clarkson/assets
chmod 755 /mnt/user/appdata/clarkson/config
chmod 755 /mnt/user/appdata/clarkson/assets
\`\`\`

### 3. Deploy Container
\`\`\`bash
docker-compose -f docker-compose.yml up -d
\`\`\`

Or use Unraid's Docker UI:
1. Go to **Docker** tab
2. Add Container
3. Repository: `clarkson:latest`
4. Container Name: `clarkson`
5. Ports: Port 3000 → 3000
6. Volumes:
   - `/config` → `/mnt/user/appdata/clarkson/config`
   - `/assets` → `/mnt/user/appdata/clarkson/assets`
7. Environment Variables:
   - `PORT`: 3000
   - `JWT_SECRET`: Generate a strong secret (e.g., `openssl rand -base64 32`)

### 4. Access Clarkson
- Web UI: `http://unraid-ip:3000`
- Database: `/mnt/user/appdata/clarkson/config/clarkson.db`
- Attachments: `/mnt/user/appdata/clarkson/assets/`

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 3000 | API server port |
| `JWT_SECRET` | (required) | JWT signing secret - must change in production |
| `CONFIG_PATH` | /config | SQLite database directory |
| `ASSETS_PATH` | /assets | File uploads directory |

### Generate JWT Secret
\`\`\`bash
# Generate strong secret
openssl rand -base64 32

# Output example:
# X8kP9mL2q4vN3zJ7fG5bH6cK8wD4eR0tY2aS5bU9vC1mN3oP7qR9tU2vW4xY6zZ
\`\`\`

## Maintenance

### Backup
\`\`\`bash
# Backup database and assets
tar -czf clarkson-backup-$(date +%Y%m%d).tar.gz \
  /mnt/user/appdata/clarkson/config/clarkson.db \
  /mnt/user/appdata/clarkson/assets/
\`\`\`

### Restore
\`\`\`bash
tar -xzf clarkson-backup-20231128.tar.gz -C /
docker-compose restart
\`\`\`

### Update
\`\`\`bash
cd clarkson
git pull
docker-compose build --no-cache
docker-compose up -d
\`\`\`

## Troubleshooting

### Container won't start
\`\`\`bash
# Check logs
docker logs clarkson

# Verify permissions
ls -la /mnt/user/appdata/clarkson/
# Should be: drwxr-xr-x root root

# Restart
docker-compose restart
\`\`\`

### Database locked
\`\`\`bash
# SQLite connection issue - restart container
docker-compose restart clarkson
\`\`\`

### Port already in use
\`\`\`bash
# Change port in docker-compose.yml
# ports:
#   - "3001:3000"  # Change 3001 to available port
docker-compose restart
\`\`\`

## Performance Tuning

The default configuration is optimized for Unraid:
- **Container Size**: ~180MB
- **Memory**: 256MB typical
- **Startup Time**: <5 seconds
- **SQLite**: Embedded, zero-config

### Optional: Enable SQLite WAL Mode (faster writes)
Edit container environment and add:
\`\`\`
SQLITE_WAL_MODE=1
\`\`\`

## Security Considerations

1. **Change JWT_SECRET** in production (required!)
2. **Use reverse proxy** (nginx) for HTTPS
3. **Set strong admin password** on first login
4. **Regular backups** to `/mnt/user/backup/`
5. **Restrict access** to `/config` directory (contains DB + secrets)

## API Documentation

### Authentication
\`\`\`bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
\`\`\`

### Vehicles
\`\`\`bash
TOKEN="your-jwt-token"

# List vehicles
curl -H "Authorization: $TOKEN" \
  http://localhost:3000/api/vehicles

# Create vehicle
curl -X POST http://localhost:3000/api/vehicles \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "year": 2023,
    "make": "Toyota",
    "model": "Camry",
    "odometer": 45000,
    "mileage_unit": "mi",
    "fuel_type": "Petrol"
  }'
\`\`\`

## Integration with Unraid Ecosystem

### Integrate with Notifiarr (optional)
1. Configure webhook in Clarkson settings
2. Receive maintenance reminder notifications
3. Optional: Send to Discord/Telegram

### Backup to Unraid Backup Drive
Add to cron job:
\`\`\`bash
# /usr/local/bin/clarkson-backup.sh
tar -czf /mnt/backup/clarkson-$(date +%Y%m%d-%H%M%S).tar.gz \
  /mnt/user/appdata/clarkson/config/
\`\`\`

## Support & Documentation

- GitHub: https://github.com/yourusername/clarkson
- Issues: Report bugs and feature requests
- Wiki: https://github.com/yourusername/clarkson/wiki

## License

Clarkson is licensed under the MIT License. See LICENSE file for details.
