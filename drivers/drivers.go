// Package drivers is a convenience package that registers all built-in
// database drivers. Import it with a blank identifier to make all drivers
// available:
//
//	import _ "github.com/nuln/dbase/drivers"
package drivers

import (
	_ "github.com/nuln/dbase/driver/bolt"
	_ "github.com/nuln/dbase/driver/gorm"
)
