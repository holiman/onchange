## `onchange`

`onchange` is a utility to do things when files change.

### Examples

```
onchange ./output.jsonl notify-send '`thing changed now'
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