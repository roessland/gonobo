# gonobo

## Nobo/Nobø Hub / Energy Control client in Go

Listens to multicast to discover the hub IP.
Then dumps the output from the G00 command.

See https://www.glendimplex.no/media/15655/nobo-hub-api-v-1-1-integration-for-advanced-users.pdf
for a description of the API.

Communicates with the hub using TCP.

## Temperature

Sadly the ovens don't provide the thermostat temperature.
You must have a Nobø Switch connected the the Nobø Hub to get access
to temperature data.

## Listening to data continuously

The program exits after receiving H00-H05 commands
since it doesn't send heartbeats, so the hub closes
the connection.
But if this is implemented you can keep listening to any changes.
Although not very interesting if you don't have any Nobø Switches...

## Example output

```
./gonobo
2021/03/14 18:19:04 Auto discovering Nobø hub by listening to multicast...
2021/03/14 18:19:05 Found hub at 192.168.1.194 with serial starting with <hubSer>>
2021/03/14 18:19:05 usage: gonobo <hubSerial>
hint: the output so far should show the first 9 digits of the serial,
      and you must provide all 12 digits
```


```
./gonobo <hubSerial>
2021/03/14 18:21:19 Auto discovering Nobø hub by listening to multicast...
2021/03/14 18:21:20 Found hub at 192.168.1.194 with serial starting with <hubSer>
2021/03/14 18:21:20 Looking for hub with serial <hubSerial>
2021/03/14 18:21:20 To hub: HELLO 1.1 <hubSerial> 20210314172120
2021/03/14 18:21:20 From hub: HELLO 1.1
2021/03/14 18:21:20 To hub: HANDSHAKE
2021/03/14 18:21:20 From hub: HANDSHAKE
2021/03/14 18:21:20 To hub: G00
H00
H01 1 Gjestesoverom 23 23 22 1 -1
H01 2 Kjøkken og stue 20 22 17 1 -1
H01 3 Andreas soverom 20 23 18 1 -1
H02 168000147085 0 Gjestesoverom 0 1 -1 -1
H02 168000147084 0 Kjøkken 0 2 -1 -1
H02 168000147080 0 Stue 0 2 -1 -1
H02 168000147076 0 Andreas soverom 0 3 -1 -1
H03 0 Default 00000,06001,08000,15001,23000,00000,06001,08000,15001,23000,00000,06001,08000,15001,23000,00000,06001,08000,15001,23000,00000,06001,08000,15001,00000,07001,00000,07001,23000
H03 20 Andreas hjemme 00000,08001,22450,00000,08001,23000,00000,08001,23000,00000,08001,23000,00000,08001,00000,08001,00000,08001,23000
H03 21 Soveromsprogrammet 00000,09001,11150,00000,09001,10300,00000,09001,10450,00000,09001,11150,00000,09001,11450,00000,09001,11000,00000,09001,11150
H03 22 Soverom hjemmekontor 00000,08151,21150,00000,08151,21150,00000,08151,21150,00000,08151,21150,00000,08151,21150,00000,10001,21150,00000,10001,21150
H03 23 Ingen hjemme 00000,00000,00000,00000,00000,00000,00000
H04 87 0 0 -1 -1 0 -1
H04 91 0 0 -1 -1 1 3
H04 92 1 0 -1 -1 1 1
V06 <snipped encryption keys>
H05 <hubSerial> My Eco Hub 2880 87 114 11123610_rev._1 20190429
```

If you keep listening after H05 and change temperature in the Nobø app
you will get events:
```
V00 3 Andreas soverom 20 22 18 1 -1
V00 3 Andreas soverom 23 22 18 1 -1
```


## Disclaimer

This is not offical Nobo / Glen Dimplex Nordic AS software,
they do not endorse this software, and I am
not a partner of Glen Dimplex Nordic AS in any way. 