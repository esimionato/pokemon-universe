/*Pokemon Universe MMORPG
Copyright (C) 2010 the Pokemon Universe Authors

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.*/
package pubattle

import (
	"pulogic/pokemon"
	pnet "network"
)

type TeamPoke struct {
	UID *UniqueId
	Nick string
	Item int
	Ability int
	Nature int
	Gender int
	Gen int
	Shiny bool
	Happiness int
	Level int
	
	Moves []int
	DVs []int
	EVs []int
}

func NewTeamPoke() *TeamPoke {
	teamPoke := TeamPoke{}
	teamPoke.UID = NewUniqueId()
	teamPoke.Nick = ""
	teamPoke.Item = 0
	teamPoke.Ability = 0
	teamPoke.Nature = 0
	teamPoke.Gender = 0
	teamPoke.Gen = 0
	teamPoke.Shiny = true
	teamPoke.Happiness = 0
	teamPoke.Level = 0
	
	teamPoke.Moves = make([]int, 4)
	teamPoke.Moves[0] = 0
	teamPoke.Moves[1] = 0
	teamPoke.Moves[2] = 0
	teamPoke.Moves[3] = 0
	
	teamPoke.DVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.DVs[i] = 0
	}
	teamPoke.EVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.EVs[i] = 0
	}	
	
	return &teamPoke
}

func NewTeamPokeFromPokemon(_pokemon *pokemon.PlayerPokemon) *TeamPoke {
	teamPoke := TeamPoke{}
	teamPoke.UID = NewUniqueIdExt(_pokemon.Base.PokemonId, 0) // TODO: Add form as SubNum
	teamPoke.Nick = _pokemon.GetNickname()
	teamPoke.Item = 0 // TODO: Add item 
	teamPoke.Ability = _pokemon.Ability.AbilityId
	teamPoke.Nature = _pokemon.Nature
	teamPoke.Gender = _pokemon.Gender
	teamPoke.Gen = 5
	teamPoke.Shiny = (_pokemon.IsShiny == 1)
	teamPoke.Happiness = _pokemon.Happiness
	teamPoke.Level = _pokemon.GetLevel()
	
	teamPoke.Moves = make([]int, 4)
	for i := 0; i < 4; i++ {
		
		if pmove := _pokemon.Moves[i]; pmove != nil {
			teamPoke.Moves[i] = pmove.Move.MoveId
		} else {
			teamPoke.Moves[i] = 0
		}
	}
	
	teamPoke.DVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.DVs[i] = _pokemon.Stats[i]
	}
	teamPoke.EVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.EVs[i] = 42
	}	
	
	return &teamPoke
}

func NewTeamPokeFromPacket(_packet *pnet.QTPacket) *TeamPoke {
	teamPoke := TeamPoke{}
	teamPoke.UID = NewUniqueIdFromPacket(_packet)
	teamPoke.Nick = _packet.ReadString()
	teamPoke.Item = int(_packet.ReadUint16())
	teamPoke.Ability = int(_packet.ReadUint16())
	teamPoke.Nature = int(_packet.ReadUint8())
	teamPoke.Gender = int(_packet.ReadUint8())
	// teamPoke.Gen = (int)_packet.ReadByte()
	teamPoke.Shiny = _packet.ReadBool()
	teamPoke.Happiness = int(_packet.ReadUint8())
	teamPoke.Level = int(_packet.ReadUint8())
	
	teamPoke.Moves = make([]int, 4)
	for i := 0; i < 4; i++ {
		teamPoke.Moves[i] = int(_packet.ReadUint32())
	}
	teamPoke.DVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.DVs[i] = int(_packet.ReadUint8())
	}
	teamPoke.EVs = make([]int, 6)
	for i := 0; i < 6; i++ {
		teamPoke.EVs[i] = int(_packet.ReadUint8())
	}
	
	return &teamPoke
}

func (t *TeamPoke) WritePacket() pnet.IPacket {
	packet := pnet.NewQTPacket()
	uIdPacket := t.UID.WritePacket()
	packet.AddBuffer(uIdPacket.GetBufferSlice())
	packet.AddString(t.Nick)
	packet.AddUint16(uint16(t.Item))
	packet.AddUint16(uint16(t.Ability))
	packet.AddUint8(uint8(t.Nature))
	packet.AddUint8(uint8(t.Gender))
	// packet.AddUint8(uint8(t.Gen)) // XXX Gen would go here
	packet.AddBool(t.Shiny)
	packet.AddUint8(uint8(t.Happiness))
	packet.AddUint8(uint8(t.Level))

	for i := 0; i < 4; i++ {
		packet.AddUint32(uint32(t.Moves[i]))
	}

	for i := 0; i < 6; i++ {
		packet.AddUint8(uint8(t.DVs[i]))
	}
	
	for i := 0; i < 6; i++ {
		packet.AddUint8(uint8(t.EVs[i]))
	}
	
	return packet
}