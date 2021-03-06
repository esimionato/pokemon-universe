package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"net"
	"strings"
	
	pos "nonamelib/pos"
	pul "pulogic"
	
	"pulogic/models"
	"nonamelib/log"
)

var AutoClientId int = 0

type Client struct {
	id int
	socket   net.Conn
	loggedIn bool

	changeTileChan chan *Packet
}

func NewClient(_socket net.Conn, _changeTileChan chan *Packet) *Client {
	AutoClientId++
	return &Client{id: AutoClientId, socket: _socket, loggedIn: false, changeTileChan: _changeTileChan}
}

func (c *Client) HandleClient() {
	for {
		packet := NewPacket()
		var headerbuffer [2]uint8
		recv, err := io.ReadFull(c.socket, headerbuffer[0:])
		if err != nil || recv == 0 {
			fmt.Printf("Disconnected: %d\n", c.id)
			break
		}
		copy(packet.Buffer[0:2], headerbuffer[0:2])
		packet.GetHeader()

		databuffer := make([]uint8, packet.MsgSize)

		reloop := false
		bytesReceived := uint16(0)
		for bytesReceived < packet.MsgSize {
			recv, err = io.ReadFull(c.socket, databuffer[bytesReceived:])
			if recv == 0 {
				reloop = true
				break
			} else if err != nil {
				fmt.Printf("Connection read error: %v\n", err)
				reloop = true
				break
			}
			bytesReceived += uint16(recv)
		}
		if reloop {
			continue
		}

		copy(packet.Buffer[2:], databuffer[:])

		header := packet.ReadUint8()
		
		if !c.loggedIn && header != HEADER_LOGIN {
			fmt.Println("Received header but user is not logged in!")
			continue
		}
		
		switch header {
		case HEADER_LOGIN: // Login
			c.ReceiveLogin(packet)
		case HEADER_REQUEST_MAP_PIECE: // Request map piece
			go c.RequestMapPiece(packet)

		case HEADER_TILE_CHANGE: // Tile changes
			c.changeTileChan <- packet
			
		case HEADER_REQUEST_MAP_LIST: // Request map list
			go c.SendMapList()
			
		case HEADER_ADD_MAP: // Add map
			go c.ReceiveAddMap(packet)
			
		case HEADER_DELETE_MAP: // Delete map
			go c.ReceiveRemoveMap(packet)
			
		case HEADER_UPDATE_TILEEVENT: // Update tile event
			go c.ReceiveTileEventUpdate(packet)
			
		case HEADER_ADD_NPC: // Add Npc
			go c.ReceiveAddNpc(packet)
			
		case HEADER_EDIT_NPC_OUTFIT: //Edit Npc Appearance
			go c.ReceiveEditNpcAppearence(packet)
			
		case HEADER_EDIT_NPC_POSITION: //Edit Npc Position
			go c.ReceiveEditNpcPosition(packet)
			
		case HEADER_DELETE_NPC: //Delete Npc
			go c.ReceiveDeleteNpc(packet)
			
		case HEADER_GET_NPC_DATA: //Retreive NPC pokemon and Events
			go c.ReceiveGetNpcPokemonAndEvents(packet)	
			
		case HEADER_GET_NPC_EVENTS:
			go c.ReceiveNpcEvents(packet)
			
		case HEADER_SET_MUSIC:
			go c.ReceiveSaveMusic(packet)
			
		case HEADER_SET_POKECENTER:
			go c.ReceiveSavePokecenter(packet)
			
		case HEADER_SET_LOCATION:
			go c.ReceiveSaveLocation(packet)
			
		default:
			fmt.Printf("Unknown header: %d\n", header)
			
		}
	}
	fmt.Printf("Client disconnected: %d\n", c.id)
}

func (c *Client) ReceiveLogin(_packet *Packet) {
	username := _packet.ReadString()
	password := _packet.ReadString()
	ver := _packet.ReadString()
	if ver == version {
		if c.checkAccount(username, password) {
			fmt.Println("- Send login")
			c.loggedIn = true
			c.SendLogin(0)
			fmt.Println("- Send map list")
			c.SendMapList()
			fmt.Println("- Send npc list")
			c.SendNpcList()
			fmt.Println("- Send location list")
			c.SendLocationList()
			fmt.Println("- Send music list")
			c.SendMusicList()
			fmt.Println("- Send pokecenter list")
			c.SendPokecenterList()
		} else {
			fmt.Println("- Send login false")
			c.SendLogin(1)
		}
	} else {
		c.SendLogin(2)
	}
}

