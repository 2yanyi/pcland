// Source code file, created by Developer@YANYINGSONG.

package connect

import (
	"fmt"
	"net/http"
	"r2/pkg/generic2/chars"
	"r2/pkg/generic2/chars/cat"
	"sort"
)

var AccessUA = make(map[string]error)

func uaAuth(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	if _, ok := AccessUA[ua]; !ok {
		fmt.Printf("%s\n", ua)
		return false
	}

	return true
}

func CallHandler(w http.ResponseWriter, r *http.Request) {
	if !uaAuth(r) {
		return
	}

	sign := r.URL.Query().Get("f")
	addr := r.URL.Query().Get("a")
	conn, ok := ConnectionPool[addr]
	if !ok {
		_, _ = w.Write(append([]byte(htmlHeader), []byte(`<h1>操作失败</h1>`)...))
		return
	}

	_, _, err := Power.Call(conn.Conn, []byte(sign), nil)
	if err != nil {
		// _, _ = w.Write([]byte(err.Error()))
		// return
	}

	_, _ = w.Write(append([]byte(htmlHeader), []byte(`<h1>操作成功</h1>`)...))
}

func DeviceList(w http.ResponseWriter, r *http.Request) {
	if !uaAuth(r) {
		return
	}

	slice := make([]*DeviceInfo, 0, len(ConnectionPool))
	for _, info := range ConnectionPool {
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
		item := fmt.Sprintf(`
<div class="ls-item">
    <div class="ls-title">%s %s</div>
    <pre><a href="" onclick="return call('/call?f=client.restart&a=%s')">重啟</a></pre>
    <pre>%s_ %s %s</pre>
</div>
<hr>
`,
			icon, info.Model, info.Tid,
			cat.SizeFormat(float64(info.Power)), info.Target,
			info.Addr,
		)
		ls.WriteString(item)
	}

	_, _ = w.Write(ls.Bytes())
}
