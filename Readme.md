XBAND

The following is a compilation of source code and files obtained over the years
that made the xband modem for the SNES and Sega Genesis what it was.

Xband Hardware
--------------
The xband hardware is a 2400 baud Rockwell chipset modem with a custom chip on
board called "Fred" which has patch vectors (10 Approximately). It worked in a
very similar fashion to game genie in terms of patching.

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

patching
--------
Xband game patches are hard to come by obviously and without them the gaming side
of xband is a lost cause. Fortunately for sega there is mk2, mk, nba jam, madden95
and nhl94/95 patches in the source provided. Some patches also have source code
as well. For snes, there is only a binary copy of super mario kart with no source
and same goes for a super street fighter 2 which i recovered from a japanese xband
and dumped,tested and confirmed working between xbands for gaming.
