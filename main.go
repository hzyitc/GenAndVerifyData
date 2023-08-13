package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"time"
)

var program = "GenAndVerifyData"
var version = "v0.0.0"

const displayInterval = 1024 * 1024 * 16
const blockSize = 4096

func getSize(path string) (uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pos, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	return uint64(pos), nil
}

func genData(buf []byte, i uint64) {
	binary.LittleEndian.PutUint64(buf, uint64(i))
	result := md5.Sum(buf[:64/8])
	copy(buf, result[:])

	l := 32 / 8

	for j := len(result); j < len(buf); j += l {
		binary.LittleEndian.PutUint32(buf[j:], uint32(j*123))

		for k := 0; k < l; k++ {
			buf[j+k] ^= buf[j-len(result)+k] + byte(k*5)
		}
	}
}

type Mode int

const (
	Mode_Write Mode = iota
	Mode_Verify
)

func do(mode Mode, path string, begin, end uint64) error {
	if begin%blockSize != 0 {
		return fmt.Errorf("begin(%d) should be aligned to %d", begin, blockSize)
	}

	if end%blockSize != 0 {
		return fmt.Errorf("end(%d) should be aligned to %d", end, blockSize)
	}

	buf := make([]byte, blockSize)
	buf2 := make([]byte, blockSize)

	fd, err := syscall.Open(path, syscall.O_RDWR, 0777)
	if err != nil {
		return err
	}

	pos, err := syscall.Seek(fd, int64(begin), io.SeekStart)
	if err != nil {
		return err
	}
	if uint64(pos) != begin {
		return fmt.Errorf("seek error: should be %d, but is %d", begin, pos)
	}

	start := time.Now()
	ot := time.Now()
	oi := end
	for i := begin; i < end; i += uint64(len(buf)) {
		genData(buf, i)

		var n int
		switch mode {
		case Mode_Write:
			n, err = syscall.Write(fd, buf)
		case Mode_Verify:
			n, err = syscall.Read(fd, buf2)
		default:
			return fmt.Errorf("unknown mode: %d", mode)
		}
		if err != nil {
			return err
		}
		if n != len(buf) {
			return fmt.Errorf("length error: should be %d, but is %d", len(buf), n)
		}
		if mode == Mode_Verify && !bytes.Equal(buf2, buf) {
			return fmt.Errorf("compare error from %d to %d", i, i+blockSize)
		}

		if i%displayInterval == 0 {
			t := time.Now()
			speed := float64(i-oi) / t.Sub(ot).Seconds()
			fmt.Printf("\rAt: %d, Progress: %.2f%%(%d/%d), Speed: %.2fMB/s ", i, (float32(i-begin) / float32(end-begin) * 100), (i - begin), (end - begin), (speed / 1024 / 1024))
			ot = t
			oi = i
		}
	}

	syscall.Close(fd)

	speed := float64(end-begin) / time.Since(start).Seconds()
	fmt.Printf("\nFinish. Avg: %.2fMB/s\n", (speed / 1024 / 1024))

	return nil
}

func feature() {
	fmt.Printf("%s can write and/or verify the file with some \"random\" data.\n", program)
	fmt.Printf("It could be used to test disk.\n")
}

func usage() {
	fmt.Printf("%s(%s)\n", program, version)
	fmt.Printf("Usage: \n")
	fmt.Printf("  %s [-write|-verify] {path} [begin [end]]\n", program)
	fmt.Printf("    -write    Write only\n")
	fmt.Printf("    -verify   Verify only\n")
	fmt.Printf("    path      Path\n")
	fmt.Printf("    begin     From 0. Align to %d. Include\n", blockSize)
	fmt.Printf("    end       From 0. Align to %d. Not include\n", blockSize)
}

func main() {
	write_only := flag.Bool("write", false, "Write only")
	verify_only := flag.Bool("verify", false, "Verify only")
	h := flag.Bool("h", false, "help")
	help := flag.Bool("help", false, "help")
	flag.Usage = usage
	flag.Parse()

	if *h || *help ||
		flag.NArg() < 1 ||
		flag.NArg() > 3 {
		feature()
		usage()
		return
	}

	path := flag.Arg(0)

	begin := uint64(0)
	if flag.NArg() >= 2 {
		var err error
		begin, err = strconv.ParseUint(flag.Arg(1), 10, 64)
		if err != nil {
			println(err.Error())
			println("")
			usage()
			return
		}
	}

	end := uint64(0)
	if flag.NArg() >= 3 {
		var err error
		end, err = strconv.ParseUint(flag.Arg(2), 10, 64)
		if err != nil {
			println(err.Error())
			println("")
			usage()
			return
		}
	}

	size, err := getSize(path)
	if err != nil {
		fmt.Printf("getSize error: %s\n", err.Error())
		return
	}

	if end == 0 {
		end = size
	}

	if begin >= size {
		fmt.Printf("Begin(%d) should less than \"%s\"'s size(%d)\n", begin, path, size)
		return
	}

	if end > size {
		fmt.Printf("End(%d) should less than or equal to \"%s\"'s size(%d)\n", end, path, size)
		return
	}

	if begin >= end {
		fmt.Printf("Begin(%d) should less than end(%d)\n", begin, end)
		return
	}

	if !*verify_only {
		var input string
		fmt.Printf("Will write to \"%s\"...\n", path)
		fmt.Printf("!!!!!!!!!!Danger!!!!!!!!!!\n")
		fmt.Printf("This will overwrite \"%s\"\n", path)
		fmt.Printf("Plaese input \"overwrite\" to make sure you known what you are doing: ")
		fmt.Scanln(&input)
		if input != "overwrite" {
			fmt.Printf("Exit.\n")
			return
		}

		err = do(Mode_Write, path, begin, end)
		if err != nil {
			fmt.Printf("Write error: %s\n", err.Error())
			return
		}
	}

	if !*write_only {
		fmt.Printf("Will verify \"%s\"...\n", path)
		err = do(Mode_Verify, path, begin, end)
		if err != nil {
			fmt.Printf("Verify error: %s\n", err.Error())
			return
		}
	}
}
