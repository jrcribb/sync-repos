package main

import (
	"cui"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var GITHUB_USERNAME = "jrcribb"

type Repository struct {
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Fork        bool   `json:"fork"`
}

const version = "2.2"

//	func splitString(input string, chunkSize int) []string {
//		var chunks []string
//		for i := 0; i < len(input); i += chunkSize {
//			end := i + chunkSize
//			if end > len(input) {
//				end = len(input)
//			}
//			chunks = append(chunks, input[i:end])
//		}
//		return chunks
//	}
func main() {
	cui.ClearScreen()
	cui.Box(1, 1, 3, 80)
	cui.XyPrintf(2, 2, 78, "%ssync-repos %s%s", cui.INVERTED, version, cui.NORMAL)
	if len(os.Args) < 2 {
		cui.XyPrintf(4, 1, 0, "%sSyncroniza repositorios forkeados de Github", cui.NORMAL)
		cui.XyPrintf(5, 3, 0, "%sModo de uso\n", cui.NORMAL)
		cui.XyPrintf(6, 3, 0, "sync-repos (%susuario_github%s)\n", cui.ITALIC, cui.NORMAL)
		return
	}

	GITHUB_USERNAME = os.Args[1]
	cui.XyPrintf(4, 1, 0, "[ ] Actualizando repositorios de %s\n", GITHUB_USERNAME)
	cui.XyPrintf(4, 1, 0, "[%s■%s]", cui.BLINK, cui.NORMAL)
	page := 1
	morePages := true
	ok := 0
	nok := 0
	synced := false
	for morePages {
		urlGH := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d", GITHUB_USERNAME, page)
		resp, err := http.Get(urlGH)
		if err != nil {
			fmt.Println("Error al realizar la solicitud HTTP:", err)
			break
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error en la respuesta HTTP: %s\n", resp.Status)
			break
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error al leer la respuesta HTTP:", err)
			break
		}

		var repos []Repository
		if err := json.Unmarshal(respBytes, &repos); err != nil {
			fmt.Println("Error al analizar la respuesta JSON:", err)
			break
		}

		for _, repo := range repos {
			if repo.Fork {
				synced = syncRepository(GITHUB_USERNAME, repo.Name, repo.Description)
				if synced {
					ok += 1
				} else {
					nok += 1
				}
				cui.XyPrintf(4, 61, 20, "[√] %d [X] %d", ok, nok)
			}
		}
		if len(repos) > 0 {
			page++
		} else {
			morePages = false
		}
	}
	cui.ClearLines(5, 20)
	cui.XyPrintf(4, 1, 0, "[√] Repositorios de %s\n actualizados", GITHUB_USERNAME)
}

func removeWorkflows(user, repo string) {
	cui.ClearLines(7, 20)
	cui.XyPrintf(7, 3, 0, "[ ] Deshabilitando workflows de %s", repo)
	cui.XyPrintf(7, 3, 0, "[%s■%s]", cui.BLINK, cui.NORMAL)
	cmd := exec.Command("gh", "workflow", "-R", fmt.Sprintf("https://github.com/%s/%s", user, repo), "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		cui.XyPrintf(7, 3, 0, "[X]")
		cui.XyPrintf(8, 3, 0, "Error: Deshabilitación de workflow fallida: %s\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			id := line[len(line)-8:]
			cui.XyPrintf(8, 5, 0, "[ ] Eliminando workflow %s", id)
			cui.XyPrintf(8, 5, 0, "[%s■%s]", cui.BLINK, cui.NORMAL)
			cmd := exec.Command("gh", "workflow", "-R", fmt.Sprintf("https://github.com/%s/%s", user, repo), "disable", id)
			_, err := cmd.CombinedOutput()
			if err != nil {
				cui.XyPrintf(8, 5, 0, "[X]")
				cui.XyPrintf(9, 5, 0, "Error: Deshabilitación fallida de %s: %s\n", id, err)
			} else {
				cui.XyPrintf(8, 5, 0, "[√]")
				cui.XyPrintf(9, 5, 0, "Status: Deshabilitación correcta de %s\n", id)
			}
		}
	}
}

func syncRepository(user, repo, description string) bool {
	cui.ClearLines(5, 20)
	cui.XyPrintf(5, 3, 0, "[ ] Sincronizando %s%s%s", cui.INVERTED, repo, cui.NORMAL)
	cui.XyPrintf(5, 3, 0, "[%s■%s]", cui.BLINK, cui.NORMAL)
	// trozos := splitString(description, 78)
	if len(description) < 78 {
		cui.XyPrintf(6, 3, 78, "%s", description)
	} else {
		cui.XyPrintf(6, 3, 78, "%s...", description[0:75])
	}
	// for _, trozo := range trozos {
	// 	cui.XyPrintf(0, 3, 200, "%s", trozo)
	// }
	cmd := exec.Command("gh", "repo", "sync", "--force", fmt.Sprintf("%s/%s", user, repo))
	_, err := cmd.CombinedOutput()
	if err != nil {
		cui.XyPrintf(5, 3, 0, "[X]")
		cui.XyPrintf(7, 3, 200, "Error: Sincronización fallida: %s\n", err)
		return false
	} else {
		cui.XyPrintf(5, 3, 0, "[√]")
		removeWorkflows(user, repo)
		return true
	}
}
