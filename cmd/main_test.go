package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestCLI(t *testing.T) {
	for _, test := range []struct {
		Name   string
		Args   []string
		Output string
	}{
		{
			Name: "Base",
			Args: []string{"./YadroTest", "../tests/test_data/test_base.txt"},
			Output: `09:00
08:48 1 client1
08:48 13 NotOpenYet
09:41 1 client1
09:48 1 client2
09:52 3 client1
09:52 13 ICanWaitNoLonger!
09:54 2 client1 1
10:25 2 client2 2
10:58 1 client3
10:59 2 client3 3
11:30 1 client4
11:35 2 client4 2
11:35 13 PlaceIsBusy
11:45 3 client4
12:33 4 client1
12:33 12 client4 1
12:43 4 client2
15:52 4 client4
19:00 11 client3
19:00
1 70 05:58
2 30 02:18
3 90 08:01
`,
		},
		{
			Name: "to_end",
			Args: []string{"./YadroTest", "../tests/test_data/test_to_end.txt"},
			Output: `09:00
08:48 1 client1
08:48 13 NotOpenYet
09:28 1 client1
09:41 1 client2
09:48 1 client3
09:52 2 client1 1
09:52 2 client2 2
09:52 2 client3 3
19:00 11 client1
19:00 11 client2
19:00 11 client3
19:00
1 100 09:08
2 100 09:08
3 100 09:08
`,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			os.Args = test.Args
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			main()

			outC := make(chan string)
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()

			w.Close()
			os.Stdout = old
			out := <-outC
			require.Equal(t, out, test.Output)
		})
	}
}
