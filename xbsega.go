package main

import "fmt"
import "log"
import "time"
import "bytes"
import "strings"
import "io/ioutil"
import "encoding/hex"
import "encoding/binary"
import "github.com/tarm/serial"

const DEBUG = true

// KNOWN BOX TYPES
const GENESIS                   string = "segb"
const SATURN			              string = "tj01"
const JSNES			                string = "sj01"
const SNES                      string = "sn07"

// OPCODES BOX SENDS TO SERVER
const msLogin				byte = 0x0b
const msGAMEIDAndPatchVersion		byte = 0x0c
const msChallengeRequest		byte = 0x0e
const msSystemVersion			byte = 0x0f
const msSendNGPVersion			byte = 0x10
const msDBIDInfo			byte = 0x11
const msSendItemFromDB			byte = 0x12
const msSendFirstItemID			byte = 0x13
const msSendNextItemID			byte = 0x14
const msSendSendQElements		byte = 0x15
const msSendAddressesToVerify		byte = 0x16
const msSendNumRankings			byte = 0x17
const msSendFirstRankingID		byte = 0x18
const msSendNextRankingID		byte = 0x19
const msSendRankingData			byte = 0x1a
const msSendInvalidPers			byte = 0x1b
const msSendOutgoingMail		byte = 0x1d
const msSendCreditDebitInfo		byte = 0x1e
const msBoxType 			byte = 0x1f
const msSendGameResults			byte = 0x20
const msSendNoGameResults		byte = 0x21
const msSendConstant			byte = 0x22
const msSendGameErrorResults		byte = 0x23
const msSendNoGameErrorResults		byte = 0x24
const msSendNetErrors			byte = 0x25
const msNoNetErrors			byte = 0x26
const msSendHiddenSerials		byte = 0x27

var options = &serial.Config{ Name: "/dev/tty.usbserial", Baud: 2400 } // CONFIGURE SERIAL PORT AND TARGET DEVICE
var port, err = serial.OpenPort(options)

var rx_buffer = []byte { 0x00 };

var crctab = [256]uint16 {
0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0,
};

var MK_PUKE = []byte { // test packet dump with mortal kombat cartridge connected to modem
0x03,0x00,0x08,0xc8,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x04,0x00,0x00,0x1f,0x73,0x65,0x67,0x62,0x0b,0x3a,0x52,0x3a,0x52,0x63,0x86,0x63,0x86,
0x00,0x00,0x00, 0x1d,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
0x00,0x00,0x00,0x00,0x00,0x08,0x67,0x53,0x09,0xff,0xff,0xff,0xff,0x00,0x00,0x12,0x4e,0x6f,0x6e,0x65,0x00,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,
0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0x24,0x24,0x24,0x24,0x24,0x24,0x24,0x24,
0x24,0x24,0x24,0x24,0x24,0x24,0x24,0x24,0x00,0x39,0x73,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x00,0x6e,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0xfa,0xce,
0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0x62,0x4a,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x00,0x80,0x00,0x00,
0x00,0x00,0x04,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x1e,0x02,0x0c,0xab,0x63,0x48,0xe9,0xff,0xff,0xff,0xff,0x0f,0x00,0x04,0x00,0x00,0x00,0x00,
0x0e,0x03,0x10,0x10,0x00,0x00,0x00,0x01,0x15,0x00,0x02,0x03,0x00,0x00,0x00,0x00,0x01,0x00,0x00,0x00,0xb0,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xcd,0x47,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x00,0xee,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0x00,0x00,0x0f,0xc7,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0x76,0x21,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x01,0x5c,0x00,0x00,0x00,0x00,0x04,0x00,
0x00,0x16,0x00,0x00,0x1d,0x00,0x00,0x21,0x24,0x25,0x00,0xff,0xff,0xff,0xff,0xff,0x00,0x00,0xff,0x02,0xff,0xff,0xff,0xff,0xff,0xff,0x00,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0x25,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0x1b,0x5d,0x00,0x1b,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,
0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x25,0x00,0x00,0x82,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0xb7,0x27,0x10,0x03,0x00,
0x08,0xc8,0x00,0x00,0x01,0xca,0x00,0x00,0x00,0x00,0x04,0x00,0x00,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,
0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,
0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,
0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x40,
0x40,0x40,0x40,0x40,0x40,0x40,0x40,0x63,0xe8,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x02,0x38,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0x40,0x40,0x40,0x40,
0x40,0x40,0x40,0x40,0x00,0xde,0x18,0x10,0x03,0x00,0x08,0xc8,0x00,0x00,0x02,0x41,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0x12,0x00,0x69,0x4c,0x10,0x03,
0x00,0x08,0xc8,0x00,0x00,0x02,0x43,0x00,0x00,0x00,0x00,0x04,0x00,0xc0,0xcc,0xc8,0x10,0x03 };

