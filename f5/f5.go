// Copyright © 2023 Sébastien Gross <seb•ɑƬ•chezwam•ɖɵʈ•org>
//
// Created: 2022-04-20
// Last changed: 2024-10-11 21:23:48
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU Affero General Public License
// as published by the Free Software Foundation, either version 3 of
// the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program. If not, see
// <http://www.gnu.org/licenses/>.

package f5

type F5Metadata struct {
	Name    string `@( Ident | QF5Name | F5Name | QString ) "{"`
	Value   string `(  "value" @( Ident | QF5Name | F5Name | QString )`
	Persist string ` | "persist" @( "true" | "false" ) )* "}"`
}
