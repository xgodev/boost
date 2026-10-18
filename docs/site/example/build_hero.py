"""Writes the hero code window of docs/site/mockup.html from the shop app.

Usage, from docs/site: python3 example/build_hero.py example/shop mockup.html

It replaces what sits between the hero:start and hero:end comments of the page.
Run go generate in the shop first, so the generated tab shows the real output.
"""
import html
import re
import sys

FX = '<a class="ref" href="https://uber-go.github.io/fx/">fx</a>'
KW = set("package import func return type struct interface var const if else for range defer go switch case default map chan nil true false select break continue".split())
TYPES = set("error string int bool any Context Repository Service Handler Options Plugin Mutex Option Server ErrorWrapper Once M datastore service newDecorator".split())

TOKEN = re.compile(r"(?P<com>//[^\n]*)|(?P<raw>`[^`]*`)|(?P<str>\"(?:\\.|[^\"\\])*\")|(?P<num>\b\d+\b)|(?P<id>[A-Za-z_]\w*)|(?P<ws>\s+)|(?P<op>.)", re.S)


def comment(text):
    if text.startswith("// Code generated"):
        return f'<span class="gen-banner">{html.escape(text)}</span>'
    out = []
    for part in re.split(r"(@\w+)", text):
        if not part:
            continue
        cls = "a" if part.startswith("@") else "c"
        out.append(f'<span class="{cls}">{html.escape(part)}</span>')
    return "".join(out)


def highlight(src):
    toks = [(m.lastgroup, m.group()) for m in TOKEN.finditer(src)]
    out = []
    i = 0
    expect_name = False  # after "func" (and its receiver), the next identifier is the function name
    depth = 0
    while i < len(toks):
        kind, text = toks[i]
        nxt = lambda k: next(((kk, tt) for kk, tt in toks[i + k:] if kk != "ws"), (None, None))
        if kind == "com":
            out.append(comment(text))
        elif kind in ("str", "raw"):
            out.append(f'<span class="st">{html.escape(text)}</span>')
        elif kind == "num":
            out.append(f'<span class="nu">{text}</span>')
        elif kind == "id":
            if text in KW:
                out.append(f'<span class="kw">{text}</span>')
                if text == "func":
                    k, t = nxt(1)
                    expect_name = k == "id" or t == "("
                    # anonymous func: "func(" directly
                    if i + 1 < len(toks) and toks[i + 1][1] == "(":
                        expect_name = False
                    depth = 0
            elif expect_name and depth == 0:
                out.append(f'<span class="fn">{text}</span>')
                expect_name = False
            elif i + 2 < len(toks) and toks[i + 1][1] == "." and toks[i + 2][0] == "id":
                pkg, name = text, toks[i + 2][1]
                after = toks[i + 3][1] if i + 3 < len(toks) else ""
                pkg_html = FX if pkg == "fx" else pkg
                if (name in TYPES or (pkg, name) == ("mongo", "Collection")) and after != "(":
                    out.append(f'<span class="ty">{pkg_html}.{name}</span>')
                else:
                    out.append(f"{pkg_html}.{name}")
                i += 3
                continue
            elif text in TYPES and (i + 1 >= len(toks) or toks[i + 1][1] != "("):
                out.append(f'<span class="ty">{text}</span>')
            else:
                out.append(text)
        else:
            if expect_name:
                if text == "(":
                    depth += 1
                elif text == ")":
                    depth -= 1
            out.append(html.escape(text))
        i += 1
    return "".join(out).rstrip("\n")


