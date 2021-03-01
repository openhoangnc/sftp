package sftp

import (
	"path/filepath"
	"runtime"
	"strings"
)

func (s *packetManager) prependRootDir(dir string) string {
	joinPath := filepath.Join(*s.rootDir, dir)
	if absPath, err := filepath.Abs(joinPath); err != nil || !strings.HasPrefix(absPath, *s.rootDir) {
		return *s.rootDir
	}
	return joinPath
}

func (s *packetManager) trimRootDir(dir string) string {
	if runtime.GOOS == "windows" && strings.HasPrefix(dir, "/") {
		dir = dir[1:]
		if absDir, err := filepath.Abs(dir); err == nil {
			dir = absDir
		}
	}

	if strings.HasPrefix(dir, *s.rootDir) {
		dir = dir[len(*s.rootDir):]
		if !strings.HasPrefix(dir, "/") {
			dir = "/" + dir
		}
	}
	return dir
}

func (s *packetManager) interceptRequest(pkt *orderedRequest) {
	if s.rootDir == nil {
		return
	}

	switch p := pkt.requestPacket.(type) {
	case *sshFxpStatPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpLstatPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpMkdirPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpRmdirPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpRemovePacket:
		p.Filename = s.prependRootDir(p.Filename)
	case *sshFxpRenamePacket:
		p.Oldpath = s.prependRootDir(p.Oldpath)
		p.Newpath = s.prependRootDir(p.Newpath)
	case *sshFxpSymlinkPacket:
		p.Targetpath = s.prependRootDir(p.Targetpath)
		p.Linkpath = s.prependRootDir(p.Linkpath)
	case *sshFxpReadlinkPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpRealpathPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpOpendirPacket:
		p.Path = s.prependRootDir(p.Path)
	case *sshFxpOpenPacket:
		p.Path = s.prependRootDir(p.Path)

	}
}

func (s *packetManager) interceptResponse(pkt *orderedResponse) {
	if s.rootDir == nil {
		return
	}

	switch p := pkt.responsePacket.(type) {
	case *sshFxpNamePacket:
		for _, a := range p.NameAttrs {
			a.Name = s.trimRootDir(a.Name)
			a.LongName = s.trimRootDir(a.LongName)
		}
	}
}
