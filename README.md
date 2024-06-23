# doit

Opinionated wrapper around `yt-dlp`.

## Why and what?

This is not a music downloader for everyone. It works for specific use cases only
(hence, _opinionated_). Let me explain.

I have a decently sized local library of music, and the library grows
**incrementally**. I organize the music in 2 ways:

- author/album/songs
  - this means that the structure looks something like this:
    ```
    author1/
      album1/
        song1.1
        sont1.2
        song1.3
      album2/
        song2.1
        song2.2
        song2.3
      singles/
        song1
        song2
    ```
- year/month/songs (this is the important one)
  - I organize the music by the time when I heard the song for the first time,
    or just when I listened to the song a bunch.
  - so the structure looks something like this:
    ```
    2023/
      January/
        song1
        song2
        song3
      February/
        song1
        song2
        song3
    ```

This means that for each month, I have a folder and inside that folder a file
named `links` in which I collect the music over a period of one full month, and
I occasionally run `doit` on that file.

The `links` file looks something like this:

```
Author - Song https://example.com
Author - Song 2 https://example2.com
Author 2 - Song A https://hi.com
-- and so on
```

You might think that it's just too much work writing every author and title by
hand, but that is **required** since there is no other way to guaranty the
authors and titles are correct and/or properly formatted. Plus, since English
input method is the only common input method on all my devices I listen to my music,
all names/titles need to be transliterated.

If your workflow is similar, `doit` is the program for you. Enjoy :)

## Is it any good?

Yes.

## Options

Available flags/options (with examples). All flags can be used with both single 
and double dash ( - and -- ).

```
index <number>             - Enables indexing and sets the starting index, or 
                             disables it if the index is negative. Default: 1
file <string>              - Set file to be used as input file. Default: "./links"
dir <string>               - Set/create directory to store downloads. Default: "./"
batch <number>             - Set number of parallel downloads. Default: 6
```

## Example usage

`./doit -index 1 -file the_prodigy.txt -dir "The Prodigy/" -batch 15`

or, if you put links for songs in file 'links' and those songs have their own
authors in that file (format: `author-song name https://example.com`) you can
also just us `./doit` and it will parse the links file and download the songs in
the current directory.

## Dependencies

Given the nature of the script (wrapper), `doit` depends on
[yt-dlp](https://github.com/yt-dlp/yt-dlp) (which, I believe, depends on `glibc`,
so as for now, `doit` won't work on `muslc` systems).

## TODO

- add options and levels of verbosity
  - better logging
- add tests for files
- add options to retry failed downloads 

Known bugs:
- it does not process additional arguments
- 

## Licence

This code is licensed under MIT licence (see LICENCE for details).
