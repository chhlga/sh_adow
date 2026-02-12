package shadow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Version struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
	Notes     string    `json:"notes"`
	Size      int64     `json:"size"`
	Hash      string    `json:"hash"`
}

type FileEntry struct {
	Path     string    `json:"path"`
	Versions []Version `json:"versions"`
}

type List struct {
	Files []FileEntry `json:"files"`
}

func LoadList(shadowPath string) (*List, error) {
	listPath := filepath.Join(shadowPath, "list.json")

	data, err := os.ReadFile(listPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &List{Files: []FileEntry{}}, nil
		}
		return nil, err
	}

	var list List
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (l *List) Save(shadowPath string) error {
	listPath := filepath.Join(shadowPath, "list.json")
	tmpPath := listPath + ".tmp"

	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, listPath)
}

func (l *List) FindFile(path string) *FileEntry {
	for i := range l.Files {
		if l.Files[i].Path == path {
			return &l.Files[i]
		}
	}
	return nil
}

func (l *List) AddVersion(path string, version Version) {
	for i := range l.Files {
		if l.Files[i].Path == path {
			l.Files[i].Versions = append([]Version{version}, l.Files[i].Versions...)
			return
		}
	}

	l.Files = append(l.Files, FileEntry{
		Path:     path,
		Versions: []Version{version},
	})
}

func (l *List) RemoveVersion(path string, versionID string) bool {
	for i := range l.Files {
		if l.Files[i].Path == path {
			for j, v := range l.Files[i].Versions {
				if v.ID == versionID {
					l.Files[i].Versions = append(l.Files[i].Versions[:j], l.Files[i].Versions[j+1:]...)

					if len(l.Files[i].Versions) == 0 {
						l.Files = append(l.Files[:i], l.Files[i+1:]...)
					}
					return true
				}
			}
		}
	}
	return false
}

func GenerateVersionID(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:4])
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func HashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
