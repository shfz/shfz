<p align="center">
  <img src="https://avatars.githubusercontent.com/u/89079069?s=200" />
</p>

<p align="center">
  <a href="https://github.com/shfz/shfz/releases">
    <img src="https://img.shields.io/github/workflow/status/shfz/shfz/goreleaser" alt="GitHub Workflow Status">
  </a>
    <a href="https://github.com/shfz/shfz/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/shfz/shfz" alt="license">
  </a>
    <a href="https://github.com/shfz/shfz/releases">
    <img src="https://img.shields.io/github/v/release/shfz/shfz" alt="release">
  </a>
    <a href="https://github.com/shfz/shfz/releases">
    <img src="https://img.shields.io/github/downloads/shfz/shfz/total" alt="downloads">
  </a>
</p>

<h1 align="center">shfz</h1>

A scenario-based web application fuzzng tool that supports fuzz generation by genetic algorithms.

<p align="center">
  <img src="https://raw.githubusercontent.com/shfz/shfz/main/image/shfz.jpg" />
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/shfz/shfz/main/image/use.jpg" />
</p>

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

### Github Actions

[A example workflow](https://github.com/shfz/demo-webapp/blob/main/.github/workflows/fuzzing.yml)

1. Create fuzzing scenario in `/fuzz` directory.

<https://github.com/shfz/demo-webapp/tree/main/fuzz>

```yml
      - uses: actions/setup-node@v2
        with:
          node-version: "16"
      - name: setup fuzzing scenario
        run: |
          cd fuzz
          npm i
          ./node_modules/typescript/bin/tsc scenario.ts
          file scenario.js
```

2. Setup webapp (by docker-compose).

```yml
      - name: setup webapp
        run: |
          docker-compose build
          docker-compose up -d
          docker-compose ps -a
```

> If this webapp is created by Python Flask, install [shfz/shfz-flask](https://github.com/shfz/shfz-flask)
>
> Note.
>
> If you use docker-compose to launch the webapp on Linux, you need to enable host.docker.internal.
>
> ```yml
>     extra_hosts:
>       - "host.docker.internal:host-gateway"
> ```
>
> And shfztrace is initialised by `fuzzUrl="http://host.docker.internal:53653"`
>
> ```python
> from flask import *
> from shfzflask import shfztrace
>
> app = Flask(__name__)
> shfztrace(app, fuzzUrl="http://host.docker.internal:53653")
> ```

3. Setup shfz

```yml
      - name: setup shfz
        run: |
          wget https://github.com/shfz/shfz/releases/download/v0.0.2/shfz_0.0.2_linux_amd64.tar.gz
          tar -zxvf shfz_0.0.2_linux_amd64.tar.gz
          sudo chmod +x shfz
          ./shfz --help
```

4. Run fuzzzing

```yml
      - name: run shfz server
        run: ./shfz server &

      - name: run fuzzing
        run: ./shfz run -f fuzz/scenario.js -n 100
```

4. (GitHub Actions) Report result in Issue

```yml
      - name: export fuzzing report
        run: >
          curl
          -F "hash=${{ github.sha }}"
          -F "repo=${{ github.repository }}"
          -F "id=${{ github.run_id }}"
          -F "job=${{ github.job }}"
          -F "number=${{ github.run_number }}"
          -F "path=/app"
          http://localhost:53653/report > report.md
      - name: create issue
        uses: peter-evans/create-issue-from-file@v3
        with:
          title: shfz result
          content-filepath: ./report.md
          labels: |
            shfz
```

5. Export fuzzing data

```yml
      - name: export fuzzing data
        run: curl http://localhost:53653/data > result.json
      - name: upload artifact
        uses: actions/upload-artifact@v2
        with:
          name: result.json
          path: ./result.json
```

6. Export application log

```yml
      - name: export application log
        run: docker logs demo-webapp_app_1 > app.log
      - name: upload artifact
        uses: actions/upload-artifact@v2
        with:
          name: app.log
          path: ./app.log
```
