// Source code file, created by Developer@YANYINGSONG.

package pages

import (
	"fmt"
	"library/console"
	"library/generic/chars"
	"net/http"
	"os"
	"path/filepath"
	"r/device"
	"r/server"
	"sort"
	"strings"
)

var AccessUA = make(map[string]error)

func uaAuth(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	if _, ok := AccessUA[ua]; !ok {
		console.INFO("%s", ua)
		return false
	}

	return true
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		help := `<!doctype html>
<pre>
Download link

# GNU/Linux

curl http://<HOST>:90/get?target=linux.386     -o PCland-@<HOST>.386
curl http://<HOST>:90/get?target=linux.amd64   -o PCland-@<HOST>.amd64
curl http://<HOST>:90/get?target=linux.amd64v1 -o PCland-@<HOST>.amd64v1
curl http://<HOST>:90/get?target=linux.arm64   -o PCland-@<HOST>.arm64
curl http://<HOST>:90/get?target=linux.loong64 -o PCland-@<HOST>.loong64
curl http://<HOST>:90/get?target=linux.riscv64 -o PCland-@<HOST>.riscv64

# Windows

curl http://<HOST>:90/get?target=windows.386     -o PCland-@<HOST>.exe
curl http://<HOST>:90/get?target=windows.amd64   -o PCland-@<HOST>.exe
curl http://<HOST>:90/get?target=windows.amd64v1 -o PCland-@<HOST>.exe
curl http://<HOST>:90/get?target=windows.arm64   -o PCland-@<HOST>.exe

# Apple

curl http://<HOST>:90/get?target=apple.amd64 -o PCland-@<HOST>.amd64
curl http://<HOST>:90/get?target=apple.arm64 -o PCland-@<HOST>.arm64
</pre>
`
		help = strings.ReplaceAll(help, "<HOST>", "211.149.130.119")
		_, _ = w.Write([]byte(help))
		return
	}

	fp := filepath.Join("bin", target)
	fmt.Printf("download %s\n", fp)

	if !chars.FileExist(fp) {
		_, _ = w.Write([]byte("target not found"))
		return
	}

	src, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", target))

	info, _ := src.Stat()
	http.ServeContent(w, r, target, info.ModTime(), src)
}

func CallHandler(w http.ResponseWriter, r *http.Request) {
	if !uaAuth(r) {
		return
	}

	sign := r.URL.Query().Get("f")
	addr := r.URL.Query().Get("a")
	conn, ok := device.Pool[addr]
	if !ok {
		_, _ = w.Write(append([]byte(htmlHeader), []byte(`<h1>操作失败</h1>`)...))
		return
	}

	_, err := server.Power.Call(conn.Conn, []byte(sign), nil)
	if err != nil {
		// _, _ = w.Write([]byte(err.Error()))
		// return
	}

	_, _ = w.Write(append([]byte(htmlHeader), []byte(`<h1>操作成功</h1>`)...))
}

func DeviceList(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%s\n", chars.ToJsonBytes(r.Header, "  "))
	if !uaAuth(r) {
		return
	}

	slice := make([]*device.DeviceInfo, 0, len(device.Pool))
	for _, info := range device.Pool {
		slice = append(slice, info)
	}
	sort.Slice(slice, func(i, j int) bool {
		if slice[i] == nil || slice[j] == nil {
			return false
		}
		return slice[i].Power > slice[j].Power
	})

	// Response to JSON.
	format := r.URL.Query().Get("format")
	if format == "json" {
		data := chars.ToJsonBytes(slice, "")
		_, _ = w.Write(data)
		return
	}

	// Response to HTML.
	icon := r.URL.Query().Get("icon")
	if icon == "" {
		icon = "🔥"
	}
	titleColor := r.URL.Query().Get("color")
	if titleColor == "" {
		titleColor = "f47983"
	}

	ls := chars.NewBuffer(1024)
	ls.WriteString(htmlHeader)
	ls.WriteString(`<h1>終端註冊面板</h1><hr>
<style>
	.ls-title { color: #` + titleColor + ` }
</style>
`)

	for _, info := range slice {

		// 保持所有的 target 长度一致
		n := 17 - len(info.Target)
		for i := 0; i < n; i++ {
			if i == 0 {
				info.Target += " "
			} else {
				info.Target += "."
			}
		}

		item := fmt.Sprintf(`
<div class="ls-item">
    <div class="ls-title">%s %s</div>
    <pre><a href="" onclick="return call('/call?f=client.restart&a=%s')">重啟</a></pre>
    <pre>%s_ %s %s</pre>
</div>
<hr>
`,
			icon, info.Model, info.Tid,
			chars.SizeFormat(float64(info.Power)), info.Target,
			info.Addr,
		)
		ls.WriteString(item)
	}

	_, _ = w.Write(ls.Bytes())
}
