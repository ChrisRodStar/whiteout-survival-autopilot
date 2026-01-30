package main

import (
	"path/filepath"

	"github.com/batazor/whiteout-survival-autopilot/internal/gift"
)

func main() {
	gift.RunRedeemer(gift.RedeemConfig{
		DevicesYAML: filepath.Join("db", "devices.yaml"),
		CodesYAML:   filepath.Join("db", "giftCodes.yaml"),
		// PythonDir: ""  // empty => redeem_code.py from internal/discordgift
	})
}
