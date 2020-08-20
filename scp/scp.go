package scp

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// the following code is a modified version of https://github.com/gnicod/goscplib
// which follows https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
//Constants

const (
	SCP_PUSH_BEGIN_FILE       = "C"
	SCP_PUSH_BEGIN_FOLDER     = "D"
	SCP_PUSH_BEGIN_END_FOLDER = "0"
	SCP_PUSH_END_FOLDER       = "E"
	SCP_PUSH_END              = "\x00"
)

type Scp struct {
	client *ssh.Client
}

func GetPerm(f *os.File) (perm string) {
	fileStat, _ := f.Stat()
	mod := fileStat.Mode()
	// if it's a directory there's high bits we want to ditch
	// only keep the low bits
	if mod > (1 << 9) {
		mod = mod % (1 << 9)
	}
	return fmt.Sprintf("%#o", uint32(mod))
}

//Initializer
func NewScp(clientConn *ssh.Client) *Scp {
	return &Scp{
		client: clientConn,
	}
}

//Pull file from remote to local
//targetFile 远端目标文件
//srcpath 本地目录
func (scp *Scp) PullFile(srcpath, targetFile string) error {
	session, err := scp.client.NewSession()
	if err != nil {
		log.Fatalln("Failed to create session: " + err.Error())
		return err
	}
	defer session.Close()

	go func() {
		iw, err := session.StdinPipe()
		if err != nil {
			log.Fatalln("Failed to create input pipe: " + err.Error())
		}
		or, err := session.StdoutPipe()
		if err != nil {
			log.Fatalln("Failed to create output pipe: " + err.Error())
		}
		fmt.Fprint(iw, "\x00")

		sr := bufio.NewReader(or)
		localFile := path.Join(srcpath, path.Base(targetFile))
		src, srcErr := os.Create(localFile)
		if srcErr != nil {
			log.Fatalln("Failed to create source file: " + srcErr.Error())
		}
		if controlString, ok := sr.ReadString('\n'); ok == nil && strings.HasPrefix(controlString, "C") {
			fmt.Fprint(iw, "\x00")
			controlParts := strings.Split(controlString, " ")
			size, _ := strconv.ParseInt(controlParts[1], 10, 64)
			/*
				bar := pb.New(int(size))
				bar.Units = pb.U_BYTES
				bar.ShowSpeed = true
				bar.Start()
				rp := io.MultiReader(sr, bar)
				if n, ok := io.CopyN(src, rp, size); ok != nil || n < size {
			*/
			if n, ok := io.CopyN(src, sr, size); ok != nil || n < size {
				fmt.Fprint(iw, "\x02")
				return
			}
			//			bar.Finish()
			sr.Read(make([]byte, 1))
		}
		fmt.Fprint(iw, "\x00")
	}()

	if err := session.Run(fmt.Sprintf("scp -f %s", targetFile)); err != nil {
		log.Fatalln("Failed to run: " + err.Error())
		return err
	}
	return nil
}

//Push one file to server
func (scp *Scp) PushFile(src string, dest string) error {
	session, err := scp.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fileSrc, srcErr := os.Open(src)
		defer fileSrc.Close()
		//fileStat, err := fileSrc.Stat()
		if srcErr != nil {
			log.Fatalln("Failed to open source file: " + srcErr.Error())
		}
		//Get file size
		srcStat, statErr := fileSrc.Stat()
		if statErr != nil {
			log.Fatalln("Failed to stat file: " + statErr.Error())
		}
		// Print the file content
		//fmt.Fprintln(w, SCP_PUSH_BEGIN_FILE+GetPerm(fileSrc), srcStat.Size(), filepath.Base(dest))
		fmt.Fprintln(w, SCP_PUSH_BEGIN_FILE+GetPerm(fileSrc), srcStat.Size(), filepath.Base(src))
		io.Copy(w, fileSrc)
		fmt.Fprint(w, SCP_PUSH_END)
	}()
	//if err := session.Run("/usr/bin/scp -rt " + filepath.Dir(dest)); err != nil {
	if err := session.Run("/usr/bin/scp -rt " + dest); err != nil {
		return err
	}
	return nil
}

//Push directory to server
func (scp *Scp) PushDir(src string, dest string) error {
	session, err := scp.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		//w = os.Stdout
		defer w.Close()
		folderSrc, _ := os.Open(src)
		fmt.Fprintln(w, SCP_PUSH_BEGIN_FOLDER+GetPerm(folderSrc), SCP_PUSH_BEGIN_END_FOLDER, filepath.Base(src))
		lsDir(w, src)
		fmt.Fprintln(w, SCP_PUSH_END_FOLDER)

	}()
	if err := session.Run("/usr/bin/scp -qrt " + dest); err != nil {
		return err
	}
	return nil
}

func prepareFile(w io.WriteCloser, src string) {
	fileSrc, srcErr := os.Open(src)
	defer fileSrc.Close()
	if srcErr != nil {
		log.Fatalln("Failed to open source file: " + srcErr.Error())
	}
	//Get file size
	srcStat, statErr := fileSrc.Stat()
	if statErr != nil {
		log.Fatalln("Failed to stat file: " + statErr.Error())
	}
	// Print the file content
	fmt.Fprintln(w, SCP_PUSH_BEGIN_FILE+GetPerm(fileSrc), srcStat.Size(), filepath.Base(src))
	io.Copy(w, fileSrc)
	fmt.Fprint(w, SCP_PUSH_END)
}

func lsDir(w io.WriteCloser, dir string) {
	fi, _ := ioutil.ReadDir(dir)
	//parcours des dossiers
	for _, f := range fi {
		if f.IsDir() {
			folderSrc, _ := os.Open(path.Join(dir, f.Name()))
			defer folderSrc.Close()
			fmt.Fprintln(w, SCP_PUSH_BEGIN_FOLDER+GetPerm(folderSrc), SCP_PUSH_BEGIN_END_FOLDER, f.Name())
			lsDir(w, path.Join(dir, f.Name()))
			fmt.Fprintln(w, SCP_PUSH_END_FOLDER)
		} else {
			prepareFile(w, path.Join(dir, f.Name()))
		}
	}
}