var MK2_PUKE = []byte { // test packet dump with mortal kombat 2 cartridge connected to modem
0x00,0x2e,0x9c,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x04,0x00,0x00,0x1f,0x73,0x65,0x67,0x62,0x0b,0x3a,0x52,0x3a,0x52,0x64,0x5c,0x64,0x5c,0x00,
0x00,0x00, 0x00,0x00,0x00,0x00,0x01,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
0x00,0x00,0x00,0x00,0x11,0x11,0x11,0x11,0xFF,0xAA,0xCC,0xEE,0x00,0x00,0x12,0x4e,0x6f,0x6e,0x65,0x00,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,
0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0x54,0x65,0x73,0x74,0x20,0x31,0x32,0x33,0x00,
0x72,0x00,0xde,0xad,0xfa,0xce,0xde,0xad,0x05,0xb6,0x10,0x03,0x00,0x2e,0x9c,0x00,0x00,0x00,0x6e,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0xfa,0xce,0xde,
0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0xde,0xad,0xfa,0xce,0x63,0xe1,0x10,0x03,0x00,0x2e,0x9c,0x00,0x00,0x00,0x80,0x00,0x00,0x00,
0x00,0x04,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x1e,0x02,0x0c,0xc4,0xcd,0xdf,0x0c,0xff,0xff,0xff,0xff,0x0f,0x00,0x04,0x00,0x00,0x00,0x00,0x0e,
0x03,0x10,0x10,0x00,0x00,0x00,0x01,0x15,0x00,0x01,0x01,0x00,0x00,0x00,0xb0,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xdb,
0x09,0x10,0x03,0x00,0x2e,0x9c,0x00,0x00,0x00,0xee,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0x00,0x00,0x0f,0xc7,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0xff,0xff,0xff,0x47,0x13,0x10,0x03,0x00,0x2e,0x9c,0x00,0x00,0x01,0x57,0x00,0x00,0x00,0x00,0x04,0x00,0x00,0x16,0x00,0x00,0x1d,0x00,
0x00,0x21,0x24,0x25,0x02,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0x01,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0x25,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
0xff,0xff,0xff,0x1b,0x5d,0x00,0x23,0x50,0x69,0x63,0x6b,0x20,0x6f,0x6e,0x20,0x6d,0x65,0x21,0x20,0x49,0x20,0x68,0x61,0x76,0x65,0x6e,0xd5,0x74,0x20,
0x67,0x6f,0x74,0x20,0x61,0x20,0x74,0x61,0x75,0x6e,0x74,0x21,0x00,0x00,0x12,0x49,0xd5,0x6d,0x61,0x08,0x10,0x03,0x00,0x2e,0x9c,0x00,0x00,0x01,0xc5,
0x00,0x00,0x00,0x00,0x04,0x00,0x40,0x20,0x61,0x20,0x6e,0x65,0x77,0x20,0x70,0x6c,0x61,0x79,0x65,0x72,0x21,0x00,0x3e,0x2e,0x10,0x03,0x00,0x2e,0x9c,
0x00,0x00,0x01,0xd4,0x00,0x00,0x00,0x00,0x04,0x00,0x40,0x12,0x00,0xf8,0xf6,0x10,0x03 };

