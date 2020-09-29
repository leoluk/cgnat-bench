# cgnat-bench

Simple benchmark that verifies CGNAT session concurrency limits.

Run the server:

    ulimit -n 65535 && go run github.com/leoluk/cgnat-bench/server -addr [::]:8081
   
Test 100 concurrent connections to server:

    go run github.com/leoluk/cgnat-bench/concurrency -n 100 -addr <server>
    
 
