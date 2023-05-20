// Source code file, created by Developer@YANYINGSONG.

package connect

const htmlHeader = `
<!doctype html>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<title>PC List</title>
<style>
    html {
        background-color: #000;
        color: #CCC;
    }
    a {
        color: #177cb0;
    }
    .ls-item {
        font-size: 12px;
    }
    .ls-title {
        font-size: 18px;
    }
    pre {
        display: inline;
    }
    hr {
        border: 0;
        padding: 3px;
        background: repeating-linear-gradient(135deg, #a2a9b6 0px, #a2a9b6 1px, transparent 1px, transparent 6px);
    }
</style>
<script>
    function call(url) {
        if (confirm('请再次确认你的操作')) {
            location.href = url
			return false
        }
    }
</script>
`
