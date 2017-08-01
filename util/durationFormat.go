// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

import (
	"time"
)

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0s.
func FormatDuration(d *time.Duration, showSecondsOver1Min bool) string {
	// Largest time is 2540400h10m10.000000000s

	if d == nil {
		return "<nil>"
	}

	var buf [32]byte
	w := len(buf)

	u := uint64(*d)
	neg := *d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			return "0s"
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u, false)
	} else {

		// Fraction of a second
		//w, u = fmtFrac(buf[:w], u, 9)
		u = u / 1000000000

		// u is now integer seconds
		sec := u
		u /= 60

		// Only includes seconds if less then 1 minutes
		if showSecondsOver1Min || u == 0 {

			w--
			buf[w] = 's'

			//w, u = fmtFrac(buf[:w], u, 9)

			leadingZero := u > 0
			w = fmtInt(buf[:w], sec%60, leadingZero)
		}

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			min := u
			u /= 60
			leadingZero := u > 0
			w = fmtInt(buf[:w], min%60, leadingZero)

			// u is now integer hours
			if u > 0 {
				w--
				buf[w] = 'h'
				hour := u
				u /= 24
				leadingZero := u > 0
				w = fmtInt(buf[:w], hour%24, leadingZero)
				// u is now integer days
				if u > 0 {
					w--
					buf[w] = ' '
					w--
					buf[w] = 'd'
					w = fmtInt(buf[:w], u, false)
				}
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64, twoDigit bool) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
		w--
		buf[w] = '0'
	} else {
		oldV := v
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
		if twoDigit && oldV < 10 {
			w--
			buf[w] = '0'
		}
	}
	return w
}
