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

const (
	BATTLEEVENT_STOPBATTLE         = 999
	BATTLEEVENT_SLEEP              = 0
	BATTLEEVENT_TEXTID             = 1
	BATTLEEVENT_TEXT               = 2
	BATTLEEVENT_CHANGEHP           = 3
	BATTLEEVENT_ANIMATION          = 4
	BATTLEEVENT_ALLOWCONTROL       = 5
	BATTLEEVENT_CHANGEPOKEMON_SELF = 6
	BATTLEEVENT_CHANGEPOKEMON      = 7
	BATTLEEVENT_CHANGESELECTION    = 8
	BATTLEEVENT_CHANGEPP           = 9
	BATTLEEVENT_CHANGESTATUS       = 10
	BATTLEEVENT_CHANGELEVELSELF    = 11
	BATTLEEVENT_CHANGELEVEL        = 12
	BATTLEEVENT_CHANGEATTACK       = 13
	BATTLEEVENT_CHANGESCREEN       = 14
	BATTLEEVENT_DIALOGUE           = 15
	BATTLEEVENT_REMOVEPLAYER       = 16
	BATTLEEVENT_CHANGEEXP          = 17
)

const (
	BATTLECONTROL_NONE               = 0
	BATTLECONTROL_CHOOSEMOVE         = 1
	BATTLECONTROL_CHOOSEPOKEMON      = 2
	BATTLECONTROL_CHOOSEPOKEMON_ITEM = 3
	BATTLECONTROL_CHOOSEATTACK_ITEM  = 4
)

const (
	BATTLESELECT_MOVE    = 0
	BATTLESELECT_POKEMON = 1
)

const (
	BATTLE_IDLE      = 0
	BATTLE_RUNNING   = 1
	BATTLE_WAITING   = 2
	BATTLE_CHANGEHP  = 3
	BATTLE_CHANGEEXP = 4
)

const (
	PLAYER  = 0
	NPC     = 1
	POKEMON = 2
	SELF    = 3
)

const (
	USER     = 0
	OPPONENT = 1
)

const (
	ONE_ONE  = 1
	ONE_NPC  = 2
	ONE_WILD = 2
	TWO_TWO  = 4
	TWO_NPC  = 5
)

const (
	MOVE_NONE          = -1
	MOVE_ATTACK        = 0
	MOVE_ATTACK1       = 0
	MOVE_ATTACK2       = 1
	MOVE_ATTACK3       = 2
	MOVE_ATTACK4       = 3
	MOVE_POKEMON       = 4
	MOVE_BAG           = 5
	MOVE_RUN           = 6
	MOVE_CHANGEPOKEMON = 7
	MOVE_FAINT         = 8
	MOVE_ANSWER        = 9
	MOVE_FINISH        = 99
)

const (
	POKEMON_NONE = -1
	POKEMON1     = 0
	POKEMON2     = 1
	POKEMON3     = 2
	POKEMON4     = 3
	POKEMON5     = 4
	POKEMON6     = 5
)

const (
	BATTLESTATUS_CONTINUE      = 0
	BATTLESTATUS_ENDWON        = 1
	BATTLESTATUS_ENDLOST       = 2
	BATTLESTATUS_CHANGEPOKEMON = 3
)

const (
	BATTLESCREEN_MOVES   = 0
	BATTLESCREEN_YESNO   = 1
	BATTLESCREEN_ATTACKS = 2
)

const (
	BATTLEANSWER_NONE = 0
	BATTLEANSWER_YES  = 1
	BATTLEANSWER_NO   = 2
)

type IBattleEvent interface {
	Execute()
}
