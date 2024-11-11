package arts

import (
	"fmt"

	"ventana.com/pkg/config"
)

func DisplayWelcomeMessage() {
	CapyAscii := `
 __      __        _                    
 \ \    / /       | |                   
  \ \  / /__ _ __ | |_ __ _ _ __   __ _ 
   \ \/ / _ \ '_ \| __/ _' | '_ \ / _' |
    \  /  __/ | | | || (_| | | | | (_| |
     \/ \___|_| |_|\__\__,_|_| |_|\__,_|
                                        
	`

	fmt.Println(config.DarkBlue(CapyAscii))
}
