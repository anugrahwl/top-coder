package animation

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func clearCMD() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func AnimateSentenceForward(str string) {
	for i := 0; i < len(str)-1; i++ {
		clearCMD()
		fmt.Print(str[:i+1])
		time.Sleep(time.Millisecond * 40)
	}
	clearCMD()
	fmt.Print(str)
}

func AnimateSentenceBackward(str string) {
	for i := len(str) - 1; i >= 0; i-- {
		time.Sleep(time.Millisecond * 30)
		clearCMD()
		fmt.Print(str[:i])
	}
}
