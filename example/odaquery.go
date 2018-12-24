package main

import (
	"fmt"

	"github.com/ch0ww/odamexgo"
)

/*
===================
MAIN FUNCTION
===================
*/
func main() {

	odasv, err := odamexgo.NewOdaURI("odamex://74.91.112.85:10669")

	if err != nil {
		fmt.Println(err)
		return
	}

	sv, err := odasv.GetServerInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	//================================================================================
	// Printing INFOS

	fmt.Println(sv.Ip + sv.Port)
	fmt.Println("")
	fmt.Println("==============================")
	fmt.Println("[SERVER INFO]")
	fmt.Println("HOSTNAME :", sv.Hostname)
	fmt.Println("Players :", sv.PlayersInGame, "(", sv.Spectators, ") /", sv.MaxClients)
	fmt.Println("Max. Clients :", sv.MaxClients)
	fmt.Println("MapName :", sv.MapName)
	fmt.Println("Website :", sv.Website)
	fmt.Println("Contact :", sv.Email)
	fmt.Println("Password-Protected :", sv.HasPassword)

	sGametype := "Cooperation"
	if sv.IsDeathmatch {
		sGametype = "Deathmatch"
	} else if sv.IsTeamDM {
		sGametype = "Team Deathmatch"
	} else if sv.IsCTF {
		sGametype = "Capture the Flag"
	}

	fmt.Println("Gamemode :", sGametype)
	fmt.Println("Skill :", sv.Skill)

	fmt.Println("")
	fmt.Println("==============================")
	fmt.Println("[WADS]")
	for iWads := 0; iWads < int(sv.WadLength); iWads++ {
		fmt.Println("-", sv.WadList[iWads].Name, "("+sv.WadList[iWads].Hash+")")
	}

	if len(sv.PlayerList) > 0 {
		fmt.Println("")
		fmt.Println("==============================")
		fmt.Println("[PLAYERS]")
		for iPlayers := 0; iPlayers < len(sv.PlayerList); iPlayers++ {
			fmt.Println("-", sv.PlayerList[iPlayers].Name, sv.PlayerList[iPlayers].Frags, sv.PlayerList[iPlayers].Deaths, sv.PlayerList[iPlayers].Points, sv.PlayerList[iPlayers].Time, sv.PlayerList[iPlayers].Ping, sv.PlayerList[iPlayers].Spectator)
		}
	}

	fmt.Println("")
	fmt.Println("==============================")
	fmt.Println("[CVARs]")
	for iCVAR := 0; iCVAR < len(sv.CVARList); iCVAR++ {
		fmt.Println(sv.CVARList[iCVAR].Name, ":", sv.CVARList[iCVAR].Value)
	}

	if len(sv.PatchList) > 0 {
		fmt.Println("")
		fmt.Println("==============================")
		fmt.Println("[PATCHs]")

		for i := 0; i < len(sv.PatchList); i++ {
			fmt.Println("-", sv.PatchList[i])
		}
	}

}