# (indent, label, file path | None, hint, hint class)
SRC = [
    (0, "internal/app/api/rest/", None, "", ""),
    (1, "handler/reservestock/", None, "", ""),
    (2, "handler.go", "internal/app/api/rest/handler/reservestock/handler.go", "@REST", "a"),
    (1, "plugins.go", "internal/app/api/rest/plugins.go", "@Provide", "a"),
    (0, "internal/pkg/domain/", None, "", ""),
    (1, "service/reservestock/", None, "", ""),
    (2, "service.go", "internal/pkg/domain/service/reservestock/service.go", "interface", ""),
    (2, "local/", None, "", ""),
    (3, "service.go", "internal/pkg/domain/service/reservestock/local/service.go", "@Middleware", "a"),
    (1, "repository/stock/", None, "", ""),
    (2, "repository.go", "internal/pkg/domain/repository/stock/repository.go", "interface", ""),
    (0, "internal/pkg/datastore/stock/", None, "", ""),
    (1, "memory/", None, "", ""),
    (2, "datastore.go", "internal/pkg/datastore/stock/memory/datastore.go", "@Default", "a"),
    (1, "mongodb/", None, "", ""),
    (2, "datastore.go", "internal/pkg/datastore/stock/mongodb/datastore.go", "@Provide", "a"),
    (2, "config.go", "internal/pkg/datastore/stock/mongodb/config.go", "@Config", "a"),
    (0, "doc.go", "doc.go", "go:generate", ""),
]
MID = "gen/builder/middleware/github.com/acme/shop/internal_/pkg/domain/service/reservestock/local/new_middleware.go"
IMPL = "gen/implementation/github.com/acme/shop/internal_/pkg/domain/repository/stock/repository_implementation.go"
G = "gen/github.com/acme/shop/internal_/"
GEN = [
    (0, "cmd/api/", None, "", ""),
    (1, "main.go", "cmd/api/main.go", "main", ""),
    (0, "gen/", None, "", ""),
    (1, "module.go", "gen/module.go", "all modules", ""),
    (1, "builder/api/", None, "", ""),
    (2, "modules.go", "gen/builder/api/modules.go", "app api", ""),
    (2, "rest_module.go", "gen/builder/api/rest_module.go", "@REST", "a"),
    (1, "builder/middleware/…/local/", None, "", ""),
    (2, "new_middleware.go", MID, "@Middleware", "a"),
    (1, "implementation/…/stock/", None, "", ""),
    (2, "repository_implementation.go", IMPL, "@Default", "a"),
    (1, "github.com/acme/shop/internal_/", None, "", ""),
    (2, "app/api/rest/", None, "", ""),
    (3, "plugins_module.go", G + "app/api/rest/plugins_module.go", "", ""),
    (3, "handler/reservestock/", None, "", ""),
    (4, "new_module.go", G + "app/api/rest/handler/reservestock/new_module.go", "", ""),
    (2, "pkg/domain/service/…/local/", None, "", ""),
    (3, "new_module.go", G + "pkg/domain/service/reservestock/local/new_module.go", "", ""),
    (2, "pkg/datastore/stock/", None, "", ""),
    (3, "memory/new_module.go", G + "pkg/datastore/stock/memory/new_module.go", "", ""),
    (3, "mongodb/new_module.go", G + "pkg/datastore/stock/mongodb/new_module.go", "", ""),
    (3, "mongodb/options_module.go", G + "pkg/datastore/stock/mongodb/options_module.go", "@Config", "a"),
]

RUN = """<span class="pr">$</span> go generate ./...
<span class="dim">5:35PM</span> <span class="ok">INF</span> scanning ~/shop (module github.com/acme/shop)
<span class="dim">5:35PM</span> <span class="ok">INF</span> generated 6 inject modules and 1 apps in ~/shop/gen, mains in ~/shop/cmd

<span class="pr">$</span> go run ./cmd/api &amp;
<span class="dim">5:36PM</span> <span class="ok">INF</span> starting Echo Server. https://echo.labstack.com/

<span class="c"># memory is the @Default datastore: 3 t-shirts in stock</span>
<span class="pr">$</span> for i in 1 2 3 4; do curl -s -w <span class="st">' %{http_code}\\n'</span> -X POST localhost:8080/stock/tshirt/reserve; done
 <span class="ok">204</span>
 <span class="ok">204</span>
 <span class="ok">204</span>
{"httpStatusCode":409,"errorCode":"409","message":"tshirt is sold out"}
 <span class="nu">409</span>

<span class="c"># same binary on MongoDB, through boost's mongo factory:</span>
<span class="c"># INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=mongodb</span>"""