func Updcrc(icrc uint16, buffer []uint8, start uint, count uint) uint16 {
for end := start + count; start < end; start++ { icrc = ((icrc << 8) & 0xff00) ^ crctab[((icrc>>8)&0x00ff)^uint16(buffer[start])] }
return icrc
}

func Send_Message(DATA []byte) {
// ----------------------------------------------------------------------------------------------
fmt.Print("PAYLOD PASSED IN: "); fmt.Printf("%x",DATA); fmt.Println("\n");
// ----------------------------------------------------------------------------------------------
ADSP_HEADER    := []byte { 0x00,0xDE,0xAD,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x04,0x00,0x80,0x01,0x00,rx_buffer[1],rx_buffer[2],0x00,0x00,0x00,0x00,0xAA,0xAA,0x10,0x03 };
var crc         = Updcrc(0xFFFF,ADSP_HEADER,0,22); crc = ^crc 		// CRC FIRST 22 BYTES OF PACKET
ADSP_HEADER[22] = byte(crc>>8); ADSP_HEADER[23] = byte(crc);    // INSERT THE NEW CRC INTO PACKET_END
// ----------------------------------------------------------------------------------------------
PAYLOD_HEADER  := []byte { 0x00,0xBE,0xEF,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x04,0x00,0x00 };
PAYLOD         := append(PAYLOD_HEADER,DATA...) 			       // LINK PAYLOAD HEADER WITH PASSED IN PAYLOAD VALUE
PAYLOD_END     := []byte { 0xAA,0xAA,0x10,0x03 };
crc 	        = Updcrc(0xFFFF,PAYLOD,0,uint(len(PAYLOD))); crc = ^crc        // CRC HEADER + PAYLOAD DATA
PAYLOD_END[0]   = byte(crc>>8); PAYLOD_END[1] = byte(crc);                     // PLUG IN NEW CRC INTO PACKET_END (FIRST TWO BYTES)
// ----------------------------------------------------------------------------------------------
DATA_PACKET    := append(ADSP_HEADER, append(PAYLOD,PAYLOD_END...)...) // MERGE EVERYTHING TOGETHER INTO 1 PACKET
// ----------------------------------------------------------------------------------------------
fmt.Printf("%X\n",DATA_PACKET)
port.Write([]byte(DATA_PACKET))
fmt.Println("\nSent Packet!\n")
return;
}

//********************************************************
//********************************************************
//*** MAIN PROGRAM LOOP                                 **
//********************************************************
//********************************************************

