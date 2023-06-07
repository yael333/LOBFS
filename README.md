# LOBFS (Library of Babel File System)

LOBFS is a filesystem implementation that brings Jorge Luis Borges' Library of Babel to your local machine. With LOBFS, you can explore the endless possibilities of the Library of Babel as if it were a file system on your computer. It's entirely written in Golang from scratch, with design choices deviating from the mainstream Javascript implementation.

## Features
- Treats the Library of Babel as a hierarchical filesystem, allowing navigation through hexagon (hex), wall, shelf, volume, and page identifiers.
- Uses the FUSE (Filesystem in User Space) interface to integrate with your Windows, Linux, \*BSD and MacOS (untested) systems.
- Allows reading of any "page" in the library determenistically using normal OS-level API.

https://github.com/konata-chan404/LOBFS/assets/42537566/6fe72509-0dee-4f1c-ab1b-a490fda8e9f2

## Installation

To install LOBFS, you first need to clone the repository:

```
git clone https://github.com/konata-chan404/LOBFS.git
cd LOBFS
```

Then build the application:

```
go build
```

## Usage

After building, you can mount the filesystem:

```
./LOBFS /mount/point
```

This will mount the Library of Babel at the directory `/mount/point`. Now you can navigate through the library using your file browser (`File Explorer`, `ranger`) or standard commands like `cd` and `ls` .

## Contributing

LOBFS is an open source project, and contributions are welcome! This project started from mere curiosity, and there's a lot to expand and improve. Feel free to create a new post on [Issues](https://github.com/konata-chan404/LOBFS/issues) about anything and it will be addressed with uttermost care 💜

Here are some roadmap changes:
- Search through the library (implemented in `babel.go`, has to be integrated in a multiplatform way in FUSE)
- Set StatFS for directory and file sizes (again, multiplatform quirks)
- Performence opimization 

## License

LOBFS is licensed under the [MIT License](LICENSE).
