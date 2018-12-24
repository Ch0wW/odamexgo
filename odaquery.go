package odamexgo

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// WadInfo : all WAD data returned from the server
type WadInfo struct {
	Name string
	Hash string
}

// PlayerInfo : all data from a player.
type PlayerInfo struct {
	Name   string
	Frags  int16
	Deaths int16
	Points int16
	Team   byte

	Spectator bool
	Time      int16
	Ping      int32
}

// CVARInfo : handles CVARs
type CVARInfo struct {
	Name  string
	Value bool
}

// TeamInfo : data for teaminfos.
type TeamInfo struct {
	Points int32
}

/*
ServerInfo : Handles all the data of an Odamex server
*/
type ServerInfo struct {
	ip, port string

	Challenge int32
	Token     int32
	Protocol  int16
	Version   int32

	Hostname    string `json:"Hostname"`
	Website     string
	Email       string
	HasPassword bool

	Fraglimit  int16
	Timelimit  int16
	TimeLeft   int16
	Scorelimit int32

	PlayersInGame byte
	Spectators    int
	MaxClients    byte
	MaxPlayers    int16
	PlayerList    []PlayerInfo
	Teams         [2]TeamInfo

	MapName      string
	IsDeathmatch bool
	Skill        byte
	IsTeamDM     bool
	IsCTF        bool

	WadLength byte
	WadList   []WadInfo
	CVARList  []CVARInfo
	PatchList []string
}

var (
	remIP   string // Remote ip
	remPort string // Remote Port
)

// ServerQuery struct handles all of the data of a server
type ServerQuery struct {
	ip     string
	port   string
	server ServerInfo

	buffer   []byte
	position int
}

const (
	dataquery = "\xA3\xDB\x0B\x00" // If we convert to Little Endian, it goes by 777123
)

/*
NewOdaURI parses an odamex:// URI and transforms it into a ServerQuery .
*/
func NewOdaURI(link string) (*ServerQuery, error) {

	var strAddr []string
	var strPort []string

	s := &ServerQuery{}

	if link != "" {
		if strings.HasPrefix(link, "odamex://") {
			fulllink := link[9:]

			// 1) Avoid repetitions and check if we have a port parameter.
			if strings.Contains(fulllink, ":") {
				strAddr = strings.Split(fulllink, ":")
				s.ip = strAddr[0]

				// 1.2) Check if there is a '/' (because QWURL used it sometimes, ODA may also used it at some time?)
				if strings.Contains(strAddr[1], "/") {
					strPort = strings.Split(strAddr[1], "/")
					s.port = strPort[0]
				} else {
					s.port = strAddr[1]
				}

				return s, nil
			}

			// if no port, assume the default port is 10666
			s.ip = fulllink
			s.port = "10666"

			return s, nil

		}
	}

	return nil, fmt.Errorf("Odamex link is not valid (argument should be \"odamex://<ip>[:<port>])\"")
}

// IsCooperation : Checks if the server is a Cooperation Game.
func IsCooperation(server ServerInfo) bool {
	if !server.IsDeathmatch && !server.IsTeamDM && !server.IsCTF {
		return true
	}
	return false
}

// IsDeathmatch : Checks if the server is a Deathmatch Game.
func IsDeathmatch(server ServerInfo) bool {
	if server.IsDeathmatch && !server.IsTeamDM && !server.IsCTF {
		return true
	}
	return false
}

// IsTeamGame : Checks if the server is a Teambased Game.
func IsTeamGame(server ServerInfo) bool {
	if server.IsTeamDM || server.IsCTF {
		return true
	}
	return false
}

// AddCVAR : Adds a CVAR to the list of CVARs after parsing it.
func (s *ServerQuery) AddCVAR(list []CVARInfo, name string) []CVARInfo {

	boolValue := s.ReadBool()

	cvar := CVARInfo{Name: name, Value: boolValue}
	list = append(list, cvar)

	return list
}

