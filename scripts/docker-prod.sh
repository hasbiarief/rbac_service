#!/bin/bash

# Production Docker Compose Management Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env exists
check_env() {
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating from .env.production template..."
        cp .env.production .env
        print_warning "Please edit .env file with your production values before continuing."
        exit 1
    fi
}

# Build and start services
start() {
    print_info "Starting production services..."
    check_env
    
    # Build and start services
    docker-compose -f docker-compose.prod.yml up -d --build
    
    print_success "Services started successfully!"
    print_info "Waiting for services to be healthy..."
    
    # Wait for services to be healthy
    sleep 10
    
    # Check service status
    docker-compose -f docker-compose.prod.yml ps
    
    print_info "Application should be available at http://localhost:8081"
    print_info "Use 'docker-compose -f docker-compose.prod.yml logs -f' to view logs"
}

# Start with Nginx
start_with_nginx() {
    print_info "Starting production services with Nginx..."
    check_env
    
    # Build and start services including nginx
    docker-compose -f docker-compose.prod.yml --profile nginx up -d --build
    
    print_success "Services with Nginx started successfully!"
    print_info "Application should be available at http://localhost"
}

# Stop services
stop() {
    print_info "Stopping production services..."
    docker-compose -f docker-compose.prod.yml down
    print_success "Services stopped successfully!"
}

# Restart services
restart() {
    print_info "Restarting production services..."
    stop
    start
}

# View logs
logs() {
    docker-compose -f docker-compose.prod.yml logs -f "${2:-}"
}

# Run migrations
migrate() {
    print_info "Running database migrations..."
    docker-compose -f docker-compose.prod.yml exec app ./server migrate
    print_success "Migrations completed!"
}

# Backup database
backup() {
    print_info "Creating database backup..."
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres huminor_rbac > "$BACKUP_FILE"
    print_success "Database backup created: $BACKUP_FILE"
}

# Show status
status() {
    print_info "Service status:"
    docker-compose -f docker-compose.prod.yml ps
    
    print_info "\nService health:"
    docker-compose -f docker-compose.prod.yml exec app wget --spider -q http://localhost:8081/health && print_success "App: Healthy" || print_error "App: Unhealthy"
}

# Clean up
clean() {
    print_warning "This will remove all containers, networks, and volumes. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_info "Cleaning up..."
        docker-compose -f docker-compose.prod.yml down -v --remove-orphans
        docker system prune -f
        print_success "Cleanup completed!"
    else
        print_info "Cleanup cancelled."
    fi
}

# Help
help() {
    echo "Production Docker Compose Management Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start         Start production services"
    echo "  start-nginx   Start services with Nginx reverse proxy"
    echo "  stop          Stop all services"
    echo "  restart       Restart all services"
    echo "  logs [service] View logs (optionally for specific service)"
    echo "  migrate       Run database migrations"
    echo "  backup        Create database backup"
    echo "  status        Show service status and health"
    echo "  clean         Remove all containers, networks, and volumes"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 logs app"
    echo "  $0 migrate"
}

# Main script
case "${1:-}" in
    start)
        start
        ;;
    start-nginx)
        start_with_nginx
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs "$@"
        ;;
    migrate)
        migrate
        ;;
    backup)
        backup
        ;;
    status)
        status
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        help
        ;;
    *)
        print_error "Unknown command: ${1:-}"
        echo ""
        help
        exit 1
        ;;
esac