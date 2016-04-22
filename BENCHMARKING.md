## Environment

1. Go 1.6
2. wrk 4.0.0
3. AWS EC2 t2.micro with Ubuntu 14.04

## Benchmarking with go-httprouting-benchmark

Test suite: [github.com/julienschmidt/go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)


| Benchmark                           | iteration| time/iteration | bytes allocated| allocations     |
| ----------------------------------- | -------- | -------------- | -------------- | --------------- |
| BenchmarkAce_GithubAll              | 10000    |  118334 ns/op  |   13792 B/op   |   167 allocs/op |
| BenchmarkBeego_GithubAll            |  5000    |  244428 ns/op  |       0 B/op   |     0 allocs/op |
| BenchmarkBone_GithubAll             |   500    | 2890918 ns/op  |  548736 B/op   |  7241 allocs/op |
| BenchmarkDenco_GithubAll            | 10000    |  111923 ns/op  |   20224 B/op   |   167 allocs/op |
| BenchmarkGin_GithubAll              | 30000    |   44235 ns/op  |       0 B/op   |     0 allocs/op |
| BenchmarkGocraftWeb_GithubAll       |  5000    |  615754 ns/op  |  131656 B/op   |  1686 allocs/op |
| BenchmarkGoji_GithubAll             |  3000    |  652552 ns/op  |   56112 B/op   |   334 allocs/op |
| BenchmarkGojiv2_GithubAll           |  2000    |  888854 ns/op  |  118864 B/op   |  3103 allocs/op |
| BenchmarkGoJsonRest_GithubAll       |  3000    |  628870 ns/op  |  134371 B/op   |  2737 allocs/op |
| **BenchmarkGolf_GithubAll**         | **20000**|**62703 ns/op** |     **0 B/op** |  **0 allocs/op**|
| BenchmarkGoRestful_GithubAll        |   100    |17627201 ns/op  |  837832 B/op   |  6913 allocs/op |
| BenchmarkGorillaMux_GithubAll       |   200    | 7464885 ns/op  |  144464 B/op   |  1588 allocs/op |
| BenchmarkHttpRouter_GithubAll       | 20000    |   74199 ns/op  |   13792 B/op   |   167 allocs/op |
| BenchmarkHttpTreeMux_GithubAll      | 10000    |  250904 ns/op  |   65856 B/op   |   671 allocs/op |
| BenchmarkKocha_GithubAll            | 10000    |  192759 ns/op  |   23304 B/op   |   843 allocs/op |
| BenchmarkLARS_GithubAll             | 30000    |   44416 ns/op  |       0 B/op   |     0 allocs/op |
| BenchmarkMacaron_GithubAll          |  2000    |  805238 ns/op  |  201138 B/op   |  1803 allocs/op |
| BenchmarkMartini_GithubAll          |   200    | 6722507 ns/op  |  228214 B/op   |  2483 allocs/op |
| BenchmarkPat_GithubAll              |   300    | 5079035 ns/op  | 1499569 B/op   | 27435 allocs/op |
| BenchmarkPossum_GithubAll           | 10000    |  309150 ns/op  |   84448 B/op   |   609 allocs/op |
| BenchmarkR2router_GithubAll         | 10000    |  278160 ns/op  |   77328 B/op   |   979 allocs/op |
| BenchmarkRevel_GithubAll            |  1000    | 1668912 ns/op  |  337424 B/op   |  5512 allocs/op |
| BenchmarkTango_GithubAll            |  3000    |  514973 ns/op  |   87076 B/op   |  2267 allocs/op |
| BenchmarkTigerTonic_GithubAll       |  2000    | 1215298 ns/op  |  233680 B/op   |  5035 allocs/op |
| BenchmarkTraffic_GithubAll          |   200    | 9580916 ns/op  | 2659331 B/op   | 21848 allocs/op |
| BenchmarkVulcan_GithubAll           |  5000    |  329675 ns/op  |   19894 B/op   |   609 allocs/op |


## Benchmarking with WRK

Test suite: [github.com/vishr/web-framework-benchmark](https://github.com/vishr/web-framework-benchmark)

![Golf benchmark result](https://cloud.githubusercontent.com/assets/1311594/14748305/fcbdc216-0886-11e6-90a4-231e78acfb60.png)

```
benchmarking golf...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   154.78ms  210.75ms 788.98ms   79.55%
    Req/Sec    13.82k     8.23k   24.70k    45.91%
  234427 requests in 10.01s, 32.42MB read
Requests/sec:  23420.23
Transfer/sec:      3.24MB

benchmarking beego...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   152.66ms  206.02ms 779.10ms   79.39%
    Req/Sec    13.67k     7.94k   21.94k    70.29%
  208864 requests in 10.01s, 28.88MB read
Requests/sec:  20856.08
Transfer/sec:      2.88MB

benchmarking echo/standard...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   162.70ms  258.59ms   1.58s    85.34%
    Req/Sec    10.67k     5.79k   21.24k    59.47%
  207126 requests in 10.01s, 28.64MB read
Requests/sec:  20688.04
Transfer/sec:      2.86MB

benchmarking gin...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   142.73ms  194.49ms 809.02ms   79.83%
    Req/Sec    10.85k     5.18k   22.25k    59.90%
  214791 requests in 10.01s, 29.70MB read
Requests/sec:  21465.80
Transfer/sec:      2.97MB

benchmarking goji...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   178.54ms  295.77ms   1.83s    85.39%
    Req/Sec    10.55k     6.45k   22.74k    59.24%
  201138 requests in 10.08s, 27.81MB read
Requests/sec:  19963.71
Transfer/sec:      2.76MB

benchmarking martini...
Running 10s test @ http://localhost:8080/teams/x-men/members/wolverine
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   119.61ms  165.13ms 689.59ms   80.80%
    Req/Sec     6.76k     3.78k   13.13k    57.71%
  125345 requests in 10.02s, 17.33MB read
Requests/sec:  12514.69
Transfer/sec:      1.73MB
```
