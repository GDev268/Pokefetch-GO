# Pokefetch-GO
Implementation of FastFetch + Pokemon-Colorscripts + PokeAPI now in the Golang language

---

## Efficiency against Pokefetch (AKA: Pokefetch V1)

Simple benchmarking result of PokeFetch VS PokeFetch GO (using hyperfine)
```
❯ hyperfine --warmup 10 -N "./pokefetch" "./pokefetch_go"
Benchmark 1: ./pokefetch
  Time (mean ± σ):     436.9 ms ±  29.0 ms    [User: 78.3 ms, System: 60.8 ms]
  Range (min … max):   399.5 ms … 478.5 ms    10 runs
 
Benchmark 2: ./pokefetch_go
  Time (mean ± σ):     206.4 ms ±  23.1 ms    [User: 70.6 ms, System: 59.8 ms]
  Range (min … max):   178.2 ms … 262.0 ms    15 runs
 
Summary
  ./pokefetch_go ran
    2.12 ± 0.28 times faster than ./pokefetch
```
