package snapman

import (
	"config"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"bufio"
	"strings"
	"time"
	"strconv"
	"sort"
)

func GetDatasets(conf *config.Config) (map[string]bool, error) {
	datasets := make(map[string]bool)
	cmd := exec.Command("zfs", "list", "-H", "-o", "name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		dataset := scanner.Text()
		dataset = strings.TrimSpace(dataset)
		if dataset == "" {
			continue
		}
		datasets[dataset] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return datasets, nil
}

func CreateSnapshot(dataset string, command string) {
	name := dataset + "@autosnap." + time.Now().Format("2006-01-02.15:04:05") + "." + command
	cmd := exec.Command("zfs", "snapshot", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "create shapshot: %v\n", err)
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(os.Stderr, "create shapshot: %s\n", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "create shapshot: %v\n", err)
	}
}

func DeleteSnapshot(snapshotName string) {
	cmd := exec.Command("zfs", "destroy", snapshotName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "delete shapshot: %v\n", err)
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(os.Stderr, "delete shapshot: %s\n", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "delete shapshot: %v\n", err)
	}
}

type Snapshot struct {
	SnapshotName string
	DatasetName  string
	CommandName  string
	CreationDate int64
}

type ByCreationDate []Snapshot

func (a ByCreationDate) Len() int {
	return len(a)
}
func (a ByCreationDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByCreationDate) Less(i, j int) bool {
	return a[i].CreationDate > a[j].CreationDate
}

type Snapshots map[string]map[string][]Snapshot
//   map[command] -> map[dataset] -> []Snapshot

func NewSnapshots() *Snapshots {
	snapshots := Snapshots(make(map[string]map[string][]Snapshot))
	return &snapshots
}

func (snapshots *Snapshots) AddSnapshot(conf *config.Config, snapshot Snapshot) {
	var command string
	if _, ok := conf.Interval[snapshot.CommandName]; ok {
		command = snapshot.CommandName
	} else {
		command = "clean"
	}
	data := (*snapshots)[command]
	if data == nil {
		data = make(map[string][]Snapshot)
		(*snapshots)[command] = data
	}
	slice := data[snapshot.DatasetName]
	if slice == nil {
		slice = make([]Snapshot, 0)
		data[snapshot.DatasetName] = slice
	}
	data[snapshot.DatasetName] = append(data[snapshot.DatasetName], snapshot)
}

func (snapshots *Snapshots) DeleteExpiredSnapshots(conf *config.Config) {
	for command := range *snapshots {
		switch command {
		case "clean":
			for dataset := range (*snapshots)[command] {
				for _, snapshot := range (*snapshots)[command][dataset] {
					DeleteSnapshot(snapshot.SnapshotName)
				}
			}
		default:
			for dataset := range (*snapshots)[command] {
				slice := (*snapshots)[command][dataset]
				sort.Sort(ByCreationDate(slice))
				leave := conf.Interval[command]
				if len(slice) > leave {
					slice = slice[leave:]
					for _, snapshot := range slice {
						DeleteSnapshot(snapshot.SnapshotName)
					}
				}
			}
		}
	}
}

func GetSnapshots(conf *config.Config) (*Snapshots, error) {
	snapshots := NewSnapshots()
	// zfs list -H -p -o name,creation -t snap
	cmd := exec.Command("zfs", "list", "-H", "-p", "-o", "name,creation", "-t", "snap")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "\t", 2)
		snapshotName, creation := split[0], split[1]
		snapshotName = strings.TrimSpace(snapshotName)
		creation = strings.TrimSpace(creation)
		secondSplit := strings.SplitN(snapshotName, "@", 2)
		datasetName, rest := secondSplit[0], secondSplit[1]
		if !strings.HasPrefix(rest, "autosnap") {
			continue
		}
		creationDate, err := strconv.ParseInt(creation, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("can't parse creation date '%s'", creation))
		}
		var commandName string
		pos := strings.LastIndex(snapshotName, ".")
		if pos == -1 {
			commandName = "unknown"
		} else {
			commandName = snapshotName[pos + 1:]
		}
		snapshot := Snapshot{
			SnapshotName: snapshotName,
			DatasetName: datasetName,
			CommandName: commandName,
			CreationDate: creationDate,
		}
		snapshots.AddSnapshot(conf, snapshot)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func Execute(conf *config.Config, command string) {
	datasets, err := GetDatasets(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: can't read datasets: %v\n", err)
		os.Exit(1)
	}
	if _, ok := conf.Interval[command]; ok {
		for dataset := range datasets {
			if conf.Included(dataset) {
				CreateSnapshot(dataset, command)
			}
		}
	}
	snapshots, err := GetSnapshots(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: can't read snapshots: %v\n", err)
		os.Exit(1)
	}
	snapshots.DeleteExpiredSnapshots(conf)
}
