package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	plist "github.com/DHowett/go-plist"
)

type Keyboard struct {
	Name      string `plist:"Product"`
	VendorID  int    `plist:"VendorID"`
	ProductID int    `plist:"ProductID"`
}

var OSXModKeys = map[string]int{
	"none":         -1,
	"caps":         0,
	"shift_l":      1,
	"control_l":    2,
	"option_l":     3,
	"command_l":    4,
	"keypad_0":     5,
	"help":         6,
	"shift_r":      9,
	"control_r":    10,
	"option_r":     11,
	"command_r":    12,
	"kernel_panic": 16,
}

type kbModMap struct {
	Src int `plist:"HIDKeyboardModifierMappingSrc"`
	Dst int `plist:"HIDKeyboardModifierMappingDst"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("You must supply exactly one remap specification in the format '<Source>:<Dest>'. Available key names are:")
		for key_name, _ := range OSXModKeys {
			fmt.Println("  -", key_name)
		}
		os.Exit(1)
	}

	desiredRemap := parseInputRemap(os.Args[1])

	// fmt.Println(newDefaultsValueTest(desiredRemap))
	keyboards := getKeyboards()

	if len(keyboards) > 0 {
		fmt.Println("Found keyboards:")
		for i, keyboard := range keyboards {
			fmt.Printf(" %3d: %s (vendor: %d; product: %d):\n", i, keyboard.Name, keyboard.VendorID, keyboard.ProductID)
			remaps := getRemap(keyboard)
			// for _, m := range remaps {
			// 	fmt.Printf("      Current map: %d => %d\n", m.Src, m.Dst)
			// }
			remaps = append(remaps, desiredRemap)
			remaps = dedupRemaps(remaps)
			setRemap(keyboard, remaps)
			// for _, m := range remaps {
			// 	fmt.Printf("      New map: %d => %d\n", m.Src, m.Dst)
			// }

			fmt.Println("      remapped successfully")
		}
		fmt.Println("\nYou'll need to logout before changes will take effect.")
	} else {
		fmt.Println("Found no keyboards!")
	}
}

func getKeyboards() []Keyboard {
	cmd := exec.Command("ioreg", "-n", "IOHIDKeyboard", "-a", "-r")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return parseIoregOutput(out)
}

func parseInputRemap(str string) (m kbModMap) {
	for i, part := range strings.SplitN(str, ":", 2) {
		key_id, ok := OSXModKeys[part]
		if !ok {
			fmt.Printf("Unknown key name: %s\n", part)
			os.Exit(1)
		}
		if i == 0 {
			m.Src = key_id
		} else {
			m.Dst = key_id
		}
	}
	return
}

func parseIoregOutput(output []byte) (keyboards []Keyboard) {
	buf := bytes.NewReader(output)
	decoder := plist.NewDecoder(buf)
	keyboards = make([]Keyboard, 0)

	err := decoder.Decode(&keyboards)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func setRemap(k Keyboard, remaps []kbModMap) {
	cmd := exec.Command("defaults", "-currentHost", "write", "-g", defaultsPropName(k), newDefaultsValue(remaps))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func getRemap(k Keyboard) (maps []kbModMap) {
	maps = make([]kbModMap, 0)

	cmd := exec.Command("defaults", "-currentHost", "read", "-g", defaultsPropName(k))
	out, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); ok {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewReader(out)
	decoder := plist.NewDecoder(buf)

	err = decoder.Decode(&maps)
	if err != nil {
		log.Fatal(err)
	}

	return
}
func dedupRemaps(input []kbModMap) (output []kbModMap) {
	remapMap := make(map[int]kbModMap)
	for _, remap := range input {
		remapMap[remap.Src] = remap
	}
	output = make([]kbModMap, 0, len(remapMap))
	for _, remap := range remapMap {
		output = append(output, remap)
	}
	return
}

func defaultsPropName(k Keyboard) string {
	return fmt.Sprintf("com.apple.keyboard.modifiermapping.%d-%d-0", k.VendorID, k.ProductID)
}

func newDefaultsValue(remaps []kbModMap) string {
	buf := new(bytes.Buffer)
	encoder := plist.NewEncoder(buf)
	if err := encoder.Encode(remaps); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
func newDefaultsValueTest(m kbModMap) string {
	buf := new(bytes.Buffer)
	encoder := plist.NewEncoderForFormat(buf, plist.OpenStepFormat)
	if err := encoder.Encode(m); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