def slug(path):
    return re.sub(r"[^a-z0-9]+", "-", path.lower()).strip("-")


def group(shop, gid, label, rows, kind):
    files, panes, first = [], [], True
    for depth, name, path, hint, hcls in rows:
        if path is None:
            files.append(f'<div class="dir" role="none" style="--d:{depth}">{html.escape(name)}</div>')
            continue
        fid = slug(path)
        sel = "true" if first else "false"
        tab = "0" if first else "-1"
        h = f'<span class="hint{" " + hcls if hcls else ""}">{html.escape(hint)}</span>' if hint else ""
        files.append(f'<button class="file" type="button" role="tab" id="ft-{fid}" aria-selected="{sel}" aria-controls="fp-{fid}" tabindex="{tab}" style="--d:{depth}"><span class="fname">{html.escape(name)}</span>{h}</button>')
        code = highlight(open(f"{shop}/{path}").read())
        hidden = "" if first else " hidden"
        panes.append(f'<div class="pane" role="tabpanel" id="fp-{fid}" aria-labelledby="ft-{fid}"{hidden}>'
                     f'<div class="pane-path"><span class="k {kind}">{"generated" if gid == "gen" else "you write"}</span>{html.escape(path)}</div>'
                     f'<pre class="code" tabindex="0">{code}</pre></div>')
        first = False
    hidden = "" if gid == "src" else " hidden"
    return (f'<div class="ide" role="tabpanel" id="p-{gid}" aria-labelledby="t-{gid}"{hidden}>\n'
            f'<div class="files" role="tablist" aria-orientation="vertical" aria-label="{label}">\n' + "\n".join(files) + "\n</div>\n"
            f'<div class="panes">\n' + "\n".join(panes) + "\n</div>\n</div>")


START = "<!-- hero:start"
END = "<!-- hero:end -->"


def main(shop, page):
    nsrc = sum(1 for r in SRC if r[2])
    ngen = sum(1 for r in GEN if r[2])
    parts = [
        '<div class="win" id="hero-win">',
        '  <div class="win-bar" role="tablist" aria-label="The shop example: what you write, what go generate writes, and the run">',
        f'    <button class="tab" type="button" role="tab" id="t-src" aria-selected="true" aria-controls="p-src"><span class="k k-ann">you write</span><span class="tab-n">{nsrc} files</span></button>',
        f'    <button class="tab" type="button" role="tab" id="t-gen" aria-selected="false" aria-controls="p-gen" tabindex="-1"><span class="k k-gen">generated</span><span class="tab-n">{ngen} files</span></button>',
        '    <button class="tab" type="button" role="tab" id="t-run" aria-selected="false" aria-controls="p-run" tabindex="-1"><span class="k k-run">run</span><span class="tab-n">terminal</span></button>',
        '    <span class="win-step">github.com/acme/shop</span>',
        '  </div>',
        group(shop, "src", "Files you write", SRC, "k-ann"),
        group(shop, "gen", "Files go generate writes", GEN, "k-gen"),
        f'<div class="ide-run" role="tabpanel" id="p-run" aria-labelledby="t-run" hidden><pre class="code" tabindex="0">{RUN}</pre></div>',
        '</div>',
    ]
    text = open(page).read()
    head, rest = text.split(START, 1)
    marker, tail = rest.split("\n", 1)
    _, tail = tail.split(END, 1)
    with open(page, "w") as f:
        f.write(head + START + marker + "\n" + "\n".join(parts) + "\n" + END + tail)


if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])
