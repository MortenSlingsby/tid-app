# tid-app

En liten CLI tool for å logge timene. Planen til videre utvikling fins i github issues.

## Installation
installer med:
```
nix profile install github:mortenslingsby/tid-app#tid
```

## Legg til AOs
For å legge til AO kjører man:
```
tid add -f <full_name> <tag> 
```
Så for en AO som er `Dataflyt Utvikling` som jeg vil bruke med tag `du`, så kjører jeg:
```
tid add -f Dataflyt-Utvikling du
```

## Start logging av tid
For å begynne med logging så kjører man:
```
tid start <tag>
```
So hvis jeg skal jobbe med utviklingsoppgaver er det `tid start du`. Da vil tiden løpe frem til du kaller den samme kommando med en annen tag

## Stop logging av tid
For å stoppe logging kjør:
```
tid stop
```

## Aggresso timeføring
For å se aggregerte timene du jobbet per AO per dag forrige uke kjør:
```
tid view
```
