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
package main

type Player struct {
	Creature // Inherit generic creature data
	
	Conn			*Connection	
	
	Location		*Location
	LastPokeCenter	*Tile
	
	Money			int
}

func NewPlayer(_name string) *Player {
	p := Player{ }
	p.uid 	= GenerateUniqueID()
	p.Conn 	= nil
	p.Outfit = NewOutfit()
	p.name = _name
	
	p.lastStep = PUSYS_TIME()
	p.moveSpeed = 280
	
	return &p
}

func (p *Player) GetType() int {
	return CTYPE_PLAYER
}

func (p *Player) SetConnection(_conn *Connection) {
	p.Conn = _conn
	p.Conn.Owner = p
	go _conn.HandleConnection()
}

func (p *Player) SetMoney(_money int) int {
	if p.Money += _money; p.Money < 0 {
		p.Money = 0
	}
	return p.Money
}

func (p *Player) OnCreatureMove(_creature ICreature, _from *Tile, _to *Tile, _teleport bool) {
	if _creature.GetUID() == p.GetUID() {
		p.lastStep = PUSYS_TIME()
	}
	
	canSeeFromTile	:= CanSeePosition(p.GetPosition(), _from.Position)
	canSeeToTile	:= CanSeePosition(p.GetPosition(), _to.Position)

	if canSeeFromTile && !canSeeToTile { 		// Leaving viewport
		p.sendCreatureMove(_creature, _from, _to)
		
		p.RemoveVisibleCreature(_creature)
		_creature.RemoveVisibleCreature(p)
	} else if canSeeToTile && !canSeeFromTile {	// Entering viewport
		p.AddVisibleCreature(_creature)
		_creature.RemoveVisibleCreature(p)
		
		p.sendCreatureMove(_creature, _from, _to)
	} else {									// Moving inside viewport
		p.AddVisibleCreature(_creature)
		_creature.AddVisibleCreature(p)
		
		p.sendCreatureMove(_creature, _from, _to)
	}
}

func (p *Player) OnCreatureTurn(_creature ICreature) {
	if _creature.GetUID() != p.GetUID() {
		p.sendCreatureTurn(_creature)
	}
}

func (p *Player) OnCreatureAppear(_creature ICreature, _isLogin bool) {
	canSeeCreature := CanSeeCreature(p, _creature)
	if !canSeeCreature {
		return
	}
	
	// We're checking inside the AddVisibleCreature method so no need to check here
	p.AddVisibleCreature(_creature)
	_creature.AddVisibleCreature(p)
}

func (p *Player) OnCreatureDisappear(_creature ICreature, _isLogout bool) {
	// TODO: Have to do something here with _isLogout
	
	p.RemoveVisibleCreature(_creature)
}

func (p *Player) AddVisibleCreature(_creature ICreature) {
	if _, found := p.VisibleCreatures[_creature.GetUID()]; !found {
		p.VisibleCreatures[_creature.GetUID()] = _creature
	}
}

func (p *Player) RemoveVisibleCreature(_creature ICreature) {
	// No need to check if the key actually exists because Go is awesome
	// http://golang.org/doc/effective_go.html#maps
	p.VisibleCreatures[_creature.GetUID()] = nil, false
}

// ------------------------------------------------------ //
func (p *Player) sendMapData(_dir int) {
	if p.Conn != nil {
		p.Conn.SendMapData(_dir, p.GetPosition())
	}
}

func (p *Player) sendCreatureMove(_creature ICreature, _from, _to *Tile) {
	if p.Conn != nil {
		p.Conn.SendCreatureMove(_creature, _from, _to)
	}
}

func (p *Player) sendCreatureTurn(_creature ICreature) {
	if p.Conn != nil {
		p.Conn.SendCreatureTurn(_creature)
	}
}

func (p *Player) sendCreatureAdd(_creature ICreature) {
	if p.Conn != nil {
		p.Conn.SendCreatureAdd(_creature)
	}
}

func (p *Player) sendCreatureRemove(_creature ICreature) {
	if p.Conn != nil {
		p.Conn.SendCreatureRemove(_creature)
	}
}

func (p *Player) sendPlayerWarp() {
	if p.Conn != nil {
		p.Conn.SendPlayerWarp(p.GetPosition())
	}
}
