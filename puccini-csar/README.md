puccini-csar
============

Tool for creating
[CSAR files](https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html#_Toc302251718).

Supported archive formats are tarballs (`.tar.gz` and `.tar`) as well as zip files (`.zip` or the
`.csar` alias).

Note that tarballs have the advantage that they can be streamed (e.g. from a HTTP URL) whereas using
the zip format would require `puccini-tosca` to first download the entire archive to the system's
temporary directory.

`create`
--------

Creates a CSAR file in the filesystem.

By default `puccini-csar` will select the archive format based on its extension, but this can be
forced via `--archive-format/-a`.

The compression level can be changed via `--compression/-c` to values from 0 (no compression) to
9 (maximum and slowest to compress/decompress). The default is 6.

If the directory already includes a `TOSCA-Metadata/TOSCA.meta` file then it will be validated and
used. Otherwise, `puccini-csar` will generate it for you. If there is only one `.yaml` file in the
root, then it will be used as the meta's `Entry-Definitions`. If there is more than one, then the
tool will emit an error unless you specify it explicitly via `--entry-definitions`. All fields in
the generated meta can be controlled via the following switches:

* `--tosca-meta-file-version`
* `--csar-version`
* `--created-by`
* `--entry-definitions` (this can be only be used once)
* `--other-definitions` (this can be repeated multiple times; order matters)

`meta`
------

Parses, validates, and extracts a CSAR's `TOSCA-Metadata/TOSCA.meta` information. Local paths as
well as URLs can be used as the argument.

The default format for output is YAML, but you can select JSON, XML, CBOR, or MessagePack instead with
`--format/-f`.
