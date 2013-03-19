# ArticleDB

**ArticleDB is not a DataBase!**

ArticleDB is a Usenet storage, intended to be a NNTP-Server Backend. It stores the
data distributet to a set of files within the folder rather than a single file.
The Articles will be splittet into two parts: The MIME-Header and the body. While the
MIME-Header is stored within the Message-ID-Index, the Body is stored in Append-Only-Files.
ArticleDB compresses the body in order to reduce disk-space usage and io-Usage. This will also
Improve Performance.

## ArticleDB's Folder structure

First of all, ArticleDB uses a flat Structure instead of a deep hierarchical Structure.

```
basis-folder/
  index.json      # The index log, written in json
  index.dbm       # The index itself. Its a godbm file.
  group.list      # The Group index. Its a godbm file.
  group.assoc     # The Group-Number index. Its a godbm file.
  data._%d_       # where %d is replaced with an integer. This are the append only files for article-body storage.
```

## This software includes

It includes godbm (https://github.com/maxymania/godbm) wich is *repackaged* into this repository.

godbm - Copyright 2011 by Christoph Hack, BSD-Licensed

# License (MIT-License) execpt godbm

Copyright (c) 2013 maxymania

Permission is hereby granted, free of charge, to any person obtaining a copy of this software
and associated documentation files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial
portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
