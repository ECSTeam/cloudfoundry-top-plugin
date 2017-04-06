// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

//
// http://csbruce.com/software/utf-8.html
//
// Most of these do not work on Windows
//
// "▁▂▃▄▅▆▇█"
// "▁▂▃▄▅▆▇█
// "█▉▊▋▌▍▎▏"
//
// ▁ ▂ ▃ ▄ ▅ ▆ ▇ █   u2580
// ▁▂▃▄▅▆▇█   This one has two odd chars display on MacOS
// ▁▂▃▅▆▇  This one looks better on MacOS
//
// ◐◑◒◓◔◕   u25C0
// ☠  ⚡ ❌ ❗
// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors

const (
	CLEAR = "\033[0m"

	// Forground color
	BLACK  = "\033[30"
	RED    = "\033[31"
	GREEN  = "\033[32"
	YELLOW = "\033[33"
	BLUE   = "\033[34"
	PURPLE = "\033[35"
	CYAN   = "\033[36"
	WHITE  = "\033[37"

	// Background colors:
	// 40m - Black
	// 41m - Red
	// 42m - Green
	// 43m - Yellow
	// 44m - Blue
	// 45m - Purple
	// 46m - Cyan
	// 47m - White

	BRIGHT    = ";1m"
	DIM       = ";2m"
	UNDERLINE = ";4m"
	//FLASH = ";5m"
	REVERSE = ";7m"

	WHITE_TEXT_SOFT_BG   = "\x1b[48;5;235m\x1b[37m"
	RED_TEXT_GREEN_BG    = "\033[31m\033[42m"
	WHITE_TEXT_RED_BG    = "\033[37m\033[41m"
	WHITE_TEXT_YELLOW_BG = "\033[37m\033[43m"
	WHITE_TEXT_BLUE_BG   = "\033[37m\033[44m"

	BRIGHT_BLACK    = BLACK + BRIGHT
	DIM_BLACK       = BLACK + DIM
	UNDERLINE_BLACK = BLACK + UNDERLINE
	REVERSE_BLACK   = BLACK + REVERSE

	BRIGHT_RED    = RED + BRIGHT
	DIM_RED       = RED + DIM
	UNDERLINE_RED = RED + UNDERLINE
	REVERSE_RED   = RED + REVERSE

	BRIGHT_GREEN    = GREEN + BRIGHT
	DIM_GREEN       = GREEN + DIM
	UNDERLINE_GREEN = GREEN + UNDERLINE
	REVERSE_GREEN   = GREEN + REVERSE

	BRIGHT_YELLOW    = YELLOW + BRIGHT
	DIM_YELLOW       = YELLOW + DIM
	UNDERLINE_YELLOW = YELLOW + UNDERLINE
	REVERSE_YELLOW   = YELLOW + REVERSE

	BRIGHT_BLUE    = BLUE + BRIGHT
	DIM_BLUE       = BLUE + DIM
	UNDERLINE_BLUE = BLUE + UNDERLINE
	REVERSE_BLUE   = BLUE + REVERSE

	BRIGHT_PURPLE    = PURPLE + BRIGHT
	DIM_PURPLE       = PURPLE + DIM
	UNDERLINE_PURPLE = PURPLE + UNDERLINE
	REVERSE_PURPLE   = PURPLE + REVERSE

	BRIGHT_WHITE    = WHITE + BRIGHT
	DIM_WHITE       = WHITE + DIM
	UNDERLINE_WHITE = WHITE + UNDERLINE
	REVERSE_WHITE   = WHITE + REVERSE

	BRIGHT_CYAN    = CYAN + BRIGHT
	DIM_CYAN       = CYAN + DIM
	UNDERLINE_CYAN = CYAN + UNDERLINE
	REVERSE_CYAN   = CYAN + REVERSE
)
