# p
Plan 9 paging implementation.

Plan 9, the operating system created by Bell Labs in the mid 1980s, had a
utility called `p`, which was a simple pager. The modern `less` command is more
powerful but I thought it would be fun to implement `p` in Go. The `less` command
is 410 lines of code in its main file. See
[here](https://github.com/vbwagner/less/blob/master/main.c). The `p` command is
89 lines of code.

I have implemented the number of lines argument handling and the list of files
to be processed. One interesting aspect of `p` is to allow a command to be entered
following a page's output. It would be interesting to know why this was
important in the original. I have implemented this in Golang.

I found an issue I was facing earlier where processing stdin then getting paged
input would result in no paging for stdin input. This was because the stdin in
this case was from the parent process and not from the terminal. I had to
redefine where to get input from.

```go
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var r = bufio.NewReader(tty)
```

## Man page

[Found here}(https://9fans.github.io/plan9port/man/man1/p.html)

## Original C code

```c
#include <u.h>
#include <libc.h>
#include <bio.h>

#define	DEF	22	/* lines in chunk: 3*DEF == 66, #lines per nroff page */

Biobuf *cons;
Biobuf bout;

int pglen = DEF;

void printfile(int);

void
main(int argc, char *argv[])
{
	int n;
	int f;

	if((cons = Bopen("/dev/tty", OREAD)) == 0) {
		fprint(2, "p: can't open /dev/tty\n");
		exits("missing /dev/tty");
	}
	Binit(&bout, 1, OWRITE);
	n = 0;
	while(argc > 1) {
		--argc; argv++;
		if(*argv[0] == '-'){
			pglen = atoi(&argv[0][1]);
			if(pglen <= 0)
				pglen = DEF;
		} else {
			n++;
			f = open(argv[0], OREAD);
			if(f < 0){
				fprint(2, "p: can't open %s\n", argv[0]);
				continue;
			}
			printfile(f);
			close(f);
		}
	}
	if(n == 0)
		printfile(0);
	exits(0);
}

void
printfile(int f)
{
	int i, j, n;
	char *s, *cmd;
	Biobuf *b;

	b = malloc(sizeof(Biobuf));
	Binit(b, f, OREAD);
	for(;;){
		for(i=1; i <= pglen; i++) {
			s = Brdline(b, '\n');
			if(s == 0){
				n = Blinelen(b);
				if(n > 0)	/* line too long for Brdline */
					for(j=0; j<n; j++)
						Bputc(&bout, Bgetc(b));
				else{		/* true EOF */
					free(b);
					return;
				}
			}else{
				Bwrite(&bout, s, Blinelen(b)-1);
				if(i < pglen)
					Bwrite(&bout, "\n", 1);
			}
		}
		Bflush(&bout);
	    getcmd:
		cmd = Brdline(cons, '\n');
		if(cmd == 0 || *cmd == 'q')
			exits(0);
		cmd[Blinelen(cons)-1] = 0;
		if(*cmd == '!'){
			if(fork() == 0){
				dup(Bfildes(cons), 0);
				execl("/bin/rc", "rc", "-c", cmd+1, 0);
			}
			waitpid();
			goto getcmd;
		}
	}
}
```

## An example of using the bang option

```c
 $ p test/long.txt
```c
#include <u.h>
#include <libc.h>
#include <bio.h>

#define	DEF	22	/* lines in chunk: 3*DEF == 66, #lines per nroff page */

Biobuf *cons;
Biobuf bout;

int pglen = DEF;

void printfile(int);

void
main(int argc, char *argv[])
{
	int n;
	int f;

	if((cons = Bopen("/dev/tty", OREAD)) == 0) {
!ls -lah
output total 72
drwxr-xr-x  14 ian  staff   448B 26 Dec 14:24 .
drwxr-xr-x  60 ian  staff   1.9K 26 Dec 12:29 ..
drwxr-xr-x  15 ian  staff   480B 26 Dec 14:24 .git
-rw-r--r--   1 ian  staff     4B 26 Dec 14:17 .gitignore
-rw-r--r--   1 ian  staff   3.1K 26 Dec 14:03 .goreleaser.yaml
-rw-r--r--   1 ian  staff    11K 26 Dec 12:29 LICENSE
-rw-r--r--   1 ian  staff   2.1K 26 Dec 14:48 README.md
-rw-r--r--   1 ian  staff   797B 26 Dec 14:19 Taskfile.yaml
drwxr-xr-x   4 ian  staff   128B 26 Dec 12:53 cmd
drwxr-xr-x  14 ian  staff   448B 26 Dec 14:24 dist
-rw-r--r--   1 ian  staff   141B 26 Dec 12:50 go.mod
-rw-r--r--   1 ian  staff   1.4K 26 Dec 12:50 go.sum
drwxr-xr-x   4 ian  staff   128B 26 Dec 13:40 test
drwxr-xr-x   4 ian  staff   128B 26 Dec 14:24 vendor
```

## Lines of code

```
$ gocloc ./README.md cmd
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Markdown                         1             25              0            148
Go                               2             24             10            135
-------------------------------------------------------------------------------
TOTAL                            3             49             10            283
-------------------------------------------------------------------------------
```