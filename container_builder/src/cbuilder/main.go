package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

const (
    Repo = "https://github.com/karan/dockerfiles"
)

// Walks all dirs in current dir and returns relative paths
// to dirs with Dockerfile in them.
func DirsWithDockerfiles() ([]string, error) {
    dockerfileDirs := []string{}
    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.Name() == "Dockerfile" {
            dockerfileDirs = append(dockerfileDirs, filepath.Dir(path))
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    fmt.Printf("dirs with dockerfiles: %+v\n", dockerfileDirs)
    return dockerfileDirs, nil
}

// Get SHA for HEAD
func headSha() (string, error) {
    cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
    var out bytes.Buffer
    defer out.Reset()
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to get HEAD SHA, %+v", err)
    }
    s := out.String()
    s = strings.TrimSpace(s)
    fmt.Printf("HEAD SHA: %s\n", s)
    return s, nil
}

// Get SHA at remote's tip
func upstreamSha() (string, error) {
    cmd := exec.Command("git", "fetch", Repo, "refs/heads/master")
    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to fetch repo, %+v", err)
    }

    cmd = exec.Command("git", "rev-parse", "--verify", "FETCH_HEAD")
    var out bytes.Buffer
    defer out.Reset()
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to get upstream SHA, %+v", err)
    }
    s := out.String()
    s = strings.TrimSpace(s)
    fmt.Printf("upstream SHA: %s\n", s)
    return s, nil
}

// Returns `git diff --name-only` output based on passed SHAs.
func gitDiff(upstream, head string) (string, error) {
    // If HEAD and upstream are the same, get uncommitted changes.
    diffRange := "FETCH_HEAD~"
    // If head and upstream are different, we want the commit diff.
    if upstream != head {
        diffRange = fmt.Sprintf("%s...$s", upstream, head)
    }
    fmt.Printf("will git diff %s\n", diffRange)

    cmd := exec.Command("git", "diff", diffRange, "--name-only", "--", ".")
    var out, errbuf bytes.Buffer
    defer out.Reset()
    defer errbuf.Reset()
    cmd.Stdout = &out
    cmd.Stderr = &errbuf
    err := cmd.Run()
    if err != nil {
        fmt.Printf("failed to git diff: %s", errbuf.String())
        return "", fmt.Errorf("failed to get diff, %+v", err)
    }
    return out.String(), nil
}

// From the passed list of paths, return a set of directories excluding $PWD.
func getOnlyDirs(paths []string) ([]string) {
    // Dir name to exists
    dirs := map[string]bool{}

    for _, path := range paths {
        dirName := filepath.Dir(path)
        if dirName != "." {
            dirs[dirName] = true
        }
    }

    dirNames := []string{}
    for k, _ := range dirs {
        dirNames = append(dirNames, k)
    }
    return dirNames
}

// Returns a list of relative paths of all dirs with git changes in them.
func DirsWithChanges() ([]string, error) {
    head, err := headSha()
    if err != nil {
        return nil, err
    }

    upstream, err := upstreamSha()
    if err != nil {
        return nil, err
    }

    changedPathsRaw, err := gitDiff(upstream, head)
    if err != nil {
        return nil, err
    }
    fmt.Printf("changed paths raw = %+v\n", changedPathsRaw)
    changedPaths := strings.Split(changedPathsRaw, "\n")
    changedDirs := getOnlyDirs(changedPaths)
    fmt.Printf("changed dirs = %+v\n", changedDirs)
    return changedDirs, nil
}

func main() {
    // Get all dirs with Dockerfile in it
    allDirs, err := DirsWithDockerfiles()
    if err != nil {
        fmt.Errorf("failed to get all dirs with Dockerfiles, %+v\n", err)
        os.Exit(1)
    }

    changedDirs, err := DirsWithChanges()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // To make go build happy
    if allDirs != nil && changedDirs != nil {}

    // For each dir in changed set, if it's in all dockerfiles set, build and
    // push it.


}