func main() {
  println("\n*************************************")
  println("**         | XBSEGA SERVER |         **")
  println("*************************************\n")

  if DEBUG == false {

  if err != nil {                                                                 // IF ERROR FROM OPENING PORT...
  println("!!!   CANNOT FIND SERIAL DEVICE   !!!\n\a\a\a")
  println("     Press CTRL+C To break/Exit")
  for { ;; }                                                                      // LOOP FOREVER UNTIL USER BREAKS OUT OF IT
  log.Fatal(err);                                                                 // LOG THE ERROR
  } else { println("SERIAL INITIALIZED!"); }                                      // OTHERWISE PRINT TEXT STATING SUCCESS

  rx_buffer := make([]byte, 128)                                                  // MAKE RX BUFFER FOR MODEM STATUS & AT COMMAND RESPONSES & CALLER ID DATA
  println("CHECKING MODEM...")

  port.Write([]byte("AT\r"))                                                      // SEND "ATTENTION" COMMAND
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
  port.Read(rx_buffer);                                                           // READ SERIAL DATA INTO THE BUFFER

  if strings.Contains(string(rx_buffer), "OK" ) {                                 // SEE IF BUFFER CONTAINS RESPONS OF "OK"
  println("MODEM DETECTED!")                                                      // PRINT TEXT
  println("CONFIGURING...");
  port.Write([]byte("AT&F\r"))                                                    // SET MODEM TO FACTORY DEFAULT
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
  port.Write([]byte("ATQ0\r"))                                                    // ENABLE STATUS CODES FROM MODEM
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
  port.Write([]byte("AT+VCID=1\r"))                                               // ENABLE CALLER ID OUTPUT
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
  port.Write([]byte("ATS0=2\r"))                                                  // AUTO-ANSWER ON 2ND RING
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
	port.Write([]byte("ATM0\r"))                                                    // TURN OFF MODEM SPEAKER
  time.Sleep(time.Millisecond * 200)                                              // DELAY FOR PROCESSING
  println("MODEM CONFIGURED!\n")                                                  // PRINT TEXT
  }

	println("---------- WAITING FOR CALL ---------")

  for {
        port.Read(rx_buffer);                                                     // READ DATA INTO BUFFER
        if strings.Contains(string(rx_buffer), "CONNECT" ) { println("------------ CONNECTING... ----------\n");
        port.Write([]byte("\r"));                                                 // SEND CARRIDGE RETURN TO XBAND (LET BOX KNOW WE ARE READY TO TALK)
        break;                                                                    // EXIT FOR LOOP
        }
        time.Sleep(time.Millisecond * 200); // DELAY BETWEEN SERIAL PORT READS
      }

} // end of debug IF statement bracket defined above

  rx_buffer = make([]byte,26)                   // make 26 byte buffer slice for receiving initial packet from box (handshake packet)
  time.Sleep(time.Second * 2)                   // give time for serial data to come in (MIGHT NEED TO ADJUST FOR PRODUCTION TO LONGER TIME)
	if DEBUG == false { port.Read(rx_buffer);	}		// READ 26 BYTES INTO BUFFER

  println("------------- CONNECTED -------------\a")

  // ------------------------------------------
  // CRC PACKET AND SEND RESPONSE
  // ------------------------------------------

  rx_buffer[0]  = 0x00                                         // MAKE SURE FIRST BYTE SHOWS 00
  rx_buffer[13] = 0x82                                         // SET ACK FLAG IN ADSP HEADER
  var crc       = Updcrc(0xFFFF,rx_buffer,0,22); crc=^crc      // CRC THE FIRST 22 BYTES (EVERYTHING UP TO THE CRC AND DLE/ETX BYTES (0x10,0x03) )
  rx_buffer[22] = byte(crc>>8); rx_buffer[23] = byte(crc);     // APPLY NEW CRC TO PACKET
  if DEBUG == false { port.Write([]byte(rx_buffer))	}      // SEND THE RESPONSE FOR HANDSHAKING TO COMPLETE.

  // ------------------------------------------
  // RECEIVE INITIAL DUMP PACKET FROM BOX
  // ------------------------------------------

  rx_buffer = make([]byte,2048)						     // DEFINE 2KB BUFFER FOR RECEIVING DATA.
  if DEBUG == false { time.Sleep(time.Second * 5); port.Read(rx_buffer); }   // WAIT 6 SECOND & READ ALL AVAILABLE SERIAL DATA INTO THE BUFFER
  if DEBUG == true { rx_buffer = MK2_PUKE }			             // IF DEBUG MODE, SPECIFY PACKET DUMP TO USE IN RX_BUFFER

  // ------------------------------------------
  // DECODE LARGE DUMP PACKET
  // ------------------------------------------

  index        := bytes.IndexByte(rx_buffer,msBoxType)+1                       // POINT TO OPCODE FOR BOXTYPE
  BOXTYPE      := string(rx_buffer[index:index+4]);			       // GET BOX TYPE AS A STRING

  index         = bytes.IndexByte(rx_buffer,msLogin)+1                         // point to msLogin opcode+1 (start of data after opcode)
  OSFREE       := binary.BigEndian.Uint16(rx_buffer[index:index+4]); index+=4; // convert bytes to decimal so we can show numbers
  DBFREE       := binary.BigEndian.Uint16(rx_buffer[index:index+4]); index+=4; // convert bytes to decimal so we can show numbers
  MISCFLAGS    := binary.BigEndian.Uint16(rx_buffer[index+2:index+4]); index+=4;
  LASTBOXSTATE := binary.BigEndian.Uint16(rx_buffer[index:index+4]); index+=4;
  PHONENUMBER  := binary.BigEndian.Uint16(rx_buffer[index:index+26]);index+=26;

  REGION       := binary.BigEndian.Uint16(rx_buffer[index:index+4]); index+=4
  SERIAL       := binary.BigEndian.Uint16(rx_buffer[index:index+4]); index+=4

  BOXID        := rx_buffer[index]; index+=1;                                    // 0-3 (PROFILE 1-4) (1 BYTE)
  CLUTID       := rx_buffer[index]; index+=1;                                    // COLOR LOOK UP TABLE INDEX VALUE (1-BYTE)
  ICONID       := rx_buffer[index]; index+=1;                                    // ICON ID NUMBER (1 byte)

  temp 	       := rx_buffer[index:index+34]; temp2 := bytes.IndexByte(temp,0x00) // Point to start and end of string with null termination
  HOMETOWN     := string(rx_buffer[index:index+temp2]); index+=34;               // Only print data from index to null byte in string

  temp 	        = rx_buffer[index:index+34]; temp2 = bytes.IndexByte(temp,0x00)  // Point to start and end of string with null termination
  USERNAME     := string(rx_buffer[index:index+temp2]);                          // Set variable but dont increment index value due to logic below

  if BOXTYPE == GENESIS { index+=35; }                                           // Increment index based on boxtype (they're slightly different)
  if BOXTYPE == SNES || BOXTYPE == JSNES { index+=34; }

  // NUM XMAILS ISNT RIGHT...
  XMAILS    := binary.BigEndian.Uint16( rx_buffer[index:index+4] ); index+=4;

  // PERSONIFICATION STUFF
  // ---------------------
  index     = bytes.IndexByte(rx_buffer,msSendInvalidPers)+1

  PASSWORD := []byte { 0x00 }; // SET PASSWORD AS NOTHING INITIALLY.
  // NOW WE SEE IF 0x5d is present after the msSendInvalidPers opcode to see if a password is even set. if so, increment index 2 bytes
  // and read the 8 following (because thats supposedly the password). if no 0x5d, then just skip to the start of the taunt
  if rx_buffer[index] == 0x5d { index+=3; } else { PASSWORD = rx_buffer[index:index+8]; index+=8; } // THIS NEEDS TO BE VETTED!!

  temp      = rx_buffer[index:index+55]; temp2 = bytes.IndexByte(temp,0x00)
  TAUNT    := string(rx_buffer[index:index+temp2]) // convert starting byte of username until null byte into string for use

  temp      = rx_buffer[index:index+55]; temp2 = bytes.IndexByte(temp,0x00)
  ABOUT    := string(rx_buffer[index:index+temp2]) // convert starting byte of username until null byte into string for use

  // ------------------------------------------
  // GET GAME ID PLUGGED ONTOP OF XB MODEM
  // ------------------------------------------

  index     = bytes.IndexByte(rx_buffer, msGAMEIDAndPatchVersion )
  GAMEID   := "00000000"
  GAMENAME := "NONE";
  var PATCH_DATA string = "";

  if BOXTYPE == GENESIS {  // KNOWN (U) SEGA GENESIS GAMES
  switch GAMEID = hex.EncodeToString(rx_buffer[index+1:index+5])
  GAMEID { case "31ed8123": GAMENAME = "Madden 95"; PATCH_DATA = "/xband/patches/segb/nfl95.mp" // POINT TO PATCH FILE
           case "ab6348e9": GAMENAME = "Mortal Kombat"; PATCH_DATA = "/xband/patches/segb/mk.mp" // POINT TO PATCH FILE
           case "c4cddf0c": GAMENAME = "Mortal Kombat II"; PATCH_DATA = "/xband/patches/segb/mk2.mp" // POINT TO PATCH FILE
           case "e30c296e": GAMENAME = "NBA JAM [Rev 1]"; PATCH_DATA = "/xband/patches/segb/nbajam.mp" // POINT TO PATCH FILE
        // case "39677bdb": GAMENAME = "NBA JAM [Rev 2]"
        // case "a61b53f8": GAMENAME = "NHL 94"
           case "8f6b9f70": GAMENAME = "NHL 95"; PATCH_DATA = "/xband/patches/segb/nhl95.mp" // POINT TO PATCH FILE
        // case "3fed23f2": GAMENAME = "Road Rash 3"
        // case "00192660": GAMENAME = "NBA Live 95"
        // case "433e2840": GAMENAME = "FIFA Soccer 95"
        // case "4a017a94": GAMENAME = "WeaponLord [Ver 1]"
        // case "bf33efc7": GAMENAME = "WeaponLord [Ver 2]"
        // case "c6906e52": GAMENAME = "Primal Rage"
        // case "51a5e383": GAMENAME = "Rampart"
        // case "067a218f": GAMENAME = "Ballz"
        // case "4d402d90": GAMENAME = "Madden 96"
        // case "afc0ce39": GAMENAME = "NHL 96"
        // case "6d14eb41": GAMENAME = "Mortal Kombat 3"
        // case "51a5e383": GAMENAME = "Rampart"
        // case "4d1c4e1d": GAMENAME = "Super Street Fighter II"
           default: 			  GAMENAME = "UNKNOWN";
        }
  }

  if BOXTYPE == SNES {  // KNOWN SUPER NINTENDO GAMES
  switch GAMEID = hex.EncodeToString(rx_buffer[index+1:index+5])
  GAMEID { // case "c4cddf0c": GAMENAME = "Mortal Kombat II [Rev 1]";
        // case "c0432172": GAMENAME = "Mortal Kombat II [Rev 2]"
        // case "127e8181": GAMENAME = "NHL 95"
        // case "1969d2af": GAMENAME = "NBA JAM Tournament Edition"
        // case "ef120a61": GAMENAME = "Super Street Fighter II"
        // case "b8958396": GAMENAME = "Madden 95"
        // case "972404cc": GAMENAME = "FIFA Int'l Soccer"
        // case "3d1c44eb": GAMENAME = "Super Mario Kart"
        // case "19a2c936": GAMENAME = "NBA Live 95"
        // case "0572a585": GAMENAME = "WeaponLord [Rev 1]"
        // case "0572dd87": GAMENAME = "WeaponLord [Rev 2]"
        // case "a8973c8c": GAMENAME = "Ken Griffey Baseball"
        // case "2d17c045": GAMENAME = "Killer Instinct"
        // case "085d3cdb": GAMENAME = "Madden 96"
        // case "25f372a5": GAMENAME = "NHL 96"
        // case "94b564b5": GAMENAME = "DOOM"
        // case "05484971": GAMENAME = "Mortal Kombat 3"
        // case "83e627ef": GAMENAME = "Kirby's Avalanche"
        // case "8f6b9f70": GAMENAME = "Zelda: A Link to the Past"
           default: 				GAMENAME = "UNKNOWN";
         }
  }

  if BOXTYPE == JSNES {  // KNOWN (J) SUPER NINTENDO GAMES
  switch GAMEID = hex.EncodeToString(rx_buffer[index+1:index+5])
  GAMEID { // case "d8222103": GAMENAME = "Super Street Fighter II"
        // case "0a2c238a": GAMENAME = "Super Mario Kart"
        // case "925b41fc": GAMENAME = "Super Fire ProWrestling X"
           default: 			  GAMENAME = "UNKNOWN";
         }
  }

  if BOXTYPE == SATURN {  // KNOWN (J) SEGA SATURN GAMES
  switch GAMEID = hex.EncodeToString(rx_buffer[index+1:index+5])
  GAMEID { // case "00000000": GAMENAME = "Decathlete"
        // case "00000001": GAMENAME = "Virtual On"
        // case "00000002": GAMENAME = "Puyo Puyo Sun"
        // case "00000003": GAMENAME = "Puzzle Bobble 3"
        // case "00000004": GAMENAME = "Saturn Bomberman"
        // case "00000005": GAMENAME = "Virtua Fighter Remix"
        // case "00000006": GAMENAME = "World Series Baseball"
        // case "00000007": GAMENAME = "Sega Worldwide Soccer '98"
        // case "00000008": GAMENAME = "Sega Rally Championship Plus"
        // case "00000009": GAMENAME = "Daytona USA Championship Circuit Edition"
           default        : GAMENAME = "UNKNOWN";
         }
  }

  // ------------------------------------------
  // PRINT OUT RECEIVED INFORMATION FROM BOX
  // ------------------------------------------

  fmt.Print("[BOX TYPE]        : "); fmt.Println(BOXTYPE);
  fmt.Print("[OS FREE MEM]     : "); fmt.Print(OSFREE); fmt.Println(" Bytes");
  fmt.Print("[DB FREE MEM]     : "); fmt.Print(DBFREE); fmt.Println(" Bytes");
  fmt.Print("[BOX FLAGS]       : "); fmt.Printf("%X",MISCFLAGS); fmt.Print("\n") // SHOW IN HEX FOR DEBUG PURPOSES
  fmt.Print("[LAST STATE]      : "); fmt.Println(LASTBOXSTATE);
  fmt.Print("[BOX PHONE NUMBER]: "); fmt.Println(PHONENUMBER);
  fmt.Print("[REGION]          : "); fmt.Println(REGION);
  fmt.Print("[SERIAL]          : "); fmt.Println(SERIAL);
  fmt.Print("[PROFILE NUMBER]  : "); fmt.Println(BOXID);
  fmt.Print("[CLUT INDEX]      : "); fmt.Println(CLUTID);
  fmt.Print("[ICON ID]         : "); fmt.Println(ICONID);
  fmt.Print("[HOME TOWN]       : "); fmt.Println(HOMETOWN);
  fmt.Print("[USER NAME]       : "); fmt.Println(USERNAME);
  fmt.Print("[XMAILS]          : "); fmt.Println(XMAILS);
  fmt.Print("[PASSWORD]        : "); fmt.Println(PASSWORD);
  fmt.Print("[TAUNT]           : "); fmt.Println(TAUNT);
  fmt.Print("[ABOUT]           : "); fmt.Println(ABOUT);
  fmt.Print("[GAME]            : "+GAMENAME); if GAMENAME == "UNKNOWN" { fmt.Println(" "+GAMEID) } else { fmt.Print("\n") }
  println("-------------------------------------\n")

  // ------------------------------------------
  // SEND TEST PACKET TO BOX
  // ------------------------------------------

  OpcodeStream := "0E0110000002" // set box to wait for call for 11000000 ticks (longword value) + 02 end stream opcode
  Send_Message(OpcodeStream); // send it...

  if DEBUG == false {
  println("----------- DISCONNECTED ------------\n\a")
  port.Write([]byte("ATH\r"));
  time.Sleep(time.Millisecond * 200)
  port.Close();
  }

// END OF PROGRAM
}