// ParseOdamex070 : Manually parses an Odamex 0.7.X server.
func (s *ServerQuery) ParseOdamex070() ServerInfo {

	var sv ServerInfo

	sv.ip = s.ip
	sv.port = s.port

	sv.Challenge = s.ReadLong()
	sv.Token = s.ReadLong()
	sv.Hostname = s.ReadString()
	sv.PlayersInGame, _ = s.ReadByte()
	sv.MaxClients, _ = s.ReadByte()
	sv.MapName = s.ReadString()
	sv.WadLength, _ = s.ReadByte()

	sv.WadList = make([]WadInfo, sv.WadLength)
	for iWads := 0; iWads < int(sv.WadLength); iWads++ {
		sv.WadList[iWads].Name = s.ReadString()
	}

	sv.IsDeathmatch = s.ReadBool()
	sv.Skill, _ = s.ReadByte()
	sv.IsTeamDM = s.ReadBool()
	sv.IsCTF = s.ReadBool()

	// =========== OUR PROBLEM STARTS NOW
	sv.PlayerList = make([]PlayerInfo, sv.PlayersInGame)
	for iPlayer := 0; iPlayer < int(sv.PlayersInGame); iPlayer++ {

		sv.PlayerList[iPlayer].Name = s.ReadString()
		sv.PlayerList[iPlayer].Frags = s.ReadShort()

		sv.PlayerList[iPlayer].Ping = s.ReadLong()

		if sv.IsCTF || sv.IsTeamDM {
			sv.PlayerList[iPlayer].Team, _ = s.ReadByte()
		} else {
			sv.PlayerList[iPlayer].Team = 3 // TEAM_NONE
		}
	}

	for iWads := 0; iWads < int(sv.WadLength); iWads++ {
		sv.WadList[iWads].Hash = s.ReadString()
	}

	sv.Website = s.ReadString()

	// Take care of the Team Infos
	if sv.IsCTF || sv.IsTeamDM {

		sv.Scorelimit = s.ReadLong()

		var bIsTeam uint8
		for iTeam := 0; iTeam < 2; iTeam++ {
			bIsTeam, _ = s.ReadByte()

			if bIsTeam == 1 {
				sv.Teams[iTeam].Points = s.ReadLong()
			}
		}
	}

	sv.Protocol = s.ReadShort()

	sv.Email = s.ReadString()

	sv.Timelimit = s.ReadShort()
	sv.TimeLeft = s.ReadShort()
	sv.Fraglimit = s.ReadShort()

	// CVARs
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_itemsrespawn")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_weaponstay")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_friendlyfire")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_allowexit")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_infiniteammo")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_nomonsters")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_monstersrespawn")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_fastmonsters")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_allowjump")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_freelook")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_waddownload")
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_emptyreset")
	_ = s.ReadBool()
	sv.CVARList = s.AddCVAR(sv.CVARList, "sv_fragexitswitch")

	for iPlayer := 0; iPlayer < int(sv.PlayersInGame); iPlayer++ {
		sv.PlayerList[iPlayer].Points = s.ReadShort()
		sv.PlayerList[iPlayer].Deaths = s.ReadShort()

		sv.PlayerList[iPlayer].Time = s.ReadShort()
	}

	_ = s.ReadLong()

	sv.MaxPlayers = s.ReadShort()

	// Spectator
	iSpectators := 0
	for iPlayer := 0; iPlayer < int(sv.PlayersInGame); iPlayer++ {
		sv.PlayerList[iPlayer].Spectator = s.ReadBool()

		if sv.PlayerList[iPlayer].Spectator == true {
			iSpectators = iSpectators + 1
		}
	}

	sv.Spectators = iSpectators

	_ = s.ReadLong()

	pass := s.ReadShort()
	if pass == 1 {
		sv.HasPassword = true
	} else {
		sv.HasPassword = false
	}

	sv.Version = s.ReadLong()

	// Check patch sizes
	PatchSize, _ := s.ReadByte()
	if PatchSize == 0 {
	} else {
		for iPatch := 0; iPatch < int(PatchSize); iPatch++ {
			sv.PatchList[iPatch] = s.ReadString()
		}
	}
	return sv
}

// GetServerInfo : Connects to an Odamex server and parses the server.
func (s *ServerQuery) GetServerInfo() (*ServerInfo, error) {

	// Translate DNS into a readable IP
	daIP, err := net.LookupIP(s.ip)
	if err != nil {
		fmt.Println("Unknown host")
	}
	s.ip = daIP[0].String()

	svlink := s.ip + ":" + s.port

	//fmt.Println(binary.LittleEndian.Uint32([]byte(dataquery)))

	//Connect udp
	conn, err := net.DialTimeout("udp", svlink, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("cannot access the server: %s", err)
	}
	defer conn.Close()

	// Query the server to check if we're a valid QW server
	_, err = conn.Write([]byte(dataquery))
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("Write Timeout: %s", err)
		}
		return nil, fmt.Errorf("write Error: %s", err)
	}

	// Read the answer and trim it, so that empty bytes won't be displayed.
	buffer := make([]byte, 8196)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	buffersize, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, fmt.Errorf("read timeout: %s", err)
		}
		return nil, fmt.Errorf("read Error: %s", err)
	}

	if buffersize <= 0 {
		return nil, fmt.Errorf("server has no data to answer with")
	}

	// Copy the whole buffer to the ServerQuery struct
	s.buffer = buffer[0:buffersize]

	sv := s.ParseOdamex070()

	return &sv, nil
}
