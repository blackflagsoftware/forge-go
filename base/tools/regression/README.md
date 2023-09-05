## Regression

This tool is to help you with regression tests.  Regression is a tool that reads tests from a json or yaml formatted file, calling your API and checking against the actual response your tests expected reponse.

Where is your data kept?  Do you pre-populate it with known data or rely on an existing data set (hoping "things" don't change).  It is really up to your API and the env var values you give it, this tool just calls that API.  As you will see in the notes, you can run this against various "environments" if desired, remember its your data that will be manipulated but the API.

#### Usage
As a golang project, this project can be compiled via the command line for any (ARCH/OS) combination you desire.  How you call it is up to you.  For local development I like to:
`go install .`
(in this directory), then make sure the PATH is pointed to your `GOPATH/bin`

command line args defaults in `[]`:

`testPathFile [""]`: path/file of tests 
`testPath [""]`: (optional) path to files of test, if present then testPathWithFile is not needed
`environment [""]`: (optional) environment to run against the test's run_environment
`dynamicValuesPath [""]`: (optional) load and save path/file for dynamic values
`printDebug [false]`: (optional) print the DynamicValues at the end of the test

usage: `regression -testPathFile=/your/path/to/file/test.json -environment=dev -printDebug=true`

If you are using this tool as part of the great project: [forge](https://github.com/blackflagsoftware/forge-go), then there is a folder called `test` and within that folder is a file called `regression.sh`, which is a good starter on how to:
- run your API
- call `regression`
- `kill` the existing process

You can add any data setup/destroy as needed.

`testPathFile` vs `testPath`: when I originally wrote this, I only had the `testPathFile` option.  I found with a lot of tests in one file, it became cumbersome to find a particular test.  `testPath` was born, it takes a directory for the value and will look at any file within that directory for `.json | .yaml | .yml` and process those files.  See `TestOutput` for file formatting.

##### env vars
If you are using this tool as part of the great project: [forge](https://github.com/blackflagsoftware/forge-go), you can set an env var called `<your_project_REGRESSION_ENVIRONMENT` that holds the environment you want to test against.  The command line `environment` will override this.

#### Tests
If you are using this tool as part of the great project: [forge](https://github.com/blackflagsoftware/forge-go), then there is a folder called `test` and within that folder is a file called `test_file.json`, that shows a sample test.

##### test payload
`name`: (string) your name for this test, (make it descriptive)
`active`: (boolean) it this active `true/false`
`host`: (string) host part of your API's uri
`path`: (string) path (including query params) of your API's uri
`test_type`: (string) acceptable values: `rest | grpc` (grpc is currently not supported)
`run_environments`: (array:string) which environments to run against, see: `environment`, if not set `dev` is used by default
`method`: (string) used only for `rest`, acceptable values: `GET | POST | PUT | PATCH | DELETE`
`auth_user`: (string) basic auth user, if empty it will not use `basic auth`
`auth_pwd`: (string) basic auth password
`request_headers`: (map) key/value headers (add any other type of authentication here)
`request_body`: (map/array) this will be the body sent to the API
`request_timeout`: (int) the request timeout to your API, if this is not set default of 30 seconds is used
`expected_response_body`: (map/array) this is the expected body to come back see `ExpectedBody`, below for more information
`expected_status`: (int) the expected status code
`wait_time`: (int) optional wait time for this test to run before running the next test, 0 is the default

example:
```
[
	{
		"name": "",
		"active": true,
		"host": "",
		"path": "",
		"test_type": "rest",
		"run_environments": ["dev"],
		"method": "",
		"auth_user": "",
		"auth_pwd": "",
		"request_headers": {
			"Content-Type": "application/json"
		},
		"request_body": {
		},
		"request_timeout": 30,
		"expected_response_body": {
			"data": {
			}
		},
		"expected_response_status": 200,
		"wait_time": 0
	}
]
```
##### ExpectedBody
When defining the `expected_request_body`, you don't have to define everything that will come back from the API payload.  The "compare" is designed to only check the `actual_response_body` elements if they are defined in the `expected_request_body`.  So if a field is dynamic or not necessary for your testing, just leave it out of the `expected_request_body`.

##### DynamicValues
Regression is designed to capture, store and reuse dynamic values from one call to the next.  Meaning, if you test a `POST` endpoint and it is designed to create a autoincrement id or a `uuid`, then there is a process to capture that value and save it for the next call, e.g.: a `PATCH` call later in your testing.

By using `dyn:<unique_name_here>` on both the capture and subsequent calls, the value will be replaced as needed.

example:
```
[
	{
		"name": "my creation call - POST",
		...
		"expected_response_body": {
			"uid": "dyn:post_uid",
			...
		}
	},
	{
		"name": "my pathc call - PATCH",
		"path": "/v1/endpoint/dyn:post_uid"
		...
	}
]
```
By default this works for `strings`.  In order for it to work other types, i.e.: int, float, bool; tack on an extra `:int|float|boolean`.  This is mostly used for the request body, so that the json works in the test file and is sent with the correct type as the request payload.

example:
```
[
	{
		"name": "my creation call - POST",
		...
		"expected_response_body": {
			"id": "dyn:post_id:int",
			...
		}
	},
	{
		"name": "my pathc call - PATCH",
		...
		"request_body": {
			"id": "dyn:post_id:int",
			...
		}
	}
]
```
There are other `DynamicValues` provided for you if you need current dates that are consistent each time you run the tests, use these like a custom DynamicValue: e.g.: `dyn:DynDateNow`
```
DynDateMinusSeven
DynDateMinusSix
DynDateMinusFive
DynDateMinusFour
DynDateMinusThree
DynDateMinusTwo
DynDateMinusOne
DynDateNow
DynDateAddOne
DynDateAddTwo
DynDateAddThree
DynDateAddFour
DynDateAddFive
DynDateAddSix
DynDateAddSeven
```

##### TestOutput
After all the test have ran, the results will be pushed to `STDOUT` in this format:
```
* Test: my first test [test_file.json]
        Status: SUCCEEDED | FAILED | SKIPPED
// if there are any errors then this lines will show for each error the test encountered
		Errors:
            Status - want: 200; got: 400
// then a summary of the counts will be last

+---------------------------------------------------+
|  Total: 2 | Successed: 1 | Failed: 1 | Skipped: 0 |
+---------------------------------------------------+

Saved file to: ./regression_results/reg_20230831203645.json
```
All the tests and results will be saved to a file, as shown above.  Each test will be printed again with an added set of fields:
`actual_response_body`: (array/map) the actual payload coming back from the endpoint call
`actual_status`: (int) the actual status code
`messages`: (array) any error messages

The results are defined here and live at the root of the json object and will look like this:
```
{
    "dateTime": "2023-08-31T20:36:45.782936-06:00",
    "total": 2,
    "succeeded": 1,
    "failed": 1,
    "skipped": 0,
    "tests": []
}
```
It is saved as `.json` file so that you can, if needed, parse via `jq` or any programmic way for reporting.