# Geospatial Client Side Cache

This package provides a client-side caching mechanism for geospatial data. It supports three primary functionalities:

- Geospatial indexing for location-based queries.
- Timestamp indexing for time-based queries.
- In-memory key-value storage for caching locations.

## Key Components

- Geospatial Interface: Manages geospatial data, allowing addition, removal, and querying of nearby locations.
- Timestamp Interface: Manages timestamped data, allowing addition, removal, and querying of locations within a time range.
- Cache Interface: Provides in-memory storage with facilities for setting, getting, and deleting cache entries, and managing their TTL (Time-to-Live).

## Cache Struct

Combines geospatial, timestamp, and cache functionalities. It also includes cleaning mechanisms to remove outdated cache entries.
