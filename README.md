# consul-docker-operator

[consul-docker-operator](https://github.com/neomantra/consul-docker-operator) watches a list of Docker Registry Images, storing information about them in Consul KV.

```
usage: consul-docker-operator [args]

  --key -k  
```

All Consul configuration is performed via environment variables, however there is an argument to load an ENV file.


## Building

```
task build

task docker:build-arm64
```

## Credits and License

Copyright (c) 2023 [Neomantra BV](https://neomantra.com).  Authored by Evan Wies.

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).