func (c *Client) checkAccount(_username string, _password string) bool {
	var accountModel models.MapchangeAccount
	err := g_orm.Where(fmt.Sprintf("%v = '%v'", models.MapchangeAccount_Username, _username)).Find(&accountModel)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	
	return c.passwordTest(_password, accountModel.Password)
}

func (c *Client) passwordTest(_plain string, _hash string) bool {
	var h hash.Hash = sha1.New()
	h.Write([]byte(_plain))

	var sha1Hash string = strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
	var original string = strings.ToUpper(_hash)

	return (sha1Hash == original)
}

func (c *Client) RequestMapPiece(_packet *Packet) {
	if c.loggedIn {
		x := int(_packet.ReadInt16())
		y := int(_packet.ReadInt16())
		z := int(_packet.ReadUint16())
		w := int(_packet.ReadUint16())
		h := int(_packet.ReadUint16())

		c.SendArea(x, y, z, w + x, h + y)
	}
}

func (c *Client) ReceiveAddMap(_packet *Packet) {
	mapName := _packet.ReadString()
	if len(mapName) > 0 {	
		mapEntity := models.Map{ Idmap: 0, 
								 Name: mapName }
		err := g_orm.Save(&mapEntity)
		if err != nil {
			log.Error("Client", "ReceiveAddMap", "Error adding map: %v", err.Error())
			return
		}
		g_map.AddMap(mapEntity.Idmap, mapName)	
		g_server.SendMapListUpdateToClients()
	}
}

func (c *Client) ReceiveRemoveMap(_packet *Packet) {
	mapId := int(_packet.ReadUint16())
	
	mapEntity := models.Map{}
	mapEntity.Idmap = mapId
	_, err := g_orm.Delete(&mapEntity)
	if err != nil {
		log.Error("Client", "ReceiveRemoveMap", "Error removing map: %v", err.Error())
		return
	}
	g_map.DeleteMap(mapId)
	g_server.SendMapListUpdateToClients()
}

func (c *Client) ReceiveTileEventUpdate(_packet *Packet) {
	x := int(_packet.ReadInt16())
	y := int(_packet.ReadInt16())
	z := int(_packet.ReadInt16())
	
	if tile, found := g_map.GetTileFromCoordinates(x, y, z); found {	
		eventType := int(_packet.ReadUint8())
		
		if eventType == TILEEVENT_NONE { // Remove event
			tile.RemoveEvent()
		} else if tile.Event != nil && tile.Event.GetEventType() == eventType { // Update
			tile.Event.UpdateFromPacket(_packet)
			tile.Event.Save()
		} else { // Add
			var newEvent ITileEvent = nil
			if eventType == TILEEVENT_WARP {
				newEvent = NewWarpFromPacket(_packet)
			}
			
			if newEvent != nil {
				tile.AddEvent(newEvent)
				//TODO Update the tile
			}
		}
	}
}

func (c *Client) ReceiveAddNpc(_packet *Packet) {
	npcName := _packet.ReadString()
	if len(npcName) > 0 {
		// New NPC object
		npc := NewNpc()
		npc.Name = npcName
				
		g_npc.AddNpc(npc)
		g_server.SendNpcToClients(npc.DbId)
	}
}

func (c *Client) ReceiveEditNpcAppearence(_packet *Packet) {
	npcId := int64(_packet.ReadUint16())
	npcName := _packet.ReadString()
	head := int(_packet.ReadUint16())
	nek := int(_packet.ReadUint16())
	upper := int(_packet.ReadUint16())
	lower := int(_packet.ReadUint16())
	feet := int(_packet.ReadUint16())
	
	if len(npcName) > 0 {
		if npc, ok := g_npc.GetNpcById(npcId); ok {
			npc.SetName(npcName)
			npc.SetOutfitPart(pul.OUTFIT_HEAD, head)
			npc.SetOutfitPart(pul.OUTFIT_NEK, nek)
			npc.SetOutfitPart(pul.OUTFIT_UPPER, upper)
			npc.SetOutfitPart(pul.OUTFIT_LOWER, lower)
			npc.SetOutfitPart(pul.OUTFIT_FEET, feet)
			
			if npc.Save() {			
				g_server.SendNpcToClients(npc.DbId)
			}
		}
	}
}

