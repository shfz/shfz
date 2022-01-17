# shfz

<p align="center">
  <img src="https://avatars.githubusercontent.com/u/89079069?s=200" />
</p>

A scenario-based web application fuzzng tool that supports fuzz generation by genetic algorithms.

## Features

- **Easy to customize** fuzzing test for web applications by scripting fuzzing scenario in JavaScript / TypeScript
- **Third-party packages** can be used in fuzzing scenario script
- **Genetic algorithms** fuzz generation can be used to efficiently increase code coverage
- **High affinity with CI**, such as GitHub Actions

## Install

Download binary from [Releases](https://github.com/shfz/shfz/releases) page, or compile from source.

#### Linux (amd64)

```
$ curl -Lo shfz.tar.gz https://github.com/shfz/shfz/releases/download/v0.0.1/shfz_0.0.1_linux_amd64.tar.gz
$ tar -zxvf shfz.tar.gz
$ sudo mv shfz /usr/local/bin/
$ sudo chmod +x /usr/local/bin/shfz
```

## Usage

To run fuzzing test with this tool, you need to create a scenario (that calls http requests for the web application, with automatically embeds the fuzz in the request parameter such as `username`, `password`).

Please refer to [shfz/shfzlib](https://github.com/shfz/shfzlib) for how to script scenarios.

And for Fuzzing generation by genetic algorithm, it is necessary to install the trace library in the web application.

Currently, the trace library is only compatible with Python Flask. (supported frameworks will be expanded in the future)

If the web application to be fuzzed uses flask, please install [shfz/shfz-flask](https://github.com/shfz/shfz-flask).

### server

In order to aggregate the results of fuzzing or generate fuzz by the genetic algorithm, it is necessary to start the server.

```
$ shfz server
```

By default, the http server starts on port `53653` on localhost.

And you can get the saved fuzz data from the `/data` endpoints.

```
$ curl -s http://localhost:53653/data | jq

{
  "status": [
    {
      "name": "login",
      "UsedFuzzs": [
        {
          "id": "0000",
          "fuzz": [
            {
              "name": "user",
              "text": "abcabc"
            }
          ],
...
```

### run

After setting up the server, specify the scenario file in another terminal and execute fuzzing.

```
$ shfz run -f scenario.js -n 100 -p 3 -t 30
```

> #### options
>
> - `-f`, `--file` scenario file (required)
> - `-n`, `--number` total number of executions (default 1)
> - `-p`, `--parallel` number of parallel executions (default 1)
> - `-t`, `--timeout` scenario execution timeout(seconds) (default 30)

You can get the result by sending a request to the server's the `/data` endpoints during or after fuzzing.

```
$ curl -s http://localhost:53653/data | jq

{
  "status": [
    {
      "name": "login",
      "UsedFuzzs": [
        {
          "id": "0000",
          "fuzz": [
            {
              "name": "user",
              "text": "abcabc"
            }
          ],
...
```

## CI integration

You can also install shfz on your local machine and run fuzzing, but we recommend integrating shfz into CI.

TBD
