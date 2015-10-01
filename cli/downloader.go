/*
 Copyright 2015 Runtime Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cli

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Downloader struct {
	Repos map[string]string
}

func NewDownloader() (*Downloader, error) {
	dl := &Downloader{}

	dl.Repos = map[string]string{}

	return dl, nil
}

func (dl *Downloader) gitClone(url string, branch string, dest string) error {
	_, err := ShellCommand(fmt.Sprintf("git clone -b %s %s %s", branch, url, dest))
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dest + "/.git/"); err != nil {
		return err
	}

	return nil
}

func (dl *Downloader) GetRepo(repoUrl string, branch string) (string, error) {
	// If repo already exists, return the temporary directory where it exists
	dir, ok := dl.Repos[repoUrl+branch]
	if ok {
		return dir, nil
	}

	dir, err := ioutil.TempDir("", "newtrepo")
	if err != nil {
		return "", err
	}

	// Otherwise, get a temporary directory and place the repo there.
	if err := dl.gitClone(repoUrl, branch, dir); err != nil {
		return "", nil
	}

	dl.Repos[repoUrl+branch] = dir

	return dir, nil
}

func (dl *Downloader) DownloadFile(repoUrl string, branch string,
	filePath string, destPath string) error {
	repoDir, err := dl.GetRepo(repoUrl, branch)
	if err != nil {
		return err
	}

	if err := CopyFile(repoDir+"/"+filePath, destPath); err != nil {
		return err
	}

	return nil
}