# audioplayer
A audio play tool running on slave mode - Slave模式音频播放器

# Description
This tool is developed for playing audios.
It's running under slave mode,so any program written by any language can call it directly.

# Dependency
`libbass.so`
`libbassflac.so`

you can find them in this repo,please copy it to /lib

# Usage

``` shell
Usage of ./audioplayer:
  -l    is playing looply.
  -loop
        is playing looply.
  -n string
        set the name of player. (default "-1")
  -name string
        set the name of player. (default "-1")
  -noquit
        set quit or not when play is over.
  -p string
        play a music by filename (default "-1")
  -play string
        play a music by filename (default "-1")
  -pos int
        specify the position of music
  -r string
        recover a music by process name (default "-1")
  -recover string
        recover a music by process name (default "-1")
  -s string
        stop a music by process name (default "-1")
  -stop string
        stop a music by process name (default "-1")
        
        
```

# Example

### Play music.

You may also need to specify a process name for it by using `-n YourName`,then you can stop,pause,recover,and change the play position with the name you have set.

you can also not specify a process name for it,insteadly,the process name will be the timestamp when the process starts.

``` shell
./audioplayer -p test.flac
```

### Stop music

Stop a audio process with the name `YourName`.
``` shell
./audioplayer -s YourName
```

### Recover music

Continue playing the music with the position when last time stopped.
``` shell
./audioplayer -r YourName
```

### Change position

change the position to 10s.
``` shell
./audioplayer -pos 10 -n YourName
```

# Warning

This tool can only be used for `none-commercial use`,I's using the library of [Un4seen Bass](http://www.un4seen.com/),they only permmited non-commercial use for free. If you want to use this tool for commercial use, you need to buy a licence from them. I'm not responsible for illegal use of this tool.

Btw, I may write a audip library myself in the future,but for current convinience of use,I will keep bass library for a period of time.