func (c *Client) ReceiveEditNpcPosition(_packet *Packet) {
	npcId := int64(_packet.ReadUint16())
	x := int(_packet.ReadUint16())
	y := int(_packet.ReadUint16())
	z := int(_packet.ReadUint16())
	
	if npc, ok := g_npc.GetNpcById(npcId); ok {
		npc.SetPositionByCoordinates(x, y, z)
		
		if npc.Save() {
			g_server.SendNpcToClients(npcId)
		}
	}		
}

func (c *Client) ReceiveDeleteNpc(_packet *Packet) {
	npcId := int64(_packet.ReadUint16())
		
	if npc, ok := g_npc.GetNpcById(npcId); ok {
		if npc.Delete() {
			g_server.SendDeleteNpcToClients(npcId)
		}
	}
}

func (c *Client) ReceiveGetNpcPokemonAndEvents(_packet *Packet) {
	npcId := int64(_packet.ReadUint16())
	
	c.SendNpcPokemon(npcId)
	c.SendNpcEvents(npcId)
}

func (c *Client) ReceiveNpcEvents(_packet *Packet) {
	npcId := int64(_packet.ReadUint16())
	events := _packet.ReadString()
	eventInitId := int(_packet.ReadUint16())


	if npc, ok := g_npc.GetNpcById(npcId); ok {
		npc.Events = events
		npc.EventInitId = eventInitId
		
		npc.Save()
	}
}

func (c *Client) ReceiveSaveMusic(_packet *Packet) {
	
}

func (c *Client) ReceiveSaveLocation(_packet *Packet) {
	
}

func (c *Client) ReceiveSavePokecenter(_packet *Packet) {

}


// //////////////////////////////////////////////
// SEND
// //////////////////////////////////////////////

func (c *Client) SendLogin(_status int) {
	packet := NewPacketExt(HEADER_SEND_LOGIN)
	packet.AddUint16(uint16(_status))
	if (_status == 2){
		packet.AddString(version)
	}
	c.Send(packet)
}

func (c *Client) SendArea(_x, _y, _z, _w, _h int) {

	packet := NewPacketExt(HEADER_SEND_TILE_AREA)
	packet.AddUint16(0)
	packet.AddUint16(uint16(_z))
	count := 0
	for x := _x; x < _w; x++ {
		for y := _y; y < _h; y++ {
			if packet.MsgSize > 8000 {
				packet.MsgSize -= 2
				packet.readPos = 3
				packet.AddUint16(uint16(count))
				c.Send(packet)
				
				packet = NewPacketExt(0x01)
				packet.AddUint16(0)
				packet.AddUint16(uint16(_z))
				count = 0
			}
			tile, found := g_map.GetTileFromCoordinates(x, y, _z)
			if found == true {
				count++

				packet.AddUint16(uint16(x))
				packet.AddUint16(uint16(y))
				packet.AddUint8(uint8(tile.Blocking))
				
				if tile.Event != nil {
					packet.AddUint8(uint8(tile.Event.GetEventType()))
					if tile.Event.GetEventType() == 1 {
						warp := tile.Event.(*Warp)
						packet.AddUint16(uint16(warp.destination.X))
						packet.AddUint16(uint16(warp.destination.Y))
						packet.AddUint16(uint16(warp.destination.Z))
					}
				} else {
					packet.AddUint8(0)
				}
				
				packet.AddUint8(uint8(len(tile.Layers)))
				for _, layer := range tile.Layers {
					if layer != nil {
						packet.AddUint8(uint8(layer.Layer))
						packet.AddUint16(uint16(layer.SpriteId))
					}
				}
			}
		}
	}
	
	packet.MsgSize -= 2
	packet.readPos = 3
	packet.AddUint16(uint16(count))
	c.Send(packet)
}

func (c *Client) SendMapList() {
	packet := NewPacketExt(HEADER_SEND_MAP_LIST)
	packet.AddUint16(uint16(len(g_map.mapNames)))
	
	for index, value := range(g_map.mapNames) {
		packet.AddUint16(uint16(index))
		packet.AddString(value)
	}
	
	c.Send(packet)
}

