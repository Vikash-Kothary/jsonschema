package render

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func build(root string) error {
	cmd := exec.Command("hugo")
	cmd.Dir = root
	return cmd.Run()
}

func RenderHTML(resume []byte, theme string) ([]byte, error) {
	sitePath, err := ioutil.TempDir(os.TempDir(), "resumic")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(sitePath)

	themePath := path.Join(sitePath, "themes")
	err = extractTheme(themePath, theme)
	if err != nil {
		return nil, err
	}

	config := map[string]string{
		"theme": theme,
	}
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}
	configPath := path.Join(sitePath, "config.json")
	err = ioutil.WriteFile(configPath, configJSON, 0600)
	if err != nil {
		return nil, err
	}

	dataPath := path.Join(sitePath, "data", "resume", "resume.json")
	err = os.MkdirAll(path.Dir(dataPath), 0700)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(dataPath, resume, 0600)
	if err != nil {
		return nil, err
	}

	contentPath := path.Join(sitePath, "content", "resume", "resume.md")
	err = os.MkdirAll(path.Dir(contentPath), 0700)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(contentPath, []byte{}, 0600)
	if err != nil {
		return nil, err
	}

	err = build(sitePath)
	if err != nil {
		return nil, err
	}

	htmlPath := path.Join(sitePath, "public", "resume", "resume", "index.html")
	return ioutil.ReadFile(htmlPath)
}