#!/bin/bash

# Test script for password generator
set -e

echo "Running password generator tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASSGEN="./passgen"
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to run test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_exit_code="${3:-0}"
    
    echo -n "Testing: $test_name... "
    
    if eval "$command" > /dev/null 2>&1; then
        actual_exit_code=$?
    else
        actual_exit_code=$?
    fi
    
    if [ $actual_exit_code -eq $expected_exit_code ]; then
        echo -e "${GREEN}PASS${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}FAIL${NC} (expected exit code $expected_exit_code, got $actual_exit_code)"
        ((TESTS_FAILED++))
    fi
}

# Helper function to test password properties
test_password_properties() {
    local test_name="$1"
    local args="$2"
    local expected_length="$3"
    local should_have_special="${4:-false}"
    
    echo -n "Testing: $test_name... "
    
    # Generate password and capture output
    output=$($PASSGEN $args 2>/dev/null | tail -1 | cut -d' ' -f2-)
    
    # Check if output exists
    if [ -z "$output" ]; then
        echo -e "${RED}FAIL${NC} (no password generated)"
        ((TESTS_FAILED++))
        return
    fi
    
    # Check length
    actual_length=${#output}
    if [ $actual_length -ne $expected_length ]; then
        echo -e "${RED}FAIL${NC} (expected length $expected_length, got $actual_length)"
        ((TESTS_FAILED++))
        return
    fi
    
    # Check character sets
    has_upper=$(echo "$output" | grep -q '[A-Z]' && echo true || echo false)
    has_lower=$(echo "$output" | grep -q '[a-z]' && echo true || echo false)
    has_number=$(echo "$output" | grep -q '[0-9]' && echo true || echo false)
    has_special=$(echo "$output" | grep -q '[!@#$%^&*()_+\-=\[\]{}|;:,.<>?]' && echo true || echo false)
    
    # Check excluded characters
    has_excluded=$(echo "$output" | grep -q '[0OIl1]' && echo true || echo false)
    
    if [ "$has_excluded" = "true" ]; then
        echo -e "${RED}FAIL${NC} (contains excluded characters: 0, O, I, l, 1)"
        ((TESTS_FAILED++))
        return
    fi
    
    if [ "$has_upper" = "false" ] || [ "$has_lower" = "false" ] || [ "$has_number" = "false" ]; then
        echo -e "${RED}FAIL${NC} (missing required character types)"
        ((TESTS_FAILED++))
        return
    fi
    
    if [ "$should_have_special" = "true" ] && [ "$has_special" = "false" ]; then
        echo -e "${RED}FAIL${NC} (should contain special characters but doesn't)"
        ((TESTS_FAILED++))
        return
    fi
    
    if [ "$should_have_special" = "false" ] && [ "$has_special" = "true" ]; then
        echo -e "${RED}FAIL${NC} (shouldn't contain special characters but does)"
        ((TESTS_FAILED++))
        return
    fi
    
    echo -e "${GREEN}PASS${NC}"
    ((TESTS_PASSED++))
}

echo "=== Basic Functionality Tests ==="

# Test help option
run_test "Help option (-h)" "$PASSGEN -h"

# Test default password generation
test_password_properties "Default password (12 chars, no special)" "" 12 false

# Test custom length
test_password_properties "Custom length (8 chars)" "-l 8" 8 false
test_password_properties "Custom length (20 chars)" "-l 20" 20 false

# Test with special characters
test_password_properties "With special chars (12 chars)" "-s" 12 true
test_password_properties "With special chars (16 chars)" "-l 16 -s" 16 true

echo ""
echo "=== Error Handling Tests ==="

# Test invalid length (too short)
run_test "Invalid length (too short)" "$PASSGEN -l 2" 1

# Test invalid length (too long)
run_test "Invalid length (too long)" "$PASSGEN -l 200" 1

# Test invalid count (too low)
run_test "Invalid count (too low)" "$PASSGEN -c 0" 1

# Test invalid count (too high)
run_test "Invalid count (too high)" "$PASSGEN -c 200" 1

# Test special chars with insufficient length
run_test "Special chars with length 3" "$PASSGEN -l 3 -s" 1

# Test invalid option
run_test "Invalid option" "$PASSGEN -x" 1

echo ""
echo "=== Multiple Password Generation Tests ==="

# Test multiple password generation
run_test "Generate 3 passwords" "$PASSGEN -c 3"
run_test "Generate 5 passwords with special chars" "$PASSGEN -c 5 -s"

echo ""
echo "=== Randomness Tests ==="

# Test that multiple runs produce different passwords
echo -n "Testing: Password randomness... "
pass1=$($PASSGEN 2>/dev/null | tail -1 | cut -d' ' -f2-)
pass2=$($PASSGEN 2>/dev/null | tail -1 | cut -d' ' -f2-)
pass3=$($PASSGEN 2>/dev/null | tail -1 | cut -d' ' -f2-)

if [ "$pass1" != "$pass2" ] && [ "$pass1" != "$pass3" ] && [ "$pass2" != "$pass3" ]; then
    echo -e "${GREEN}PASS${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${RED}FAIL${NC} (generated identical passwords)"
    ((TESTS_FAILED++))
fi

echo ""
echo "=== Test Summary ==="
echo "Tests passed: $TESTS_PASSED"
echo "Tests failed: $TESTS_FAILED"
echo "Total tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi