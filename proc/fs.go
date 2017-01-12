package proc

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
	"syscall"
)

//ProcMount where to find the information in os
const ProcMount = "/proc/mounts"

//ValidFSList List of valid fs format to export
var ValidFSList = []string{"ext2", "ext3", "ext4", "btrfs", "xfs", "gfs", "ntfs", "vfat"}

//AllFs struct to contain all fs fnformation
type AllFs struct {
	List map[string]FsInfo
}

//Init prepare list of Fs
func (fs *AllFs) Init() {
	sort.Strings(ValidFSList)
	fs.Update()
}

//Update refresh list of Fs
func (fs *AllFs) Update() {
	fs.List = readValidFs()
}

//FsInfo Information on Fs
type FsInfo struct {
	Dev   string
	Mount string
	Type  string
	Size  FsSize
}

//FsSize Information on space of Fs
type FsSize struct {
	BlockSize uint64
	Avail     uint64
	Used      uint64
	Total     uint64
}

func readFsSize(path string) FsSize {
	var stat syscall.Statfs_t
	_ = syscall.Statfs(path, &stat)
	fsSize := FsSize{}
	fsSize.BlockSize = uint64(stat.Bsize)
	fsSize.Total = uint64(stat.Blocks) * uint64(stat.Bsize)
	fsSize.Used = (uint64(stat.Blocks) - uint64(stat.Bfree)) * uint64(stat.Bsize)
	fsSize.Avail = uint64(stat.Bavail) * uint64(stat.Bsize)
	return fsSize
}

func readValidFs() map[string]FsInfo {
	all := map[string]FsInfo{}
	inFile, _ := os.Open(ProcMount)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		res := regexp.MustCompile(`[ \t]+`).Split(line, -1)
		fsMount := res[1]
		fsType := res[2]
		s := sort.SearchStrings(ValidFSList, fsType)
		if s < len(ValidFSList) && ValidFSList[s] == fsType && !strings.Contains(fsMount, "docker") {
			//all = append(all, FsInfo{Dev: res[0], Mount: res[1], Type: fsType, Size: readFsSize(fsMount)})
			all[strings.Replace(res[0], "/", "-", -1)+"@"+strings.Replace(res[1], "/", "-", -1)] = FsInfo{Dev: res[0], Mount: res[1], Type: fsType, Size: readFsSize(fsMount)}
		}
	}
	return all
}
