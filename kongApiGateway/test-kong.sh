#!/bin/bash
# Test script for Kong API Gateway

echo "======================================"
echo "Kong API Gateway Health Check"
echo "======================================"
echo ""

# Check if Kong container is running
echo "1. Checking if Kong container is running..."
if docker compose ps | grep -q "kong.*Up"; then
    echo "   ✓ Kong container is running"
else
    echo "   ✗ Kong container is not running"
    echo "   Run: docker compose up -d"
    exit 1
fi
echo ""

# Check Kong health
echo "2. Checking Kong health status..."
if docker compose ps | grep -q "kong.*healthy"; then
    echo "   ✓ Kong is healthy"
else
    echo "   ⚠ Kong may not be fully ready yet"
fi
echo ""

# Test Kong proxy endpoint
echo "3. Testing Kong proxy endpoint (http://localhost:8000)..."
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8000/ 2>&1)
if [ "$response" = "404" ]; then
    echo "   ✓ Kong is responding (404 expected for root path)"
elif [ "$response" = "000" ]; then
    echo "   ✗ Cannot connect to Kong on port 8000"
    echo "   Check if port is accessible: netstat -an | grep 8000"
    exit 1
else
    echo "   ⚠ Unexpected response code: $response"
fi
echo ""

# Test configured route
echo "4. Testing configured route (POST /auth/create-account)..."
response=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8000/auth/create-account 2>&1)
if [ "$response" = "500" ]; then
    echo "   ⚠ Kong routes to backend, but auth service is not running on port 5002"
    echo "   Start auth service: cd ../auth_service && make run"
elif [ "$response" = "000" ]; then
    echo "   ✗ Cannot connect to Kong"
    exit 1
elif [ "$response" = "404" ]; then
    echo "   ✗ Route not configured properly"
    echo "   Check Kong logs: docker compose logs kong"
else
    echo "   ✓ Route is configured (response code: $response)"
fi
echo ""

# Check backend services
echo "5. Checking backend services..."
echo "   Checking auth service (port 5002)..."
if nc -z localhost 5002 2>/dev/null; then
    echo "   ✓ Auth service is listening on port 5002"
else
    echo "   ✗ Auth service is not running on port 5002"
    echo "   Start with: cd ../auth_service && make run"
fi

echo "   Checking user service (port 5001)..."
if nc -z localhost 5001 2>/dev/null; then
    echo "   ✓ User service is listening on port 5001"
else
    echo "   ✗ User service is not running on port 5001"
    echo "   Start with: cd ../user_service && npm run start:dev"
fi
echo ""

echo "======================================"
echo "Summary"
echo "======================================"
echo "Kong API Gateway: http://localhost:8000"
echo "Kong Admin API: http://localhost:8001"
echo ""
echo "Next steps:"
echo "1. Start auth service: cd ../auth_service && make run"
echo "2. Start user service: cd ../user_service && npm run start:dev"
echo "3. Test API: curl -X POST http://localhost:8000/auth/create-account -d '{...}'"
echo ""
