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
package network
//Datastructures to be sent between server and client

//Datastructures used to transfer tiles
type Tile struct {
	X        int
	Y        int
	Blocking int
	Layers   []*Layer
}

func NewTile() *Tile {
	return &Tile{Layers: make([]*Layer, 0)}
}

func (t *Tile) AddLayer(_index int, _sprite int) {
	layer := NewLayer()
	layer.Index = _index
	layer.Sprite = _sprite
	t.Layers = append(t.Layers, layer)
}

type Layer struct {
	Index  int
	ID     int
	Sprite int
}

func NewLayer() *Layer {
	return &Layer{}
}

//===============================================
// Server -> Client

//Tiles (HEADER_TILES)
type Data_Tiles struct {
	Tiles []*Tile
}

func NewData_Tiles() (msg *Message) {
	msg = NewMessage(HEADER_TILES)
	msg.Tiles = &Data_Tiles{Tiles: make([]*Tile, 0)}
	return
}

func (d *Data_Tiles) AddNewTile(_x int, _y int, _blocking int) {
	tile := NewTile()
	tile.X = _x
	tile.Y = _y
	tile.Blocking = _blocking
	d.Tiles = append(d.Tiles, tile)
}

func (d *Data_Tiles) AddTile(_tile *Tile) {
	d.Tiles = append(d.Tiles, _tile)
}
