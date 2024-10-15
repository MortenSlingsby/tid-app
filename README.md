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

## Legg til custom tid eller fix feil tid
Man kan ikke endre eksisterende logger, men du kan legge til custom timer som kan være både + og - timer med kommando: `tid fix <tag> <tid_i_minutt>`. Så hvis du har glemt å bytte tag fra A til B for en time siden så kan du legge til 1 minus time for A og 1 plus time for B for å fikse det. Det gjør man slik:
```
tid fix A -60
tid fix B 60
```

## Aggresso timeføring
For å se aggregerte timene du jobbet per AO per dag kjør:
```
tid view <relative_week>
```
Relative week defaulter til uke - 1 (i andre ord, forrige uke), så man trenger ikke det argumentet. Men hvis du vil se f.eks denne uke eller for flere uker siden så kan du gi antall uker tilbake i tid her. Så oversikt for 2 uker siden kjøres med `tid view 2`

## Oversikt over AOs
For å se hvilke AOs som du har lagt til kjør du `tid list`

## Oversikt over log til i dag
For å se hva du har logget så langt i dat kjør `tid log`

## Drop log
Man kan droppe en log med kommando `tid drop <log_id>`. Log id kan finnes ved å kjøre `tid log`
