#!/bin/bash

# API Test Script for Task Management Application
# Tests all API endpoints with proper authentication and error handling

# Don't exit on errors - we'll handle them explicitly
set -u

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
ADMIN_USERNAME="${ADMIN_USERNAME:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123}"
COOKIES_FILE="/tmp/test-api-cookies.txt"
DB_FILE="${DATABASE_PATH:-./data/tasks.db}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
print_test() {
    echo -e "${BLUE}▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Cleanup function
cleanup() {
    rm -f "$COOKIES_FILE"
}

trap cleanup EXIT

# Check if server is running
check_server() {
    print_test "Checking if server is running..."
    if curl -s -f "$BASE_URL/api/health" > /dev/null 2>&1; then
        print_success "Server is running"
        return 0
    else
        print_error "Server is not running at $BASE_URL"
        echo "Please start the server first:"
        echo "  cd backend/src && ./task-management-app"
        exit 1
    fi
}

# Test health check endpoint
test_health() {
    print_test "Testing GET /api/health"
    RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/health" 2>&1)
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ -z "$HTTP_CODE" ]; then
        print_error "Health check failed - no HTTP code received. Response: $RESPONSE"
        return 1
    fi
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        if echo "$BODY" | grep -q "ok\|status"; then
            print_success "Health check passed"
        else
            print_error "Health check returned unexpected body: $BODY"
        fi
    else
        print_error "Health check failed with HTTP $HTTP_CODE. Body: $BODY"
    fi
}

# Test login
test_login() {
    print_test "Testing POST /api/login"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"$ADMIN_PASSWORD\"}" \
        -c "$COOKIES_FILE" 2>&1) || {
        print_error "Login curl command failed"
        return 1
    }
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ -z "$HTTP_CODE" ] || ! [[ "$HTTP_CODE" =~ ^[0-9]+$ ]]; then
        print_error "Login failed - invalid HTTP code: '$HTTP_CODE'. Full response: $RESPONSE"
        return 1
    fi
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        if echo "$BODY" | grep -q "\"username\".*\"$ADMIN_USERNAME\""; then
            print_success "Login successful"
            USER_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
            print_info "Logged in as user ID: $USER_ID"
            return 0
        else
            print_error "Login returned unexpected response: $BODY"
            return 1
        fi
    else
        print_error "Login failed with HTTP $HTTP_CODE: $BODY"
        return 1
    fi
}

