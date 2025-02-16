# ChonkFS

FUSE file-system that split files in chunks and save them remotely.
Made for protocols such as torrent based systems. 

<p align="center">
<img src="./assets/chonker.png" alt="avatar" width="300"/>
</p>

<p align="center">
<i>Chonkers are cats, and all cats are liquids that fit anywhere, disregard of their chonkiness.</i>
</p>

<p align="center">
<i>Just like this file-system.</i>
</p>

*Based on the incredible [go-fuse](https://github.com/hanwen/go-fuse) library.*

## Architecture

The file-system is divided in three main components:

```mermaid
flowchart LR
    User --File Operations<br/>through FUSE--> FuseWrapper
    subgraph "ChonkFS"
        FuseWrapper[FuseWrapper<br><i>Abstract FS specifics</i>] --> Chonker[Chonker<br><i>Split files/Join chunks</i>]
        Chonker --Read/Store--> RAM

        subgraph "Storage"
            subgraph LayerRAM["Layer"]
                RAM 
                subgraph LayerDisk["Layer"]
                    LocalDisk[Local Disk]
                    FTP
                end
            end
        end
    end

    RAM --> LocalDisk
    LocalDisk --> FTP
```

## Storage

* Memory (RAM): Implemented
* Disk: Implemented
* FTP: TODO
* S3: TODO