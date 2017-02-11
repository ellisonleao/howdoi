howdoi
======

This is a work in progress Go port of the awesome Python [howdoi](https://github.com/gleitz/howdoi) lib

## Installing

```
go get -u github.com/ellisonleao/howdoi
```

## Usage

```
howdoi [-h|--help] [-p|--pos POS] [-a|--all] [-l|--link] [-n|--num-answers NUM_ANSWERS] [-v|--version]
              [QUERY [QUERY ...]]

instant coding answers via the command line

positional arguments:
  QUERY                 the question to answer

optional arguments:
  -h, --help            show this help message and exit
  -p POS, --pos POS     select answer in specified position (default: 1)
  -a, --all             display the full text of the answer
  -l, --link            display only the answer link
  -n NUM_ANSWERS, --num-answers NUM_ANSWERS
                        number of answers to return
  -v, --version         displays the current version of howdoi`
```
