/*
Copyright © 2021 Marco Ostaska

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package squad

// graphQL most primitive data for squad resturns
type graphQL struct {
	Data data `json:"data"`
}

// squad's user information
type users struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

// squad's Assignees information for cards
type assignees struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// squad's cards information
type cards struct {
	Identifier     string      `json:"identifier"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	PrimaryLabels  []string    `json:"primaryLabels"`
	SecondaryLabel string      `json:"secondaryLabel"`
	DueAt          string      `json:"dueAt"`
	Swimlane       string      `json:"swimlane"`
	WorkstateType  string      `json:"workstateType"`
	Assignees      []assignees `json:"assignees"`
}

// squad's SwimlaneWorkstates information
type swimlaneWorkstates struct {
	Name string `json:"name"`
}

// Squad is an abstraction to squad
type squad struct {
	Name               string               `json:"name"`
	Users              []users              `json:"users"`
	Description        string               `json:"description"`
	Geography          string               `json:"geography"`
	SquadUsersCount    int                  `json:"squadUsersCount"`
	Cards              []cards              `json:"cards"`
	SwimlaneWorkstates []swimlaneWorkstates `json:"swimlaneWorkstates"`
}

// Data is the squad data
type data struct {
	Squad squad `json:"squad"`
}
