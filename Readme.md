XBAND

The following is a compilation of source code and files obtained over the years
that made the xband modem for the SNES and Sega Genesis what it was.

Xband Hardware
--------------
The xband hardware is a 2400 baud Rockwell chipset modem with a custom chip on
board called "Fred" which has patch vectors (10 Approximately). It worked in a
very similar fashion to game genie in terms of patching.

How did it work?
----------------
Basically you would dial up to a central server which would then handshake and
your modem would puke its guts so to speak to the server. It would grab mail, get
player profiles, and any system patches (if there were any). It would also report to
the server how much free space in SRAM it had. Sram was 64kb in size and held game
patches and your profile and all your stats and custom player icons and such and news
and mail.

For matchmaking, if you chose to challenge a player or get a random match, it would
dial the central server, and spent a minute or two downloading a 10kb or smaller game
patch which would be used with the FRED chip to patch in its own network game code
into the existing game on the fly and interact with the modem.

The game networking was quite an impressive feat. It would send controller data back
and fourth between consoles and also sync vblanks between machines. It would switch between
interlaced and non interlaced video modes to slow down the faster consoles to they would
be running in lock-stop as much as possible or at most, a few frames apart (up to about 6-8 max).

However, due to latency, this device is sensitive. Copper phone lines produce stable
latency that was much lower and stable, however today over the internet its much less-predictable and
latency varies quite a bit. In some of the source code docs you'll find reference to about 35-40ms
of latency being the absolute max it can handle (has to be stable). Ive tested this over voip
many a time and this absolutely holds true that it cannot handle anything over that, otherwise
you repeatdly trip the resync code in the game patch and it either lags out or plays
horribly.

Xband Protocol
--------------
The xband network protocol is a modified early version of ADSP Appletalk.
Most if not all packets end with  10 03 (In hex) which equates to DLE/ETX characters.
The first packet received from xband when it initially connects is 26 bytes in length
Eesentially you need to puke the same packet back to the xband but with modifications
to initiate the handshake.

See xbsega.go (sample handshake in go) for an example of how to do this
along with functions to see how it parses whats called the "puke" packet
which comes after the handshake completes after initial connect.

Beyond this, most of the xband protocol once handshakes has to layer the adsp
packet header on top of the actual payload packet you are sending. The XBAND
actual core packet format uses opcodes (byte) that determines the action you want to
do. There are actual server side opcodes and client side opcodes. These can all be
found within the source code provided in the catapult.tar.gz file listed here but
may require some digging. Ensure that when you use/send these values, that you find
the sizes of the data you are going to send as well to the xband. If the packet
length and sizing of the data isnt exactly right, you can spend a bunch of time
waiting for the xband to timeout or crash entirely and do nothing, Its very much
a trial and error thing in some cases. Also the snes is little endian and Sega
is big endian format. Keep this in mind when crafting packets.

See packets.txt file for samples on crafted packets that you have to attach/append
the adsp header packet on top of and also adjust the crc value in the adsp header
to include the full length of the payload.

patching
--------
Xband game patches are hard to come by obviously and without them the gaming side
of xband is a lost cause. Fortunately for sega there is mk2, mk, nba jam, madden95
and nhl94/95 patches in the source provided. Some patches also have source code
as well. For snes, there is only a binary copy of super mario kart with no source
and same goes for a super street fighter 2 which i recovered from a japanese xband
and dumped,tested and confirmed working between xbands for gaming.
