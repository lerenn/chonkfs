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


## Architecture

The file-system is divided in three main components:

```mermaid
flowchart LR
    User --File Operations<br/>through FUSE--> Wrapper
    subgraph "ChonkFS"
        Wrapper[Wrapper<br><i>Abstract FS specifics</i>] --> Chonker[Chonker<br><i>Split in chunks</i>]
        Chonker --Read/Store--> Storage

        subgraph "Storage"
            RAM --> LocalDisk[Local Disk] 
            LocalDisk --> Remote
        end
    end
```