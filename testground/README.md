# testground plans

check out first https://github.com/testground/testground to install `testground` on your host.

then, on this directory you can see the list:

```sh
export TESTGROUND_HOME=`pwd`
testground plan list --testcases
```

to run a test, for example ipfssync:
```
export TESTGROUND_HOME=`pwd`
testground daemon &
testground run single --plan=ipfssync --testcase=accept --builder=docker:go --runner=local:docker --instances=10
```