# Test get current user
test_get_current_user() {
    print_test "Testing GET /api/me"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/me" \
        -b "$COOKIES_FILE")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        if echo "$BODY" | grep -q "\"username\""; then
            print_success "Get current user successful"
            IS_ADMIN=$(echo "$BODY" | grep -o '"user_type":"[^"]*"' | cut -d: -f2 | tr -d '"')
            print_info "User type: $IS_ADMIN"
        else
            print_error "Get current user returned unexpected response: $BODY"
        fi
    else
        print_error "Get current user failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test create task
test_create_task() {
    print_test "Testing POST /api/tasks"
    TASK_DATA='{
        "description": "Test task from API test script",
        "status": "not_started",
        "progress": 0
    }'
    
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/tasks" \
        -H "Content-Type: application/json" \
        -b "$COOKIES_FILE" \
        -d "$TASK_DATA")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 201 ]; then
        TASK_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
        if [ -n "$TASK_ID" ]; then
            print_success "Task created with ID: $TASK_ID"
            echo "$TASK_ID" > /tmp/test-task-id.txt
        else
            print_error "Task created but ID not found in response: $BODY"
        fi
    else
        print_error "Create task failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test create second task (for dependency testing)
test_create_second_task() {
    print_test "Testing POST /api/tasks (second task for dependencies)"
    TASK_DATA='{
        "description": "Second test task for dependency testing",
        "status": "not_started",
        "progress": 0
    }'
    
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/tasks" \
        -H "Content-Type: application/json" \
        -b "$COOKIES_FILE" \
        -d "$TASK_DATA")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 201 ]; then
        TASK_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
        if [ -n "$TASK_ID" ]; then
            print_success "Second task created with ID: $TASK_ID"
            echo "$TASK_ID" > /tmp/test-task-id-2.txt
        else
            print_error "Second task created but ID not found in response: $BODY"
        fi
    else
        print_error "Create second task failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test get tasks
test_get_tasks() {
    print_test "Testing GET /api/tasks"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/tasks" \
        -b "$COOKIES_FILE")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        TASK_COUNT=$(echo "$BODY" | grep -o '"id":[0-9]*' | wc -l)
        print_success "Get tasks successful (found $TASK_COUNT tasks)"
    else
        print_error "Get tasks failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test get tasks with filter
test_get_tasks_filtered() {
    print_test "Testing GET /api/tasks?status=not_started"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/tasks?status=not_started" \
        -b "$COOKIES_FILE")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        print_success "Get filtered tasks successful"
    else
        print_error "Get filtered tasks failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test get specific task
test_get_task() {
    if [ -f /tmp/test-task-id.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing GET /api/tasks/$TASK_ID"
        RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/tasks/$TASK_ID" \
            -b "$COOKIES_FILE")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            if echo "$BODY" | grep -q "\"id\":$TASK_ID"; then
                print_success "Get task $TASK_ID successful"
            else
                print_error "Get task returned unexpected response: $BODY"
            fi
        else
            print_error "Get task failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test get task - task ID not found"
    fi
}

# Test update task
test_update_task() {
    if [ -f /tmp/test-task-id.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing PUT /api/tasks/$TASK_ID"
        UPDATE_DATA='{
            "description": "Updated test task description",
            "status": "in_progress",
            "progress": 50
        }'
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/tasks/$TASK_ID" \
            -H "Content-Type: application/json" \
            -b "$COOKIES_FILE" \
            -d "$UPDATE_DATA")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            if echo "$BODY" | grep -q "\"progress\":50"; then
                print_success "Update task successful"
            else
                print_error "Update task returned unexpected response: $BODY"
            fi
        else
            print_error "Update task failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test update task - task ID not found"
    fi
}

# Test add dependency
test_add_dependency() {
    if [ -f /tmp/test-task-id.txt ] && [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        DEP_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing POST /api/tasks/$TASK_ID/dependencies"
        DEP_DATA="{\"dependency_id\":$DEP_ID}"
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/tasks/$TASK_ID/dependencies" \
            -H "Content-Type: application/json" \
            -b "$COOKIES_FILE" \
            -d "$DEP_DATA")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 200 ]; then
            print_success "Add dependency successful (Task $TASK_ID depends on Task $DEP_ID)"
        else
            print_error "Add dependency failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test add dependency - task IDs not found"
    fi
}

# Test get dependencies
test_get_dependencies() {
    if [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        print_test "Testing GET /api/tasks/$TASK_ID/dependencies"
        RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/tasks/$TASK_ID/dependencies" \
            -b "$COOKIES_FILE")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            print_success "Get dependencies successful"
        else
            print_error "Get dependencies failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test get dependencies - task ID not found"
    fi
}

# Test completing task with incomplete dependency (should fail)
test_complete_task_with_dependency() {
    if [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        print_test "Testing PUT /api/tasks/$TASK_ID (try to complete with incomplete dependency - should fail)"
        UPDATE_DATA='{
            "status": "completed"
        }'
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/tasks/$TASK_ID" \
            -H "Content-Type: application/json" \
            -b "$COOKIES_FILE" \
            -d "$UPDATE_DATA")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 400 ]; then
            if echo "$BODY" | grep -qi "dependency\|cannot mark\|completed"; then
                print_success "Correctly rejected completion of task with incomplete dependency"
            else
                print_error "Rejected but with unexpected error message: $BODY"
            fi
        else
            print_error "Should have rejected completion but got HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test dependency validation - task ID not found"
    fi
}

# Test complete dependency first
test_complete_dependency() {
    if [ -f /tmp/test-task-id.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing PUT /api/tasks/$TASK_ID (complete dependency task)"
        UPDATE_DATA='{
            "status": "completed",
            "progress": 100
        }'
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/tasks/$TASK_ID" \
            -H "Content-Type: application/json" \
            -b "$COOKIES_FILE" \
            -d "$UPDATE_DATA")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            print_success "Completed dependency task"
        else
            print_error "Complete dependency task failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test complete dependency - task ID not found"
    fi
}

# Test complete task (should succeed now)
test_complete_task() {
    if [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        print_test "Testing PUT /api/tasks/$TASK_ID (complete task - should succeed now)"
        UPDATE_DATA='{
            "status": "completed",
            "progress": 100
        }'
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/tasks/$TASK_ID" \
            -H "Content-Type: application/json" \
            -b "$COOKIES_FILE" \
            -d "$UPDATE_DATA")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | head -n-1)
        
        if [ "$HTTP_CODE" -eq 200 ]; then
            print_success "Completed task successfully"
        else
            print_error "Complete task failed with HTTP $HTTP_CODE: $BODY"
        fi
    else
        print_error "Cannot test complete task - task ID not found"
    fi
}

# Test remove dependency
test_remove_dependency() {
    if [ -f /tmp/test-task-id.txt ] && [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        DEP_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing DELETE /api/tasks/$TASK_ID/dependencies/$DEP_ID"
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/api/tasks/$TASK_ID/dependencies/$DEP_ID" \
            -b "$COOKIES_FILE")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        
        if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
            print_success "Remove dependency successful"
        else
            print_error "Remove dependency failed with HTTP $HTTP_CODE"
        fi
    else
        print_error "Cannot test remove dependency - task IDs not found"
    fi
}

# Test delete task
test_delete_task() {
    if [ -f /tmp/test-task-id.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id.txt)
        print_test "Testing DELETE /api/tasks/$TASK_ID"
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/api/tasks/$TASK_ID" \
            -b "$COOKIES_FILE")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        
        if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
            print_success "Delete task successful"
            rm -f /tmp/test-task-id.txt
        else
            print_error "Delete task failed with HTTP $HTTP_CODE"
        fi
    else
        print_error "Cannot test delete task - task ID not found"
    fi
}

# Test delete second task
test_delete_second_task() {
    if [ -f /tmp/test-task-id-2.txt ]; then
        TASK_ID=$(cat /tmp/test-task-id-2.txt)
        print_test "Testing DELETE /api/tasks/$TASK_ID (second task)"
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/api/tasks/$TASK_ID" \
            -b "$COOKIES_FILE")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        
        if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
            print_success "Delete second task successful"
            rm -f /tmp/test-task-id-2.txt
        else
            print_error "Delete second task failed with HTTP $HTTP_CODE"
        fi
    else
        print_error "Cannot test delete second task - task ID not found"
    fi
}

# Test create user (admin only)
test_create_user() {
    print_test "Testing POST /api/users (admin only)"
    USER_DATA='{
        "username": "testuser_'$(date +%s)'",
        "password": "testpass123",
        "email": "testuser@example.com",
        "display_name": "Test User",
        "user_type": "regular"
    }'
    
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users" \
        -H "Content-Type: application/json" \
        -b "$COOKIES_FILE" \
        -d "$USER_DATA")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 201 ]; then
        print_success "Create user successful"
    else
        print_error "Create user failed with HTTP $HTTP_CODE: $BODY"
    fi
}

# Test logout
test_logout() {
    print_test "Testing POST /api/logout"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/logout" \
        -b "$COOKIES_FILE")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 204 ]; then
        print_success "Logout successful"
    else
        print_error "Logout failed with HTTP $HTTP_CODE"
    fi
}

# Test unauthorized access
test_unauthorized_access() {
    print_test "Testing GET /api/tasks (without authentication - should fail)"
    RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/tasks")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" -eq 401 ] || [ "$HTTP_CODE" -eq 403 ]; then
        print_success "Correctly rejected unauthorized access"
    else
        print_error "Should have rejected unauthorized access but got HTTP $HTTP_CODE"
    fi
}

# Main test execution
main() {
    echo "=========================================="
    echo "  Task Management API Test Suite"
    echo "=========================================="
    echo ""
    echo "Base URL: $BASE_URL"
    echo "Admin Username: $ADMIN_USERNAME"
    echo ""
    
    # Cleanup from previous runs
    rm -f "$COOKIES_FILE" /tmp/test-task-id*.txt
    
    # Run tests
    echo "Starting tests..."
    if ! check_server; then
        echo "ERROR: Server check failed"
        exit 1
    fi
    
    echo "Running health check..."
    if ! test_health; then
        print_error "Health check failed, stopping tests"
        exit 1
    fi
    
    echo "Running login test..."
    if ! test_login; then
        print_error "Login failed, stopping tests"
        exit 1
    fi
    
    echo "Continuing with remaining tests..."
    test_get_current_user
    test_create_task
    test_create_second_task
    test_get_tasks
    test_get_tasks_filtered
    test_get_task
    test_update_task
    test_add_dependency
    test_get_dependencies
    test_complete_task_with_dependency
    test_complete_dependency
    test_complete_task
    test_remove_dependency
    test_delete_task
    test_delete_second_task
    test_create_user
    test_logout
    test_unauthorized_access
    
    # Summary
    echo ""
    echo "=========================================="
    echo "  Test Summary"
    echo "=========================================="
    echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Failed: $TESTS_FAILED${NC}"
    echo ""
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}Some tests failed.${NC}"
        exit 1
    fi
}

# Run main function
main
