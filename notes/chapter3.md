## Chapter 3 - Write a Log Package

- Logs can be used to store, share and process ordered data:
  - replicate databases (write-ahead log)
  - coordinate distributed services (consensus algorithms, e.g. Raft)
  - manage state in front-end applications (e.g. Redux)
- many distributed systems problems can be easier solved by breaking down system changes into single, atomic operations that can be stored, shared and processed with a log

### How Logs work
- A **log** is an append-only sequence of **records**:
  - text lines for humans to read
  - binary-encoded messages for other programs to read
  - orders records by time and indexes each record by its offset and time created
- **Segments**
  - a log is split into a list of segments
  - free up disk by deleting old segments
  - when the active segment is filled, create a new segment and make it the active segment
- Each segment comprises a **store** file and an **index** file
  - store file contain the record data - continually append records to this file
  - index file map each record's offset to its position in the store file - index files are small enough to be memory-mapped
- How to read a record given its offset:
  - get the entry from the index file -> position of the record in the store file
  - read the record at that position in the store file

### Build a Log
- Data Modeling
  - Record: the data stored in the log
  - Store: the file that records are stored in
  - Index: the file that index entries are stored in
  - Segment: the abstraction that ties a store and an index together
  - Log: the abstraction that ties all the segments together
- Write to a buffered writer instead of directly to the file to reduce the number of system calls and improve performance
- use `*Width` constants to specify the number of bytes that make up each index entry
  - an index file comprises a persisted file and a memory-mapped file
- The `segment` type wraps the `index` and `store` types to coordinate operations
  - when the log appends a record to the active segment, the segment needs to write the data to its store and add a new entry in the index
  - for reads, the segment needs to look up the entry from the index and then fetch the data from the store
- The log consists of a list of segments and a pointer to the active segment to append writes to
  - The directory is where the segments are stored
- use a RWMutex to grant access to reads when there isn't a write holding the lock
- expose methods which return offset range, to determine which nodes have oldest and newest data, and which need replication
- truncate segments due to limited disk space
- expose method `Reader()` to read the whole log to enable coordinated consensus, support snapshots and to restore a log