func (c *Client) SendNpcList() {
	packet := NewPacketExt(HEADER_SEND_NPC_LIST)
	packet.AddUint16(uint16(len(g_npc.Npcs)))
	
	for _id, npc := range(g_npc.Npcs) {
		packet.AddUint16(uint16(_id))
		packet.AddString(npc.Name)
		packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_HEAD)))
		packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_NEK)))
		packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_UPPER)))
		packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_LOWER)))
		packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_FEET)))
		packet.AddUint16(uint16(npc.Position.X))
		packet.AddUint16(uint16(npc.Position.Y))
		packet.AddUint16(uint16(npc.Position.Z))
	}
	
	c.Send(packet)
}

func (c *Client) SendNpc(_npcid int64) {
	packet := NewPacketExt(HEADER_SEND_NPC)
	
	npc, _ := g_npc.GetNpcById(_npcid)
	packet.AddUint16(uint16(_npcid))
	packet.AddString(npc.Name)
	packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_HEAD)))
	packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_NEK)))
	packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_UPPER)))
	packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_LOWER)))
	packet.AddUint16(uint16(npc.GetOutfitPart(pul.OUTFIT_FEET)))
	packet.AddUint16(uint16(npc.Position.X))
	packet.AddUint16(uint16(npc.Position.Y))
	packet.AddUint16(uint16(npc.Position.Z))
	
	c.Send(packet)
}

func (c *Client) SendDeleteNpc(_id int64) {
	packet := NewPacketExt(HEADER_SEND_DELETE_NPC)
	packet.AddUint16(uint16(_id))
	
	c.Send(packet)
}

func (c *Client) SendNpcPokemon(_npcid int64) {
	packet := NewPacketExt(HEADER_SEND_NPC_POKEMON)
	packet.AddUint16(uint16(_npcid))
	
	npc, _ := g_npc.GetNpcById(_npcid)
	packet.AddUint16(uint16(len(npc.Pokemons)))
	
	for _id, pokemon := range(npc.Pokemons) {
		packet.AddUint16(uint16(_id))
		packet.AddUint16(uint16(pokemon.pokId))
		packet.AddString(pokemon.Name)
		packet.AddUint16(uint16(pokemon.Hp))
		packet.AddUint16(uint16(pokemon.Att))
		packet.AddUint16(uint16(pokemon.Att_spec))
		packet.AddUint16(uint16(pokemon.Def))
		packet.AddUint16(uint16(pokemon.Def_spec))
		packet.AddUint16(uint16(pokemon.Speed))
		packet.AddUint16(uint16(pokemon.Gender))
		packet.AddUint16(uint16(pokemon.Held_item))
	}
	
	c.Send(packet)
}

func (c *Client) SendNpcEvents(_id int64) {
	npc, _ := g_npc.GetNpcById(_id)
	
	packet := NewPacketExt(HEADER_SEND_NPC_PANELS)
	packet.AddUint16(uint16(_id))
	packet.AddString(npc.Events)
	packet.AddUint16(uint16(npc.EventInitId))
	
	c.Send(packet)
}

func (c *Client) SendLocationList() {
	packet := NewPacketExt(HEADER_SEND_LOCATION_LIST)
	packet.AddUint16(uint16(len(g_locations.locations)))
	
	for index, location := range(g_locations.locations) {
		packet.AddUint16(uint16(index))
		packet.AddString(location.Name)
		packet.AddUint16(uint16(location.Idpokecenter))
		packet.AddUint16(uint16(location.Idmusic))
	}
	
	c.Send(packet)
}

func (c *Client) SendPokecenterList() {
	packet := NewPacketExt(HEADER_SEND_POKECENTER_LIST)
	packet.AddUint16(uint16(len(g_locations.locations)))
	
	for index, pokecenter := range(g_locations.pokecenters) {
		packet.AddUint16(uint16(index))
		position := pos.NewPositionFromHash(pokecenter.Position)
		packet.AddUint16(uint16(position.X))
		packet.AddUint16(uint16(position.Y))
		packet.AddUint16(uint16(position.Z))
		packet.AddString(pokecenter.Description)
	}
	
	c.Send(packet)
}

func (c *Client) SendMusicList() {
	packet := NewPacketExt(HEADER_SEND_MUSIC_LIST)
	packet.AddUint16(uint16(len(g_locations.musics)))
	
	for index, music := range(g_locations.musics) {
		packet.AddUint16(uint16(index))
		packet.AddString(music.Title)
		packet.AddString(music.Filename)
	}
	
	c.Send(packet)
}

func (c *Client) Send(_packet *Packet) {
	_packet.SetHeader()
	c.socket.Write(_packet.Buffer[0:_packet.MsgSize])
}