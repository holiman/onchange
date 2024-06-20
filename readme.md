## `onchange`

`onchange` is a utility to do things when files change. And not only files, it
can react to changes in ports, or webpages.

So, if you want to be notified (via the excellent [ntfy](https://ntfy.sh/)) if a
server shuts down, taking down the ssh port:
```
onchange tcp://myserver.my.domain:22 curl -d "Someting happened" ntfy.sh/mytopic
```

Or same thing for websites

```
onchange https://webpage.domain/index.html curl -d "Webpage updated" ntfy.sh/mytopic
```
Or, of course, for local files

```
onchange ./output.jsonl notify-send 'file was changed now'
```

Note, this utility does not do shell expansion, so things like pipes
does _not_ work:
```
onchange ./generated.testcase "echo 1 > a.txt"
```
However, the outputs from the command are forwarded by default, so
the command `stdout` goes to `stdout`, and vice versa for `stderr`. Therfore, this
works:
```
 ./onchange ./output.jsonl echo "this is one time" >> a.txt
```
To prevent the output from `onchange` to disturb the output, you can explicitly
redirect it thus:
```
./onchange -1=stdout.txt -2=stderr.txt ./output.jsonl echo "this is one time"
```
Or all output into one:
```
./onchange --output=output.txt ./output.jsonl echo "this is one time"
```