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

### Advanced

The `command` portion, that is, the "what do do on change" also obtains access to the `status`.

Example:
```
$ go run ./cmd/onchange https://www.dn.se/ cat
2024/06/20 17:06:27 INFO Watching subject=https://www.dn.se/ cmd=cat args=[]
2024/06/20 17:06:27 INFO HTTP initial status subject=https://www.dn.se/ status=200:0xdca7be6f0db81965337eaec99afd5ea50703ce53b151046f1f0e604a560b1071
2024/06/20 17:06:27 INFO Running cmd=/usr/bin/cat
200:0xdca7be6f0db81965337eaec99afd5ea50703ce53b151046f1f0e604a560b1071
```

Here, we're watching `https://www.dn.se/` for changes. Any time the website changes,
the command `cat` will be invoked. Now, in addition, the `status`, in this case
`200:0xdca7be6f0db81965337eaec99afd5ea50703ce53b151046f1f0e604a560b1071`, will be
written to the standard input, `stdin`, of the `cat` process.

Note: The format for the `status`:es of the various watchers are a bit up in the air, and subject
to change. Please don't build your house on them.
