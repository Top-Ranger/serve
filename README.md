# serve

A simple program to share single files over the network. Inspired by [pshs](https://github.com/mgorny/pshs/)

# Usage

Assume you want to share 'dir/file' and 'dir/file2'. Just use these files as command-line arguments:

```
serve dir/file dir/file2
```

You will get an output similar to:

```
2017/12/28 23:25:37 Adding dir%2Ffile
2017/12/28 23:25:37 Adding dir%2Ffile2
2017/12/28 23:25:37 Server reachable at http://localhost:8080/
2017/12/28 23:25:37 Server publicly reachable at http://XXX.XXX.XXX.XXX:8080/
```

Your files are now accessible at http://XXX.XXX.XXX.XXX:8080/dir%2Ffile and http://XXX.XXX.XXX.XXX:8080/dir%2Ffile2

For more options like switching ports use

```
serve --help
```

# License

MIT. See LICENSE
