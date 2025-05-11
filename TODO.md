# TODO Items

## Command Test Improvements

To fully implement and complete the command tests, the following improvements are needed:

1. **Make Command Functions Mockable**:
   - Convert direct functions like `listDevices`, `getDevice`, etc. to variable functions in the main code
   - Example conversion:
   ```go
   // Change from this:
   func listDevices(token string) ([]Device, error) {
       // implementation
   }

   // To this:
   type ListDevicesFunc func(token string) ([]Device, error)
   
   var listDevices ListDevicesFunc = listDevicesImpl
   
   func listDevicesImpl(token string) ([]Device, error) {
       // implementation
   }
   ```

2. **Fix Skipped Tests**:
   - After functions are made mockable, uncomment and re-enable skipped tests in:
     - `/cmd/devices/devices_test.go`
     - `/cmd/events/events_test.go`
     - `/cmd/error_handling_test.go`
     - `/cmd/flags_test.go`

3. **Fix Error Message Assertions**:
   - Update tests to check for the actual error messages produced by the Cobra library
   - Some tests are failing because they're checking for "requires 1 arg" but Cobra outputs "accepts 1 arg(s)"

4. **Fix Completion Tests**:
   - The completion tests are failing due to file handling issues
   - Need to investigate and fix the pipe closure error in `TestCommandCompletions`

## Code Coverage Improvements

Once all tests are passing:

1. Generate a coverage report to identify areas needing additional tests:
   ```
   go test ./... -cover
   ```

2. Generate a detailed coverage report:
   ```
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

3. Add tests for any code paths that aren't currently covered

## Integration Test Improvements

1. Create an environment for running integration tests:
   ```
   INTEGRATION_TESTS=true go test ./pkg/tests/integration/... -v
   ```

2. Consider adding more realistic test scenarios with mocked API responses